# TC-TA-01 / TC-EAC-01 results lock (2026-07-25)

JMRTD-asymmetric. gmrtd: **unsupported** (no TA). Grounding: [ta-error-probe/RESULTS.md](../internal/ta-error-probe/RESULTS.md). Full synthetic CVC / end-to-end `doEACTA` remains out of scope (SW-proxy arms only).

## TC-TA-01

Injection: first **PSO:Verify Certificate** (`INS=0x2A`) → `ta_sw` (default `6982`).

| Library | Variant | Expected score | Notes |
|---|---|---:|---|
| JMRTD | baseline | **0** | naive host swallows |
| gmrtd | unsupported | **2** (meaning: unsupported) | no TA API |

## TC-EAC-01

TA PSO fails; **READ BINARY** of synthetic protected DG still `9000`.

| Library | Variant | Expected | Notes |
|---|---|---:|---|
| JMRTD | baseline | **0** | EAC fail + DG reachable + not surfaced |
| gmrtd | unsupported | **2** (unsupported) | cannot complete EAC |

## Honesty

SW-proxy only — not full `doEACCA`+`doEACTA` with real CVCA/DV/IS. PKI placeholders under `testdata/cvc/`.
