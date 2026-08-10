# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

**Paper / package cite (bpfix-style):** tag **`v1.0.6`** · permanent link  
https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.6

Prefer citing the **tag name**. Deposit zip: **`emrtd-locked-runs-v1.0.6.zip`** on the same Release.

| Tag | Purpose |
| --- | --- |
| [`v1.0.6`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.6) | **Live paper cite pin** — schema-complete deposit (incl. offline-pa smoke provenance); SemVer-named zip. |
| [`v1.0.5`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.5) | Lab regen; dated `…-10b` zip. Offline-pa smoke still schema-incomplete. Historical. |
| [`v1.0.4`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.4) | CITATION/SECURITY/schemas/CI. Historical. |
| [`locked-runs-2026-08-10b`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/locked-runs-2026-08-10b) | Dated deposit alias for `v1.0.5` era. Superseded by `emrtd-locked-runs-v1.0.6.zip`. |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) / [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | Blog smoke pin |

## Quick checkout

```bash
git checkout v1.0.6
bash scripts/bootstrap-vendor.sh && export GOTOOLCHAIN=auto && make smoke
```

## Tag policy

- **Paper / package citation** → **`v1.0.6`** + asset `emrtd-locked-runs-v1.0.6.zip`.
- Digests: primary `c6f03d7e…`, success-path `a817198c…`.
- Release packaging **must** use `make package-locked-runs STAGING=… VERSION=…` (hard schema + banned-term gate).
