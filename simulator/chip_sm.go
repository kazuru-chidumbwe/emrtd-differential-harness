package simulator

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/gmrtd/gmrtd/cryptoutils"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/tlv"
	"github.com/gmrtd/gmrtd/utils"
)

// ChipSM is a BAC secure-messaging endpoint for the synthetic chip.
// SSC advances once per inbound SM command and once per outbound SM response,
// matching gmrtd NfcSession Encode+Decode pairing.
type ChipSM struct {
	ksEnc     []byte
	ksMac     []byte
	ssc       []byte
	encCipher cipher.Block
}

// EstablishBacSm derives post-BAC session keys (K.IFD⊕K.ICC → KDF) and SSC
// as in ICAO Doc 9303 Part 11 / gmrtd bac.setupSecureMessaging.
func EstablishBacSm(kIfd, kIcc, rndIcc, rndIfd []byte) (*ChipSM, error) {
	if len(kIfd) != 16 || len(kIcc) != 16 {
		return nil, fmt.Errorf("EstablishBacSm: K.IFD/K.ICC must be 16 bytes")
	}
	if len(rndIcc) != 8 || len(rndIfd) != 8 {
		return nil, fmt.Errorf("EstablishBacSm: RND.ICC/RND.IFD must be 8 bytes")
	}
	kXor := utils.XorBytes(kIfd, kIcc)
	ksEnc := cryptoutils.KDF(kXor, cryptoutils.KDF_COUNTER_KSENC, cryptoutils.TDES, 112)
	ksMac := cryptoutils.KDF(kXor, cryptoutils.KDF_COUNTER_KSMAC, cryptoutils.TDES, 112)
	encCipher, err := cryptoutils.CipherForKey(cryptoutils.TDES, ksEnc)
	if err != nil {
		return nil, fmt.Errorf("EstablishBacSm: enc cipher: %w", err)
	}
	ssc := make([]byte, 8)
	copy(ssc[0:4], rndIcc[4:8])
	copy(ssc[4:8], rndIfd[4:8])
	return &ChipSM{
		ksEnc:     ksEnc,
		ksMac:     ksMac,
		ssc:       ssc,
		encCipher: encCipher,
	}, nil
}

// PlainCommand is the unprotected command recovered from an SM C-APDU.
type PlainCommand struct {
	INS  byte
	P1   byte
	P2   byte
	Data []byte
	Le   int
}

func (sm *ChipSM) sscIncrement() {
	pre := new(big.Int).SetBytes(sm.ssc)
	post := new(big.Int).Add(pre, big.NewInt(1))
	if len(post.Bytes()) > len(sm.ssc) {
		sm.ssc = make([]byte, len(sm.ssc))
	} else {
		post.FillBytes(sm.ssc)
	}
}

func (sm *ChipSM) pad(data []byte) []byte {
	return cryptoutils.ISO9797Method2Pad(data, sm.encCipher.BlockSize())
}

func (sm *ChipSM) unpad(data []byte) ([]byte, error) {
	return cryptoutils.ISO9797Method2Unpad(data)
}

func (sm *ChipSM) cbc(data []byte, encrypt bool) ([]byte, error) {
	iv := make([]byte, sm.encCipher.BlockSize())
	return cryptoutils.CryptCBC(sm.encCipher, iv, data, encrypt)
}

func (sm *ChipSM) mac(data []byte) ([]byte, error) {
	return cryptoutils.ISO9797RetailMacDes(sm.ksMac, data)
}

