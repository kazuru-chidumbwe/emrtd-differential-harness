# EMRTD Differential Harness

Differential test harness for open-source eMRTD (electronic passport) reader libraries. Identical synthetic chip profiles are replayed against multiple stacks; each run captures APDU traces and classifies whether negotiation failures or downgrades reach a typical application caller.

  
**Synthetic only:** no physical passport, no NFC hardware required for smoke tests.

## Related publication

**Jul 2026 — case study (DEV Community):** [Differential Testing Revealed What Conformance Testing Missed: A Case Study with Open-Source eMRTD Libraries](https://dev.to/kazuru_73322ef9a7d6ed2b18/differential-testing-revealed-what-conformance-testing-missed-a-case-study-with-open-source-emrtd-1nie)

The essay’s reproducibility claims match tag **`blog-b10-2026-07`** (`ef15b10`). `main` / paper branches may advance for suite work without changing what the post describes.

**Disclosure (2026-07-26):** harness **source** is already public. What remains pending public deposit under informal publication timing is the locked full-run JSON trees for independent hash checks (not a first-time source release). See [`docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md`](docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md). Paper evidence pin: tag **`paper-manifest-2026-07-25`** @ `110c546`.

```bash
git checkout blog-b10-2026-07   # blog reproduction pin
```

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

### Docker (smoke gate)

One-file reproduction without a hand-tuned host (network required at build). **Option A:** image clones gmrtd and resolves `org.jmrtd:jmrtd:0.8.6` from Maven Central (SHA-256 checked); it does **not** clone E3V3A.

```bash
docker build -t emrtd-harness .
docker run --rm emrtd-harness
# optional offline PA fixtures (incl. TC-PA-04a/04b chained):
docker run --rm emrtd-harness bash scripts/run_offline_pa.sh
```

CA matrix profiles: `python3 profiles/generate_ca_sweep.py` then `bash scripts/run_ca_sweep_gmrtd.sh`. See [`docs/CA-RESULTS-MATRIX-2026-07-25.md`](docs/CA-RESULTS-MATRIX-2026-07-25.md) and [`docs/PA-RESULTS-2026-07-25.md`](docs/PA-RESULTS-2026-07-25.md).

Expected output includes `logs/suite-*/artifact-manifest.json` — the **canonical** published object. Cite `FIG-01`…`FIG-04` (not manuscript figure numbers). Derived `summary-*.md` tables reference the manifest.

---

## Blog reproduction anchor (B10)

The Dev.to essay (*Differential Testing Revealed What Conformance Testing Missed…*) reproduces **TC-AC-01 smoke** (N=1 per library) at:

| Field | Value |
| --- | --- |
| Tag | [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/blog-b10-2026-07) |
| Commit | `ef15b10` |
| Profile | `profiles/pace-then-bac-downgrade.json` |
| Libraries | gmrtd + JMRTD (baseline drivers only) |

```bash
git clone https://github.com/kazuru-chidumbwe/emrtd-differential-harness.git
cd emrtd-differential-harness
git checkout blog-b10-2026-07
bash scripts/bootstrap-vendor.sh
export GOTOOLCHAIN=auto
make smoke
```

Independent lab reproduction (Ubuntu 24.04, July 2026): both drivers green, `observability_score: 0`, ~6 s wall clock. Paper-grade N=100 manifests use `make suite` / `make paper` on `main` (manifest freeze at `e31d945`).

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

See `docs/PROVENANCE.md` and `docs/ARCHITECTURE.md`. Release anchors: [`docs/TAGS.md`](docs/TAGS.md).

---

## Primary wire sweep (AC-01 profile grid, Jul 2026)

**200-run wire-tier sweep** across 50 injection-point profiles, gmrtd + JMRTD, baseline + mitigated:

| Library | Baseline (score 0) | Mitigated (score 2) |
| --- | --- | --- |
| gmrtd | 50 / 50 | 50 / 50 |
| JMRTD | 50 / 50 | 50 / 50 |

Lab-verified 2026-07-09 on `test-server` (64 s). Suite `artifact-manifest.json` SHA-256 (review pin):

`d8afa16137b79877de1e8dc42de5b2136be2d832651b0aee0bd6071b499f6b67`

An earlier note recorded `fa84b9cdef61f2eb4dbb4d01816e949464d1ce12801b42119cb177aa9c0777e0`; that digest does **not** match any retained file under the Jul-9 log directory and is **not** the review pin. Locked raw run-trees remain pending public deposit under informal publication timing (earlier on maintainer reply); see [`docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md`](docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md) and [`docs/SWEEP-RESULTS-2026-07-09.md`](docs/SWEEP-RESULTS-2026-07-09.md). Option A corroboration (0.8.6, 2026-07-20): `d505d5212480a68b2253020a6c701c4f50fa71ee8a7da23000f970c573c60051`.

```bash
python3 classifier/run_suite.py --manifest suites/ac-01-sweep-full.json
```

CI: push to `jmrtd-sweep-*` or manual dispatch of workflow **JMRTD sweep and CA-mitigated verification**.

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
| `JMRTD_VERSION` | `0.8.6` | Maven Central `org.jmrtd:jmrtd` pin (Option A) |
| `PROFILE` | profile-specific | Override chip profile path |
| `LOG_DIR` | `logs` | Trace output directory |

JMRTD is resolved from Maven Central (`drivers/jmrtd/pom.xml`); see [docs/JMRTD-PIN.md](docs/JMRTD-PIN.md). Do not clone archived E3V3A/JMRTD.

---

## License

MIT — see [LICENSE](LICENSE). Upstream libraries (gmrtd, JMRTD) retain their own licenses.

---

## Citation

If you use this harness in research, cite the repository URL and include the `run_id` from your trace JSON.
