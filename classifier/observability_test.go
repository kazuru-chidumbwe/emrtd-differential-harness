package classifier

import "testing"

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
		ChipAuthFailed: true, ChipAuthSuccess: false, FailureSurfacedToCaller: false,
	})
	if score != ScoreSilent {
		t.Fatalf("got %d want silent", score)
	}
}

func TestClassifyTCCA01_logged(t *testing.T) {
	// Failed but ChipAuthSuccess still true → logged branch (not the silent triple).
	score, _ := ClassifyTCCA01(TCCA01Input{
		ChipAuthFailed: true, ChipAuthSuccess: true, FailureSurfacedToCaller: false,
	})
	if score != ScoreLogged {
		t.Fatalf("got %d want logged", score)
	}
}

func TestClassifyTCCA01_surfaced(t *testing.T) {
	score, _ := ClassifyTCCA01(TCCA01Input{
		ChipAuthFailed: true, ChipAuthSuccess: false, FailureSurfacedToCaller: true,
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
