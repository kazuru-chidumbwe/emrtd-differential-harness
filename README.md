# EMRTD Differential Harness

Differential test harness for open-source eMRTD (electronic passport) reader libraries. Identical synthetic chip profiles are replayed against multiple stacks; each run captures APDU traces and classifies whether negotiation failures or downgrades reach a typical application caller.

  
**Synthetic only:** no physical passport, no NFC hardware required for smoke tests.

## Related publication

**Jul 2026 — case study (DEV Community):** [Differential Testing Revealed What Conformance Testing Missed: A Case Study with Open-Source eMRTD Libraries](https://dev.to/kazuru_73322ef9a7d6ed2b18/differential-testing-revealed-what-conformance-testing-missed-a-case-study-with-open-source-emrtd-1nie)

The essay’s reproducibility claims match tags **`v0.1.0`** / **`blog-b10-2026-07`** (same tree). `main` / paper branches may advance for suite work without changing what the post describes.

**Disclosure (2026-08-10):** harness **source** is already public. Locked full-run JSON trees for independent hash verification are on GitHub Release (Zenodo DOI pending) — see [`docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md`](docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md). **Paper cite pin:** **`v1.0.3`** — https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.3 · See [`CHANGELOG.md`](CHANGELOG.md) and [`docs/TAGS.md`](docs/TAGS.md).

```bash
git checkout v1.0.3
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
| SemVer | [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) |
| Essay tag | [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/blog-b10-2026-07) |
| Profile | `profiles/pace-then-bac-downgrade.json` |
| Libraries | gmrtd + JMRTD (baseline drivers only) |

```bash
git clone https://github.com/kazuru-chidumbwe/emrtd-differential-harness.git
cd emrtd-differential-harness
git checkout v0.1.0
bash scripts/bootstrap-vendor.sh
export GOTOOLCHAIN=auto
make smoke
```

Independent lab reproduction (Ubuntu 24.04, July 2026): both drivers green, `observability_score: 0`, ~6 s wall clock. Paper-grade work cites **`v1.0.3`** (historical freeze tags `paper-manifest-*` / earlier SemVer remain aliases).

---

## What this harness measures (and what it does not)

**Object under test:** open-source **reader-library** negotiation and error-surfacing logic (call-boundary Observability Score), not the chip RF layer and not physical NFC silicon.

**Real-world grounding (not RF):**

- ICAO Doc 9303 Part 11 permits PACE→BAC fallback when CardAccess/PACE parameters are unavailable (inspection systems SHOULD try BAC).
- Upstream gmrtd ships `--skipPace` as an explicit integrator control over whether PACE runs.
- Public JMRTD integrator survey documents catch-and-continue around `doPACE` ([survey](docs/JMRTD-PUBLIC-CATCH-CONTINUE-SURVEY-2026-07-26.md)).
- Coordinated disclosure contacts for gmrtd, JMRTD, and pymrtd are on record ([disclosure note](docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md)).

**Synthetic only:** no physical passport, no NFC hardware, no live PKD. ADV profiles are lab APDU-boundary channel faults — **Not RF / silicon**.

---

## Paper arm set (on `main`)

| Arm | Role | Where |
| --- | --- | --- |
| **TC-AC-01** 50-profile × 2 libs × 2 variants (**200 runs**) | Primary PACE→BAC observability factorial | `suites/ac-01-sweep-full.json` · pin `6c6cbbd1…` |
| **Option A** JMRTD **0.8.6** re-run | Same score table corroboration | pin `2efd190e…` |
| **TC-AC-ADV** (4 channel-fault profiles, *n*=16) | Adversarial corroborating; Not RF | `profiles/adv/` · tag `tc-ac-adv-2026-07-28` · digest `99d38845…` |
| **CA / AA / PA** sweeps | Mechanism corroboration grids | `docs/CA-RESULTS-*`, `AA-RESULTS-*`, `PA-RESULTS-*` |
| **TA-EAC** + success-path FP controls | Unsupported-path / non-false-reject checks | `docs/TA-EAC-*`, `docs/SUCCESS-PATH-*` · control `31aa96db…` |
| **Middleware** explicit-reject | Raises Score 0→2 without upstream patches | `middleware/` · mitigated runners |

**N=100** (AC-01 wire suite): repeating each deterministic profile demonstrates **harness stability and reproducibility**, not behavioural variance / statistical power.

Evidence pin: **`v1.0.3`**. Locked run-trees: GitHub Release (see disclosure note).

---

## Middleware (explicit-reject)

`middleware/` enforces explicit-reject on PACE→BAC downgrade (and EAC-CA failure). Mitigated drivers:

- `cmd/tc-ac-01-mitigated` (gmrtd)
- `cmd/tc-ca-01-mitigated` (gmrtd)
- `TcAc01MitigatedRunner` (JMRTD — no catch-and-continue)

Before/after comparison is a diff of run artifacts, not a rewritten test.

---

## Limitations (simulator, not silicon)

This harness replays **synthetic chip profiles** through in-process APDU transceivers. It does **not** exercise physical NFC hardware, real eMRTD silicon, or live PKD/CRL infrastructure.

- **N=100 repetitions:** Repeating each deterministic profile N=100 demonstrates harness stability and result reproducibility rather than estimating behavioural variance.
- **Blog scope:** `make smoke` / TC-AC-01 single-profile.
- **Paper scope:** full arm set above (`make suite-paper` / tagged evidence + Zenodo locked trees).

See `docs/PROVENANCE.md` and `docs/ARCHITECTURE.md`. Release anchors: [`docs/TAGS.md`](docs/TAGS.md).

---

## Primary wire sweep (AC-01 profile grid, Jul 2026)

**200-run wire-tier sweep** across 50 injection-point profiles, gmrtd + JMRTD, baseline + mitigated:

| Library | Baseline (score 0) | Mitigated (score 2) |
| --- | --- | --- |
| gmrtd | 50 / 50 | 50 / 50 |
| JMRTD | 50 / 50 | 50 / 50 |

**Equal scores ≠ identical library defects.** Both baselines score 0 under the stated naive-caller models, but silence is reached differently:

| Stack | How baseline silence arises | Epistemic status |
| --- | --- | --- |
| **gmrtd** | Reference client / driver path records `PaceErr` at the session layer but does not branch on it before BAC — close to an unchecked-error omission in ordinary use | Demonstrated property of the shipped usage pattern |
| **JMRTD** | `PassportService.doPACE` **throws**; the baseline harness `catch`es and continues to BAC (the API permits unconstrained catch). Mitigated runners do **not** catch-and-continue | API-permitted integrator pattern; attested in public third-party clients ([survey](docs/JMRTD-PUBLIC-CATCH-CONTINUE-SURVEY-2026-07-26.md)), not a shipped JMRTD demo defect |

Do not read the 50/50 table as “both libraries behave identically.” Manuscript §5.3.1 spells out the same split.

Lab-verified 2026-07-09 on `test-server` (64 s). Score table above is stable under regenerate; the **canonical suite pin** (SHA-256 of the locked Jul-9 `artifact-manifest.json`) and Option A corroboration live only in:

[`docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md`](docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md)

Do not cite older draft digests from blog-era notes. Locked raw run-trees are deposited on Zenodo for independent hash verification (see disclosure note). Details of the Jul-9 lab pass: [`docs/SWEEP-RESULTS-2026-07-09.md`](docs/SWEEP-RESULTS-2026-07-09.md).

```bash
python3 classifier/run_suite.py --manifest suites/ac-01-sweep-full.json
```

CI: push to `main` or manual dispatch of workflow **JMRTD sweep and CA-mitigated verification**.

---

## TC-AC-ADV (channel-abort corroboration)

Shallow adversarial arm (*n*=16): four `profiles/adv/` modalities (`timeout` at MSE:Set AT / GENERAL AUTHENTICATE, `no_response`, `transport_abort`) × gmrtd+JMRTD × baseline+mitigated. Baseline Score **0** / mitigated Score **2** on all cells. Does **not** change primary AC-01 digest `6c6cbbd1…`. Tag: `tc-ac-adv-2026-07-28`. **Not RF / silicon.**

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
