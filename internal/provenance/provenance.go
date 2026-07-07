package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Record is embedded in every run artifact for paper-grade reproducibility.
type Record struct {
	HarnessCommit   string `json:"harness_commit"`
	HarnessDirty    bool   `json:"harness_dirty"`
	GoVersion       string `json:"go_version,omitempty"`
	JavaVersion     string `json:"java_version,omitempty"`
	PythonVersion   string `json:"python_version,omitempty"`
	ProfilePath     string `json:"profile_path"`
	ProfileSHA256   string `json:"profile_sha256"`
	SuiteID         string `json:"suite_id,omitempty"`
	SuiteSeed       int    `json:"suite_seed,omitempty"`
	SuiteN          int    `json:"suite_n,omitempty"`
	RunIndex        int    `json:"run_index,omitempty"`
	Driver          string `json:"driver"`
	Variant         string `json:"variant"`
	Middleware      string `json:"middleware,omitempty"`
	CapturedAtUTC   string `json:"captured_at_utc"`
}

// Options configures provenance collection for a single run.
type Options struct {
	ProfilePath string
	SuiteID     string
	SuiteSeed   int
	SuiteN      int
	RunIndex    int
	Driver      string
	Variant     string
	Middleware  string
	JavaVersion string
	PythonVer   string
}

// Collect builds a provenance record from the environment and profile file.
func Collect(opts Options) (Record, error) {
	commit, dirty, err := gitHead()
	if err != nil {
		commit = "unknown"
	}

	profSHA, err := fileSHA256(opts.ProfilePath)
	if err != nil {
		return Record{}, fmt.Errorf("profile hash: %w", err)
	}

	return Record{
		HarnessCommit: commit,
		HarnessDirty:  dirty,
		GoVersion:     runtime.Version(),
		JavaVersion:   opts.JavaVersion,
		PythonVersion: opts.PythonVer,
		ProfilePath:   filepath.ToSlash(opts.ProfilePath),
		ProfileSHA256: profSHA,
		SuiteID:       opts.SuiteID,
		SuiteSeed:     opts.SuiteSeed,
		SuiteN:        opts.SuiteN,
		RunIndex:      opts.RunIndex,
		Driver:        opts.Driver,
		Variant:       opts.Variant,
		Middleware:    opts.Middleware,
		CapturedAtUTC: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func gitHead() (string, bool, error) {
	root, err := findModuleRoot()
	if err != nil {
		return "", false, err
	}
	out, err := exec.Command("git", "-C", root, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", false, err
	}
	status, _ := exec.Command("git", "-C", root, "status", "--porcelain").Output()
	return strings.TrimSpace(string(out)), len(strings.TrimSpace(string(status))) > 0, nil
}

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func fileSHA256(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}
