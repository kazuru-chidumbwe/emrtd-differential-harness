# TC-CA-01 CA skew matrix — profile lock + 40-run pin (2026-07-25)

**Suite:** `ca-01-sweep` (10 profiles × 2 libraries × 2 variants = **40 runs**)

**Axes (harness-consumed):**
| Axis | Values |
| --- | --- |
| `skew_direction` | `v2chip_v1term` (`ca_v2_chip_v1_terminal_skew`), `v1chip_v2term` (`ca_v1_chip_v2_terminal_skew`) |
| `ca_sw` | `6FFF`, `6300`, `6982`, `6A88`, `6A80` |

**Injection:** first CA MSE:Set AT (`INS=0x22`, `P1=0x41`, `P2=0xA4`) returns `ca_sw`. `ca_fail_on` is documentation-only (simulators do not branch on it). DG14 fixture `testdata/dg14/ca-v2-sample.hex` advertises CA for both direction labels; true CA version negotiation is not simulated.

**Generator:** `python3 profiles/generate_ca_sweep.py` → `profiles/ca-sweep/` + `index.json`

**Smoke anchors (unchanged paths):**
- `profiles/ca-v1-v2-skew.json` — forward skew, SW `6FFF`
- `profiles/ca-v2-terminal-v1.json` — reverse skew label, SW `6985`

**Observed Observability Score (lab host Docker `emrtd-harness:paper-2026-07-25`, 2026-07-25):**

| Library | Variant | Score | Cells |
| --- | --- | ---: | ---: |
| gmrtd | baseline | 0 (silent) | 10 / 10 |
| gmrtd | mitigated | 2 (surfaced) | 10 / 10 |
| JMRTD | baseline | 0 (silent) | 10 / 10 |
| JMRTD | mitigated | 2 (surfaced) | 10 / 10 |

**Lab pin:**

| Field | Value |
| --- | --- |
| Log dir | `logs/ca-01-sweep-full-20260725T204638Z` |
| Wall clock | ~148 s (container) |
| `bundle_sha256` | `ce12ee2ca83e0c7f6137802a4386df01a7a2996ac5f883441d190d1d8ef9b569` |
| Runner | `scripts/run_ca_sweep_full.sh` |

**Reproduce:**

```bash
bash scripts/run_ca_sweep_full.sh
# or gmrtd-only: bash scripts/run_ca_sweep_gmrtd.sh
```

**Notes:**
- Option A pin: `org.jmrtd:jmrtd:0.8.6` (see `docs/JMRTD-PIN.md`).
- Corroborating depth for the paper is this 40-run matrix, not a second 200-run PACE factorial.
