# TC-PA offline combinatorial grid — R1 (≥24 cells)

**Tier:** offline (pymrtd) — never pooled with wire \(N\).

**Axes (deterministic; \(n=1\) per fixture):**

| Axis | Levels |
| --- | --- |
| Digest / LDS hash | SHA-1, MD5, SHA-256 |
| Signer validity vs inspection date (2026-07-07) | fresh; expired (−1y); expired (−5y); not-yet-valid |
| Chain shape | self-signed; CSCA→DSC |

**Grid size:** \(3 \times 4 \times 2 = 24\) fixtures. Claim is **combinatorial coverage of policy conditions**, not statistical variance. Live PKD/CRL out of scope.

**Generator:** `profiles/generate_pa_sweep.py` → `testdata/sod/pa-sweep/` + `testdata/sod/fixtures/pa-sweep-*`

**Control cells:** `sha256` + `fresh` (both chain shapes) set `expect_policy_rejection=false`.

**Classifier:** when `expect_policy_rejection=true` and `verify()` succeeds with no error → Observability Score **0** (silent policy gap). When `verify()` raises → Score **2** (surfaced reject — library already enforces that cell).

**Smoke anchors (unchanged):** TC-PA-01 / 03 / 04a / 04b under `testdata/sod/tc-pa-0*.json`.

## Lab pin (test-server)

| Field | Value |
| --- | --- |
| Host | `test-server` (Ubuntu; Go 1.25.0; OpenJDK 17; Python 3.12.3; pymrtd vendor) |
| Log dir | `logs/pa-sweep-full-r1-test-server` |
| `bundle_sha256` | `c231482386ab15422b17331f7021432dce01cbc4b52c4c7e644c8e2f2488ecc2` |
| Runner | `scripts/run_pa_sweep_full.sh` |

**Outcome counts (this pin):**

| Class | Count | Meaning |
| --- | ---: | --- |
| Silent policy cells (score 0) | 14 | SHA-1 and expiry / not-yet under naive `verify()` |
| Surfaced reject (score 2) | 8 | All MD5 cells — pymrtd raises on MD5 |
| Control (score 0 or 1) | 2 | SHA-256 + fresh |

## Reproduce

```bash
python3 profiles/generate_pa_sweep.py
LOG_DIR=logs/pa-sweep-full-r1-test-server bash scripts/run_pa_sweep_full.sh
```

Catalog pointer: `profiles/catalog.json` → `offline_cases` (smoke anchors); sweep index at `testdata/sod/pa-sweep/index.json`.
