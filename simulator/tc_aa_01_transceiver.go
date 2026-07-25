package simulator

import (
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
)

// TcAa01Transceiver: BAC mutual auth succeeds; first INTERNAL AUTHENTICATE fails with aaSW.
type TcAa01Transceiver struct {
	bac    *TcAc01Transceiver
	aaSW   []byte
	bacOK  bool
	aaFail bool
}

func NewTcAa01Transceiver(aaStatusWord string, pass *password.Password) *TcAa01Transceiver {
	return &TcAa01Transceiver{
		bac:  NewTcAc01Transceiver("6FFF", pass),
		aaSW: utils.HexToBytes(aaStatusWord),
	}
}

var _ iso7816.Transceiver = (*TcAa01Transceiver)(nil)

func (t *TcAa01Transceiver) Transceive(cla, ins, p1, p2 int, data []byte, le int, encodedData []byte) []byte {
	// Active Authentication: INTERNAL AUTHENTICATE (INS 0x88)
	if byte(ins) == iso7816.INS_INTERNAL_AUTHENTICATE {
		if !t.aaFail {
			t.aaFail = true
			return append([]byte(nil), t.aaSW...)
		}
		return []byte{0x6D, 0x00}
	}
	// Decline PACE / CA MSE paths cleanly
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

func (t *TcAa01Transceiver) BacEstablished() bool {
	return t.bacOK
}

func (t *TcAa01Transceiver) AaFailed() bool {
	return t.aaFail
}
