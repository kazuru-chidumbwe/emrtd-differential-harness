# Middleware (§VIII) — explicit-reject on negotiation downgrade

Original engineering contribution: a thin wrapper that prevents **silent PACE→BAC fallback** (and analogous CA failures) unless the integrator explicitly opts in.

## Problem

gmrtd records PACE failure on the session object but still proceeds to BAC when secure messaging is not established. A naive caller that checks only `ReadDocument()` success observes Observability Score **0** (silent downgrade).

## Solution

`middleware/negotiate.go` enforces:

```
PACE fails → return ErrPaceFailed (score 2) unless AllowBACFallback=true
```

`middleware/ca.go` applies the same contract to EAC-CA:

```
Chip auth fails → return ErrChipAuthFailed unless AllowContinue=true
```

## Drivers

| Driver | Variant | Middleware |
| --- | --- | --- |
| `cmd/tc-ac-01` | baseline | none (library default) |
| `cmd/tc-ac-01-mitigated` | mitigated | `NegotiatePACEBAC{AllowBACFallback:false}` |
| `cmd/tc-ca-01-mitigated` | mitigated | `PerformChipAuth{AllowContinue:false}` |
| JMRTD `TcAc01MitigatedRunner` | mitigated | rethrow `PACEException` (no catch-and-continue) |

## Before/after evidence

Run the wire-tier suite manifest (`suites/ac-01-wire.json`):

- **baseline:** gmrtd/jmrtd → 100% silent (N=100, fixed synthetic profile)
- **mitigated:** gmrtd → 100% surfaced (0% silent)

Aggregated summaries link each published percentage to per-run JSON artifacts via `artifact_manifest` (see `docs/PROVENANCE.md`).

## Integration example

```go
sess := middleware.NegotiatePACEBAC(nfc, doc, pass, middleware.Options{
    AllowBACFallback: false, // explicit reject — recommended default for PACE-capable flows
})
if sess.SurfacedError != nil {
    return sess.SurfacedError
}
```

Opt-in fallback (documents integrator acceptance of downgrade risk):

```go
sess := middleware.NegotiatePACEBAC(nfc, doc, pass, middleware.Options{
    AllowBACFallback: true,
})
```
