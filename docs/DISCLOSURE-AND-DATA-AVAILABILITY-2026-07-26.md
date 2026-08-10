# Disclosure and data availability

**Date:** 2026-08-10 (venue-agnostic; no fixed maintainer-facing embargo)  
**Evidence pin (cite):** tag `v1.0.7` — https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.7  
**Locked-run deposit:** Release `v1.0.7` / asset `emrtd-locked-runs-v1.0.7.zip`  
**Historical SemVer:** `v1.0.5` / `v1.0.4` / `v1.0.3` / `v1.0.2` / `v1.0.1` / `v1.0.0`  
**ADV corroboration pin:** tag `tc-ac-adv-2026-07-28` (also present on tip)

This note records pin/disclosure facts for independent verification. It is intentionally **venue-agnostic**. An earlier alternate-venue draft for related manuscript text was explored and is now abandoned; it is not under concurrent consideration.

## What is already public

Harness **source** (profiles, generators, classifier, drivers, CI workflows, and documentation) has been public on GitHub since the July 2026 blog release:

- Tag `blog-b10-2026-07` / `v0.1.0` (`96fa6a4` after history scrub)
- Repository: https://github.com/kazuru-chidumbwe/emrtd-differential-harness

Regenerating suites via `make paper` (or equivalent) reproduces **Observability Scores** and suite structure. Fresh runs embed new UTC timestamps, so they do **not** reproduce an identical historical suite hash. Cite locked-run digests (below) or the Zenodo deposit for byte-identical verification.

## Locked-run public deposit

- The **locked full-run trees** needed for independent hash verification of the paper pins (primary PACE suite and success-path control (AC-01 factorial + FP control)) are **public** at GitHub Release [v1.0.7](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.7) (asset `emrtd-locked-runs-v1.0.7.zip`). A Zenodo DOI will be minted from **this** archive when available; until then cite the Release URL for hash verification.
- **Deposit gate:** Release zips are built only via `make package-locked-runs` (banned-term + abs-path + **full** run-artifact schema scan). As of `v1.0.7` that scan is **206/206** OK (AC-01 factorial + success-path control). Do not mint Zenodo from ``v1.0.6`` or older dated `…-10` / `…-10b` zips.
- Those trees are **not** in the public git tree (`logs/` is gitignored).
- Public deposit proceeds under author judgment after a reasonable interval from first effective contact per library (informally ~2–3 weeks), or earlier on any maintainer reply — consistent with informal notice rather than a formal embargo schedule. As of **2026-08-10** that interval has elapsed for all three libraries.

## Maintainer disclosure contacts (as of 2026-08-10)

| Stack | Contact date | Status |
| --- | --- | --- |
| JMRTD | 2026-07-08 (first notice, informal, no deadline stated); follow-up 2026-07-19 | **No email/advisory reply** as of 2026-08-10; notices on record |
| gmrtd | 2026-07-19 (first effective; 8 Jul bounced; informal, no deadline stated) | **Maintainer email replies** 2026-07-29 and 2026-08-03 (Oscar): confirms record-and-continue; constructive discussion of optional PR https://github.com/gmrtd/gmrtd/pull/446 (open; maintainer may ask to move surfacing into shared `reader`). **Paper evaluates upstream pin `8fea245` (release 0.45.0)** — pre–PR-446 CLI surfacing — see [`GMRTD-PIN.md`](GMRTD-PIN.md). |
| pymrtd (ZeroPass) | 2026-07-26 (GitHub Security vulnerability report) | **No advisory reply** as of 2026-08-10; report on record |

**Publication timing:** no fixed date was committed to maintainers. Public deposit of locked full-run trees + Zenodo DOI proceeds under author judgment after a reasonable interval for this severity class (informally on the order of two to three weeks from first effective contact per library), or earlier on any maintainer reply.

## Research pin (do not confuse with “first public release”)

| Object | Value |
| --- | --- |
| Annotated tag (live cite) | `v1.0.7` |
| Locked-run asset | `emrtd-locked-runs-v1.0.7.zip` on Release `v1.0.7` |
| gmrtd library pin | `8fea245048d3b4e76483d048b202ff7f5269728c` (release 0.45.0; [`GMRTD-PIN.md`](GMRTD-PIN.md); verified in `make paper`) |
| ADV corroboration tag | `tc-ac-adv-2026-07-28` |
| Primary PACE cite | `e15f4b57f0226566d2257cf8912c8521493d678a8e866524150ff7d925e67006` |
| Option A (JMRTD 0.8.6) | Same as primary — tip regen uses Central `0.8.6` only |
| Success-path control | `b029a9dcefb8277b5790c19fa6b13c19cfca9a9dc6e1786e72a3be4e04fda245` |

**Note:** Raw `artifact-manifest.json` trees are **not** in this git repository (`logs/` is gitignored). Verify digests against the GitHub Release zip (or Zenodo DOI once minted).

### Retired digests (do not cite)

| Digest | Why retired |
| --- | --- |
| `b98b354de6276bf08b3a75dac4ef37593a3ef728be65cb8218b21e8e4e787273` | Success-path control before offline-pa included in fingerprint (`v1.0.5`) |
| `6c6cbbd1959ec047980c1021c5d3ddf196c10b088d7884ab8ed80f56ae81c9b5` | Path-sanitized strip reissue (superseded by lab regen `e15f4b57…`) |
| `2efd190e80d90fd3c3793b3b481805a567f3630329cf434f8a4cd33fbe9bddd4` | Path-sanitized Option A strip (superseded; tip uses single primary digest) |
| `31aa96db837a5d2d947da4f8d01d3263dc5a44dbb73a46bf4cddd3195da7e726` | Prior success-path control |
| `d8afa16137b79877de1e8dc42de5b2136be2d832651b0aee0bd6071b499f6b67` | Pre–path-sanitization primary manifest |
| `d505d5212480a68b2253020a6c701c4f50fa71ee8a7da23000f970c573c60051` | Pre–path-sanitization Option A manifest |
| `fa84b9cdef61f2eb4dbb4d01816e949464d1ce12801b42119cb177aa9c0777e0` | Draft/lab note only; never matched retained Jul-9 tree |

Retired primary `c6f03d7e…` / control `a817198c…` (v1.0.6). Cite only `e15f4b57…` (primary) or `b029a9dc…` (success-path control) above.

Tag `v1.0.7` is the **paper cite pin**. New evidence semantics or deposit changes require a new SemVer.
