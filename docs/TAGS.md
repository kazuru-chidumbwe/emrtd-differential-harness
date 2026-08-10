# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

**Paper / package cite (bpfix-style):** tag **`v1.0.5`** · permanent link  
https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.5

Prefer citing the **tag name**. Short hashes below are illustrative; peel with `git rev-parse <tag>^{}`.

| Tag | Purpose |
| --- | --- |
| [`v1.0.5`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.5) | **Live paper cite pin** — lab regen digests; relative `profile_path`; schema-clean deposit `locked-runs-2026-08-10b`. |
| [`v1.0.4`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.4) | CITATION/SECURITY/schemas/CI. Digests superseded by `v1.0.5`. |
| [`v1.0.3`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.3) | Path-sanitized strip digests. Historical. |
| [`v1.0.2`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.2) | PA fixture regen. Historical. |
| [`v1.0.1`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.1) | Soft alternate-venue wording. Historical. |
| [`v1.0.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.0) | First paper-ready SemVer. Historical. |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) | First public **smoke-tier** SemVer (same tree as `blog-b10-2026-07`) |
| [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | **Dev.to B10 essay** — TC-AC-01 smoke (N=1) |
| [`locked-runs-2026-08-10b`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/locked-runs-2026-08-10b) | **Live** lab-regen locked-run zip |
| [`locked-runs-2026-08-10`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/locked-runs-2026-08-10) | Path-sanitized strip zip. Superseded by `…-10b`. |
| [`tc-ac-adv-2026-07-28`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/tc-ac-adv-2026-07-28) | Historical ADV feature pin |

## Quick checkout

```bash
git checkout v1.0.5
bash scripts/bootstrap-vendor.sh && export GOTOOLCHAIN=auto && make smoke
```

## Tag policy

- **Paper / package citation** → **`v1.0.5`**. See [`CHANGELOG.md`](../CHANGELOG.md).
- **Blog citations** → `blog-b10-2026-07` / `v0.1.0` only.
- **Digest vs tag:** Suite digests (`c6f03d7e…`, control `b98b354d…`) are hashes of locked evidence trees. Regenerating suites reproduces Observability Scores but not identical timestamped hashes.
