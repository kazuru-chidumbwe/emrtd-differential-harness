# Disclosure and data availability (aligned with JISA manuscript)

**Date:** 2026-07-26  
**Branch:** `jmrtd-sweep-2026-07-09`  
**Evidence pin:** tag `paper-manifest-2026-07-25` @ commit `110c546`

This note matches the manuscript Data Availability / §9 rewrite after the Gates review. It supersedes wording that implied the *entire* harness was still embargoed until 7 August 2026.

## What is already public

Harness **source** (profiles, generators, classifier, drivers, CI workflows, and documentation) has been public on GitHub since the July 2026 blog release:

- Tag `blog-b10-2026-07` / `v0.1.0` (`ef15b10`)
- Repository: https://github.com/kazuru-chidumbwe/emrtd-differential-harness

Regenerating suites via `make paper` (or equivalent) reproduces **Observability Scores** and suite structure. Fresh runs embed new UTC timestamps, so they do **not** reproduce an identical historical `bundle_sha256` / suite hash.

## What remains pending public deposit under informal publication timing

- The **locked full-run trees** needed for independent hash verification of the paper pins (primary PACE suite `d8afa161…`, Option A `d505d521…`, plus corroborating CA/AA/PA bundles).
- Those trees currently live on the lab host (`logs/suite-ac-01-sweep-full-20260709T052334Z/`, Option A re-run `…20260720T163501Z/`, and R1 arm dirs). They are **not** in the public git tree (`logs/` is gitignored).
- Zenodo DOI for the archived package (planned at camera-ready or immediately after the gate).

Reviewers may request the locked primary package under journal confidential-review access before the public gate.

## Research pin (do not confuse with “first public release”)

| Object | Value |
| --- | --- |
| Branch | `jmrtd-sweep-2026-07-09` |
| Annotated tag | `paper-manifest-2026-07-25` @ `110c546` |
| Primary PACE cite (manuscript) | `d8afa16137b79877de1e8dc42de5b2136be2d832651b0aee0bd6071b499f6b67` (Jul-9; earlier `fa84b9…` is not the review pin) |
| Option A corroboration | `d505d5212480a68b2253020a6c701c4f50fa71ee8a7da23000f970c573c60051` |
| Success-path control | `31aa96db837a5d2d947da4f8d01d3263dc5a44dbb73a46bf4cddd3195da7e726` |

Tag `paper-manifest-2026-07-25` remains the **evidence pin**. Docs-only commits may land on the branch tip after that tag without retagging.

## Related

- [`TAGS.md`](TAGS.md) — checkout policy  
- [`R1-RESPONSE-EXPAND-EVIDENCE-2026-07-25.md`](R1-RESPONSE-EXPAND-EVIDENCE-2026-07-25.md) — evidence depth  
- [`R5-RESPONSE-REMAINING-GAPS-2026-07-26.md`](R5-RESPONSE-REMAINING-GAPS-2026-07-26.md) — remaining-gaps triage  
- [`SWEEP-RESULTS-2026-07-09.md`](SWEEP-RESULTS-2026-07-09.md) — Jul-9 lab sweep summary  
