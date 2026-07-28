# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Citation / essay pins (`blog-*`, `paper-manifest-*`) remain valid reproducibility anchors.
Prefer **SemVer** (`vX.Y.Z`) for SoftwarX / package citations; see [`docs/TAGS.md`](docs/TAGS.md).

## [Unreleased]

### Changed

- (none pending)

## [tc-ac-adv-2026-07-28] — 2026-07-28

Paper corroboration pin for the shallow TC-AC-ADV arm. **Does not** replace `paper-manifest-2026-07-26` for the primary AC-01 factorial.

### Added

- TC-AC-ADV: synthetic APDU-boundary `pace_channel` modalities (`timeout`, `no_response`, `transport_abort`) + profiles under `profiles/adv/` (feature commit `10137e5`). Shallow Lab corroboration only (*n*=16); does **not** change primary AC-01 digest `d8afa161…`. Manuscript cites suite digest `99d38845…`. Not RF / silicon.

### Changed

- Disclosure docs / README: removed hard dated maintainer-facing public-release gate language. Public locked-run deposit follows informal publication timing aligned with IEEE Access §9.1 (no fixed deadline communicated to maintainers). Suite digests unchanged.
- `docs/TAGS.md`: cite `tc-ac-adv-2026-07-28` for ADV; primary pin remains `paper-manifest-2026-07-26`.

## [0.1.0] — 2026-07-07

### Added

- First public smoke-tier SemVer release (TC-AC-01, N=1 per library).
- Tagged at the same commit as essay pin `blog-b10-2026-07` (`96fa6a4`).
- Keep a Changelog + SemVer discipline documented on the default branch.

[tc-ac-adv-2026-07-28]: https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/tc-ac-adv-2026-07-28
[0.1.0]: https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0
