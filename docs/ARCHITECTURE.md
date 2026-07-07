# Harness architecture

```
profiles/ + suites/*.json manifests
        ↓
classifier/run_suite.py → drivers (baseline + mitigated)
        ↓
middleware/ (§VIII explicit-reject) — contribution layer
        ↓
logs/{run_id}.json  (provenance block per run)
        ↓
classifier/aggregate.py → summary + artifact-manifest
        ↓
Published figure_id → artifact_refs (SHA-256 linked)
```

## Middleware first

The harness separates **observation** (baseline drivers) from **mitigation** (`middleware/` + mitigated drivers). Baseline proves silent downgrade; mitigated proves the fix surfaces failure (score 0→2).

See [`middleware/README.md`](../middleware/README.md).

## Suite manifests

| Manifest | Scope | Command |
| --- | --- | --- |
| `suites/ac-01-wire.json` | Blog wire tier (AC-01 × gmrtd/jmrtd × baseline/mitigated) | `make suite` |
| `suites/paper-matrix.json` | Paper matrix (+ CA, offline PA) | `make suite-paper` |

## Observability Score

Scored per `(test_case × library × variant)`; never mean-aggregated across libraries. Finding threshold: **≥95% over N=100** (deterministic fixed profile = reproducibility proof).

## Provenance

Every run records commit, profile SHA-256, suite seed, driver, variant. See [`docs/PROVENANCE.md`](PROVENANCE.md).
