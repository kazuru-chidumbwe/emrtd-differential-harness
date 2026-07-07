.PHONY: smoke suite suite-paper repro test test-middleware

smoke:
	bash scripts/quick_test.sh

suite:
	MANIFEST=suites/ac-01-wire.json bash scripts/run_suite.sh

suite-paper:
	MANIFEST=suites/paper-matrix.json bash scripts/run_suite.sh

repro:
	bash scripts/independent-repro.sh

test-classifier:
	go test ./classifier/...
	python3 -c "import sys; sys.path.insert(0,'classifier'); from observability import TCAC01Outcome, classify_tc_ac_01, ObservabilityScore; assert classify_tc_ac_01(TCAC01Outcome(True,True,None,False))==ObservabilityScore.SILENT"

test-middleware:
	go test ./middleware/...

test: test-classifier test-middleware
