package classifier

import (
	"encoding/json"
	"os"
	"testing"
)

func TestClassifyTCAC01_gmrtdSilentDowngrade(t *testing.T) {
	score, _ := ClassifyTCAC01(TCAC01Input{
		PaceFailed:           true,
		BacSuccess:           true,
		BacErr:               "",
		PaceSurfacedToCaller: false,
	})
	if score != ScoreSilent {
		t.Fatalf("got %d want silent", score)
	}
}

func TestClassifyTCAC01_middlewareSurfaced(t *testing.T) {
	score, _ := ClassifyTCAC01(TCAC01Input{
		PaceFailed:           true,
		BacSuccess:           false,
		BacErr:               "",
		PaceSurfacedToCaller: true,
	})
	if score != ScoreSurfaced {
		t.Fatalf("got %d want surfaced", score)
	}
}

func TestClassifyTCAC01_paceFailBacFailLogged(t *testing.T) {
	score, _ := ClassifyTCAC01(TCAC01Input{
		PaceFailed:           true,
		BacSuccess:           false,
		BacErr:               "bac failed",
		PaceSurfacedToCaller: false,
	})
	if score != ScoreLogged {
		t.Fatalf("got %d want logged", score)
	}
}

func TestClassifyTCCA01_silent(t *testing.T) {
	score, _ := ClassifyTCCA01(TCCA01Input{
		ChipAuthFailed: true, ChipAuthSuccess: false, SessionContinueOK: true, FailureSurfacedToCaller: false,
	})
	if score != ScoreSilent {
		t.Fatalf("got %d want silent", score)
	}
}

func TestClassifyTCCA01_logged(t *testing.T) {
	// Failed but session did not continue → logged branch (not the silent quadruple).
	score, _ := ClassifyTCCA01(TCCA01Input{
		ChipAuthFailed: true, ChipAuthSuccess: false, SessionContinueOK: false, FailureSurfacedToCaller: false,
	})
	if score != ScoreLogged {
		t.Fatalf("got %d want logged", score)
	}
}

func TestClassifyTCCA01_surfaced(t *testing.T) {
	score, _ := ClassifyTCCA01(TCCA01Input{
		ChipAuthFailed: true, ChipAuthSuccess: false, SessionContinueOK: true, FailureSurfacedToCaller: true,
	})
	if score != ScoreSurfaced {
		t.Fatalf("got %d want surfaced", score)
	}
}

func TestClassifyTCAA01_silent(t *testing.T) {
	score, _ := ClassifyTCAA01(TCAA01Input{
		ActiveAuthFailed: true, ActiveAuthSuccess: false, FailureSurfacedToCaller: false,
	})
	if score != ScoreSilent {
		t.Fatalf("got %d want silent", score)
	}
}

func TestClassifyTCAA01_logged(t *testing.T) {
	score, _ := ClassifyTCAA01(TCAA01Input{
		ActiveAuthFailed: true, ActiveAuthSuccess: true, FailureSurfacedToCaller: false,
	})
	if score != ScoreLogged {
		t.Fatalf("got %d want logged", score)
	}
}

func TestClassifyTCAA01_surfaced(t *testing.T) {
	score, _ := ClassifyTCAA01(TCAA01Input{
		ActiveAuthFailed: true, ActiveAuthSuccess: false, FailureSurfacedToCaller: true,
	})
	if score != ScoreSurfaced {
		t.Fatalf("got %d want surfaced", score)
	}
}

func TestClassifyTCTA01_peerUnsupported(t *testing.T) {
	score, meaning := ClassifyTCTA01(TCTA01Input{PeerUnsupported: true})
	if score != ScoreSurfaced {
		t.Fatalf("got %d want surfaced(unsupported)", score)
	}
	if meaning == "" {
		t.Fatal("expected unsupported meaning")
	}
}

func TestClassifyTCTA01_silent(t *testing.T) {
	score, _ := ClassifyTCTA01(TCTA01Input{
		TerminalAuthFailed: true, TerminalAuthSuccess: false, FailureSurfacedToCaller: false,
	})
	if score != ScoreSilent {
		t.Fatalf("got %d want silent", score)
	}
}

