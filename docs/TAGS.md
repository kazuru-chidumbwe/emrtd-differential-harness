# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

**Paper / package cite (bpfix-style):** tag **`v1.0.3`** · permanent link  
https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.3

Prefer citing the **tag name**. Short hashes below are illustrative; peel with `git rev-parse <tag>^{}`.

| Tag | Purpose |
| --- | --- |
| [`v1.0.3`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.3) | **Live paper cite pin** — path-sanitized locked-run digests; PA fixtures DER-clean; venue-agnostic tip. Cite this for manuscripts. |
| [`v1.0.2`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.2) | PA fixture regen. Superseded for *cite* by `v1.0.3`. |
| [`v1.0.1`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.1) | Soft alternate-venue wording. Historical. |
| [`v1.0.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.0) | First paper-ready SemVer. Historical. |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) | First public **smoke-tier** SemVer (same tree as `blog-b10-2026-07`) |
| [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | **Dev.to B10 essay** — TC-AC-01 smoke (N=1) |
| [`paper-manifest-2026-07-26`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-26) | **Historical** evidence freeze tag name. Live digests are on Release `locked-runs-2026-08-10`. |
| [`paper-manifest-2026-07-25`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-25) | Historical; superseded |
| [`paper-manifest-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07) | Stale (7 Jul); do not cite |
| [`tc-ac-adv-2026-07-28`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/tc-ac-adv-2026-07-28) | Historical ADV feature pin (`profiles/adv/`; digest `99d38845…`). Covered by tip. |
| [`locked-runs-2026-08-10`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/locked-runs-2026-08-10) | **Live** path-sanitized locked-run zip (also on `v1.0.3`) |
| [`locked-runs-2026-07-26`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/locked-runs-2026-07-26) | Historical tag; asset replaced with sanitized zip |
| [`ac01-n100-lab`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/ac01-n100-lab) | Lab-verified AC-01 wire suite N=100 |

## Quick checkout

```bash
# Paper reproduction (cite this)
git checkout v1.0.3
bash scripts/bootstrap-vendor.sh && export GOTOOLCHAIN=auto && make smoke

# Blog / smoke only
git checkout v0.1.0 && make smoke
```

## Tag policy

- **Paper / package citation** → **`v1.0.3`** (permanent tree URL above). See [`CHANGELOG.md`](../CHANGELOG.md).
- **Blog citations** → `blog-b10-2026-07` / `v0.1.0` only.
- **Dated `paper-manifest-*` / `locked-runs-*` / `tc-ac-adv-*` / earlier SemVer** → historical aliases; suite digests remain the hash-verification objects.
- **Freeze rule** → suite digests are immutable. New evidence semantics → new SemVer.

**Digest vs tag:** Suite digests (`6c6cbbd1…`, `2efd190e…`, CA `43f9bd1d…`) are hashes of locked evidence trees. `v1.0.3` pins harness *source semantics* that match those claims. Regenerating suites reproduces Observability Scores but not identical timestamped hashes.
