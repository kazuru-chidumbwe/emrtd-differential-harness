# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Citation / essay pins (`blog-*`, `paper-manifest-*`) remain valid historical anchors.
Prefer **SemVer** (`vX.Y.Z`) for package / paper citations; live pin is **`v1.0.1`** — see [`docs/TAGS.md`](docs/TAGS.md).

## [Unreleased]

### Changed

- (none pending)

## [1.0.1] — 2026-08-10

### Changed

- Disclosure: alternate-venue draft wording (now abandoned); no concurrent-venue implication.
- Live paper cite pin is `v1.0.1` (suite digests unchanged from `v1.0.0`).

## [1.0.0] — 2026-08-10

### Added

- SemVer paper cite pin `v1.0.0` (bpfix-style permanent tree URL).
- Public locked-run deposit (Release asset; digests `d8afa161…` / `d505d521…` / `31aa96db…`).

### Changed

- Default branch is the research tip (`main`); Java package `org.emrtd.harness.jmrtd`.
- Disclosure / README / TAGS: venue-agnostic wording (prior alternate-venue draft abandoned).
- README: full paper arm set (TC-AC-ADV, CA/AA/PA, TA-EAC, grounding).
- Paper citations prefer SemVer over dated `paper-manifest-*` tag names.

## [tc-ac-adv-2026-07-28] — 2026-07-28

Paper corroboration pin for the shallow TC-AC-ADV arm. **Does not** replace `paper-manifest-2026-07-26` for the primary AC-01 factorial.

### Added

- TC-AC-ADV: synthetic APDU-boundary `pace_channel` modalities (`timeout`, `no_response`, `transport_abort`) + profiles under `profiles/adv/` (feature commit `10137e5`). Shallow Lab corroboration only (*n*=16); does **not** change primary AC-01 digest `d8afa161…`. Manuscript cites suite digest `99d38845…`. Not RF / silicon.

### Changed

- Disclosure docs / README: removed hard dated maintainer-facing public-release gate language. Public locked-run deposit follows informal publication timing (no fixed deadline communicated to maintainers). Suite digests unchanged.
- `docs/TAGS.md`: cite `tc-ac-adv-2026-07-28` for ADV; primary pin remains `paper-manifest-2026-07-26`.

## [0.1.0] — 2026-07-07

### Added

- First public smoke-tier SemVer release (TC-AC-01, N=1 per library).
- Tagged at the same commit as essay pin `blog-b10-2026-07` (`96fa6a4`).
- Keep a Changelog + SemVer discipline documented on the default branch.

[1.0.1]: https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.1
[1.0.0]: https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.0
[tc-ac-adv-2026-07-28]: https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/tc-ac-adv-2026-07-28
[0.1.0]: https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0
