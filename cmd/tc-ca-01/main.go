// TC-CA-01 smoke: BAC session then EAC-CA fails at MSE:Set AT; gmrtd records ChipAuthErr.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gmrtd/gmrtd/bac"
	"github.com/gmrtd/gmrtd/chipauth"
	"github.com/gmrtd/gmrtd/document"
	"github.com/gmrtd/gmrtd/iso7816"
	"github.com/gmrtd/gmrtd/password"
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
	Tier            string       `json:"tier"`
	ChipAuthErr     string       `json:"chip_auth_err"`
	ChipAuthSuccess bool         `json:"chip_auth_success"`
	StepErr         string       `json:"step_err"`
	Observability   int          `json:"observability_score"`
	ObservabilityMe string       `json:"observability_meaning"`
	Trace           []traceEntry `json:"trace"`
}

func main() {
	profilePath := flag.String("profile", "profiles/ca-v1-v2-skew.json", "synthetic chip profile JSON")
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

	caSW := p.CAInjection.CaSW
	if caSW == "" {
		caSW = "6FFF"
	}

	transceiver := simulator.NewTcCa01Transceiver(caSW, pass)
	nfc := iso7816.NewNfcSession(transceiver)

	doc := &document.Document{}
	if p.CardAccessHex != "" {
		doc.Mf.CardAccess, err = document.NewCardAccess(utils.HexToBytes(p.CardAccessHex))
		if err != nil {
			fatal("card access", err)
		}
	}

	if p.Dg14HexPath != "" {
		dg14Hex, err := profile.LoadHexFile(p.Dg14HexPath)
		if err != nil {
			fatal("dg14", err)
		}
		if err := doc.NewDG(14, dg14Hex); err != nil {
			fatal("dg14 parse", err)
		}
	}

	docEx := &document.DocumentEx{Document: *doc}

	// BAC only (no PACE in this profile path)
	docEx.Session.BacResult, docEx.Session.BacErr = bac.NewBAC(nfc, doc, pass).DoBAC()
	if docEx.Session.BacErr != nil {
		fatal("bac", docEx.Session.BacErr)
	}

	// Chip Authentication slice (same as performChipAuthentication CA branch)
	docEx.Session.ChipAuthResult, docEx.Session.ChipAuthErr = chipauth.NewChipAuth(nfc, doc).DoChipAuth()

	chipErrStr := errString(docEx.Session.ChipAuthErr)
	chipOK := docEx.Session.ChipAuthResult != nil && docEx.Session.ChipAuthResult.Success
	stepErr := "" // performChipAuthentication returns nil even when ChipAuthErr set
	obs, obsMeaning := classifyCA(chipErrStr, chipOK, stepErr)

	runID := fmt.Sprintf("%s-gmrtd-%s", p.ID, time.Now().UTC().Format("20060102T150405Z"))
	result := smokeResult{
		RunID:           runID,
		TestCase:        p.ID,
		Library:         "gmrtd",
		Mechanism:       p.Mechanism,
		Condition:       p.Condition,
		Tier:            p.Tier,
		ChipAuthErr:     chipErrStr,
		ChipAuthSuccess: chipOK,
		StepErr:         stepErr,
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

	if chipErrStr == "" || chipOK {
		fmt.Fprintf(os.Stderr, "TC-CA-01 gate failed: chip_auth_err=%q chip_auth_success=%v\n", chipErrStr, chipOK)
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

func classifyCA(chipErr string, chipOK bool, stepErr string) (int, string) {
	if chipErr != "" && !chipOK && stepErr == "" {
		return 0, "silent — ChipAuthErr on session, performChipAuthentication-equivalent step returns nil"
	}
	if chipErr != "" {
		return 1, "logged — ChipAuth error present in session/trace"
	}
	return 2, "surfaced — explicit error on step return"
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
