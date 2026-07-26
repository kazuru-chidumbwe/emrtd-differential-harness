# AC-01 profile sweep — lab verification (2026-07-09)

**Suite:** `ac-01-sweep-full` (50 profiles × 2 libraries × 2 variants = **200 runs**)

**Host:** lab host (`test-server`), Ubuntu, Go 1.25, Java 17, Maven 3.8.7

**Wall clock:** 64.2 s (suite only)

**Manifest SHA-256 (review pin):** `d8afa16137b79877de1e8dc42de5b2136be2d832651b0aee0bd6071b499f6b67`

Earlier drafts recorded `fa84b9cdef61f2eb4dbb4d01816e949464d1ce12801b42119cb177aa9c0777e0`. That digest does not match retained Jul-9 bytes and is **not** the review pin (see [`DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md`](DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md)).

**Log directory (on lab):** `logs/suite-ac-01-sweep-full-20260709T052334Z/`

## Observability score counts

| Library | Variant | Score | Runs |
| --- | --- | --- | --- |
| gmrtd | baseline | 0 (silent) | 50 / 50 |
| gmrtd | mitigated | 2 (surfaced) | 50 / 50 |
| JMRTD | baseline | 0 (silent) | 50 / 50 |
| JMRTD | mitigated | 2 (surfaced) | 50 / 50 |

**Total:** 200 / 200 runs completed; all baseline cells 100% silent; all mitigated cells 100% surfaced.

## Reproduce

```bash
bash scripts/bootstrap-vendor.sh
bash scripts/install-jmrtd-local.sh
export GOTOOLCHAIN=auto
python3 classifier/run_suite.py --manifest suites/ac-01-sweep-full.json
```

GitHub Actions workflow `.github/workflows/jmrtd-sweep.yml` runs the same suite when billing allows.

## Notes

- PA fixture (`TC-PA-01`) not run on this lab pass (PEP 668 pip restriction).
- Full per-run JSON artifacts remain on the lab host; cite manifest SHA-256 above for paper tables.
