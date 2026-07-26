# JMRTD measurement pin — Option A (locked 19 Jul 2026)

## Decision

**Option A is required, not optional.** Measure the **live** canonical JMRTD host API from Maven Central — same maintainer line (`martijno` / Martijn Oostdijk), same SCM (SourceForge SVN / jmrtd.org). Do **not** measure archived E3V3A/GitHub.

| Field | Value |
| --- | --- |
| Coordinate | `org.jmrtd:jmrtd:**0.8.6**` |
| Verified | 19 Jul 2026 — **literal** `curl` `repo1.maven.org/.../maven-metadata.xml` → `<latest>0.8.6</latest>` / `<release>0.8.6</release>` / `lastUpdated` `20260421074743` |
| Jar SHA-256 (same fetch) | `5C303D7BA0DB892411E739A9920B3E0FB3C62416344CD7F220F359BDD91C0C5B` |
| POM developer | `martijno` — Martijn Oostdijk |
| Upstream SCM (from 0.8.6 POM) | `scm:svn:https://svn.code.sf.net/p/jmrtd/code` · http://jmrtd.org · https://sourceforge.net/projects/jmrtd/ |
| Retired pin | E3V3A/JMRTD tag `0.5.2` @ `0b71be7…` (archived GitHub mirror; last push 2021-06-23) |

**Re-check before each paper-grade resolve:** re-read `maven-metadata.xml` — do not inherit version from search UI alone.

## Why not Option B (alternate pin)

- E3V3A is **not** a plausible stand-in for “JMRTD as an integrator would encounter it in 2026.”
- Reviewer can pull Central in ~30s; measuring a 2021 archive while live releases continue is a **credibility / cherry-pick risk**.
- Disclosure symmetry: I.10 pairs disclosure + PR for both libraries. Option A aligns **measurement target = disclosure target**.

## Class B + S_API — re-verify before counting re-runs

0.5.2 → 0.8.6 is roughly a decade. **Do not** carry Class B / S_API assumptions forward.

See in-repo Class B re-verify notes under `docs/` (Option A / 0.8.6 pin).

| Check | Preliminary (19 Jul) |
| --- | --- |
| Shipped demo in Central / sources | **None** — Class B “no official demo” **still holds** |
| Public `PassportService.doPACE` | Declares `throws CardServiceException` (integrator boundary) |
| Internal `PACEProtocol` | May throw `CardServiceProtocolException`; **do not** pin catch from this layer alone |
| Empirical gate | `catch (Exception)` + per-cell class/SW; suite continues after cell fail — **2a design closed**; fill table in 2b |

## Harness wiring

- `drivers/jmrtd/pom.xml` → dependency version **0.8.6**
- `scripts/bootstrap-vendor.sh` → clones **gmrtd only**; documents JMRTD via Maven
- `scripts/install-jmrtd-local.sh` → `mvn dependency:get` for Central artifact

## API migration notes (0.5.2 → 0.8.x)

- `PassportService(CardService)` removed → `PassportServices.open(card)` helper
- `doPACE(key, oid, params)` → `doPACE(key, oid, params, parameterId)`
- `JMRTDSecurityProvider` removed → BouncyCastle provider directly
- `Util.computeKeySeedForBAC` → `Util.computeKeySeed(..., "SHA-1", true)`
- **S_API:** diagnostic `Exception` (not `Throwable`/`Error`); per-cell class in manifest; narrow after distribution; SW×class = finding if present
- **Isolation:** subprocess-per-cell; `run_suite.py` continues after non-zero cell

## Evidence sequence (priority #1)

1. Finish Class B / harness S_API alignment on 0.8.6 — **done**  
2. Re-run smoke TC-AC-01 (baseline + mitigated) — **done** (20 Jul)  
3. Re-run headline N=40 — **done** (test-server 20 Jul) → see [`OPTION-A-HEADLINE-40-2026-07-20.md`](OPTION-A-HEADLINE-40-2026-07-20.md)  
4. New headline manifest SHA — **`3A8D18ADE8A66DDFFAF3344606AD5382F80F544A6831A08A47DBE1202775176B`** (replaces `fa84b9cd…` for Option A headline claims)  
5. Paper methods: coordinate **0.8.6** + jar SHA-256 + fetch date  
6. Full factorial N=200 (`ac-01-sweep-full`) — **done** (20 Jul) → [`OPTION-A-SWEEP-FULL-2026-07-20.md`](OPTION-A-SWEEP-FULL-2026-07-20.md) · SHA `D505D5212480A68B2253020A6C701C4F50FA71EE8A7DA23000F970C573C60051`

**Pre-commit:** diverge from 0.5.2-era = **finding**. Headline soft deadline **~29 Jul 2026** — **met** (20 Jul). Full factorial also met same day.

## Disclosure contact (JMRTD side)

- https://jmrtd.org · https://sourceforge.net/projects/jmrtd/ · `martijno` / `info@jmrtd.org`  
- **Not** E3V3A

## Provenance of the retired pin

See in-repo provenance notes under `docs/` (retired E3V3A pin rationale).
