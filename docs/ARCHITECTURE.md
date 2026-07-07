# Harness architecture

```
JMRTD chip simulator (synthetic profiles in profiles/)
        ↕ APDU trace capture → logs/{run_id}/
Negotiation driver → drivers/jmrtd | drivers/gmrtd | drivers/pymrtd
        ↓
classifier/ → Observability Score 0/1/2
        ↓
Aggregate consistency % per tuple (finding if ≥95% over N=100)
```

## Observability Score

| Score | Meaning |
| --- | --- |
| 0 | Silent — failure/downgrade not surfaced to caller |
| 1 | Logged — visible in trace/logs only |
| 2 | Surfaced — explicit error/result to integrator |

Scored per `(library × mechanism × condition)`; never mean-aggregated across libraries.

## Middleware (§VIII)

Go wrapper on gmrtd enforcing explicit-reject on targeted downgrade; re-run full suite N=100 before/after.
