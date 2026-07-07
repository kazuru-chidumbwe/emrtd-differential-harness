# Simulator profiles

Synthetic chip behaviour definitions — **no physical passport required.**

Catalog: [`catalog.json`](catalog.json) lists profile IDs, paths, seeds, and tiers.

| Profile | Test case | Tier | Status |
| --- | --- | --- | --- |
| `pace-then-bac-downgrade.json` | TC-AC-01 | wire | smoke + suite |
| `pace-then-bac-downgrade-alt-mrz.json` | TC-AC-01 | wire | alternate MRZ |
| `ca-v1-v2-skew.json` | TC-CA-01 | wire | gmrtd + jmrtd |
| `ca-v2-terminal-v1.json` | TC-CA-01 | wire | variant SW 6985 |
| `fi-cardaccess-truncated.json` | TC-FI-01 | wire | scaffold |

Offline PA fixtures live under `testdata/sod/` (pymrtd tier — stratified, not pooled with wire).
