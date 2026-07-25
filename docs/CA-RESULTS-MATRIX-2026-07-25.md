# TC-CA-01 CA skew matrix — profile lock (2026-07-25)

**Suite:** `ca-01-sweep` (10 profiles × 2 libraries × 2 variants = **40 runs** when fully executed)

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

**Expected Observability Score (by design, deterministic APDU replay):**

| Library | Variant | Expected score | Cells |
| --- | --- | ---: | ---: |
| gmrtd | baseline | 0 (silent) | 10 / 10 |
| gmrtd | mitigated | 2 (surfaced) | 10 / 10 |
| JMRTD | baseline | 0 (silent) | 10 / 10 |
| JMRTD | mitigated | 2 (surfaced) | 10 / 10 |

**Reproduce (gmrtd matrix):**

```bash
bash scripts/bootstrap-vendor.sh
export GOTOOLCHAIN=auto
bash scripts/run_ca_sweep_gmrtd.sh
```

**Reproduce (JMRTD Option A 0.8.6 — single smoke cell):**

```bash
bash scripts/install-jmrtd-local.sh
( cd drivers/jmrtd && mvn -q -DskipTests package )
java -cp drivers/jmrtd/target/jmrtd-tc-ac-01-0.1.0.jar org.emrtd.harness.jmrtd.TcCa01MitigatedRunner \
  -profile profiles/ca-v1-v2-skew.json -log-dir logs
```

**Notes:**
- Full 40-run lab wall-clock and suite `bundle_sha256` are filled when `run_ca_sweep_gmrtd.sh` + JMRTD loop complete on `test-server`; until then cite this profile lock + smoke anchors.
- Option A pin: `org.jmrtd:jmrtd:0.8.6` (see `docs/JMRTD-PIN.md`).
