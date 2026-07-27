# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

Prefer citing the **tag name**. Short hashes below are the commits the tags currently peel to (`git rev-parse <tag>^{}`).

| Tag | Commit | Purpose |
| --- | --- | --- |
| [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | `96fa6a4` | **Dev.to B10 essay** — TC-AC-01 smoke (N=1), gmrtd + JMRTD baseline |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) | `96fa6a4` | First public smoke-tier release (same tree as blog tag) |
| [`paper-manifest-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07) | `9062a08` | **Stale (7 Jul)** — predates 9 Jul 200-run sweep; do not hand to reviewers |
| [`paper-manifest-2026-07-25`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-25) | `42b4f1a` | **Historical** — superseded by `paper-manifest-2026-07-26`. Kept so old links do not 404. At this commit, some docs still cite retired digest `fa84b9…`; do not hand to reviewers as the live pin. |
| [`paper-manifest-2026-07-26`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-26) | `0c5e3df` | **Live Access / paper review pin (FROZEN)** — AC-01 primary digest `d8afa161…` (immutable suite object), Option A `d505d521…`, CA SM-session continue-check + bundle `43f9bd1d…`, digest correction docs. Source public; locked full-run JSON trees remain confidential Supplementary Material pending public deposit (see [`DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md`](DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md)). Force-retargeted **once** on 2026-07-26 after a history scrub (commit-hash rename only; suite digests unchanged). **Do not force-move this tag again.** If the pin must advance, cut a new dated tag (e.g. `paper-manifest-2026-07-27`). |
| [`ac01-n100-lab`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/ac01-n100-lab) | `9e3edcb` | Lab-verified AC-01 wire suite N=100 (400 runs, 4 cells) |

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

- **SemVer / SoftwarX C1** → `v0.1.0` (same tree as `blog-b10-2026-07`). See [`CHANGELOG.md`](../CHANGELOG.md).
- **Blog citations** → `blog-b10-2026-07` only (not `main`, not `paper-manifest-*`). Dual-cite with `v0.1.0` when a package version is required.
- **Paper evidence (reviewers)** → `paper-manifest-2026-07-26` (not `paper-manifest-2026-07-25`, not `paper-manifest-2026-07`, not `main`).
- **N=100 rates in paper §VI** → `ac01-n100-lab` until re-verified at manifest freeze.
- **Public locked-run deposit** → locked full-run packages / Zenodo, **not** a first-time source release (source already public via blog tags).
- **Freeze rule** → never force-move a paper-manifest tag once handed to reviewers. Docs-only commits may land on the branch tip **without** retagging. New evidence semantics → new dated tag. SemVer advances (`v0.1.1`, `v0.2.0`) when the public release boundary changes.

**Digest vs tag:** Suite digests (`d8afa161…`, `d505d521…`, CA `43f9bd1d…`) are hashes of locked evidence trees — the immutable reproduction objects. The review tag pins harness *source semantics* that match those claims. Regenerating suites reproduces Observability Scores but not identical timestamped hashes.

New tags are added when a reproducibility boundary changes — not on every doc commit.
