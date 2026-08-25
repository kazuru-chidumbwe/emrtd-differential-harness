# gmrtd pins (paper reproducibility)

## Defect pin (historical primary factorial)

**Pinned commit:** `8fea245048d3b4e76483d048b202ff7f5269728c`  
**Short:** `8fea245`  
**Upstream tag/release:** `0.45.0` (`chore(main): release 0.45.0 (#417)`)  
**Harness cite:** `v1.0.7` · digest `e15f4b57…`  
**Repository:** https://github.com/gmrtd/gmrtd

At this commit:

- `reader/reader.go` `performPace` records `Session.PaceErr` and **returns nil**, so `ReadDocument` can succeed after PACE-fail ∧ BAC-ok.
- The bundled reference client `cmd/gmrtd-reader` does **not** inspect `Session.PaceErr`.

Default `scripts/gmrtd-pin.sh` and `make paper` on tag `v1.0.7` use this pin.

## Remediation remeasurement pin (v1.1.3)

**Pinned commit:** `64bd6ab8fbf8802c718a6da0dcc6f6312a3404ca`  
**Short:** `64bd6ab`  
**Upstream tag:** `v1.1.3` (24 Aug 2026)  
**PR #446 merge commit:** `1701e74a746a260a5e1707f0c5ef34e100feb32b`  
**Harness cite:** `v1.0.8` · gmrtd-only digest `04cdd3dd…`  
**Lab deposit:** 2026-08-25

At this commit:

- Default is **fail-closed**: a recorded PACE error stops `ReadDocument` before BAC unless `AllowBacFallbackOnPaceError()` is set.
- Locked remeasurement (`ac-01-sweep-gmrtd-only`, 100 runs): baseline **50/50 Score 2**; mitigated **50/50 Score 2**.

```bash
GMRTD_COMMIT=64bd6ab8fbf8802c718a6da0dcc6f6312a3404ca bash scripts/bootstrap-vendor.sh
# requires harness ≥ v1.0.8 (ReaderStatus API + TC-AC-01 gate accepts fail-closed)
python3 classifier/run_suite.py --manifest suites/ac-01-sweep-gmrtd-only.json
```

## Bootstrap

```bash
bash scripts/bootstrap-vendor.sh   # checks out GMRTD_COMMIT (default: defect pin)
# or override: GMRTD_COMMIT=<sha> bash scripts/bootstrap-vendor.sh
```

Per-run JSON may include `gmrtd_commit` resolved from the vendored tree (or `GMRTD_COMMIT` env).
