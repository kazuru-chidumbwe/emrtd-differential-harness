# Disclosure and data availability (aligned with manuscript)

**Date:** 2026-07-26  
**Branch:** `jmrtd-sweep-2026-07-09`  
**Evidence pin:** tag `paper-manifest-2026-07-26` (supersedes `paper-manifest-2026-07-25` @ `110c546`)  
**Venue (current):** IEEE Access retarget in progress (this note stays venue-agnostic for pin/disclosure facts)

This note matches the manuscript Data Availability / §9 rewrite after review. It supersedes wording that implied the *entire* harness was still embargoed until 7 August 2026.

## What is already public

Harness **source** (profiles, generators, classifier, drivers, CI workflows, and documentation) has been public on GitHub since the July 2026 blog release:

- Tag `blog-b10-2026-07` / `v0.1.0` (`ef15b10`)
- Repository: https://github.com/kazuru-chidumbwe/emrtd-differential-harness

Regenerating suites via `make paper` (or equivalent) reproduces **Observability Scores** and suite structure. Fresh runs embed new UTC timestamps, so they do **not** reproduce an identical historical suite hash.

## What remains pending public deposit under informal publication timing

- The **locked full-run trees** needed for independent hash verification of the paper pins (primary PACE suite, Option A corroboration, plus CA/AA/PA bundles).
- Those trees currently live on the lab host (`logs/suite-ac-01-sweep-full-20260709T052334Z/`, Option A re-run `…20260720T163501Z/`, and R1 arm dirs). They are **not** in the public git tree (`logs/` is gitignored).
- Zenodo DOI for the archived package (planned at camera-ready or immediately after the gate).

Reviewers may request the locked primary package under journal confidential-review access before the public gate.

## Maintainer disclosure contacts (as of 2026-07-26)

| Stack | Status |
| --- | --- |
| gmrtd | First effective notice 2026-07-19; PR https://github.com/gmrtd/gmrtd/pull/446 open |
| JMRTD | First notice 2026-07-08; follow-up 2026-07-19 |
| pymrtd (ZeroPass) | First notice **2026-07-26** via GitHub Security vulnerability report |

## Research pin (do not confuse with “first public release”)

| Object | Value |
| --- | --- |
| Branch | `jmrtd-sweep-2026-07-09` |
| Annotated tag (live) | `paper-manifest-2026-07-26` |
| Historical tag | `paper-manifest-2026-07-25` @ `110c546` (do not hand to reviewers) |
| Primary PACE cite (manuscript) | `d8afa16137b79877de1e8dc42de5b2136be2d832651b0aee0bd6071b499f6b67` |
| Option A corroboration | `d505d5212480a68b2253020a6c701c4f50fa71ee8a7da23000f970c573c60051` |
| CA SM-session bundle | `43f9bd1db0c26709d6a33015fdb82979a0e4f097911123cbac641c1d3eb3c050` |
| Success-path control | `31aa96db837a5d2d947da4f8d01d3263dc5a44dbb73a46bf4cddd3195da7e726` |

**Note:** Raw `artifact-manifest.json` trees are **not** in this git repository (`logs/` is gitignored). Reviewers verify digests against the confidential supplementary package / lab drop until the 7 Aug public deposit.

### Retired draft digest (do not cite)

Some July 2026 lab notes and early manuscript drafts recorded suite digest `fa84b9cdef61f2eb4dbb4d01816e949464d1ce12801b42119cb177aa9c0777e0`. That string does **not** match any retained file under the Jul-9 log directory. It is **not** the review pin. Cite only `d8afa161…` (primary) or `d505d521…` (Option A) above.

Tag `paper-manifest-2026-07-26` is the **live evidence pin** (SM-session CA continue-check + digest correction). Tag `paper-manifest-2026-07-25` is historical only.

## Related

- [`TAGS.md`](TAGS.md) — checkout policy  
- [`R1-RESPONSE-EXPAND-EVIDENCE-2026-07-25.md`](R1-RESPONSE-EXPAND-EVIDENCE-2026-07-25.md) — evidence depth  
- [`R5-RESPONSE-REMAINING-GAPS-2026-07-26.md`](R5-RESPONSE-REMAINING-GAPS-2026-07-26.md) — remaining-gaps triage  
- [`SWEEP-RESULTS-2026-07-09.md`](SWEEP-RESULTS-2026-07-09.md) — Jul-9 lab sweep summary  
