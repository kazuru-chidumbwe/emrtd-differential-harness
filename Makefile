.PHONY: smoke suite suite-paper paper repro test test-middleware verify-manifest preflight-locked-runs

smoke:
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
	cd classifier && python3 -m unittest test_schemas -v

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
