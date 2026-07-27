package simulator

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gmrtd/gmrtd/cryptoutils"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
)

// TcAc01Transceiver: PACE APDUs fail (SW or synthetic channel abort); BAC uses dynamic chip-side mutual auth.
//
// paceFailOn selects which PACE negotiation step fails:
//   - "mse_set_at" (default): MSE:Set AT itself fails immediately — PACE never begins key
//     agreement.
//   - "general_authenticate": MSE:Set AT succeeds (returns 9000), but the first General
//     Authenticate call (PACE step 1, encrypted-nonce exchange) fails instead. This models
//     a chip that accepts the PACE mechanism selection but fails partway through the
//     key-agreement handshake, rather than rejecting PACE outright.
//
// paceChannel selects how that failure is delivered at the APDU boundary:
//   - "" / "sw": return configured paceSW (honest synthetic chip).
//   - "timeout": brief delay then empty response (timed-exchange abstraction; not wall-clock RF).
//   - "no_response": empty frame (jammed / incomplete exchange abstraction).
//   - "transport_abort": empty response (hard transport failure abstraction).
//
// Channel modes keep BAC capable afterward so E = PACE-fail ∧ BAC-ok remains well-defined.
// They are deliberately distinct from returning a configured status word.
//
// This is a deliberate two-point simplification of the real four-exchange PACE state
// machine (MSE:Set AT, then General Authenticate steps 1-4). We do not simulate failure at
// each of the four GA sub-steps individually; "general_authenticate" fails at the first GA
// call encountered.
type TcAc01Transceiver struct {
	paceFailed  bool
	paceSW      []byte
	paceFailOn  string
	paceChannel string
	pass        *password.Password
	rndIcc      []byte
	// Post-BAC session material (set on successful EXTERNAL AUTHENTICATE).
	lastKIfd   []byte
	lastRndIfd []byte
	bacSession bool
}

func NewTcAc01Transceiver(paceStatusWord string, pass *password.Password) *TcAc01Transceiver {
	return NewTcAc01TransceiverWithInjection(paceStatusWord, "mse_set_at", pass)
}

func NewTcAc01TransceiverWithInjection(paceStatusWord string, paceFailOn string, pass *password.Password) *TcAc01Transceiver {
	return NewTcAc01TransceiverWithChannel(paceStatusWord, paceFailOn, "", pass)
}

func NewTcAc01TransceiverWithChannel(paceStatusWord, paceFailOn, paceChannel string, pass *password.Password) *TcAc01Transceiver {
	// Only "general_authenticate" selects the alternate injection point. Every other value —
	// including the empty string and the published blog profile's "first_pace_apdu" — maps to
	// the original, pinned behavior (fail at MSE:Set AT). This preserves byte-for-byte
	// reproduction of tag blog-b10-2026-07 while adding a genuinely distinct second point.
	if paceFailOn != "general_authenticate" {
		paceFailOn = "mse_set_at"
	}
	switch paceChannel {
	case "", "sw", "timeout", "no_response", "transport_abort":
	default:
		paceChannel = "sw"
	}
	if paceChannel == "" {
		paceChannel = "sw"
	}
	return &TcAc01Transceiver{
		paceSW:      utils.HexToBytes(paceStatusWord),
		paceFailOn:  paceFailOn,
		paceChannel: paceChannel,
		pass:        pass,
		rndIcc:      utils.HexToBytes("4608F91988702212"), // ICAO 9303 Part 11 D.3
	}
}

var _ iso7816.Transceiver = (*TcAc01Transceiver)(nil)

func (t *TcAc01Transceiver) deliverPaceFailure() []byte {
	switch t.paceChannel {
	case "timeout":
		// Symbolic delay: in-process APDU boundary has no real NFC deadline; documents
		// timed-exchange class without RF fidelity.
		time.Sleep(50 * time.Millisecond)
		return []byte{}
	case "no_response", "transport_abort":
		return []byte{}
	default:
		return append([]byte(nil), t.paceSW...)
	}
}

