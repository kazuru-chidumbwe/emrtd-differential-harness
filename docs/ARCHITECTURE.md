# Harness architecture

```
JMRTD chip simulator (synthetic profiles in profiles/)
        ↕ APDU trace capture → logs/{run_id}/
Negotiation driver → cmd/tc-ac-01 | cmd/tc-ac-01-mitigated | drivers/jmrtd
        ↓
classifier/ → Observability Score 0/1/2 (shared Go + Python + Java contract)
        ↓
scripts/run_suite.sh → classifier/aggregate.py
        ↓
summary-{timestamp}.json + .md (finding if ≥95% over N=100)
```

## Observability Score

| Score | Meaning |
| --- | --- |
| 0 | Silent — failure/downgrade not surfaced to caller |
| 1 | Logged — visible in trace/logs only |
| 2 | Surfaced — explicit error/result to integrator |

Scored per `(library × mechanism × condition × variant)`; never mean-aggregated across libraries.

TC-AC-01 uses `ClassifyTCAC01` / `classify_tc_ac_01` with fields `pace_failed`, `bac_success`, `bac_err`, `pace_surfaced_to_caller`.

## Suite runner

```bash
make suite              # SUITE_N=100 default
SUITE_N=10 make suite   # quick local check
```

Runs gmrtd baseline, gmrtd mitigated (middleware), and jmrtd baseline — each N times — then aggregates.

## Middleware (§VIII)

`middleware/negotiate.go` wraps gmrtd PACE+BAC: when `AllowBACFallback=false`, PACE failure returns `ErrPaceFailed` instead of silent BAC fallback. Driver: `cmd/tc-ac-01-mitigated`.

## Reproducibility note

Fixed synthetic profile → deterministic outcomes across N runs. N=100 validates harness stability, not input-population variance. See README limitations.
