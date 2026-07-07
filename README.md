# EMRTD Differential Harness

Differential test harness for open-source eMRTD (electronic passport) reader libraries. Identical synthetic chip profiles are replayed against multiple stacks; each run captures APDU traces and classifies whether negotiation failures or downgrades reach a typical application caller.

  
**Synthetic only:** no physical passport, no NFC hardware required for smoke tests.

---

## Quick start

Requirements: Go 1.25+ (`GOTOOLCHAIN=auto` works), Java 17, Maven, Python 3.

```bash
git clone https://github.com/kazuru-chidumbwe/emrtd-differential-harness.git
cd emrtd-differential-harness
bash scripts/bootstrap-vendor.sh
export GOTOOLCHAIN=auto
make smoke          # single-run smoke gate
make suite          # AC-01 wire manifest, N=100 + provenance-linked summary
make suite-paper    # full paper matrix manifest
make paper          # CI: tests → smoke → suite → verify → artifacts/
```

Expected output includes `logs/suite-*/artifact-manifest.json` — the **canonical** published object. Cite `FIG-01`…`FIG-04` (not manuscript figure numbers). Derived `summary-*.md` tables reference the manifest.

---

## Middleware (§VIII contribution)

`middleware/` enforces explicit-reject on PACE→BAC downgrade (and EAC-CA failure). Mitigated drivers:

- `cmd/tc-ac-01-mitigated` (gmrtd)
- `cmd/tc-ca-01-mitigated` (gmrtd)
- `TcAc01MitigatedRunner` (JMRTD — no catch-and-continue)

Before/after comparison is a diff of run artifacts, not a rewritten test.

---

## Limitations (simulator, not silicon)

This harness replays **synthetic chip profiles** through in-process APDU transceivers. It does **not** exercise physical NFC hardware, real eMRTD silicon, or live PKD/CRL infrastructure.

- **N=100 repetitions:** Repeating each deterministic profile N=100 demonstrates harness stability and result reproducibility rather than estimating behavioural variance.
- **Blog scope:** `make suite` (TC-AC-01 wire tier).
- **Paper scope:** `make suite-paper` (adds TC-CA-01, offline PA scaffold).

See `docs/PROVENANCE.md` and `docs/ARCHITECTURE.md`.

---

## TC-AC-01 (smoke gate)

Profile: `profiles/pace-then-bac-downgrade.json`

The synthetic chip advertises PACE, returns `6FFF` on the first PACE APDU, and keeps BAC mutual authentication working. Both gmrtd and JMRTD complete BAC; Observability Score **0** means a naive caller checking only session success is not told PACE failed.

---

## Layout

```
profiles/     Chip profiles + catalog.json
suites/       Suite manifests (ac-01-wire, paper-matrix)
simulator/    APDU transceivers
cmd/          Go drivers (baseline + mitigated)
middleware/   §VIII explicit-reject (PACE + CA)
drivers/      JMRTD Java drivers; pymrtd offline
classifier/   manifest.py (canonical), run_suite.py, aggregate.py, verify_manifest.py
scripts/      bootstrap-vendor.sh, quick_test.sh, run_suite.sh, paper.sh, independent-repro.sh
docs/         SCHEMA.md, ARCHITECTURE.md, PROVENANCE.md
schemas/      Frozen v1 schema docs
artifacts/    Verified manifest copies from `make paper` (gitignored except .gitkeep)
logs/         Run output (gitignored)
Makefile      smoke | suite | suite-paper | paper | repro | test
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
