# Option A headline N=40 — 20 Jul 2026 (test-server)

**Status:** **SUITE OK** · paper-grade Option A evidence for live Central `org.jmrtd:jmrtd:0.8.6`  
**Supersedes for Option A claims:** 0.5.2-era manifest `fa84b9cd…` (full factorial re-run separate / in progress)

## Pins

| Pin | Value |
| --- | --- |
| JMRTD | `org.jmrtd:jmrtd:0.8.6` · jar SHA-256 `5C303D7BA0DB892411E739A9920B3E0FB3C62416344CD7F220F359BDD91C0C5B` |
| gmrtd vendor | `8fea245048d3b4e76483d048b202ff7f5269728c` |
| Suite | `suites/ac-01-headline-40.json` · seed 1 · MRZ=`orig` |
| Host | test-server · Go 1.25.0 · OpenJDK 17.0.19 · Python 3.12.3 |

## Canonical object

| Field | Value |
| --- | --- |
| Log dir | `logs/suite-ac-01-headline-40-20260720T150542Z-test-server/` |
| Published copy | `artifacts/ac-01-headline-40-option-a-0.8.6-20260720T150542Z-artifact-manifest.json` |
| Manifest SHA-256 | `3A8D18ADE8A66DDFFAF3344606AD5382F80F544A6831A08A47DBE1202775176B` |
| Generated at | `2026-07-20T15:05:55Z` |

## Counts (40/40)

| Library | baseline (n=10) | mitigated (n=10) |
| --- | ---: | ---: |
| gmrtd | S_INT **0** | S_INT **2** |
| jmrtd | S_INT **0** | S_INT **2** |

- Cells: 5 SW × 2 injection × 2 libs × 2 variants (MRZ fixed `orig`)
- Within-JMRTD: single `pace_exception_class` = `org.jmrtd.CardServiceProtocolException` all SW×injection
- **No cross-library S_INT diverge** under Option A headline cells

## Build note (JMRTD shade)

Maven shade of signed `bcprov` left `META-INF/*.{SF,RSA,…}` → all 20 JMRTD cells failed with `Invalid signature file digest`. Fixed in `drivers/jmrtd/pom.xml` (exclude signature entries). Rebuild required before Option A suite.

## Scope

This is the **headline** Option A re-verify (sponsor soft deadline ~29 Jul). Full factorial N=200 (`ac-01-sweep-full`) is the robustness / Evaluation §6.1 scale claim — re-run separately on the same 0.8.6 pin; do not cite `fa84b9cd…` as live Option A evidence.
