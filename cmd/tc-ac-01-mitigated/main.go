// TC-AC-01 mitigated: middleware §VIII explicit-reject on PACE failure.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/classifier"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/output"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/profile"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/provenance"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/runid"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/middleware"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/simulator"
)

type traceEntry struct {
	Label   string `json:"label"`
	CApdu   string `json:"capdu"`
	RApdu   string `json:"rapdu"`
	Success bool   `json:"success"`
}

type smokeResult struct {
	output.Meta
	PaceErr         string       `json:"pace_err"`
	BacErr          string       `json:"bac_err"`
	BacSuccess      bool         `json:"bac_success"`
	MiddlewareErr   string       `json:"middleware_err,omitempty"`
	Observability   int          `json:"observability_score"`
	ObservabilityMe string       `json:"observability_meaning"`
	Trace           []traceEntry `json:"trace"`
}

func main() {
	profilePath := flag.String("profile", "profiles/pace-then-bac-downgrade.json", "profile JSON")
	logDir := flag.String("log-dir", "logs", "log directory")
	variant := flag.String("variant", "mitigated", "variant label")
	suiteID := flag.String("suite-id", "", "suite id")
	suiteSeed := flag.Int("suite-seed", 1, "suite seed")
	suiteN := flag.Int("suite-n", 1, "suite N")
	runIndex := flag.Int("run-index", 0, "run index")
	figureID := flag.String("figure-id", "", "figure id")
	flag.Parse()

	p, err := profile.Load(*profilePath)
	if err != nil {
		fatal("profile", err)
	}

	pass, err := password.NewPasswordMrzi(p.MRZ.DocumentNumber, p.MRZ.DateOfBirth, p.MRZ.DateOfExpiry)
	if err != nil {
		fatal("password", err)
	}

	paceSW := p.Injection.PaceSW
	if paceSW == "" {
		paceSW = "6FFF"
	}

	nfc := iso7816.NewNfcSession(simulator.NewTcAc01Transceiver(paceSW, pass))
	doc := &document.Document{}
	doc.Mf.CardAccess, err = document.NewCardAccess(utils.HexToBytes(p.CardAccessHex))
	if err != nil {
		fatal("card access", err)
	}

	sess := middleware.NegotiatePACEBAC(nfc, doc, pass, middleware.Options{AllowBACFallback: false})

	paceErrStr := errString(sess.PaceErr)
	bacErrStr := errString(sess.BacErr)
	mwErrStr := errString(sess.SurfacedError)
	obs, obsMeaning := classifier.ClassifyTCAC01(classifier.TCAC01Input{
		PaceFailed:           paceErrStr != "" || mwErrStr != "",
		BacSuccess:           sess.BacSuccess,
		BacErr:               bacErrStr,
		PaceSurfacedToCaller: mwErrStr != "",
	})

	prov, err := provenance.Collect(provenance.Options{
		ProfilePath: *profilePath,
		SuiteID:     *suiteID,
		SuiteSeed:   *suiteSeed,
		SuiteN:      *suiteN,
		RunIndex:    *runIndex,
		Driver:      "go/tc-ac-01-mitigated",
		Variant:     *variant,
		Middleware:  "explicit-reject-pace",
	})
	if err != nil {
		fatal("provenance", err)
	}

	runID := runid.New(fmt.Sprintf("%s-gmrtd-mitigated", p.ID))
	result := smokeResult{
		Meta: output.Meta{
			RunID: runID, TestCase: p.ID, Library: "gmrtd",
			Mechanism: p.Mechanism, Condition: p.Condition, Tier: p.Tier,
			Variant: *variant, FigureID: *figureID, Provenance: prov,
		},
		PaceErr: paceErrStr, BacErr: bacErrStr, BacSuccess: sess.BacSuccess,
		MiddlewareErr: mwErrStr, Observability: obs.Int(), ObservabilityMe: obsMeaning,
		Trace: buildTrace(nfc.ApduLog()),
	}

	if err := output.WriteJSON(*logDir, runID, result); err != nil {
		fatal("write", err)
	}
	_ = json.NewEncoder(os.Stdout).Encode(result)

	if paceErrStr == "" || mwErrStr == "" {
		os.Exit(1)
	}
}

func fatal(step string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", step, err)
	os.Exit(2)
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func buildTrace(log *iso7816.ApduLog) []traceEntry {
	if log == nil {
		return nil
	}
	out := make([]traceEntry, 0, len(log.AllEntries()))
	for _, e := range log.AllEntries() {
		ok := len(e.Rx) >= 2 && e.Rx[len(e.Rx)-2] == 0x90 && e.Rx[len(e.Rx)-1] == 0x00
		out = append(out, traceEntry{e.Desc, utils.BytesToHex(e.Tx), utils.BytesToHex(e.Rx), ok})
	}
	return out
}
