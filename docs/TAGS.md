# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

Prefer citing the **tag name**. Short hashes below are the commits the tags currently peel to (`git rev-parse <tag>^{}`).

| Tag | Commit | Purpose |
| --- | --- | --- |
| [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | `96fa6a4` | **Dev.to B10 essay** — TC-AC-01 smoke (N=1), gmrtd + JMRTD baseline |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) | `96fa6a4` | First public smoke-tier release (same tree as blog tag) |
| [`paper-manifest-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07) | `9062a08` | **Stale (7 Jul)** — predates 9 Jul 200-run sweep; do not hand to reviewers |
| [`paper-manifest-2026-07-25`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-25) | `42b4f1a` | **Historical** — superseded by `paper-manifest-2026-07-26`. Kept so old links do not 404. At this commit, some docs still cite retired digest `fa84b9…`; do not hand to reviewers as the live pin. |
| [`paper-manifest-2026-07-26`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/paper-manifest-2026-07-26) | *(see `git rev-parse paper-manifest-2026-07-26^{}`)* | **Live IEEE Access / paper review pin** — AC-01 primary digest `d8afa161…` (immutable suite object), Option A `d505d521…`, CA SM-session continue-check + bundle `43f9bd1d…`. Source public; locked full-run JSON trees are confidential Supplementary Material for Access review; public deposit follows informal publication timing (no fixed maintainer-facing embargo — [`DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md`](DISCLOSURE-AND-DATA-AVAILABILITY-2026-07-26.md)). Force-retargeted for history scrub(s) only (commit-hash rename; suite digests unchanged). **Cite the tag name**, not a short hash. New evidence semantics → new dated tag. |
| [`tc-ac-adv-2026-07-28`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/tc-ac-adv-2026-07-28) | *(tip at tag time; peels to ADV docs + `10137e5` feature)* | **Shallow ADV corroboration pin** — `profiles/adv/`, suite digest `99d38845…` in manuscript §6.1.1. Does **not** replace `paper-manifest-2026-07-26`. Not RF / silicon. |
| [`ac01-n100-lab`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/ac01-n100-lab) | `9e3edcb` | Lab-verified AC-01 wire suite N=100 (400 runs, 4 cells) |

## Quick checkout

```bash
# Blog / smoke reproduction
git checkout blog-b10-2026-07 && make smoke

# Paper evidence pin (source public; locked run-trees confidential until public deposit)
git checkout paper-manifest-2026-07-26
docker build -t emrtd-harness . && docker run --rm emrtd-harness

# N=100 AC-01 wire evidence (longer run)
git checkout ac01-n100-lab && make suite
```

## Tag policy

- **SemVer / SoftwarX C1** → `v0.1.0` (same tree as `blog-b10-2026-07`). See [`CHANGELOG.md`](../CHANGELOG.md).
- **Blog citations** → `blog-b10-2026-07` only (not `main`, not `paper-manifest-*`). Dual-cite with `v0.1.0` when a package version is required.
- **Paper evidence (reviewers)** → `paper-manifest-2026-07-26` (not `paper-manifest-2026-07-25`, not `paper-manifest-2026-07`, not `main`).
- **ADV channel-abort corroboration (shallow)** → annotated tag [`tc-ac-adv-2026-07-28`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/tc-ac-adv-2026-07-28) (feature landed at `10137e5`, `profiles/adv/`). Primary factorial pin remains `paper-manifest-2026-07-26`; ADV uses manuscript suite digest `99d38845…`. Not RF / silicon.
- **N=100 rates in paper §VI** → `ac01-n100-lab` until re-verified at manifest freeze.
- **Public locked-run / Zenodo deposit** → judgment under informal publication timing; **not** a first-time source release (source already public via blog tags); **not** a hard maintainer-facing embargo date.
- **Freeze rule** → suite digests are immutable. History scrubs that only rename commits (no digest change) may retarget this tag; new evidence semantics → new dated tag. SemVer advances (`v0.1.1`, `v0.2.0`) when the public release boundary changes.

**Digest vs tag:** Suite digests (`d8afa161…`, `d505d521…`, CA `43f9bd1d…`) are hashes of locked evidence trees — the immutable reproduction objects. The review tag pins harness *source semantics* that match those claims. Regenerating suites reproduces Observability Scores but not identical timestamped hashes.

New tags are added when a reproducibility boundary changes — not on every doc commit.
