# Disclosure and data availability (aligned with manuscript)

**Date:** 2026-07-27 (aligned with Access §9.1 rewrite — no fixed maintainer-facing embargo)  
**Branch:** `jmrtd-sweep-2026-07-09`  
**Evidence pin:** tag `paper-manifest-2026-07-26` (supersedes historical `paper-manifest-2026-07-25`)  
**Venue (current):** IEEE Access retarget in progress (this note stays venue-agnostic for pin/disclosure facts)

This note matches the manuscript Data Availability / §9 rewrite after sponsor clarification: **no fixed deadline was communicated to maintainers**. It supersedes wording that implied a hard dated public-release gate or that the *entire* harness remained embargoed.

## What is already public

Harness **source** (profiles, generators, classifier, drivers, CI workflows, and documentation) has been public on GitHub since the July 2026 blog release:

- Tag `blog-b10-2026-07` / `v0.1.0` (`96fa6a4` after history scrub)
- Repository: https://github.com/kazuru-chidumbwe/emrtd-differential-harness

Regenerating suites via `make paper` (or equivalent) reproduces **Observability Scores** and suite structure. Fresh runs embed new UTC timestamps, so they do **not** reproduce an identical historical suite hash.

## What remains pending public deposit

- The **locked full-run trees** needed for independent hash verification of the paper pins (primary PACE suite, Option A corroboration, plus CA/AA/PA bundles).
- Those trees currently live on the lab host (`logs/suite-ac-01-sweep-full-20260709T052334Z/`, Option A re-run `…20260720T163501Z/`, and R1 arm dirs). They are **not** in the public git tree (`logs/` is gitignored).
- Zenodo DOI for the archived package (planned once authors judge a reasonable interval has elapsed — informally ~2–3 weeks from first effective contact per library — or earlier on any maintainer reply; see manuscript §9.1).

Internal planning may use a target date; that date is **not** a maintainer-facing commitment and must not appear in the paper as a hard embargo trigger.

Reviewers may request the locked primary package under journal confidential-review access before public deposit.

## Maintainer disclosure contacts (as of 2026-07-26)

| Stack | Contact date | Status |
| --- | --- | --- |
| JMRTD | 2026-07-08 (first notice, informal, no deadline stated); follow-up 2026-07-19 | **No email/advisory reply** as of 2026-07-26; notices on record |
| gmrtd | 2026-07-19 (first effective; 8 Jul bounced; informal, no deadline stated) | **No email/advisory reply** as of 2026-07-26; notices on record; PR https://github.com/gmrtd/gmrtd/pull/446 open (not an email reply) |
| pymrtd (ZeroPass) | 2026-07-26 (GitHub Security vulnerability report) | **No advisory reply** as of 2026-07-26; report on record |

**Publication timing:** no fixed date was committed to maintainers. Public deposit of locked full-run trees + Zenodo DOI proceeds under author judgment after a reasonable interval for this severity class (informally on the order of two to three weeks from first effective contact per library), or earlier on any maintainer reply — consistent with informal notice rather than a formal embargo schedule.

## Research pin (do not confuse with “first public release”)

| Object | Value |
| --- | --- |
| Branch | `jmrtd-sweep-2026-07-09` |
| Annotated tag (live) | `paper-manifest-2026-07-26` |
| Historical tag | `paper-manifest-2026-07-25` @ `42b4f1a` (do not hand to reviewers) |
| Primary PACE cite (manuscript) | `d8afa16137b79877de1e8dc42de5b2136be2d832651b0aee0bd6071b499f6b67` |
| Option A corroboration | `d505d5212480a68b2253020a6c701c4f50fa71ee8a7da23000f970c573c60051` |
| Success-path control | `31aa96db837a5d2d947da4f8d01d3263dc5a44dbb73a46bf4cddd3195da7e726` |

**Note:** Raw `artifact-manifest.json` trees are **not** in this git repository (`logs/` is gitignored). Reviewers verify digests against the confidential supplementary package / lab drop until public deposit.

### Retired draft digest (do not cite)

Some July 2026 lab notes and early manuscript drafts recorded suite digest `fa84b9cdef61f2eb4dbb4d01816e949464d1ce12801b42119cb177aa9c0777e0`. That string does **not** match any retained file under the Jul-9 log directory. It is **not** the review pin. Cite only `d8afa161…` (primary) or `d505d521…` (Option A) above.

Tag `paper-manifest-2026-07-26` remains the **evidence pin**. Docs-only commits may land on the branch tip after that tag without retagging.
