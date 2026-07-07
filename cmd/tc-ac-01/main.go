// TC-AC-01 smoke run: PACE failure recorded on session, BAC proceeds (gmrtd wire tier).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gmrtd/gmrtd/bac"
	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
	"github.com/gmrtd/gmrtd/pace"
	"github.com/gmrtd/gmrtd/utils"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/profile"
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
	PaceErr         string       `json:"pace_err"`
	BacErr          string       `json:"bac_err"`
	BacSuccess      bool         `json:"bac_success"`
	StepErr         string       `json:"step_err,omitempty"`
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

	docEx := &document.DocumentEx{Document: *doc}

	// PACE (records error, does not return error to caller — gmrtd reader.go performPace)
	docEx.Session.PaceResult, docEx.Session.PaceCamResult, docEx.Session.PaceErr =
		pace.NewPace(nfc, doc, pass).DoPACE()

	// BAC only if secure messaging not established (PACE failed)
	if nfc.SM() == nil {
		docEx.Session.BacResult, docEx.Session.BacErr = bac.NewBAC(nfc, doc, pass).DoBAC()
	}

	paceErrStr := errString(docEx.Session.PaceErr)
	bacErrStr := errString(docEx.Session.BacErr)
	bacOK := docEx.Session.BacResult != nil && docEx.Session.BacResult.Success
	obs, obsMeaning := classify(paceErrStr, bacErrStr, bacOK)

	runID := fmt.Sprintf("%s-gmrtd-%s", p.ID, time.Now().UTC().Format("20060102T150405Z"))
	result := smokeResult{
		RunID:           runID,
		TestCase:        p.ID,
		Library:         "gmrtd",
		Mechanism:       p.Mechanism,
		Condition:       p.Condition,
		PaceErr:         paceErrStr,
		BacErr:          bacErrStr,
		BacSuccess:      bacOK,
		Observability:   obs,
		ObservabilityMe: obsMeaning,
		Trace:           buildTrace(nfc.ApduLog()),
	}

	if err := os.MkdirAll(*logDir, 0o755); err != nil {
		fatal("log dir", err)
	}
	outPath := filepath.Join(*logDir, runID+".json")
	raw, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fatal("marshal", err)
	}
	if err := os.WriteFile(outPath, raw, 0o644); err != nil {
		fatal("write log", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)

	if paceErrStr == "" || !bacOK {
		fmt.Fprintf(os.Stderr, "TC-AC-01 gate failed: pace_err=%q bac_success=%v\n", paceErrStr, bacOK)
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

func classify(paceErr, bacErr string, bacOK bool) (int, string) {
	if paceErr != "" && bacOK && bacErr == "" {
		return 0, "silent — PACE failure on session, BAC succeeded; negotiation steps return nil"
	}
	if paceErr != "" {
		return 1, "logged — PACE error present in session/trace"
	}
	return 2, "surfaced — explicit error path"
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
