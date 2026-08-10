package simulator

import (
	_ "embed"
	"encoding/binary"

	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
)

//go:embed fixtures/ac01-ef-sod.hex
var ac01EfSodHex string

// ICAO Doc 9303 Part 10 Table 31 EF.DIR sample.
var ac01EfDir = utils.HexToBytes("61094F07A000000247100161094F07A000000247200161094F07A000000247200261094F07A0000002472003")

// ICAO Doc 9303 Part 10 EF.COM sample (gmrtd document/com_test.go).
var ac01EfCom = utils.HexToBytes("60145F0104303130365F36063034303030305C026175")

// File IDs used by gmrtd reader.ReadDocument (reader/reader.go).
const (
	fileIDCardAccess = uint16(0x011C)
	fileIDEFSOD      = uint16(0x011D)
	fileIDEFCOM      = uint16(0x011E)
	fileIDEFDIR      = uint16(0x2F00)
)

// TcAc01DocumentTransceiver wraps TcAc01Transceiver so gmrtd reader.ReadDocument
// can complete with err==nil after PACE-fail ∧ BAC-ok: Select MF / CardAccess /
// Select AID, then post-BAC SM reads of minimal EF.DIR / EF.SOD / EF.COM.
// Other elementary files return file-not-found (6A82) so DG reads are skipped
// via NewDG(nil) (same path as gmrtd reader_test TestReadLDS1DgsFilesNotFound).
type TcAc01DocumentTransceiver struct {
	bac          *TcAc01Transceiver
	cardAccess   []byte
	sm           *ChipSM
	selectedFile uint16
	haveFile     bool
	files        map[uint16][]byte
}

// NewTcAc01DocumentTransceiver builds a ReadDocument-capable AC-01 chip.
// cardAccessHex is the profile CardAccess EF body (plaintext, pre-PACE).
func NewTcAc01DocumentTransceiver(paceStatusWord, paceFailOn, paceChannel, cardAccessHex string, pass *password.Password) *TcAc01DocumentTransceiver {
	sod := utils.HexToBytes(ac01EfSodHex)
	return &TcAc01DocumentTransceiver{
		bac:        NewTcAc01TransceiverWithChannel(paceStatusWord, paceFailOn, paceChannel, pass),
		cardAccess: utils.HexToBytes(cardAccessHex),
		files: map[uint16][]byte{
			fileIDCardAccess: bytesClone(utils.HexToBytes(cardAccessHex)),
			fileIDEFDIR:      bytesClone(ac01EfDir),
			fileIDEFSOD:      bytesClone(sod),
			fileIDEFCOM:      bytesClone(ac01EfCom),
		},
	}
}

func bytesClone(b []byte) []byte {
	if b == nil {
		return nil
	}
	out := make([]byte, len(b))
	copy(out, b)
	return out
}

var _ iso7816.Transceiver = (*TcAc01DocumentTransceiver)(nil)

func (t *TcAc01DocumentTransceiver) Transceive(cla, ins, p1, p2 int, data []byte, le int, encodedData []byte) []byte {
	claB, insB, p1B, p2B := byte(cla), byte(ins), byte(p1), byte(p2)

	// Post-BAC: unprotected READ BINARY must fail (SM required).
	if t.sm != nil && insB == iso7816.INS_READ_BINARY && (claB&iso7816.CLA_MASK) != iso7816.CLA_MASK {
		return []byte{0x69, 0x87}
	}

	// SM-wrapped commands after BAC.
	if t.sm != nil && (claB&iso7816.CLA_MASK) == iso7816.CLA_MASK {
		cmd, err := t.sm.UnwrapSmCommand(claB, insB, p1B, p2B, data)
		if err != nil {
			return []byte{0x69, 0x88}
		}
		return t.dispatchPlainSM(cmd)
	}

	// BAC EXTERNAL AUTHENTICATE — establish chip SM.
	if insB == iso7816.INS_EXTERNAL_AUTHENTICATE {
		rsp := t.bac.Transceive(cla, ins, p1, p2, data, le, encodedData)
		if len(rsp) >= 2 && rsp[len(rsp)-2] == 0x90 && rsp[len(rsp)-1] == 0x00 {
			if kIfd, rndIfd, rndIcc, kIcc, ok := t.bac.BacSessionMaterial(); ok {
				if sm, err := EstablishBacSm(kIfd, kIcc, rndIcc, rndIfd); err == nil {
					t.sm = sm
				}
			}
		}
		return rsp
	}

	// Plaintext SELECT / READ before or without SM.
	return t.dispatchPlain(claB, insB, p1B, p2B, data, le, encodedData)
}

