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

// TCCA01Input is the shared outcome contract for TC-CA-01 wire-tier runs.
type TCCA01Input struct {
	ChipAuthFailed          bool
	ChipAuthSuccess         bool
	SessionContinueOK       bool // emergent: post-CA DG/session probe still succeeds
	FailureSurfacedToCaller bool
}

// ClassifyTCCA01 scores whether an EAC-CA failure is visible to a naive caller.
// Silent requires an emergent continue-check (SessionContinueOK), parallel to BacSuccess / ProtectedDGAccessible.
func ClassifyTCCA01(in TCCA01Input) (Score, string) {
	if in.ChipAuthFailed && !in.ChipAuthSuccess && in.SessionContinueOK && !in.FailureSurfacedToCaller {
		return ScoreSilent, "silent — chip auth failure; session/DG still usable without surfaced error"
	}
	if in.ChipAuthFailed && !in.FailureSurfacedToCaller {
		return ScoreLogged, "logged — chip auth failure visible in session/trace only"
	}
	return ScoreSurfaced, "surfaced — explicit error at caller boundary"
}

// TCAA01Input is the shared outcome contract for TC-AA-01 wire-tier runs.
type TCAA01Input struct {
	ActiveAuthFailed        bool
	ActiveAuthSuccess       bool
	FailureSurfacedToCaller bool
}

// ClassifyTCAA01 scores whether an Active Authentication failure is visible to a naive caller.
func ClassifyTCAA01(in TCAA01Input) (Score, string) {
	if in.ActiveAuthFailed && !in.ActiveAuthSuccess && !in.FailureSurfacedToCaller {
		return ScoreSilent, "silent — AA failure on session; step returns nil / caller swallows"
	}
	if in.ActiveAuthFailed && !in.FailureSurfacedToCaller {
		return ScoreLogged, "logged — AA failure visible in session/trace only"
	}
	return ScoreSurfaced, "surfaced — explicit error at caller boundary"
}

// TCTA01Input is the shared outcome contract for TC-TA-01 (JMRTD-asymmetric) runs.
type TCTA01Input struct {
	TerminalAuthFailed      bool
	TerminalAuthSuccess     bool
	FailureSurfacedToCaller bool
	PeerUnsupported         bool
}

// ClassifyTCTA01 scores TA failure visibility. PeerUnsupported forces an explicit non-score path.
func ClassifyTCTA01(in TCTA01Input) (Score, string) {
	if in.PeerUnsupported {
		return ScoreSurfaced, "unsupported — peer library has no TA implementation"
	}
	if in.TerminalAuthFailed && !in.TerminalAuthSuccess && !in.FailureSurfacedToCaller {
		return ScoreSilent, "silent — TA failure not surfaced to naive caller"
	}
	if in.TerminalAuthFailed && !in.FailureSurfacedToCaller {
		return ScoreLogged, "logged — TA failure visible in session/trace only"
	}
	return ScoreSurfaced, "surfaced — explicit error at caller boundary"
}

// TCEAC01Input is the shared outcome contract for TC-EAC-01 (CA+TA) runs.
type TCEAC01Input struct {
	EACFailed               bool
	EACSuccess              bool
	ProtectedDGAccessible   bool
	FailureSurfacedToCaller bool
	PeerUnsupported         bool
}

// ClassifyTCEAC01 scores EAC workflow / protected-DG access observability.
func ClassifyTCEAC01(in TCEAC01Input) (Score, string) {
	if in.PeerUnsupported {
		return ScoreSurfaced, "unsupported — peer library cannot complete EAC (no TA)"
	}
	// Silent downgrade: EAC/TA failed but protected DG still appears accessible / caller continues.
	if in.EACFailed && !in.EACSuccess && in.ProtectedDGAccessible && !in.FailureSurfacedToCaller {
		return ScoreSilent, "silent — EAC/TA failure; protected DG still reachable without surfaced error"
	}
	if in.EACFailed && !in.FailureSurfacedToCaller {
		return ScoreLogged, "logged — EAC failure visible in session/trace only"
	}
	return ScoreSurfaced, "surfaced — explicit error at caller boundary"
}
