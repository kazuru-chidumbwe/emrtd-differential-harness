package simulator

import (
	"bytes"
	"fmt"

	"github.com/gmrtd/gmrtd/cryptoutils"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
)

// TcAc01Transceiver: PACE APDUs fail with paceSW; BAC uses dynamic chip-side mutual auth.
type TcAc01Transceiver struct {
	paceFailed bool
	paceSW     []byte
	pass       *password.Password
	rndIcc     []byte
}

func NewTcAc01Transceiver(paceStatusWord string, pass *password.Password) *TcAc01Transceiver {
	return &TcAc01Transceiver{
		paceSW: utils.HexToBytes(paceStatusWord),
		pass:   pass,
		rndIcc: utils.HexToBytes("4608F91988702212"), // ICAO 9303 Part 11 D.3
	}
}

var _ iso7816.Transceiver = (*TcAc01Transceiver)(nil)

func (t *TcAc01Transceiver) Transceive(cla, ins, p1, p2 int, data []byte, le int, encodedData []byte) []byte {
	switch byte(ins) {
	case iso7816.INS_MANAGE_SE, iso7816.INS_GENERAL_AUTHENTICATE:
		if !t.paceFailed {
			t.paceFailed = true
			return append([]byte(nil), t.paceSW...)
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
	_ = plain[8:16]

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
	return append(eIcc, mIcc...), nil
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