func (t *TcAc01Transceiver) Transceive(cla, ins, p1, p2 int, data []byte, le int, encodedData []byte) []byte {
	switch byte(ins) {
	case iso7816.INS_MANAGE_SE:
		if t.paceFailOn == "mse_set_at" && !t.paceFailed {
			t.paceFailed = true
			return t.deliverPaceFailure()
		}
		// paceFailOn == "general_authenticate": MSE:Set AT is accepted.
		return []byte{0x90, 0x00}
	case iso7816.INS_GENERAL_AUTHENTICATE:
		if t.paceFailOn == "general_authenticate" && !t.paceFailed {
			t.paceFailed = true
			return t.deliverPaceFailure()
		}
		if t.paceFailOn == "mse_set_at" {
			// MSE:Set AT already failed; PACE should not reach GA in a correct client,
			// but return the same rejection defensively if it does.
			return t.deliverPaceFailure()
		}
	case iso7816.INS_GET_CHALLENGE:
		return append(bytes.Clone(t.rndIcc), 0x90, 0x00)
	case iso7816.INS_EXTERNAL_AUTHENTICATE:
		rsp, err := t.bacMutualAuthResponse(encodedData)
		if err != nil {
			return []byte{0x6F, 0x00}
		}
		return append(rsp, 0x90, 0x00)
	}
	return []byte{0x6D, 0x00}
}

func (t *TcAc01Transceiver) bacMutualAuthResponse(cApdu []byte) ([]byte, error) {
	if len(cApdu) < 5 || cApdu[1] != iso7816.INS_EXTERNAL_AUTHENTICATE {
		return nil, fmt.Errorf("not external authenticate")
	}
	lc := int(cApdu[4])
	if lc != 40 || len(cApdu) < 5+lc {
		return nil, fmt.Errorf("unexpected command data length %d", lc)
	}
	cmd := cApdu[5 : 5+lc]

	kSeed, err := t.bacKseed()
	if err != nil {
		return nil, err
	}
	kEnc := cryptoutils.KDF(kSeed, cryptoutils.KDF_COUNTER_KSENC, cryptoutils.TDES, 112)
	kMac := cryptoutils.KDF(kSeed, cryptoutils.KDF_COUNTER_KSMAC, cryptoutils.TDES, 112)

	eIfd := cmd[0:32]
	mIfd := cmd[32:40]

	expMac, err := iso9797RetailMac(kMac, iso9797Pad(eIfd))
	if err != nil || !bytes.Equal(mIfd, expMac) {
		return nil, fmt.Errorf("mac verify failed")
	}

	block, err := cryptoutils.CipherForKey(cryptoutils.TDES, kEnc)
	if err != nil {
		return nil, err
	}
	plain, err := cryptoutils.CryptCBC(block, make([]byte, 8), eIfd, false)
	if err != nil {
		return nil, err
	}
	rndIfd := plain[0:8]
	_ = plain[8:16] // echoed RND.ICC
	kIfd := bytes.Clone(plain[16:32])

	kIcc := utils.HexToBytes("0B4F80323EB3191CB04970CB4052790B")
	s := make([]byte, 32)
	copy(s[0:8], t.rndIcc)
	copy(s[8:16], rndIfd)
	copy(s[16:32], kIcc)

	eIcc, err := cryptoutils.CryptCBC(block, make([]byte, 8), s, true)
	if err != nil {
		return nil, err
	}
	mIcc, err := iso9797RetailMac(kMac, iso9797Pad(eIcc))
	if err != nil {
		return nil, err
	}
	t.lastKIfd = kIfd
	t.lastRndIfd = bytes.Clone(rndIfd)
	t.bacSession = true
	return append(eIcc, mIcc...), nil
}

// BacSessionMaterial returns K.IFD / RND.IFD after a successful EXTERNAL AUTHENTICATE.
func (t *TcAc01Transceiver) BacSessionMaterial() (kIfd, rndIfd, rndIcc, kIcc []byte, ok bool) {
	if !t.bacSession || len(t.lastKIfd) != 16 || len(t.lastRndIfd) != 8 {
		return nil, nil, nil, nil, false
	}
	return bytes.Clone(t.lastKIfd), bytes.Clone(t.lastRndIfd), bytes.Clone(t.rndIcc),
		utils.HexToBytes("0B4F80323EB3191CB04970CB4052790B"), true
}

func (t *TcAc01Transceiver) bacKseed() ([]byte, error) {
	key, err := t.pass.Key()
	if err != nil {
		return nil, err
	}
	return key[0:16], nil
}

func iso9797Pad(data []byte) []byte {
	return cryptoutils.ISO9797Method2Pad(data, cryptoutils.DES_BLOCK_SIZE_BYTES)
}

func iso9797RetailMac(kMac, data []byte) ([]byte, error) {
	return cryptoutils.ISO9797RetailMacDes(kMac, data)
}
