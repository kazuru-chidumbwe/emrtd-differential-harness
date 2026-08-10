# gmrtd pin (paper reproducibility)

**Pinned commit:** `8fea245048d3b4e76483d048b202ff7f5269728c`  
**Short:** `8fea245`  
**Upstream tag/release:** `0.45.0` (`chore(main): release 0.45.0 (#417)`)  
**Repository:** https://github.com/gmrtd/gmrtd

## Why this pin

The primary AC-01 baseline drives gmrtd’s shipped `reader.ReadDocument` API. At this commit:

- `reader/reader.go` `performPace` records `Session.PaceErr` and **returns nil** (`// NB errors are just recorded at this point`), so `ReadDocument` can succeed after PACE-fail ∧ BAC-ok.
- The bundled reference client `cmd/gmrtd-reader` does **not** inspect `Session.PaceErr` (unlike Passive Authentication session fields). Later local patches that surface PaceErr in the CLI are out of scope for this pin.

JMRTD remains separately locked under Option A (Maven Central `0.8.6`); see [JMRTD-PIN.md](JMRTD-PIN.md). “Option A (locked)” refers to JMRTD, not to a floating gmrtd HEAD.

## Bootstrap

```bash
bash scripts/bootstrap-vendor.sh   # checks out GMRTD_COMMIT
# or override: GMRTD_COMMIT=<sha> bash scripts/bootstrap-vendor.sh
```

`make paper` / `scripts/verify_gmrtd_pin.sh` fail if `$GMRTD_PATH` (default `../_vendor/gmrtd`) is not at this commit.

## Provenance

Per-run JSON may include `gmrtd_commit` resolved from the vendored tree (or `GMRTD_COMMIT` env).
