# emrtd-differential-harness

Differential test harness for open-source eMRTD (electronic passport) reader libraries. Identical synthetic chip profiles are replayed against multiple stacks; each run captures APDU traces and classifies whether negotiation failures or downgrades reach a typical application caller.

**Author:** Kazuru  
**Synthetic only:** no physical passport, no NFC hardware required for smoke tests.

---

## Quick start

Requirements: Go 1.25+ (`GOTOOLCHAIN=auto` works), Java 17, Maven, Python 3.

```bash
git clone https://github.com/kazuru-chidumbwe/emrtd-differential-harness.git
cd emrtd-differential-harness
bash scripts/bootstrap-vendor.sh
export GOTOOLCHAIN=auto
bash scripts/quick_test.sh
```

Expected output ends with `SMOKE OK — traces written under logs/`. JSON traces include `run_id`, `observability_score`, and full APDU arrays.

First JMRTD run builds `jmrtd-0.5.2` from vendor sources (`scripts/install-jmrtd-local.sh`) because the Maven Central artifact is empty.

---

## TC-AC-01 (smoke gate)

Profile: `profiles/pace-then-bac-downgrade.json`

The synthetic chip advertises PACE, returns `6FFF` on the first PACE APDU, and keeps BAC mutual authentication working. Both gmrtd and JMRTD complete BAC; Observability Score **0** means a naive caller checking only session success is not told PACE failed.

---

## Layout

```
profiles/     Chip behaviour definitions (JSON)
simulator/    APDU transceivers for synthetic profiles
cmd/          Go drivers (tc-ac-01, tc-ca-01)
drivers/      JMRTD Java driver; pymrtd offline scaffold
classifier/   Trace → Observability Score (0/1/2)
scripts/      bootstrap-vendor.sh, quick_test.sh, install-jmrtd-local.sh
testdata/     Synthetic EF/DG/SOD fixtures
logs/         Run output (gitignored)
```

---

## Observability Score

| Score | Meaning |
| --- | --- |
| 0 | Silent — downgrade/failure not surfaced to naive caller |
| 1 | Logged — visible in trace or session fields only |
| 2 | Surfaced — caller must handle explicit error/result |

---

## Environment variables

| Variable | Default | Purpose |
| --- | --- | --- |
| `GMRTD_PATH` | `../_vendor/gmrtd` | gmrtd source checkout |
| `JMRTD_PATH` | `../_vendor/JMRTD/jmrtd` | JMRTD source tree |
| `PROFILE` | profile-specific | Override chip profile path |
| `LOG_DIR` | `logs` | Trace output directory |

---

## License

MIT — see [LICENSE](LICENSE). Upstream libraries (gmrtd, JMRTD) retain their own licenses.

---

## Citation

If you use this harness in research, cite the repository URL and include the `run_id` from your trace JSON.
