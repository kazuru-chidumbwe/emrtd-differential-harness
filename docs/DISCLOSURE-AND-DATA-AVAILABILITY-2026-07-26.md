# Disclosure and data availability

**Date:** 2026-08-25 (venue-agnostic; no fixed maintainer-facing embargo)  
**Zenodo version DOI:** [10.5281/zenodo.22097289](https://doi.org/10.5281/zenodo.22097289) (concept [10.5281/zenodo.22095366](https://doi.org/10.5281/zenodo.22095366))  
**Evidence pin (live cite):** tag `v1.0.9` — https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.9  
**Defect factorial deposit:** Release `v1.0.7` / asset `emrtd-locked-runs-v1.0.7.zip` (digest `e15f4b57…`)  
**Remeasurement deposit:** Release `v1.0.9` / asset `emrtd-gmrtd-v1.1.3-remeasurement-v1.0.8.zip` (digest `04cdd3dd…`; same bytes as on Release `v1.0.8`)  
**Historical SemVer:** `v1.0.8` / `v1.0.7` / `v1.0.6` / …  
**ADV corroboration pin:** tag `tc-ac-adv-2026-07-28` (also present on tip)

This note records pin/disclosure facts for independent verification. It is intentionally **venue-agnostic**.

## What is already public

Harness **source** (profiles, generators, classifier, drivers, CI workflows, and documentation) has been public on GitHub since the July 2026 blog release:

- Tag `blog-b10-2026-07` / `v0.1.0` (`96fa6a4` after history scrub)
- Repository: https://github.com/kazuru-chidumbwe/emrtd-differential-harness

Regenerating suites via `make paper` (or equivalent) reproduces **Observability Scores** and suite structure. Fresh runs embed new UTC timestamps, so they do **not** reproduce an identical historical suite hash. Cite locked-run digests (below) or the GitHub Release assets for byte-identical verification. Zenodo version DOI **[10.5281/zenodo.22097289](https://doi.org/10.5281/zenodo.22097289)** (concept [10.5281/zenodo.22095366](https://doi.org/10.5281/zenodo.22095366); GitHub→Zenodo on Release `v1.0.9`). Locked-run JSON trees remain on the GitHub Release assets listed below.

## Locked-run public deposit

- **Defect factorial (200-run primary + FP control):** GitHub Release [v1.0.7](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.7) (asset `emrtd-locked-runs-v1.0.7.zip`). Zip SHA-256 `25140e4e29f3568340f997b23e9077af701126a920d26bf2d6072cb5626144a3`. Primary manifest `e15f4b57f0226566d2257cf8912c8521493d678a8e866524150ff7d925e67006`; success-path `b029a9dcefb8277b5790c19fa6b13c19cfca9a9dc6e1786e72a3be4e04fda245`.
- **gmrtd `v1.1.3` remeasurement (100-run gmrtd-only):** GitHub Release [v1.0.9](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.9) (asset `emrtd-gmrtd-v1.1.3-remeasurement-v1.0.8.zip`). Manifest SHA-256 `04cdd3dd61f983e90140cc2bc16913aa4f0f3fc5dd040cd271fd2f72387f93e3`. Evaluation pin `64bd6ab` / tag `v1.1.3` (merge `1701e74`). See [`GMRTD-PIN.md`](GMRTD-PIN.md).
- **Deposit gate:** Release zips are built only via `make package-locked-runs` (banned-term + abs-path + **full** run-artifact schema scan) when packaging full locked-run trees. Do not mint Zenodo from ``v1.0.6`` or older dated `…-10` / `…-10b` zips.
- Those trees are **not** in the public git tree (`logs/` is gitignored).
- Public deposit proceeds under author judgment after a reasonable interval from first effective contact per library (informally ~2–3 weeks), or earlier on any maintainer reply. As of **2026-08-25** that interval has elapsed for all three libraries.

## Maintainer disclosure contacts (as of 2026-08-25)

| Stack | Contact date | Status |
| --- | --- | --- |
| JMRTD | 2026-07-08 (first notice, informal, no deadline stated); follow-up 2026-07-19 | **No email/advisory reply** as of 2026-08-17; notices on record |
| gmrtd | 2026-07-19 (first effective; 8 Jul bounced; informal, no deadline stated) | **Maintainer email replies** 2026-07-29, 2026-08-03, 2026-08-17 (Oscar). Author-opened PR https://github.com/gmrtd/gmrtd/pull/446 **merged 2026-08-24** as release `v1.1.3` (merge commit `1701e74a746a260a5e1707f0c5ef34e100feb32b`; tag tip `64bd6ab`). Default fail-closed after a recorded PACE error unless `AllowBacFallbackOnPaceError` is set. **Defect factorial evaluates upstream pin `8fea245` (release 0.45.0)** — pre–merge — see [`GMRTD-PIN.md`](GMRTD-PIN.md). Remeasurement on `v1.1.3` scores 2 on 50/50 baseline cells (digest `04cdd3dd…`). |
| pymrtd (ZeroPass) | 2026-07-26 (GitHub Security vulnerability report) | **No advisory reply** as of 2026-08-17; report on record |

**Publication timing:** no fixed date was committed to maintainers. Public deposit of locked full-run trees proceeds under author judgment after a reasonable interval for this severity class (informally on the order of two to three weeks from first effective contact per library), or earlier on any maintainer reply. Software archive DOI: [10.5281/zenodo.22097289](https://doi.org/10.5281/zenodo.22097289). Locked-run trees remain on GitHub Releases (also linked from the Zenodo record).

## Research pin (do not confuse with “first public release”)

| Object | Value |
| --- | --- |
| Annotated tag (live cite) | `v1.0.9` |
| Zenodo version DOI | `10.5281/zenodo.22097289` |
| Zenodo concept DOI | `10.5281/zenodo.22095366` |
| Defect locked-run asset | `emrtd-locked-runs-v1.0.7.zip` on Release `v1.0.7` |
| Remeasurement asset | `emrtd-gmrtd-v1.1.3-remeasurement-v1.0.8.zip` on Release `v1.0.9` |
| gmrtd defect pin | `8fea245048d3b4e76483d048b202ff7f5269728c` (release 0.45.0; [`GMRTD-PIN.md`](GMRTD-PIN.md)) |
| gmrtd remeasurement pin | `64bd6ab8fbf8802c718a6da0dcc6f6312a3404ca` / tag `v1.1.3` |
| ADV corroboration tag | `tc-ac-adv-2026-07-28` |
| Primary PACE cite (defect) | `e15f4b57f0226566d2257cf8912c8521493d678a8e866524150ff7d925e67006` |
| Success-path control | `b029a9dcefb8277b5790c19fa6b13c19cfca9a9dc6e1786e72a3be4e04fda245` |
| Remeasurement (gmrtd `v1.1.3`) | `04cdd3dd61f983e90140cc2bc16913aa4f0f3fc5dd040cd271fd2f72387f93e3` |

**Note:** Raw `artifact-manifest.json` trees are **not** in this git repository (`logs/` is gitignored). Verify digests against the GitHub Release zips. Software archive: [https://doi.org/10.5281/zenodo.22097289].

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

Retired primary `c6f03d7e…` / control `a817198c…` (v1.0.6). Cite only `e15f4b57…` (primary), `b029a9dc…` (success-path control), or `04cdd3dd…` (remeasurement) above.

Tag `v1.0.9` is the **live paper cite pin**. New evidence semantics or deposit changes require a new SemVer.
