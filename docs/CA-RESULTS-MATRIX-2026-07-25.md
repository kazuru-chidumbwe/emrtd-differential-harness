# TC-CA-01 CA skew matrix — 40-run pin + SM-session continue-check

**Suite:** `ca-01-sweep` (10 profiles × 2 libraries × 2 variants = **40 runs**)

**Axes (harness-consumed):**
| Axis | Values |
| --- | --- |
| `skew_direction` | `v2chip_v1term` (`ca_v2_chip_v1_terminal_skew`), `v1chip_v2term` (`ca_v1_chip_v2_terminal_skew`) |
| `ca_sw` | `6FFF`, `6300`, `6982`, `6A88`, `6A80` |

**Injection:** first CA MSE:Set AT (`INS=0x22`, `P1=0x41`, `P2=0xA4`) returns `ca_sw`. `ca_fail_on` is documentation-only (simulators do not branch on it). DG14 fixture `testdata/dg14/ca-v2-sample.hex` advertises CA for both direction labels; true CA version negotiation is not simulated.

**SM-session continue-check:** After BAC, chip and reader share session keys. After CA MSE reject, drivers probe via **SM-wrapped** READ BINARY (`nfc.ReadBinaryFromOffset` / JMRTD `SecureMessagingWrapper`). Unprotected B0 returns `6987` (Expected SM data objects missing). `session_continue_ok` requires MAC-valid SM success. This is crypto-real under post-BAC keys but still **harness-probed** (libraries do not spontaneously READ BINARY after CA fail) — distinct from AC-01 `BacSuccess` (library-path emergent) and from EAC `ProtectedDGAccessible` (by construction). Mitigated arms still force `FailureSurfacedToCaller=true`.

**Generator:** `python3 profiles/generate_ca_sweep.py` → `profiles/ca-sweep/` + `index.json`

**Smoke anchors (unchanged paths):**
- `profiles/ca-v1-v2-skew.json` — forward skew, SW `6FFF`
- `profiles/ca-v2-terminal-v1.json` — reverse skew label, SW `6985`

**Observed Observability Score (test-server):**

| Library | Variant | Score | Cells |
| --- | --- | ---: | ---: |
| gmrtd | baseline | 0 (silent) | 10 / 10 |
| gmrtd | mitigated | 2 (surfaced) | 10 / 10 |
| JMRTD | baseline | 0 (silent) | 10 / 10 |
| JMRTD | mitigated | 2 (surfaced) | 10 / 10 |

**Lab pin:**

| Field | Value |
| --- | --- |
| Host | `test-server` (Ubuntu; Go 1.25.0; OpenJDK 17; Python 3.12.3) |
| Log dir | `logs/ca-01-sweep-full-sm-20260726T-test-server` |
| `bundle_sha256` | `43f9bd1db0c26709d6a33015fdb82979a0e4f097911123cbac641c1d3eb3c050` |
| Runner | `scripts/run_ca_sweep_full.sh` |
| Retired stub-era | `logs/ca-01-sweep-full-r1-test-server` / `52af54e5…` |

**Reproduce:**

```bash
bash scripts/run_ca_sweep_full.sh
# or gmrtd-only: bash scripts/run_ca_sweep_gmrtd.sh
```

**Notes:**
- Option A pin: `org.jmrtd:jmrtd:0.8.6` (see `docs/JMRTD-PIN.md`).
- Corroborating depth for the paper is this 40-run matrix, not a second 200-run PACE factorial.
- Supersedes prior by-construction-only pin `ce12ee2c…` (pre–continue-check).
