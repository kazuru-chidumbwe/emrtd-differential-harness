package middleware

import (
	"testing"

	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/simulator"
)

func TestNegotiatePACEBAC_explicitReject(t *testing.T) {
	pass, err := password.NewPasswordMrzi("L898902C", "690806", "940623")
	if err != nil {
		t.Fatal(err)
	}
	nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver("6FFF", pass))
	doc := &document.Document{}
	doc.Mf.CardAccess, err = document.NewCardAccess(utils.HexToBytes("31143012060A04007F0007020204020402010202010D"))
	if err != nil {
		t.Fatal(err)
	}

	sess := NegotiatePACEBAC(nfc, doc, pass, Options{AllowBACFallback: false})
	if sess.SurfacedError == nil {
		t.Fatal("expected surfaced error on PACE failure")
	}
	if sess.BacSuccess {
		t.Fatal("BAC should not run after explicit reject")
	}
}

func TestNegotiatePACEBAC_allowFallback(t *testing.T) {
	pass, err := password.NewPasswordMrzi("L898902C", "690806", "940623")
	if err != nil {
		t.Fatal(err)
	}
	nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver("6FFF", pass))
	doc := &document.Document{}
	doc.Mf.CardAccess, err = document.NewCardAccess(utils.HexToBytes("31143012060A04007F0007020204020402010202010D"))
	if err != nil {
		t.Fatal(err)
	}

	sess := NegotiatePACEBAC(nfc, doc, pass, Options{AllowBACFallback: true})
	if sess.SurfacedError != nil {
		t.Fatalf("unexpected surfaced error: %v", sess.SurfacedError)
	}
	if !sess.BacSuccess {
		t.Fatal("expected BAC success with fallback allowed")
	}
}
