# Harness architecture

```
Scenario (suite manifest)
    ↓
Execution (drivers + middleware policy)
    ↓
Observation (classifier → score 0/1/2)
    ↓
Finding (≥95% over N=100)
    ↓
Artifact (artifact-manifest.json — canonical)
```

See [SCHEMA.md](SCHEMA.md) for frozen v1 schemas and figure ID rules (`FIG-01`, not “Figure 2”).

## Canonical artifact

**`artifact-manifest.json`** is the primary published object. Each `FIG-xx` entry bundles scenario, execution, observation, finding, `artifact_refs`, and `bundle_sha256`. `TABLE-01` and `SUMMARY-01` are derived views.

## Middleware (contribution layer)

Baseline drivers observe library behaviour. Mitigated drivers apply `middleware/` explicit-reject policy. Compare artifacts — not rewritten tests.

## Commands

| Command | Purpose |
| --- | --- |
| `make smoke` | Single-run gate |
| `make suite` | Blog wire manifest (`FIG-01`…`FIG-04`) |
| `make suite-paper` | Full paper matrix |
| `make paper` | CI pipeline — fails if tests/smoke/suite incomplete |

## Methodology (frozen)

Repeating each deterministic profile N=100 demonstrates harness stability and result reproducibility rather than estimating behavioural variance.