func (t *TcAc01DocumentTransceiver) dispatchPlainSM(cmd *PlainCommand) []byte {
	rsp := t.dispatchPlain(0x00, cmd.INS, cmd.P1, cmd.P2, cmd.Data, cmd.Le, nil)
	if len(rsp) < 2 {
		out, err := t.sm.WrapSmResponse(0x6F00, nil)
		if err != nil {
			return []byte{0x6F, 0x00}
		}
		return out
	}
	sw := uint16(rsp[len(rsp)-2])<<8 | uint16(rsp[len(rsp)-1])
	payload := rsp[:len(rsp)-2]
	out, err := t.sm.WrapSmResponse(sw, payload)
	if err != nil {
		return []byte{0x6F, 0x00}
	}
	return out
}

func (t *TcAc01DocumentTransceiver) dispatchPlain(cla, ins, p1, p2 byte, data []byte, le int, encodedData []byte) []byte {
	switch ins {
	case iso7816.INS_SELECT:
		return t.selectFile(p1, p2, data)
	case iso7816.INS_READ_BINARY:
		return t.readBinary(p1, p2, le)
	case iso7816.INS_MANAGE_SE, iso7816.INS_GENERAL_AUTHENTICATE,
		iso7816.INS_GET_CHALLENGE, iso7816.INS_EXTERNAL_AUTHENTICATE:
		// Delegate PACE/BAC (EXTERNAL AUTH handled above when not SM-wrapped).
		return t.bac.Transceive(int(cla), int(ins), int(p1), int(p2), data, le, encodedData)
	default:
		return []byte{0x6D, 0x00}
	}
}

func (t *TcAc01DocumentTransceiver) selectFile(p1, p2 byte, data []byte) []byte {
	// Select MF: P1=00 P2=0C data=3F00 (or empty).
	if p1 == 0x00 && (len(data) == 0 || (len(data) == 2 && data[0] == 0x3f && data[1] == 0x00)) {
		t.haveFile = false
		return []byte{0x90, 0x00}
	}
	// Select by AID: P1=04.
	if p1 == 0x04 {
		t.haveFile = false
		return []byte{0x90, 0x00}
	}
	// Select EF: P1=02, file ID in data.
	if p1 == 0x02 && len(data) == 2 {
		fid := binary.BigEndian.Uint16(data)
		if _, ok := t.files[fid]; ok {
			t.selectedFile = fid
			t.haveFile = true
			return []byte{0x90, 0x00}
		}
		t.haveFile = false
		return []byte{0x6A, 0x82} // file not found
	}
	return []byte{0x6A, 0x86}
}

func (t *TcAc01DocumentTransceiver) readBinary(p1, p2 byte, le int) []byte {
	if !t.haveFile {
		return []byte{0x6A, 0x82}
	}
	body := t.files[t.selectedFile]
	offset := int(p1)<<8 | int(p2)
	if offset < 0 || offset > len(body) {
		return []byte{0x6A, 0x86}
	}
	if le <= 0 {
		le = 256
	}
	end := offset + le
	if end > len(body) {
		end = len(body)
	}
	out := append(bytesClone(body[offset:end]), 0x90, 0x00)
	return out
}

// ChipSM returns the post-BAC secure messaging endpoint (tests).
func (t *TcAc01DocumentTransceiver) ChipSM() *ChipSM {
	return t.sm
}
