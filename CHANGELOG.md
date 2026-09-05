# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Citation / essay pins (`blog-*`, `paper-manifest-*`) remain valid historical anchors.
Prefer **SemVer** (`vX.Y.Z`) for package / paper citations; live pin is **`v1.0.9`** · Zenodo `10.5281/zenodo.22097289` — see [`docs/TAGS.md`](docs/TAGS.md) and [`docs/ZENODO.md`](docs/ZENODO.md).

## [Unreleased]

## [1.0.11] — 2026-09-05

### Changed

- Tip metadata hygiene: `.zenodo.json` notes and README cite surfaces are venue-neutral (no target journal names on the public tree).
- Added [`scripts/wire_zenodo_doi.py`](scripts/wire_zenodo_doi.py) with venue-neutral notes template for future DOI wires.
- Live **evidence** pin remains **`v1.0.9`** / Zenodo `10.5281/zenodo.22097289`. This tag archives tip metadata only.

## [1.0.10] — 2026-08-25

### Changed

- Cite surfaces record minted Zenodo version DOI `10.5281/zenodo.22097289` (concept `10.5281/zenodo.22095366`) for tag `v1.0.9` (GitHub→Zenodo).
- Added [`docs/ZENODO.md`](docs/ZENODO.md). Live evidence pin remains **`v1.0.9`**.

## [1.0.9] — 2026-08-25

### Changed

- Cite-surface sync: `CITATION.cff`, `.zenodo.json`, `README.md`, and `docs/DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md` now declare live pin **`v1.0.9`**.
- DAS: PR #446 recorded as **merged** (`v1.1.3` / `1701e74`); remeasurement digest `04cdd3dd…` listed; Zenodo described as **pending** (GitHub Releases are the canonical deposits).
- `scripts/preflight_banned_terms.py`: skip self when scanning a tree that includes this script.

### Note

- Defect factorial deposit remains Release **`v1.0.7`**. Remeasurement bytes remain digest `04cdd3dd…` (asset also on `v1.0.8`).

## [1.0.8] — 2026-08-25

### Added

- Remeasurement of locked TC-AC-01 (gmrtd-only, 100 runs) against upstream `v1.1.3` (`64bd6ab`): baseline and mitigated each **50/50 Score 2** (fail-closed). Manifest SHA-256 `04cdd3dd61f983e90140cc2bc16913aa4f0f3fc5dd040cd271fd2f72387f93e3`.

### Changed

- `cmd/tc-ac-01`: `ReaderStatus` implements `Status(reader.Status)` for gmrtd ≥ typed-status API; gate accepts Score-0 silent path **or** Score-2 fail-closed path.
- Unit test expects fail-closed Score 2 under the remeasurement pin.
- Docs: dual gmrtd pins (defect `8fea245` / remeasurement `v1.1.3`).

### Note

- Historical 200-run defect factorial remains under tag `v1.0.7` + pin `8fea245` (digest `e15f4b57…`). Do not conflate the two pins.

## [1.0.7] — 2026-08-10

### Fixed (artifact / baseline review)

- Baseline gmrtd arm drives shipped `reader.ReadDocument` (not hand-composed DoPACE/DoBAC); `PaceSurfacedToCaller` measured from return error.
- Pin gmrtd to `8fea245` (release 0.45.0); bootstrap/paper verify; provenance `gmrtd_commit`.
- Drop vacuous wire N=100 → N=1 (deterministic simulator); shared observability vectors; `SOURCE_DATE_EPOCH`; fence TC-FI-01 / fail-closed pymrtd.

### Changed

- Live cite pin `v1.0.7`. Deposit `emrtd-locked-runs-v1.0.7.zip`. Primary `e15f4b57f0226566d2257cf8912c8521493d678a8e866524150ff7d925e67006`; success-path `b029a9dcefb8277b5790c19fa6b13c19cfca9a9dc6e1786e72a3be4e04fda245`. Schema preflight **206/206** OK (AC-01 factorial + success-path control).

## [1.0.6] — 2026-08-10

### Fixed

- `drivers/pymrtd-offline/run_smoke.py`: emit `variant` + `provenance` so offline PA smoke validates under `schemas/run-artifact-v1.json` (closes the last schema gap in the deposit).
- Success-path control: offline PA smoke runs **before** directory fingerprint.

### Added

- `scripts/package_locked_runs.sh` + `make package-locked-runs` — hard release gate (banned-term + abs-path + full run-artifact schema scan) → SemVer-named `emrtd-locked-runs-vX.Y.Z.zip`.

### Changed

- Live cite pin `v1.0.6`. Deposit asset `emrtd-locked-runs-v1.0.6.zip` (SemVer naming; replaces dated `…-10b`). Primary digest unchanged `c6f03d7e…`; success-path control `a817198c…` (includes schema-valid offline-pa). Package schema preflight **294/294** OK.

## [1.0.5] — 2026-08-10

### Fixed

- Provenance: never emit empty `harness_commit` (Java exit-code check); 1-based `run_index` defaults; AA/CA scripts pass `-run-index`.
- Store `profile_path` relative to harness root (Go/Java/Python) so deposits do not embed host absolute paths.

### Changed

- Fresh lab regen on Lab Test Server (`4510c47` tip): primary AC-01 digest `c6f03d7e…`; success-path control `b98b354d…`.
- Option A is no longer a separate tree — tip uses JMRTD 0.8.6 only; cite the primary digest for the factorial.
- Locked-run deposit: Release `locked-runs-2026-08-10b` (zip `emrtd-locked-runs-2026-08-10b.zip`). Primary factorial schema-OK; offline-pa smoke under success-path still lacked provenance until `v1.0.6`.
- Live paper cite pin was `v1.0.5` (superseded by `v1.0.6`).

## [1.0.4] — 2026-08-10

### Added

- `CITATION.cff` + `.zenodo.json` — author/ORCID/keywords for GitHub “Cite this repository” and Zenodo GitHub-integration metadata.
- `SECURITY.md` — harness vulnerability reporting channel; points to disclosure log for third-party library notices.
- Frozen JSON Schema files under `schemas/` with `verify_manifest.py` enforcement.
- CI on `main` / PRs; `scripts/preflight_banned_terms.py`.

### Changed

- Live paper cite pin was `v1.0.4` (path-sanitized digests from strip-and-rezip; superseded by `v1.0.5` lab regen).
