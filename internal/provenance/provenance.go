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
	HarnessCommit string `json:"harness_commit"`
	HarnessDirty  bool   `json:"harness_dirty"`
	GoVersion     string `json:"go_version,omitempty"`
	JavaVersion   string `json:"java_version,omitempty"`
	PythonVersion string `json:"python_version,omitempty"`
	ProfilePath   string `json:"profile_path"`
	ProfileSHA256 string `json:"profile_sha256"`
	SuiteID       string `json:"suite_id,omitempty"`
	SuiteSeed     int    `json:"suite_seed,omitempty"`
	SuiteN        int    `json:"suite_n,omitempty"`
	RunIndex      int    `json:"run_index"`
	Driver        string `json:"driver"`
	Variant       string `json:"variant"`
	Middleware    string `json:"middleware,omitempty"`
	CapturedAtUTC string `json:"captured_at_utc"`
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
	commit, dirty := resolveHarnessCommit()

	profSHA, err := fileSHA256(opts.ProfilePath)
	if err != nil {
		return Record{}, fmt.Errorf("profile hash: %w", err)
	}

	runIndex := opts.RunIndex
	if runIndex < 1 {
		runIndex = 1
	}

	return Record{
		HarnessCommit: commit,
		HarnessDirty:  dirty,
		GoVersion:     runtime.Version(),
		JavaVersion:   opts.JavaVersion,
		PythonVersion: opts.PythonVer,
		ProfilePath:   repoRelativePath(opts.ProfilePath),
		ProfileSHA256: profSHA,
		SuiteID:       opts.SuiteID,
		SuiteSeed:     opts.SuiteSeed,
		SuiteN:        opts.SuiteN,
		RunIndex:      runIndex,
		Driver:        opts.Driver,
		Variant:       opts.Variant,
		Middleware:    opts.Middleware,
		CapturedAtUTC: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// repoRelativePath stores profile_path relative to the harness module root when possible.
func repoRelativePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	root, err := findModuleRoot()
	if err != nil {
		return filepath.ToSlash(abs)
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return filepath.ToSlash(abs)
	}
	return filepath.ToSlash(rel)
}

func resolveHarnessCommit() (string, bool) {
	if env := strings.TrimSpace(os.Getenv("EMRTD_HARNESS_COMMIT")); env != "" {
		return env, false
	}
	commit, dirty, err := gitHead()
	if err != nil || commit == "" {
		return "unknown", false
	}
	return commit, dirty
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
	commit := strings.TrimSpace(string(out))
	if commit == "" {
		return "", false, fmt.Errorf("empty git rev-parse HEAD")
	}
	status, _ := exec.Command("git", "-C", root, "status", "--porcelain").Output()
	return commit, len(strings.TrimSpace(string(status))) > 0, nil
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
