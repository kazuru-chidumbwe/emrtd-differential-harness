package output

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kazuru-chidumbwe/emrtd-differential-harness/internal/provenance"
)

// Meta is common run metadata for all drivers.
type Meta struct {
	RunID      string             `json:"run_id"`
	TestCase   string             `json:"test_case"`
	Library    string             `json:"library"`
	Mechanism  string             `json:"mechanism"`
	Condition  string             `json:"condition"`
	Tier       string             `json:"tier"`
	Variant    string             `json:"variant"`
	FigureID   string             `json:"figure_id,omitempty"`
	Provenance provenance.Record  `json:"provenance"`
}

// WriteJSON writes a run artifact to logDir/{runID}.json.
func WriteJSON(logDir, runID string, v any) error {
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(logDir, runID+".json"), raw, 0o644)
}
