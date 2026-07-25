# TC-AA-01 Active Authentication — 20-run SW grid (2026-07-25 / R1)

Forced AA failure after BAC. Grounding: [internal/aa-error-probe/RESULTS.md](../internal/aa-error-probe/RESULTS.md).

**Suite:** `aa-01-sweep` (5 profiles × 2 libraries × 2 variants = **20 runs**)

**Axes:**
| Axis | Values |
| --- | --- |
| `aa_sw` | `6300`, `6982`, `6A80`, `6A88`, `6985` |

`6FFF` is **not** in this grid: JMRTD soft-succeeds INTERNAL AUTHENTICATE on `6FFF`, so that SW is not a forced-failure cell under Option A.

**Injection:** BAC succeeds; first INTERNAL AUTHENTICATE (`INS=0x88`) returns `aa_sw`.

**Observability Score (by construction):** Baseline drivers set `FailureSurfacedToCaller=false` (naive host). Mitigated drivers set it true via explicit-reject middleware. Unlike TC-AC-01 (`BacSuccess`), TC-CA-01 (post-CA `session_continue_ok`), and TC-EAC-01 (`ProtectedDGAccessible`), AA has **no emergent continue-check** — scores are the designed naive-host / mitigated outcome, not a discovered library difference.

| Library | Variant | Score | Cells |
| --- | --- | ---: | ---: |
| gmrtd | baseline | 0 | 5 / 5 |
| gmrtd | mitigated | 2 | 5 / 5 |
| JMRTD | baseline | 0 | 5 / 5 |
| JMRTD | mitigated | 2 | 5 / 5 |

**Lab pin (test-server):**

| Field | Value |
| --- | --- |
| Host | `test-server` (Ubuntu; Go 1.25.0; OpenJDK 17; Python 3.12.3) |
| Log dir | `logs/aa-01-sweep-full-r1-test-server` |
| `bundle_sha256` | `ba02987997bcd18b2ef8d5c440ce5c051716060b69b6b8a29bcaa03ecfa3272c` |
| Runner | `scripts/run_aa_sweep_full.sh` |
| Generator | `profiles/generate_aa_sweep.py` → `profiles/aa-sweep/` |

**Reproduce:**

```bash
python3 profiles/generate_aa_sweep.py
LOG_DIR=logs/aa-01-sweep-full-r1-test-server bash scripts/run_aa_sweep_full.sh
```

**Notes:**
- Bilateral corroborating depth for the paper — not a second 200-run PACE factorial.
- Smoke anchors remain under `profiles/aa-internal-auth-reject.json` / `suites/aa-01-smoke.json`.
