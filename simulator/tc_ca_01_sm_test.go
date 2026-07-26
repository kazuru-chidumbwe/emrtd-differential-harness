package simulator

import (
	"bytes"
	"testing"

	"github.com/gmrtd/gmrtd/bac"
	"github.com/gmrtd/gmrtd/cryptoutils"
	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
)

func TestEstablishBacSmICAOWorkedExample(t *testing.T) {
	kIfd := utils.HexToBytes("0B795240CB7049B01C19B33E32804F0B")
	kIcc := utils.HexToBytes("0B4F80323EB3191CB04970CB4052790B")
	rndIcc := utils.HexToBytes("4608F91988702212")
	rndIfd := utils.HexToBytes("781723860C06C226")

	sm, err := EstablishBacSm(kIfd, kIcc, rndIcc, rndIfd)
	if err != nil {
		t.Fatal(err)
	}
	expEnc := utils.HexToBytes("979EC13B1CBFE9DCD01AB0FED307EAE5")
	expMac := utils.HexToBytes("F1CB1F1FB5ADF208806B89DC579DC1F8")
	expSSC := utils.HexToBytes("887022120C06C226")
	if !bytes.Equal(sm.ksEnc, expEnc) {
		t.Fatalf("ksEnc mismatch\nexp %x\ngot %x", expEnc, sm.ksEnc)
	}
	if !bytes.Equal(sm.ksMac, expMac) {
		t.Fatalf("ksMac mismatch\nexp %x\ngot %x", expMac, sm.ksMac)
	}
	if !bytes.Equal(sm.ssc, expSSC) {
		t.Fatalf("SSC mismatch\nexp %x\ngot %x", expSSC, sm.ssc)
	}
}

func TestChipSMRoundTripWithReader(t *testing.T) {
	kIfd := utils.HexToBytes("0B795240CB7049B01C19B33E32804F0B")
	kIcc := utils.HexToBytes("0B4F80323EB3191CB04970CB4052790B")
	rndIcc := utils.HexToBytes("4608F91988702212")
	rndIfd := utils.HexToBytes("781723860C06C226")

	chip, err := EstablishBacSm(kIfd, kIcc, rndIcc, rndIfd)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := iso7816.NewSecureMessaging(cryptoutils.TDES, bytes.Clone(chip.ksEnc), bytes.Clone(chip.ksMac))
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.SetSSC(bytes.Clone(chip.ssc)); err != nil {
		t.Fatal(err)
	}

	plain := iso7816.NewCApdu(0x00, iso7816.INS_READ_BINARY, 0x00, 0x00, nil, 5)
	encCmd, err := reader.Encode(plain)
	if err != nil {
		t.Fatal(err)
	}
	wire := encCmd.Encode()
	smData := shortCapduData(t, wire)
	cmd, err := chip.UnwrapSmCommand(wire[0], wire[1], wire[2], wire[3], smData)
	if err != nil {
		t.Fatal(err)
	}
	if cmd.INS != iso7816.INS_READ_BINARY {
		t.Fatalf("ins %02x", cmd.INS)
	}

	payload := []byte{0x61, 0x03, 0x5F, 0x2E, 0x00}
	wireRsp, err := chip.WrapSmResponse(0x9000, payload)
	if err != nil {
		t.Fatal(err)
	}
	rapdu, err := reader.Decode(wireRsp)
	if err != nil {
		t.Fatal(err)
	}
	if !rapdu.IsSuccess() || !bytes.Equal(rapdu.Data, payload) {
		t.Fatalf("rapdu %+v", rapdu)
	}
}

func shortCapduData(t *testing.T, enc []byte) []byte {
	t.Helper()
	if len(enc) < 5 {
		t.Fatalf("short capdu %x", enc)
	}
	lc := int(enc[4])
	if len(enc) < 5+lc {
		t.Fatalf("capdu len %x", enc)
	}
	return enc[5 : 5+lc]
}

func TestTcCa01UnprotectedB0FailsSMReadSucceeds(t *testing.T) {
	pass, err := password.NewPasswordMrzi("L898902C", "690806", "940623")
	if err != nil {
		t.Fatal(err)
	}
	tr := NewTcCa01Transceiver("6982", pass)
	nfc := iso7816.NewNfcSession(tr)
	doc := &document.Document{}

	bacResult, err := bac.NewBAC(nfc, doc, pass).DoBAC()
	if err != nil {
		t.Fatalf("bac: %v", err)
	}
	if bacResult == nil || !bacResult.Success {
		t.Fatal("bac not successful")
	}
	if !tr.BacEstablished() {
		t.Fatal("chip SM not established after BAC")
	}

	raw := tr.Transceive(0x00, 0xB0, 0x00, 0x00, nil, 8, nil)
	if len(raw) >= 2 && raw[len(raw)-2] == 0x90 && raw[len(raw)-1] == 0x00 {
		t.Fatalf("unprotected B0 must not succeed, got %x", raw)
	}
	if len(raw) < 2 || raw[len(raw)-2] != 0x69 || raw[len(raw)-1] != 0x87 {
		t.Fatalf("expected 6987 (SM objects missing) on unprotected B0, got %x", raw)
	}

	err = nfc.MseSetAT(0x41, 0xA4, nil)
	if err == nil {
		t.Fatal("expected CA MSE failure")
	}
	if !tr.CaFailed() {
		t.Fatal("caFail flag not set")
	}

	data, err := nfc.ReadBinaryFromOffset(0, 5)
	if err != nil {
		t.Fatalf("SM READ BINARY after CA fail: %v", err)
	}
	if !bytes.Equal(data, []byte{0x61, 0x03, 0x5F, 0x2E, 0x00}) {
		t.Fatalf("payload %x", data)
	}
}

func TestChipSMTamperedMACFails(t *testing.T) {
	kIfd := utils.HexToBytes("0B795240CB7049B01C19B33E32804F0B")
	kIcc := utils.HexToBytes("0B4F80323EB3191CB04970CB4052790B")
	rndIcc := utils.HexToBytes("4608F91988702212")
	rndIfd := utils.HexToBytes("781723860C06C226")
	chip, err := EstablishBacSm(kIfd, kIcc, rndIcc, rndIfd)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := iso7816.NewSecureMessaging(cryptoutils.TDES, bytes.Clone(chip.ksEnc), bytes.Clone(chip.ksMac))
	if err != nil {
		t.Fatal(err)
	}
	_ = reader.SetSSC(bytes.Clone(chip.ssc))
	encCmd, err := reader.Encode(iso7816.NewCApdu(0x00, iso7816.INS_READ_BINARY, 0, 0, nil, 5))
	if err != nil {
		t.Fatal(err)
	}
	wire := encCmd.Encode()
	smData := bytes.Clone(shortCapduData(t, wire))
	if len(smData) > 0 {
		smData[len(smData)-1] ^= 0xFF // corrupt MAC tail
	}
	_, err = chip.UnwrapSmCommand(wire[0], wire[1], wire[2], wire[3], smData)
	if err == nil {
		t.Fatal("expected MAC failure")
	}
}
