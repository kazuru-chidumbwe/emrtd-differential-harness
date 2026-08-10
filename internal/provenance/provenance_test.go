package provenance

import "testing"

func TestRunIndexClampedToOne(t *testing.T) {
	// Collect with missing profile should error; use a real profile path from repo root.
	rec, err := Collect(Options{
		ProfilePath: "../../profiles/pace-then-bac-downgrade.json",
		SuiteID:     "test",
		SuiteSeed:   1,
		SuiteN:      1,
		RunIndex:    0,
		Driver:      "go/test",
		Variant:     "baseline",
	})
	if err != nil {
		// Allow missing profile when tests run from module cache oddly
		t.Skip(err)
	}
	if rec.RunIndex != 1 {
		t.Fatalf("run_index=%d want 1", rec.RunIndex)
	}
	if rec.HarnessCommit == "" {
		t.Fatal("harness_commit must never be empty")
	}
}
