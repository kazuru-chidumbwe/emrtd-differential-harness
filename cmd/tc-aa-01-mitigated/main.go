// TC-AA-01 mitigated: middleware explicit-reject on Active Authentication failure.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/gmrtd/gmrtd/bac"
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
	ActiveAuthErr     string       `json:"active_auth_err"`
	ActiveAuthSuccess bool         `json:"active_auth_success"`
	MiddlewareErr     string       `json:"middleware_err,omitempty"`
	Observability     int          `json:"observability_score"`
	ObservabilityMe   string       `json:"observability_meaning"`
	Trace             []traceEntry `json:"trace"`
}

func main() {
	profilePath := flag.String("profile", "profiles/aa-internal-auth-reject.json", "profile JSON")
	logDir := flag.String("log-dir", "logs", "log directory")
	variant := flag.String("variant", "mitigated", "variant")
	suiteID := flag.String("suite-id", "", "suite id")
	suiteSeed := flag.Int("suite-seed", 1, "seed")
	suiteN := flag.Int("suite-n", 1, "N")
	runIndex := flag.Int("run-index", 0, "index")
	figureID := flag.String("figure-id", "", "figure")
	flag.Parse()

	p, err := profile.Load(*profilePath)
	if err != nil {
		fatal("profile", err)
	}
	pass, err := password.NewPasswordMrzi(p.MRZ.DocumentNumber, p.MRZ.DateOfBirth, p.MRZ.DateOfExpiry)
	if err != nil {
		fatal("password", err)
	}

	aaSW := p.AAInjection.AaSW
	if aaSW == "" {
		aaSW = "6982"
	}

	nfc := iso7816.NewNfcSession(simulator.NewTcAa01Transceiver(aaSW, pass))
	doc := &document.Document{}
	if p.Dg15HexPath != "" {
		dg15Hex, err := profile.LoadHexFile(p.Dg15HexPath)
		if err != nil {
			fatal("dg15", err)
		}
		if err := doc.NewDG(15, dg15Hex); err != nil {
			fatal("dg15 parse", err)
		}
	}

	if _, err := bac.NewBAC(nfc, doc, pass).DoBAC(); err != nil {
		fatal("bac", err)
	}

	aa := middleware.PerformActiveAuth(nfc, doc, middleware.AAOptions{AllowContinue: false})
	aaErrStr := errString(aa.ActiveAuthErr)
	mwErrStr := errString(aa.SurfacedError)
	obs, obsMeaning := classifier.ClassifyTCAA01(classifier.TCAA01Input{
		ActiveAuthFailed:        aaErrStr != "",
		ActiveAuthSuccess:       aa.ActiveAuthOK,
		FailureSurfacedToCaller: mwErrStr != "",
	})

	prov, err := provenance.Collect(provenance.Options{
		ProfilePath: *profilePath, SuiteID: *suiteID, SuiteSeed: *suiteSeed,
		SuiteN: *suiteN, RunIndex: *runIndex, Driver: "go/tc-aa-01-mitigated",
		Variant: *variant, Middleware: "explicit-reject-aa",
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
		ActiveAuthErr: aaErrStr, ActiveAuthSuccess: aa.ActiveAuthOK, MiddlewareErr: mwErrStr,
		Observability: obs.Int(), ObservabilityMe: obsMeaning,
		Trace: buildTrace(nfc.ApduLog()),
	}
	if err := output.WriteJSON(*logDir, runID, result); err != nil {
		fatal("write", err)
	}
	_ = json.NewEncoder(os.Stdout).Encode(result)
	if aaErrStr == "" || mwErrStr == "" {
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
