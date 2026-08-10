# Disclosure and data availability

**Date:** 2026-08-10 (venue-agnostic; no fixed maintainer-facing embargo)  
**Evidence pin (cite):** tag `v1.0.2` — https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.2  
**Historical SemVer:** `v1.0.1` / `v1.0.0` (AC-01 suite digests unchanged)  
**Historical evidence freeze:** `paper-manifest-2026-07-26` (suite digests unchanged)  
**ADV corroboration pin:** tag `tc-ac-adv-2026-07-28` (also present on `v1.0.2` tip)

This note records pin/disclosure facts for independent verification. It is intentionally **venue-agnostic**. An earlier alternate-venue draft for related manuscript text was explored and is now abandoned; it is not under concurrent consideration.

## What is already public

Harness **source** (profiles, generators, classifier, drivers, CI workflows, and documentation) has been public on GitHub since the July 2026 blog release:

- Tag `blog-b10-2026-07` / `v0.1.0` (`96fa6a4` after history scrub)
- Repository: https://github.com/kazuru-chidumbwe/emrtd-differential-harness

Regenerating suites via `make paper` (or equivalent) reproduces **Observability Scores** and suite structure. Fresh runs embed new UTC timestamps, so they do **not** reproduce an identical historical suite hash. Cite locked-run digests (below) or the Zenodo deposit for byte-identical verification.

## Locked-run public deposit

- The **locked full-run trees** needed for independent hash verification of the paper pins (primary PACE suite, Option A corroboration, success-path control, plus CA/AA/PA bundles) are **public** at GitHub Release [locked-runs-2026-07-26](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/locked-runs-2026-07-26) (asset emrtd-locked-runs-2026-07-26.zip). A Zenodo DOI will be minted from the same archive when available; until then cite the Release URL for hash verification.
- Those trees are **not** in the public git tree (`logs/` is gitignored).
- Public deposit proceeds under author judgment after a reasonable interval from first effective contact per library (informally ~2–3 weeks), or earlier on any maintainer reply — consistent with informal notice rather than a formal embargo schedule. As of **2026-08-10** that interval has elapsed for all three libraries.

## Maintainer disclosure contacts (as of 2026-08-10)

| Stack | Contact date | Status |
| --- | --- | --- |
| JMRTD | 2026-07-08 (first notice, informal, no deadline stated); follow-up 2026-07-19 | **No email/advisory reply** as of 2026-08-10; notices on record |
| gmrtd | 2026-07-19 (first effective; 8 Jul bounced; informal, no deadline stated) | **Maintainer email replies** 2026-07-29 and 2026-08-03 (Oscar): confirms record-and-continue; constructive discussion of optional PR https://github.com/gmrtd/gmrtd/pull/446 (open; maintainer may ask to move surfacing into shared `reader`) |
| pymrtd (ZeroPass) | 2026-07-26 (GitHub Security vulnerability report) | **No advisory reply** as of 2026-08-10; report on record |

**Publication timing:** no fixed date was committed to maintainers. Public deposit of locked full-run trees + Zenodo DOI proceeds under author judgment after a reasonable interval for this severity class (informally on the order of two to three weeks from first effective contact per library), or earlier on any maintainer reply.

## Research pin (do not confuse with “first public release”)

| Object | Value |
| --- | --- |
| Annotated tag (live cite) | `v1.0.2` |
| Historical SemVer | `v1.0.1` / `v1.0.0` (AC-01 suite digests unchanged) |
| Historical evidence freeze | `paper-manifest-2026-07-26` (suite digests; superseded for cite by `v1.0.2`) |
| ADV corroboration tag | `tc-ac-adv-2026-07-28` |
| Primary PACE cite | `d8afa16137b79877de1e8dc42de5b2136be2d832651b0aee0bd6071b499f6b67` |
| Option A corroboration | `d505d5212480a68b2253020a6c701c4f50fa71ee8a7da23000f970c573c60051` |
| Success-path control | `31aa96db837a5d2d947da4f8d01d3263dc5a44dbb73a46bf4cddd3195da7e726` |

**Note:** Raw `artifact-manifest.json` trees are **not** in this git repository (`logs/` is gitignored). Verify digests against the GitHub Release zip (or Zenodo DOI once minted).

### Retired draft digest (do not cite)

Some July 2026 lab notes and early manuscript drafts recorded suite digest `fa84b9cdef61f2eb4dbb4d01816e949464d1ce12801b42119cb177aa9c0777e0`. That string does **not** match any retained file under the Jul-9 log directory. It is **not** the review pin. Cite only `d8afa161…` (primary) or `d505d521…` (Option A) above.

Tag `v1.0.2` is the **paper cite pin**. Suite digests remain immutable; docs-only commits after `v1.0.2` require a new SemVer for cite updates.
