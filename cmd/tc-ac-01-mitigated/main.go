// TC-AC-01 mitigated run: middleware §VIII explicit-reject on PACE failure (no silent BAC fallback).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/utils"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/classifier"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/profile"
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
	RunID           string       `json:"run_id"`
	TestCase        string       `json:"test_case"`
	Library         string       `json:"library"`
	Mechanism       string       `json:"mechanism"`
	Condition       string       `json:"condition"`
	Variant         string       `json:"variant"`
	PaceErr         string       `json:"pace_err"`
	BacErr          string       `json:"bac_err"`
	BacSuccess      bool         `json:"bac_success"`
	MiddlewareErr   string       `json:"middleware_err,omitempty"`
	Observability   int          `json:"observability_score"`
	ObservabilityMe string       `json:"observability_meaning"`
	Trace           []traceEntry `json:"trace"`
}

func main() {
	profilePath := flag.String("profile", "profiles/pace-then-bac-downgrade.json", "synthetic chip profile JSON")
	logDir := flag.String("log-dir", "logs", "output directory for run traces")
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

	transceiver := simulator.NewTcAc01Transceiver(paceSW, pass)
	nfc := iso7816.NewNfcSession(transceiver)

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

	runID := runid.New(fmt.Sprintf("%s-gmrtd-mitigated", p.ID))
	result := smokeResult{
		RunID:           runID,
		TestCase:        p.ID,
		Library:         "gmrtd",
		Mechanism:       p.Mechanism,
		Condition:       p.Condition,
		Variant:         "mitigated",
		PaceErr:         paceErrStr,
		BacErr:          bacErrStr,
		BacSuccess:      sess.BacSuccess,
		MiddlewareErr:   mwErrStr,
		Observability:   obs.Int(),
		ObservabilityMe: obsMeaning,
		Trace:           buildTrace(nfc.ApduLog()),
	}

	if err := writeResult(*logDir, runID, result); err != nil {
		fatal("write log", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)

	if paceErrStr == "" || mwErrStr == "" {
		fmt.Fprintf(os.Stderr, "TC-AC-01 mitigated gate failed: pace_err=%q middleware_err=%q\n", paceErrStr, mwErrStr)
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

func writeResult(logDir, runID string, result smokeResult) error {
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return err
	}
	outPath := filepath.Join(logDir, runID+".json")
	raw, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, raw, 0o644)
}

func buildTrace(log *iso7816.ApduLog) []traceEntry {
	if log == nil {
		return nil
	}
	entries := log.AllEntries()
	out := make([]traceEntry, 0, len(entries))
	for _, e := range entries {
		success := len(e.Rx) >= 2 && e.Rx[len(e.Rx)-2] == 0x90 && e.Rx[len(e.Rx)-1] == 0x00
		out = append(out, traceEntry{
			Label:   e.Desc,
			CApdu:   utils.BytesToHex(e.Tx),
			RApdu:   utils.BytesToHex(e.Rx),
			Success: success,
		})
	}
	return out
}
