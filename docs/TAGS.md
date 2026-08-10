# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

**Paper / package cite (bpfix-style):** tag **`v1.0.0`** · permanent link  
https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.0

Prefer citing the **tag name**. Short hashes below are illustrative; peel with `git rev-parse <tag>^{}`.

| Tag | Purpose |
| --- | --- |
| [`v1.0.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.0) | **Live paper cite pin** — venue-agnostic tip; full arm set (incl. TC-AC-ADV); `org.emrtd.harness.jmrtd`; locked-run deposit documented. Cite this for manuscripts. |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) | First public **smoke-tier** SemVer (same tree as `blog-b10-2026-07`) |
| [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | **Dev.to B10 essay** — TC-AC-01 smoke (N=1) |
| [`paper-manifest-2026-07-26`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-26) | **Historical** evidence freeze (AC-01 digests `d8afa161…` / Option A `d505d521…`). Superseded for *cite* by `v1.0.0`; kept for digests / old links. |
| [`paper-manifest-2026-07-25`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-25) | Historical; superseded |
| [`paper-manifest-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07) | Stale (7 Jul); do not cite |
| [`tc-ac-adv-2026-07-28`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/tc-ac-adv-2026-07-28) | Historical ADV feature pin (`profiles/adv/`; digest `99d38845…`). Covered by `v1.0.0` tip. |
| [`locked-runs-2026-07-26`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/locked-runs-2026-07-26) | Public locked-run zip deposit (also attached to `v1.0.0` Release) |
| [`ac01-n100-lab`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/ac01-n100-lab) | Lab-verified AC-01 wire suite N=100 |

## Quick checkout

```bash
# Paper reproduction (cite this)
git checkout v1.0.0
bash scripts/bootstrap-vendor.sh && export GOTOOLCHAIN=auto && make smoke

# Blog / smoke only
git checkout v0.1.0 && make smoke
```

## Tag policy

- **Paper / package citation** → **`v1.0.0`** (permanent tree URL above). See [`CHANGELOG.md`](../CHANGELOG.md).
- **Blog citations** → `blog-b10-2026-07` / `v0.1.0` only.
- **Dated `paper-manifest-*` / `locked-runs-*` / `tc-ac-adv-*`** → historical aliases; suite digests remain the hash-verification objects.
- **Freeze rule** → suite digests are immutable. New evidence semantics → new SemVer (`v1.0.1`, `v1.1.0`, …).

**Digest vs tag:** Suite digests (`d8afa161…`, `d505d521…`, CA `43f9bd1d…`) are hashes of locked evidence trees. `v1.0.0` pins harness *source semantics* that match those claims. Regenerating suites reproduces Observability Scores but not identical timestamped hashes.
