# AA forced-failure error propagation (2026-07-25)

Empirical probes under `internal/aa-error-probe/`. Libraries: **gmrtd** (vendored) and **JMRTD 0.8.6**.

**Harness wiring (2026-07-25):** TC-AA-01 profiles/drivers live — see [docs/AA-RESULTS-2026-07-25.md](../../docs/AA-RESULTS-2026-07-25.md).

## gmrtd

| Probe | Forced failure | Library return | Integration boundary |
|---|---|---|---|
| A | INTERNAL AUTHENTICATE SW=`6982` | `DoActiveAuth` → `err` non-nil, `result.Success=false` | — |
| B | Empty AA response / bad RSA ciphertext | `ValidateActiveAuthSignature` → `err` non-nil, `Success=false` | — |
| C | Same as A, reader-shaped | `Session.ActiveAuthErr` set | `performChipAuthentication` still returns `nil` (“errors are just recorded”) |
| D | CLI | — | `surfacePaceErr` exists; **no** `surfaceActiveAuthErr` |

**Verdict:** Protocol layer surfaces AA failure correctly. Host/reader layer can still silently continue — same record-only pattern as CA/PACE.

## JMRTD 0.8.6

| Probe | Forced failure | Behavior |
|---|---|---|
| J1 | `sendInternalAuthenticate` throws `CardServiceException(SW=6982)` | `AAProtocol.doAA` throws `CardServiceProtocolException` (propagates) |
| J2 | APDU OK, garbage response bytes | `doAA` returns `AAResult` — **no signature verification** in protocol |
| J3 | Challenge length ≠ 8 | `IllegalArgumentException` wrapped by catch-all into `CardServiceProtocolException` |

**Verdict:** Chip/APDU failure propagates as exception. Cryptographic AA failure (bad signature with SW=9000) is **caller-side** — silent at `doAA` unless host verifies `AAResult`.

## Implication for “now” list

AA is **same shape as CA** for harness work (wire `DoActiveAuth`/`doAA`, inject failure, score swallow vs surface), **not** a greenfield protocol. Grounding supports that classification.

**Still separate:** CA→AA assurance downgrade (caller thinks CA-grade when only AA, or vice versa) is not covered by “does AA fail loudly.”
