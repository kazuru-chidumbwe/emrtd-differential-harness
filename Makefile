.PHONY: smoke suite suite-paper paper repro test test-middleware verify-manifest

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

test-middleware:
	go test ./middleware/...

test: test-classifier test-middleware

verify-manifest:
	@test -n "$(LOG_DIR)" || (echo "usage: make verify-manifest LOG_DIR=logs/suite-..." && exit 1)
	python3 classifier/verify_manifest.py "$(LOG_DIR)" --manifest suites/ac-01-wire.json