// UnwrapSmCommand verifies the SM command MAC and returns the plaintext command.
// Advances SSC once (pairs with reader Encode).
func (sm *ChipSM) UnwrapSmCommand(cla, ins, p1, p2 byte, smData []byte) (*PlainCommand, error) {
	if cla&iso7816.CLA_MASK != iso7816.CLA_MASK {
		return nil, fmt.Errorf("UnwrapSmCommand: expected SM CLA mask 0x0C (got %02x)", cla)
	}
	sm.sscIncrement()

	nodes, err := tlv.Decode(smData)
	if err != nil {
		return nil, fmt.Errorf("UnwrapSmCommand: tlv: %w", err)
	}
	tag8E := nodes.NodeByTag(0x8E)
	if !tag8E.IsValidNode() {
		return nil, fmt.Errorf("UnwrapSmCommand: missing DO'8E'")
	}
	actMAC := tag8E.Value()

	macNodes := &tlv.TlvNodes{}
	for _, n := range nodes.Nodes() {
		if n.Tag() == 0x8E {
			continue
		}
		macNodes.AddNode(tlv.NewTlvSimpleNode(n.Tag(), n.Value()))
	}

	header := []byte{iso7816.CLA_MASK, ins, p1, p2}
	headerPad := sm.pad(header)
	macData := append(bytes.Clone(sm.ssc), headerPad...)
	macData = append(macData, macNodes.Encode()...)
	expMAC, err := sm.mac(sm.pad(macData))
	if err != nil {
		return nil, fmt.Errorf("UnwrapSmCommand: mac: %w", err)
	}
	if !bytes.Equal(expMAC, actMAC) {
		return nil, fmt.Errorf("UnwrapSmCommand: MAC mismatch")
	}

	cmd := &PlainCommand{INS: ins, P1: p1, P2: p2}
	if n := nodes.NodeByTag(0x87); n.IsValidNode() {
		val := n.Value()
		if len(val) < 1 || val[0] != 0x01 {
			return nil, fmt.Errorf("UnwrapSmCommand: DO'87' version")
		}
		plain, err := sm.cbc(val[1:], false)
		if err != nil {
			return nil, err
		}
		cmd.Data, err = sm.unpad(plain)
		if err != nil {
			return nil, err
		}
	} else if n := nodes.NodeByTag(0x85); n.IsValidNode() {
		val := n.Value()
		if len(val) < 1 || val[0] != 0x01 {
			return nil, fmt.Errorf("UnwrapSmCommand: DO'85' version")
		}
		plain, err := sm.cbc(val[1:], false)
		if err != nil {
			return nil, err
		}
		cmd.Data, err = sm.unpad(plain)
		if err != nil {
			return nil, err
		}
	}
	if n := nodes.NodeByTag(0x97); n.IsValidNode() {
		leBytes := n.Value()
		if len(leBytes) == 1 {
			if leBytes[0] == 0 {
				cmd.Le = 256
			} else {
				cmd.Le = int(leBytes[0])
			}
		} else if len(leBytes) == 2 {
			cmd.Le = int(binary.BigEndian.Uint16(leBytes))
			if cmd.Le == 0 {
				cmd.Le = 65536
			}
		}
	}
	return cmd, nil
}

// WrapSmResponse builds an SM R-APDU (data||SW) the reader Decode accepts.
// Advances SSC once (pairs with reader Decode).
func (sm *ChipSM) WrapSmResponse(status uint16, data []byte) ([]byte, error) {
	sm.sscIncrement()

	nodes := &tlv.TlvNodes{}
	if len(data) > 0 {
		enc, err := sm.cbc(sm.pad(data), true)
		if err != nil {
			return nil, err
		}
		val := append([]byte{0x01}, enc...)
		nodes.AddNode(tlv.NewTlvSimpleNode(0x87, val))
	}
	sw := make([]byte, 2)
	binary.BigEndian.PutUint16(sw, status)
	nodes.AddNode(tlv.NewTlvSimpleNode(0x99, sw))

	macData := make([]byte, 0)
	macData = append(macData, sm.ssc...)
	macData = append(macData, nodes.NodeByTag(0x85).Encode()...)
	macData = append(macData, nodes.NodeByTag(0x87).Encode()...)
	macData = append(macData, nodes.NodeByTag(0x99).Encode()...)
	mac, err := sm.mac(sm.pad(macData))
	if err != nil {
		return nil, err
	}
	nodes.AddNode(tlv.NewTlvSimpleNode(0x8E, mac))

	out := nodes.Encode()
	out = append(out, sw...)
	return out, nil
}
