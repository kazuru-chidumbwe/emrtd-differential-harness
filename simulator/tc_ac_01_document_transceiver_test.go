package simulator_test

import (
	"testing"

	"github.com/gmrtd/gmrtd/cms"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/reader"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/classifier"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/profile"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/simulator"
)

type silentStatus struct{}

func (silentStatus) Status(reader.Status) {}

// At gmrtd v1.1.3+, ReadDocument fails closed after a recorded PACE error unless
// AllowBacFallbackOnPaceError is set. Historical Score-0 behaviour at pin 8fea245
// is preserved under harness tag v1.0.7.
func TestTcAc01ReadDocumentBaselineScore2FailClosed(t *testing.T) {
	p, err := profile.Load("../profiles/pace-then-bac-downgrade.json")
	if err != nil {
		t.Fatal(err)
	}
	pass, err := password.NewPasswordMrzi(p.MRZ.DocumentNumber, p.MRZ.DateOfBirth, p.MRZ.DateOfExpiry)
	if err != nil {
		t.Fatal(err)
	}
	tr := simulator.NewTcAc01DocumentTransceiver("6FFF", "mse_set_at", "sw", p.CardAccessHex, pass)
	nfc := iso7816.NewNfcSession(tr)
	r := reader.NewReader(silentStatus{}, nfc, &cms.GenericCertPool{})
	docEx, _, readErr := r.ReadDocument(pass, nil, nil)
	if readErr == nil {
		t.Fatal("expected ReadDocument error (fail-closed after PACE)")
	}
	if docEx == nil || docEx.Session.PaceErr == nil {
		t.Fatal("expected Session.PaceErr set on partial DocumentEx")
	}
	bacOK := docEx.Session.BacResult != nil && docEx.Session.BacResult.Success
	if bacOK {
		t.Fatal("expected BAC not completed under fail-closed default")
	}
	obs, _ := classifier.ClassifyTCAC01(classifier.TCAC01Input{
		PaceFailed:           true,
		BacSuccess:           bacOK,
		BacErr:               "",
		PaceSurfacedToCaller: readErr != nil,
	})
	if obs.Int() != 2 {
		t.Fatalf("observability=%d want 2", obs.Int())
	}
}
