# Release tags

Annotated tags mark reproducible anchors. **`main` may advance** after a tag — always `git checkout <tag>` when reproducing a cited result.

**Paper / package cite (live):** tag **`v1.0.9`** · permanent link  
https://github.com/kazuru-chidumbwe/emrtd-differential-harness/tree/v1.0.9

Prefer citing the **tag name**. Defect deposit zip: **`emrtd-locked-runs-v1.0.7.zip`**. Remeasurement: Release `v1.0.9` asset `emrtd-gmrtd-v1.1.3-remeasurement-v1.0.8.zip` / digest `04cdd3dd…`.

| Tag | Purpose |
| --- | --- |
| [`v1.0.9`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.9) | **Live paper cite pin** — cite-surface sync (CITATION/DAS/README); dual-pin docs; remeasurement asset. |
| [`v1.0.8`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.8) | Historical — first remeasurement code tag; superseded as live cite by `v1.0.9`. |
| [`v1.0.7`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.7) | Defect factorial deposit — pin `8fea245`; AC-01+control `e15f4b57…` / `b029a9dc…`. |
| [`v1.0.6`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.6) | Historical — pre–ReadDocument baseline; superseded by `v1.0.7`. |
| [`v1.0.5`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.5) | Lab regen; dated `…-10b` zip. Historical. |
| [`v1.0.4`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v1.0.4) | CITATION/SECURITY/schemas/CI. Historical. |
| [`locked-runs-2026-08-10b`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/locked-runs-2026-08-10b) | Dated deposit alias for `v1.0.5` era. |
| [`v0.1.0`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/v0.1.0) / [`blog-b10-2026-07`](https://github.com/kazuru-chidumbwe/emrtd-differential-harness/releases/tag/blog-b10-2026-07) | Blog smoke pin |

## Quick checkout

```bash
# Live cite (remeasurement-capable tip)
git checkout v1.0.9
bash scripts/bootstrap-vendor.sh && export GOTOOLCHAIN=auto

# Defect factorial (Score 0 at 8fea245) — use the deposit tag
git checkout v1.0.7
bash scripts/bootstrap-vendor.sh && export GOTOOLCHAIN=auto && make smoke

# Remeasurement (Score 2 at gmrtd v1.1.3)
git checkout v1.0.9
GMRTD_COMMIT=64bd6ab8fbf8802c718a6da0dcc6f6312a3404ca bash scripts/bootstrap-vendor.sh
python3 classifier/run_suite.py --manifest suites/ac-01-sweep-gmrtd-only.json
```

## Tag policy

- **Live paper citation** → **`v1.0.9`**; cite **`v1.0.7`** digests for the original 200-run defect claim; cite remeasurement digest `04cdd3dd…` from Release `v1.0.9`.
- Digests: primary defect `e15f4b57…`, success-path `b029a9dc…`, remeasurement `04cdd3dd…`.
- Release packaging **must** use `make package-locked-runs STAGING=… VERSION=…` (hard schema + banned-term gate) when refreshing full locked-run deposit zips.
- Zenodo DOI **pending**; GitHub Release URLs are the canonical deposits until a DOI is minted.
