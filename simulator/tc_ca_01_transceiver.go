package simulator

import (
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
)

// TcCa01Transceiver: BAC mutual auth succeeds and establishes chip SM; CA MSE:Set AT fails
// with caSW (SM-wrapped). Post-CA continue-check requires SM-validated READ BINARY.
type TcCa01Transceiver struct {
	bac    *TcAc01Transceiver
	caSW   []byte
	bacOK  bool
	caFail bool
	sm     *ChipSM
}

func NewTcCa01Transceiver(caStatusWord string, pass *password.Password) *TcCa01Transceiver {
	return &TcCa01Transceiver{
		bac:  NewTcAc01Transceiver("6FFF", pass),
		caSW: utils.HexToBytes(caStatusWord),
	}
}

var _ iso7816.Transceiver = (*TcCa01Transceiver)(nil)

func (t *TcCa01Transceiver) Transceive(cla, ins, p1, p2 int, data []byte, le int, encodedData []byte) []byte {
	claB, insB, p1B, p2B := byte(cla), byte(ins), byte(p1), byte(p2)

	// After BAC SM is established: unprotected READ BINARY must fail (SM dependence).
	// ISO 7816-4: 6987 = Expected SM data objects missing (not generic 6982).
	if t.sm != nil && insB == iso7816.INS_READ_BINARY && (claB&iso7816.CLA_MASK) != iso7816.CLA_MASK {
		return []byte{0x69, 0x87}
	}

	// SM-wrapped commands (post-BAC).
	if t.sm != nil && (claB&iso7816.CLA_MASK) == iso7816.CLA_MASK {
		cmd, err := t.sm.UnwrapSmCommand(claB, insB, p1B, p2B, data)
		if err != nil {
			return []byte{0x69, 0x88} // SM data objects incorrect
		}
		return t.dispatchPlain(cmd)
	}

	// BAC EXTERNAL AUTHENTICATE — establish chip SM from session material.
	if insB == iso7816.INS_EXTERNAL_AUTHENTICATE {
		rsp := t.bac.Transceive(cla, ins, p1, p2, data, le, encodedData)
		if len(rsp) >= 2 && rsp[len(rsp)-2] == 0x90 && rsp[len(rsp)-1] == 0x00 {
			if kIfd, rndIfd, rndIcc, kIcc, ok := t.bac.BacSessionMaterial(); ok {
				if sm, err := EstablishBacSm(kIfd, kIcc, rndIcc, rndIfd); err == nil {
					t.sm = sm
					t.bacOK = true
				}
			}
		}
		return rsp
	}

	// Pre-SM CA MSE (should not occur after DoBAC); keep plain reject for robustness.
	if insB == iso7816.INS_MANAGE_SE && p1B == 0x41 && (p2B == 0xA4 || p2B == 0xA6) {
		if !t.caFail {
			t.caFail = true
			return append([]byte(nil), t.caSW...)
		}
	}

	// PACE path — not used in TC-CA-01
	if insB == iso7816.INS_MANAGE_SE || insB == iso7816.INS_GENERAL_AUTHENTICATE {
		return []byte{0x6D, 0x00}
	}

	return t.bac.Transceive(cla, ins, p1, p2, data, le, encodedData)
}

func (t *TcCa01Transceiver) dispatchPlain(cmd *PlainCommand) []byte {
	// CA MSE:Set AT / Set KAT — inject failure once, SM-wrapped status.
	if cmd.INS == iso7816.INS_MANAGE_SE && cmd.P1 == 0x41 && (cmd.P2 == 0xA4 || cmd.P2 == 0xA6) {
		if !t.caFail {
			t.caFail = true
			sw := uint16(t.caSW[0])<<8 | uint16(t.caSW[1])
			out, err := t.sm.WrapSmResponse(sw, nil)
			if err != nil {
				return []byte{0x6F, 0x00}
			}
			return out
		}
		out, err := t.sm.WrapSmResponse(0x6D00, nil)
		if err != nil {
			return []byte{0x6F, 0x00}
		}
		return out
	}

	// SM-session continue-check: after CA reject, SM READ BINARY still succeeds.
	if t.caFail && cmd.INS == iso7816.INS_READ_BINARY {
		payload := []byte{0x61, 0x03, 0x5F, 0x2E, 0x00}
		out, err := t.sm.WrapSmResponse(0x9000, payload)
		if err != nil {
			return []byte{0x6F, 0x00}
		}
		return out
	}

	out, err := t.sm.WrapSmResponse(0x6D00, nil)
	if err != nil {
		return []byte{0x6F, 0x00}
	}
	return out
}

func (t *TcCa01Transceiver) BacEstablished() bool {
	return t.bacOK && t.sm != nil
}

func (t *TcCa01Transceiver) CaFailed() bool {
	return t.caFail
}

func (t *TcCa01Transceiver) ChipSM() *ChipSM {
	return t.sm
}
