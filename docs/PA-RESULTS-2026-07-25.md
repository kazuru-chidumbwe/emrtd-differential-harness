# TC-PA offline matrix — chained trust path (2026-07-25)

**Tier:** offline (pymrtd) — never pooled with wire N.

| Case | Condition | Chain | `expect_policy_rejection` | Expected score (naive `verify()`) |
| --- | --- | --- | :---: | ---: |
| TC-PA-01 | Weak digest (SHA-1 DG hash) | self-signed | true | 0 |
| TC-PA-03 | Expired DSC | self-signed | true | 0 |
| TC-PA-04a | Fresh DSC | CSCA→DSC | false | 0 (CMS OK; no trust-store policy) |
| TC-PA-04b | Expired DSC | CSCA→DSC | true | 0 |

**Fixtures:** `testdata/sod/` (+ `fixtures/`). Generator for 04: `python profiles/generate_pa04_chained_fixture.py`.

**Classifier:** cases with `expect_policy_rejection: true` that return success without error score Observability Score 0 (silent policy gap). Documented in manuscript §5.3.2.

**Reproduce:**

```bash
bash scripts/run_offline_pa.sh
```

Catalog pointer: `profiles/catalog.json` → `offline_cases`.
