# Success-path control — test-server runbook

Host: **test-server** only.

```bash
cd /path/to/emrtd-differential-harness
go test ./middleware/ -count=1
mkdir -p logs/suite-ac-01-success-path-control-$(date -u +%Y%m%dT%H%M%SZ)-test-server

# BAC-only × allow false/true
go run ./cmd/tc-ac-01-mitigated \
  -profile profiles/success-path/bac-only.json \
  -allow-bac-fallback=false -success-path \
  -suite-id ac-01-success-path-control -log-dir logs/sp-control

go run ./cmd/tc-ac-01-mitigated \
  -profile profiles/success-path/bac-only.json \
  -allow-bac-fallback=true -success-path \
  -suite-id ac-01-success-path-control -log-dir logs/sp-control

# PACE-fail × allow true
go run ./cmd/tc-ac-01-mitigated \
  -profile profiles/success-path/pace-fail-allow.json \
  -allow-bac-fallback=true -success-path \
  -suite-id ac-01-success-path-control -log-dir logs/sp-control
```

Then compute manifest SHA-256 of the log dir and update paper Evaluation §6.5 (private manuscript pack; not in this repo) and any matching manuscript block.
JMRTD success-path entries in the suite JSON may need runner flags for empty CardAccess — confirm on server before citing Java cells.
