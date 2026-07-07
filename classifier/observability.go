// Package classifier maps run outcomes to Observability Score (0/1/2).
// Contract is shared with classifier/observability.py and drivers/jmrtd Observability.java.
package classifier

// Score is the Observability Score for a single (library × mechanism × condition) run.
type Score int

const (
	ScoreSilent   Score = 0 // failure/downgrade not surfaced to naive caller
	ScoreLogged   Score = 1 // visible in trace/session only
	ScoreSurfaced Score = 2 // explicit error/result at integrator boundary
)

// TCAC01Input is the shared outcome contract for TC-AC-01 wire-tier runs.
type TCAC01Input struct {
	PaceFailed           bool
	BacSuccess           bool
	BacErr               string
	PaceSurfacedToCaller bool
}

// ClassifyTCAC01 scores whether a PACE→BAC downgrade is visible to a naive caller.
func ClassifyTCAC01(in TCAC01Input) (Score, string) {
	if in.PaceFailed && in.BacSuccess && in.BacErr == "" && !in.PaceSurfacedToCaller {
		return ScoreSilent, "silent — PACE failure not surfaced to caller; BAC succeeded"
	}
	if in.PaceFailed && !in.PaceSurfacedToCaller {
		return ScoreLogged, "logged — PACE failure visible in session/trace only"
	}
	return ScoreSurfaced, "surfaced — explicit error at caller boundary"
}

// Int returns the JSON-serializable score value.
func (s Score) Int() int { return int(s) }
