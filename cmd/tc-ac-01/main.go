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
	"github.com/gmrtd/gmrtd/pace"
	"github.com/gmrtd/gmrtd/utils"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/classifier"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/output"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/profile"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/provenance"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/runid"
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
	Observability   int          `json:"observability_score"`
	ObservabilityMe string       `json:"observability_meaning"`
	Trace           []traceEntry `json:"trace"`
}

type suiteFlags struct {
	profilePath string
	logDir      string
	variant     string
	suiteID     string
	suiteSeed   int
	suiteN      int
	runIndex    int
	figureID    string
}

func parseFlags() suiteFlags {
	profilePath := flag.String("profile", "profiles/pace-then-bac-downgrade.json", "synthetic chip profile JSON")
	logDir := flag.String("log-dir", "logs", "output directory for run traces")
	variant := flag.String("variant", "baseline", "run variant label")
	suiteID := flag.String("suite-id", "", "suite manifest id")
	suiteSeed := flag.Int("suite-seed", 1, "suite PRNG seed (metadata)")
	suiteN := flag.Int("suite-n", 1, "suite repetition count (metadata)")
	runIndex := flag.Int("run-index", 0, "1-based run index within suite entry")
	figureID := flag.String("figure-id", "", "published figure identifier")
	flag.Parse()
	return suiteFlags{
		profilePath: *profilePath,
		logDir:      *logDir,
		variant:     *variant,
		suiteID:     *suiteID,
		suiteSeed:   *suiteSeed,
		suiteN:      *suiteN,
		runIndex:    *runIndex,
		figureID:    *figureID,
	}
}

func main() {
	flags := parseFlags()

	p, err := profile.Load(flags.profilePath)
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

	transceiver := simulator.NewTcAc01TransceiverWithInjection(paceSW, p.Injection.PaceFailOn, pass)
	nfc := iso7816.NewNfcSession(transceiver)

	doc := &document.Document{}
	doc.Mf.CardAccess, err = document.NewCardAccess(utils.HexToBytes(p.CardAccessHex))
	if err != nil {
		fatal("card access", err)
	}

	docEx := &document.DocumentEx{Document: *doc}
	docEx.Session.PaceResult, docEx.Session.PaceCamResult, docEx.Session.PaceErr =
		pace.NewPace(nfc, doc, pass).DoPACE()

	if nfc.SM() == nil {
		docEx.Session.BacResult, docEx.Session.BacErr = bac.NewBAC(nfc, doc, pass).DoBAC()
	}

	paceErrStr := errString(docEx.Session.PaceErr)
	bacErrStr := errString(docEx.Session.BacErr)
	bacOK := docEx.Session.BacResult != nil && docEx.Session.BacResult.Success
	obs, obsMeaning := classifier.ClassifyTCAC01(classifier.TCAC01Input{
		PaceFailed:           paceErrStr != "",
		BacSuccess:           bacOK,
		BacErr:               bacErrStr,
		PaceSurfacedToCaller: false,
	})

	prov, err := provenance.Collect(provenance.Options{
		ProfilePath: flags.profilePath,
		SuiteID:     flags.suiteID,
		SuiteSeed:   flags.suiteSeed,
		SuiteN:      flags.suiteN,
		RunIndex:    flags.runIndex,
		Driver:      "go/tc-ac-01",
		Variant:     flags.variant,
	})
	if err != nil {
		fatal("provenance", err)
	}

	runID := runid.New(fmt.Sprintf("%s-gmrtd", p.ID))
	result := smokeResult{
		Meta: output.Meta{
			RunID:      runID,
			TestCase:   p.ID,
			Library:    "gmrtd",
			Mechanism:  p.Mechanism,
			Condition:  p.Condition,
			Tier:       p.Tier,
			Variant:    flags.variant,
			FigureID:   flags.figureID,
			Provenance: prov,
		},
		PaceErr:         paceErrStr,
		BacErr:          bacErrStr,
		BacSuccess:      bacOK,
		Observability:   obs.Int(),
		ObservabilityMe: obsMeaning,
		Trace:           buildTrace(nfc.ApduLog()),
	}

	if err := output.WriteJSON(flags.logDir, runID, result); err != nil {
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
