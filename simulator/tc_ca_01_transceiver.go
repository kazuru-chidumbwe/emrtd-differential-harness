package simulator

import (
	"bytes"

	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
)

// TcCa01Transceiver: BAC mutual auth succeeds; CA MSE:Set AT fails with caSW.
type TcCa01Transceiver struct {
	bac    *TcAc01Transceiver
	caSW   []byte
	bacOK  bool
	caFail bool
}

func NewTcCa01Transceiver(caStatusWord string, pass *password.Password) *TcCa01Transceiver {
	return &TcCa01Transceiver{
		bac:  NewTcAc01Transceiver("6FFF", pass),
		caSW: utils.HexToBytes(caStatusWord),
	}
}

var _ iso7816.Transceiver = (*TcCa01Transceiver)(nil)

func (t *TcCa01Transceiver) Transceive(cla, ins, p1, p2 int, data []byte, le int, encodedData []byte) []byte {
	// Chip Authentication MSE:Set AT (0x41A4) or MSE:Set KAT (0x41A6)
	if byte(ins) == iso7816.INS_MANAGE_SE && p1 == 0x41 && (p2 == 0xA4 || p2 == 0xA6) {
		if !t.caFail {
			t.caFail = true
			return append([]byte(nil), t.caSW...)
		}
	}
	// Emergent continue-check: after CA reject, READ BINARY (INS 0xB0) still succeeds.
	if t.caFail && byte(ins) == 0xB0 {
		payload := []byte{0x61, 0x03, 0x5F, 0x2E, 0x00}
		return append(payload, 0x90, 0x00)
	}
	// PACE path — not used in TC-CA-01; decline cleanly
	if byte(ins) == iso7816.INS_MANAGE_SE || byte(ins) == iso7816.INS_GENERAL_AUTHENTICATE {
		return []byte{0x6D, 0x00}
	}
	if byte(ins) == iso7816.INS_EXTERNAL_AUTHENTICATE {
		t.bacOK = true
	}
	rsp := t.bac.Transceive(cla, ins, p1, p2, data, le, encodedData)
	if len(rsp) >= 2 && rsp[len(rsp)-2] == 0x90 && byte(ins) == iso7816.INS_EXTERNAL_AUTHENTICATE {
		t.bacOK = true
	}
	return rsp
}

func (t *TcCa01Transceiver) BacEstablished() bool {
	return t.bacOK
}

func (t *TcCa01Transceiver) CaFailed() bool {
	return t.caFail
}

// Placeholder for future SM-aware CA responses.
func (t *TcCa01Transceiver) cloneSW(sw []byte) []byte {
	return bytes.Clone(sw)
}
