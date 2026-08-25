# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

**Paper / package cite (live):** tag **`v1.0.8`** · permanent link  
https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.8

Prefer citing the **tag name**. Defect deposit zip remains **`emrtd-locked-runs-v1.0.7.zip`**. Remeasurement: Release `v1.0.8` asset / digest `04cdd3dd…`.

| Tag | Purpose |
| --- | --- |
| [`v1.0.8`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.8) | **Live paper cite pin** — gmrtd `v1.1.3` remeasurement (100-run gmrtd-only); ReaderStatus API; dual-pin docs. |
| [`v1.0.7`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.7) | Historical defect factorial — pin `8fea245`; AC-01+control deposit `e15f4b57…` / `b029a9dc…`. |
| [`v1.0.6`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.6) | Historical — pre–ReadDocument baseline; superseded by `v1.0.7`. |
| [`v1.0.5`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.5) | Lab regen; dated `…-10b` zip. Historical. |
| [`v1.0.4`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.4) | CITATION/SECURITY/schemas/CI. Historical. |
| [`locked-runs-2026-08-10b`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/locked-runs-2026-08-10b) | Dated deposit alias for `v1.0.5` era. |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) / [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | Blog smoke pin |

## Quick checkout

```bash
# Defect factorial (Score 0 at 8fea245)
git checkout v1.0.7
bash scripts/bootstrap-vendor.sh && export GOTOOLCHAIN=auto && make smoke

# Remeasurement (Score 2 at gmrtd v1.1.3)
git checkout v1.0.8
GMRTD_COMMIT=64bd6ab8fbf8802c718a6da0dcc6f6312a3404ca bash scripts/bootstrap-vendor.sh
python3 classifier/run_suite.py --manifest suites/ac-01-sweep-gmrtd-only.json
```

## Tag policy

- **Live paper citation** → **`v1.0.8`** (includes remeasurement code + docs); cite **`v1.0.7`** digests for the original 200-run defect claim.
- Digests: primary defect `e15f4b57…`, success-path `b029a9dc…`, remeasurement `04cdd3dd…`.
- Release packaging **must** use `make package-locked-runs STAGING=… VERSION=…` (hard schema + banned-term gate) when refreshing deposit zips.
