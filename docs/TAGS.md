# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

**Paper / package cite (bpfix-style):** tag **`v1.0.7`** · permanent link  
https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.7

Prefer citing the **tag name**. Deposit zip: **`emrtd-locked-runs-v1.0.7.zip`** on the same Release.

| Tag | Purpose |
| --- | --- |
| [`v1.0.7`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.7) | **Live paper cite pin** — C&S / CoSe M1–M4: ReadDocument baseline + gmrtd pin; AC-01+control deposit; SemVer-named zip. |
| [`v1.0.6`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.6) | Historical — pre–ReadDocument baseline; superseded by `v1.0.7`. |
| [`v1.0.5`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.5) | Lab regen; dated `…-10b` zip. Offline-pa smoke still schema-incomplete. Historical. |
| [`v1.0.4`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.4) | CITATION/SECURITY/schemas/CI. Historical. |
| [`locked-runs-2026-08-10b`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/locked-runs-2026-08-10b) | Dated deposit alias for `v1.0.5` era. Superseded by `emrtd-locked-runs-v1.0.7.zip`. |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) / [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | Blog smoke pin |

## Quick checkout

```bash
git checkout v1.0.7
bash scripts/bootstrap-vendor.sh && export GOTOOLCHAIN=auto && make smoke
```

## Tag policy

- **Paper / package citation** → **`v1.0.7`** + asset `emrtd-locked-runs-v1.0.7.zip`.
- Digests: primary `e15f4b57…`, success-path `b029a9dc…`.
- Release packaging **must** use `make package-locked-runs STAGING=… VERSION=…` (hard schema + banned-term gate).
