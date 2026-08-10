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

func (silentStatus) Status(string) {}

func TestTcAc01ReadDocumentBaselineScore0(t *testing.T) {
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
	if readErr != nil {
		t.Fatalf("ReadDocument err=%v (want nil — library swallows PaceErr)", readErr)
	}
	if docEx == nil || docEx.Session.PaceErr == nil {
		t.Fatal("expected Session.PaceErr set")
	}
	bacOK := docEx.Session.BacResult != nil && docEx.Session.BacResult.Success
	if !bacOK {
		t.Fatalf("expected BAC success; bacErr=%v", docEx.Session.BacErr)
	}
	obs, _ := classifier.ClassifyTCAC01(classifier.TCAC01Input{
		PaceFailed:           true,
		BacSuccess:           true,
		BacErr:               "",
		PaceSurfacedToCaller: readErr != nil,
	})
	if obs.Int() != 0 {
		t.Fatalf("observability=%d want 0", obs.Int())
	}
}
