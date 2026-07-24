package middleware

import (
	"testing"

	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/simulator"
)

const cardAccessPaceHex = "31143012060A04007F0007020204020402010202010D"

func mustPassword(t *testing.T) *password.Password {
	t.Helper()
	pass, err := password.NewPasswordMrzi("L898902C", "690806", "940623")
	if err != nil {
		t.Fatal(err)
	}
	return pass
}

func mustPaceDoc(t *testing.T) *document.Document {
	t.Helper()
	doc := &document.Document{}
	var err error
	doc.Mf.CardAccess, err = document.NewCardAccess(utils.HexToBytes(cardAccessPaceHex))
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestNegotiatePACEBAC_explicitReject(t *testing.T) {
	pass := mustPassword(t)
	nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver("6FFF", pass))
	doc := mustPaceDoc(t)

	sess := NegotiatePACEBAC(nfc, doc, pass, Options{AllowBACFallback: false})
	if !sess.PaceOffered {
		t.Fatal("expected PaceOffered")
	}
	if sess.SurfacedError == nil {
		t.Fatal("expected surfaced error on PACE failure")
	}
	if sess.BacSuccess {
		t.Fatal("BAC should not run after explicit reject")
	}
}

func TestNegotiatePACEBAC_allowFallback(t *testing.T) {
	pass := mustPassword(t)
	nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver("6FFF", pass))
	doc := mustPaceDoc(t)

	sess := NegotiatePACEBAC(nfc, doc, pass, Options{AllowBACFallback: true})
	if sess.SurfacedError != nil {
		t.Fatalf("unexpected surfaced error: %v", sess.SurfacedError)
	}
	if !sess.BacSuccess {
		t.Fatal("expected BAC success with fallback allowed")
	}
}

func TestNegotiatePACEBAC_paceSuccessNoFalsePositive(t *testing.T) {
	pass := mustPassword(t)
	nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver("6FFF", pass))
	doc := mustPaceDoc(t)

	old := runPace
	runPace = func(*iso7816.NfcSession, *document.Document, *password.Password) error {
		return nil
	}
	defer func() { runPace = old }()

	sess := NegotiatePACEBAC(nfc, doc, pass, Options{AllowBACFallback: false})
	if sess.SurfacedError != nil {
		t.Fatalf("false positive on PACE success: %v", sess.SurfacedError)
	}
	if sess.PaceErr != nil {
		t.Fatalf("unexpected PaceErr: %v", sess.PaceErr)
	}
}

func TestNegotiatePACEBAC_bacOnlyNoFalsePositive(t *testing.T) {
	pass := mustPassword(t)
	doc := &document.Document{} // no CardAccess → PACE not advertised

	for _, allow := range []bool{false, true} {
		// Fresh transceiver each pass — BAC mutual-auth state is single-use.
		nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver("6FFF", pass))
		sess := NegotiatePACEBAC(nfc, doc, pass, Options{AllowBACFallback: allow})
		if sess.PaceOffered {
			t.Fatalf("AllowBACFallback=%v: PaceOffered should be false", allow)
		}
		if sess.SurfacedError != nil {
			t.Fatalf("AllowBACFallback=%v: false positive on BAC-only chip: %v", allow, sess.SurfacedError)
		}
		if !sess.BacSuccess {
			t.Fatalf("AllowBACFallback=%v: expected BAC success on BAC-only chip (bacErr=%v)", allow, sess.BacErr)
		}
	}
}
