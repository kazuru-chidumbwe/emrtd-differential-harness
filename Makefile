.PHONY: smoke suite repro test-classifier

smoke:
	bash scripts/quick_test.sh

suite:
	bash scripts/run_suite.sh

repro:
	bash scripts/sponsor-repro.sh

test-classifier:
	go test ./classifier/...
	python3 -c "import sys; sys.path.insert(0,'classifier'); from observability import TCAC01Outcome, classify_tc_ac_01, ObservabilityScore, consistency_pct; o=TCAC01Outcome(True,True,None,False); assert classify_tc_ac_01(o)==ObservabilityScore.SILENT; assert consistency_pct([0,0,0], ObservabilityScore.SILENT)==100.0; print('python classifier OK')"
