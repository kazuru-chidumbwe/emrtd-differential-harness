# Synthetic TA/EAC PKI placeholders (2026-07-25)

This directory holds **requirements and placeholders** for CVCA→DV→IS material.
Smoke cells (TC-TA-01 / TC-EAC-01) use **APDU SW-proxy** and do not yet consume real CVC DER.

## Planned artifacts (follow-on)

| File | Role |
|---|---|
| `cvca.cvcert` | Country Verifying CA |
| `dv.cvcert` | Document Verifier |
| `is.cvcert` | Inspection System (role IS) |
| `is-private.pkcs8` | IS private key matching `is.cvcert` |
| `README` (this file) | Build steps via `cert-cvc` / JMRTD `CVCertificateBuilder` |

## Current status

- **Not generated** — full CVC crypto gated on programme need for `doEACTA` end-to-end.
- SW-proxy profiles are sufficient for observability scoring of TA APDU failure / EAC access downgrade.
- Full CVC generation remains a follow-on; not required for current paper cells.
