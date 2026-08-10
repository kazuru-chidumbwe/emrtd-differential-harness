.PHONY: smoke suite suite-paper paper repro test test-middleware verify-manifest verify-gmrtd-pin preflight-locked-runs package-locked-runs

verify-gmrtd-pin:
	bash scripts/verify_gmrtd_pin.sh

smoke: verify-gmrtd-pin
	bash scripts/quick_test.sh

suite:
	MANIFEST=suites/ac-01-wire.json bash scripts/run_suite.sh

suite-paper:
	MANIFEST=suites/paper-matrix.json bash scripts/run_suite.sh

paper:
	bash scripts/paper.sh

repro:
	bash scripts/independent-repro.sh

test-classifier:
	go test ./classifier/...
	python3 -c "import sys; sys.path.insert(0,'classifier'); from observability import TCAC01Outcome, classify_tc_ac_01, ObservabilityScore; assert classify_tc_ac_01(TCAC01Outcome(True,True,None,False))==ObservabilityScore.SILENT"
	python3 -m pip install -q -r classifier/requirements-verify.txt
	cd classifier && python3 -m unittest test_schemas test_observability test_observability_vectors -v

test-middleware:
	go test ./middleware/...

test: test-classifier test-middleware

verify-manifest:
	@test -n "$(LOG_DIR)" || (echo "usage: make verify-manifest LOG_DIR=logs/suite-..." && exit 1)
	python3 classifier/verify_manifest.py "$(LOG_DIR)" --manifest suites/ac-01-wire.json

# Fail loudly if a staged locked-run tree still contains programme/lab provenance strings.
# Usage: make preflight-locked-runs STAGING=path/to/extracted-or-built-deposit
preflight-locked-runs:
	@test -n "$(STAGING)" || (echo "usage: make preflight-locked-runs STAGING=path/to/staging" && exit 1)
	python3 scripts/preflight_banned_terms.py "$(STAGING)"

# Hard release-build gate: banned-term + abs-path + schema scan, then SemVer-named zip.
# Usage: make package-locked-runs STAGING=path/to/staging VERSION=1.0.6
package-locked-runs:
	@test -n "$(STAGING)" || (echo "usage: make package-locked-runs STAGING=... VERSION=1.0.6" && exit 1)
	@test -n "$(VERSION)" || (echo "usage: make package-locked-runs STAGING=... VERSION=1.0.6" && exit 1)
	bash scripts/package_locked_runs.sh "$(STAGING)" "$(VERSION)"
