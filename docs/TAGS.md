# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

| Tag | Commit | Purpose |
| --- | --- | --- |
| [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | `ef15b10` | **Dev.to B10 essay** — TC-AC-01 smoke (N=1), gmrtd + JMRTD baseline |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) | `ef15b10` | First public smoke-tier release (same tree as blog tag) |
| [`paper-manifest-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07) | `e31d945` | **Stale (7 Jul)** — predates 9 Jul 200-run sweep; do not hand to reviewers |
| [`paper-manifest-2026-07-25`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-25) | `110c546` | **Historical** — superseded by `paper-manifest-2026-07-26`. Kept so old links do not 404. At this commit, some docs still cite retired digest `fa84b9…`; do not hand to reviewers as the live pin. |
| [`paper-manifest-2026-07-26`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-26) | `20695ba` | **Live Access / paper review pin** — AC-01 primary digest `d8afa161…` (immutable suite object), Option A `d505d521…`, CA SM-session continue-check + bundle `43f9bd1d…`, digest correction docs. Source public; locked full-run JSON trees remain confidential Supplementary Material pending public deposit (see [`DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md`](DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md)). |
| [`ac01-n100-lab`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/ac01-n100-lab) | `50154a7` | Lab-verified AC-01 wire suite N=100 (400 runs, 4 cells) |

## Quick checkout

```bash
# Blog / smoke reproduction
git checkout blog-b10-2026-07 && make smoke

# Paper evidence pin (source public; locked run-trees gated — see disclosure note)
git checkout paper-manifest-2026-07-26
docker build -t emrtd-harness . && docker run --rm emrtd-harness

# N=100 AC-01 wire evidence (longer run)
git checkout ac01-n100-lab && make suite
```

## Tag policy

- **Blog citations** → `blog-b10-2026-07` only (not `main`, not `paper-manifest-*`).
- **Paper evidence (reviewers)** → `paper-manifest-2026-07-26` (not `paper-manifest-2026-07-25`, not `paper-manifest-2026-07`, not `main`).
- **N=100 rates in paper §VI** → `ac01-n100-lab` until re-verified at manifest freeze.
- **Public locked-run deposit** → locked full-run packages / Zenodo, **not** a first-time source release (source already public via blog tags).

**Digest vs tag:** Suite digests (`d8afa161…`, `d505d521…`, CA `43f9bd1d…`) are hashes of locked evidence trees. The review tag pins harness *source semantics* that match those claims. Regenerating suites reproduces Observability Scores but not identical timestamped hashes.

New tags are added when a reproducibility boundary changes — not on every doc commit.
