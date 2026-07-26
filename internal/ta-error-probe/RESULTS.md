# TA forced-failure error propagation (2026-07-25)

Libraries: **JMRTD 0.8.6** (has TA). **gmrtd**: no TA — not comparable.

## JMRTD APDU-level (`APDULevelEACTACapable`)

| Probe | Forced failure | Behavior |
|---|---|---|
| T1 | `sendPSOExtendedLengthMode` throws SW=`6982` | `CardServiceException` propagates |
| T2 | `sendMutualAuthenticate` throws SW=`6982` | `CardServiceException` propagates |
| T4 | gmrtd | **unsupported** — no TA package |

Methods: `sendMSESetDST`, `sendPSOExtendedLengthMode`, `sendMSESetATExtAuth`, `sendGetChallenge`, `sendMutualAuthenticate`.

## Implication

TA APDU failures are loud at the JMRTD capable interface (same family as CA MSE fail). Harness TC-TA-01 uses SW-proxy on first PSO (honest CA-style), not full CVC crypto. Full `doEACTA` / synthetic CVC chain remains out of scope.
