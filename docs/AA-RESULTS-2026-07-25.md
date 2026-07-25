# TC-AA-01 Active Authentication — profile lock (2026-07-25)

Forced AA failure after BAC. Grounding: [internal/aa-error-probe/RESULTS.md](../internal/aa-error-probe/RESULTS.md). Suite: `suites/aa-01-smoke.json`.

## Injection

- BAC succeeds (MRZ).
- First **INTERNAL AUTHENTICATE** (`INS=0x88`) returns `aa_sw` (default `6982`).
- gmrtd: `DoActiveAuth()` → non-nil error, `Success=false`.
- JMRTD: `doAA(...)` throws; baseline catches (naive host); mitigated surfaces.

## Observability Score (by construction)

Baseline drivers set `FailureSurfacedToCaller=false` (naive host). Mitigated drivers set it true via explicit-reject middleware. Unlike TC-AC-01 (`BacSuccess`) and TC-EAC-01 (`ProtectedDGAccessible`), CA/AA have no emergent “session continues” check: identical scores across gmrtd and JMRTD corroborate the naive-host model applied consistently, not a discovered library difference.

| Library | Variant | Score |
|---|---|---:|
| gmrtd | baseline | **0** |
| gmrtd | mitigated | **2** |
| JMRTD | baseline | **0** |
| JMRTD | mitigated | **2** |

## Reproduce

```bash
# gmrtd
./scripts/run_tc_aa_01.sh

# JMRTD (from drivers/jmrtd after package)
java -cp target/jmrtd-tc-ac-01-0.1.0.jar org.emrtd.harness.jmrtd.TcAa01Runner \
  --profile profiles/aa-internal-auth-reject.json --log-dir logs --variant baseline
```

## Notes

- CA↔AA assurance-downgrade is a **separate** claim (not this cell).
- Included in the one-paper evidence-depth ladder as bilateral corroborating (not a second 200-run).
