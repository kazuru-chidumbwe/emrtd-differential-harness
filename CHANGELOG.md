# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Citation / essay pins (`blog-*`, `paper-manifest-*`) remain valid historical anchors.
Prefer **SemVer** (`vX.Y.Z`) for package / paper citations; live pin is **`v1.0.6`** — see [`docs/TAGS.md`](docs/TAGS.md).

## [Unreleased]

### Changed

- (none pending)

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
