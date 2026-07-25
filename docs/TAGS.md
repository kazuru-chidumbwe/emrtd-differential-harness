# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

| Tag | Commit | Purpose |
| --- | --- | --- |
| [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | `ef15b10` | **Dev.to B10 essay** — TC-AC-01 smoke (N=1), gmrtd + JMRTD baseline |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) | `ef15b10` | First public smoke-tier release (same tree as blog tag) |
| [`paper-manifest-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07) | `e31d945` | **Stale (7 Jul)** — predates 9 Jul 200-run sweep; do not hand to reviewers |
| [`paper-manifest-2026-07-25`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-25) | tip of `jmrtd-sweep-2026-07-09` | **Private paper pin** — 200-run AC-01 sweep, CI workflow, PA-04 chain fixtures, Dockerfile |
| [`ac01-n100-lab`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/ac01-n100-lab) | `50154a7` | Lab-verified AC-01 wire suite N=100 (400 runs, 4 cells) |

## Quick checkout

```bash
# Blog / smoke reproduction
git checkout blog-b10-2026-07 && make smoke

# Paper evidence pin (private handoff until 7 Aug public gate)
git checkout paper-manifest-2026-07-25
docker build -t emrtd-harness . && docker run --rm emrtd-harness

# N=100 AC-01 wire evidence (longer run)
git checkout ac01-n100-lab && make suite
```

## Tag policy

- **Blog citations** → `blog-b10-2026-07` only (not `main`, not `paper-manifest-*`).
- **Paper evidence (paper reviewers)** → `paper-manifest-2026-07-25` (not `paper-manifest-2026-07`, not `main`).
- **N=100 rates in paper §VI** → `ac01-n100-lab` until re-verified at manifest freeze.

New tags are added when a reproducibility boundary changes — not on every doc commit.
