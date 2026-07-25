// TC-TA-01 / TC-EAC-01 gmrtd peer stub: Terminal Authentication not implemented upstream.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kazuru-chidumbwe/emrtd-differential-harness/classifier"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/output"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/profile"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/provenance"
	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/runid"
)

type result struct {
	output.Meta
	PeerUnsupported   bool   `json:"peer_unsupported"`
	Observability     int    `json:"observability_score"`
	ObservabilityMe   string `json:"observability_meaning"`
	Note              string `json:"note"`
}

func main() {
	profilePath := flag.String("profile", "profiles/ta-pso-verify-reject.json", "profile JSON")
	logDir := flag.String("log-dir", "logs", "log directory")
	variant := flag.String("variant", "unsupported", "variant")
	suiteID := flag.String("suite-id", "", "suite id")
	suiteSeed := flag.Int("suite-seed", 1, "seed")
	suiteN := flag.Int("suite-n", 1, "N")
	runIndex := flag.Int("run-index", 0, "index")
	figureID := flag.String("figure-id", "", "figure")
	flag.Parse()

	p, err := profile.Load(*profilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "profile: %v\n", err)
		os.Exit(2)
	}

	var obs classifier.Score
	var meaning string
	switch p.Mechanism {
	case "EAC":
		obs, meaning = classifier.ClassifyTCEAC01(classifier.TCEAC01Input{PeerUnsupported: true})
	default:
		obs, meaning = classifier.ClassifyTCTA01(classifier.TCTA01Input{PeerUnsupported: true})
	}

	prov, err := provenance.Collect(provenance.Options{
		ProfilePath: *profilePath, SuiteID: *suiteID, SuiteSeed: *suiteSeed,
		SuiteN: *suiteN, RunIndex: *runIndex, Driver: "go/tc-ta-eac-unsupported", Variant: *variant,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "provenance: %v\n", err)
		os.Exit(2)
	}

	runID := runid.New(fmt.Sprintf("%s-gmrtd-unsupported", p.ID))
	out := result{
		Meta: output.Meta{
			RunID: runID, TestCase: p.ID, Library: "gmrtd",
			Mechanism: p.Mechanism, Condition: p.Condition, Tier: p.Tier,
			Variant: *variant, FigureID: *figureID, Provenance: prov,
		},
		PeerUnsupported: true,
		Observability:   obs.Int(),
		ObservabilityMe: meaning,
		Note:            "gmrtd does not implement Terminal Authentication; EAC incomplete (CA-only)",
	}
	if err := output.WriteJSON(*logDir, runID, out); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(2)
	}
	_ = json.NewEncoder(os.Stdout).Encode(out)
}
