# Stack drivers

| Directory | Stack | Interface |
| --- | --- | --- |
| `jmrtd/` | JMRTD + passport applet simulator | Java — APDU transceiver |
| `gmrtd/` | gmrtd | Go — `iso7816.NfcSession` mock |
| `pymrtd-offline/` | pymrtd | Python — `run_smoke.py` for TC-PA-* (offline tier) |

Each driver emits a normalized trace format consumed by `classifier/`.