func TestClassifyTCTA01_logged(t *testing.T) {
	score, _ := ClassifyTCTA01(TCTA01Input{
		TerminalAuthFailed: true, TerminalAuthSuccess: true, FailureSurfacedToCaller: false,
	})
	if score != ScoreLogged {
		t.Fatalf("got %d want logged", score)
	}
}

func TestClassifyTCTA01_surfaced(t *testing.T) {
	score, _ := ClassifyTCTA01(TCTA01Input{
		TerminalAuthFailed: true, TerminalAuthSuccess: false, FailureSurfacedToCaller: true,
	})
	if score != ScoreSurfaced {
		t.Fatalf("got %d want surfaced", score)
	}
}

func TestClassifyTCEAC01_peerUnsupported(t *testing.T) {
	score, _ := ClassifyTCEAC01(TCEAC01Input{PeerUnsupported: true})
	if score != ScoreSurfaced {
		t.Fatalf("got %d want surfaced(unsupported)", score)
	}
}

func TestClassifyTCEAC01_silentProtectedDG(t *testing.T) {
	score, _ := ClassifyTCEAC01(TCEAC01Input{
		EACFailed: true, EACSuccess: false, ProtectedDGAccessible: true, FailureSurfacedToCaller: false,
	})
	if score != ScoreSilent {
		t.Fatalf("got %d want silent", score)
	}
}

func TestClassifyTCEAC01_loggedNoDG(t *testing.T) {
	score, _ := ClassifyTCEAC01(TCEAC01Input{
		EACFailed: true, EACSuccess: false, ProtectedDGAccessible: false, FailureSurfacedToCaller: false,
	})
	if score != ScoreLogged {
		t.Fatalf("got %d want logged", score)
	}
}

func TestClassifyTCEAC01_surfaced(t *testing.T) {
	score, _ := ClassifyTCEAC01(TCEAC01Input{
		EACFailed: true, EACSuccess: false, ProtectedDGAccessible: true, FailureSurfacedToCaller: true,
	})
	if score != ScoreSurfaced {
		t.Fatalf("got %d want surfaced", score)
	}
}

func TestObservabilityVectorsJSON(t *testing.T) {
	raw, err := os.ReadFile("../testdata/observability-vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var rows []struct {
		ID       string `json:"id"`
		Mech     string `json:"mechanism"`
		Input    map[string]any `json:"input"`
		Expected int    `json:"expected_score"`
	}
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		row := row
		t.Run(row.ID, func(t *testing.T) {
			var score Score
			switch row.Mech {
			case "TC-AC-01":
				bacErr, _ := row.Input["bac_err"].(string)
				score, _ = ClassifyTCAC01(TCAC01Input{
					PaceFailed:           row.Input["pace_failed"].(bool),
					BacSuccess:           row.Input["bac_success"].(bool),
					BacErr:               bacErr,
					PaceSurfacedToCaller: row.Input["pace_surfaced_to_caller"].(bool),
				})
			case "TC-CA-01":
				score, _ = ClassifyTCCA01(TCCA01Input{
					ChipAuthFailed:          row.Input["chip_auth_failed"].(bool),
					ChipAuthSuccess:         row.Input["chip_auth_success"].(bool),
					SessionContinueOK:       row.Input["session_continue_ok"].(bool),
					FailureSurfacedToCaller: row.Input["failure_surfaced_to_caller"].(bool),
				})
			case "TC-AA-01":
				score, _ = ClassifyTCAA01(TCAA01Input{
					ActiveAuthFailed:        row.Input["active_auth_failed"].(bool),
					ActiveAuthSuccess:       row.Input["active_auth_success"].(bool),
					FailureSurfacedToCaller: row.Input["failure_surfaced_to_caller"].(bool),
				})
			default:
				t.Fatalf("unsupported %s", row.Mech)
			}
			if int(score) != row.Expected {
				t.Fatalf("got %d want %d", score, row.Expected)
			}
		})
	}
}
