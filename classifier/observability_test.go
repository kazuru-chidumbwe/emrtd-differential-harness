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
