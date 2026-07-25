# Frozen schemas (v1)

Do not change field semantics without bumping the relevant `*_version` integer and documenting migration in this file.

## Pipeline abstraction

```
Scenario  →  Execution  →  Observation  →  Finding  →  Artifact
```

| Layer | Meaning | Where |
| --- | --- | --- |
| **Scenario** | What is being tested (`scenario_id`, profile, library, variant) | Suite manifest + manifest `entries[FIG-xx].scenario` |
| **Execution** | How it runs (driver, middleware policy) | `entries[FIG-xx].execution` |
| **Observation** | Raw scored outcomes (N, silent/logged/surfaced %) | `entries[FIG-xx].observation` |
| **Finding** | Threshold decision (≥95% over N) | `entries[FIG-xx].finding` |
| **Artifact** | Canonical manifest entry + per-run JSON + SHA-256 | `artifact-manifest.json` |

Drivers implement **Execution**. Classifiers implement **Observation**. Aggregators implement **Finding**. The manifest is the **Artifact**.

## Figure IDs

Repository IDs are stable: `FIG-01`, `FIG-02`, …, `TABLE-01`, `SUMMARY-01`, `MANIFEST-01`.

The manuscript maps `FIG-01` → “Figure 5” without renaming repository objects.

Never use prose numbering (“Figure 2”) in the repository.

## Canonical object

**`artifact-manifest.json`** in each suite log directory is the primary published object. Summaries and markdown tables are derived views.

```json
{
  "artifact_version": 1,
  "suite": "ac-01-wire",
  "commit": "<git-sha>",
  "entries": {
    "FIG-01": { "type": "finding", "scenario": {}, "execution": {}, "observation": {}, "finding": {}, "artifact_refs": [], "bundle_sha256": "..." },
    "TABLE-01": { "type": "table", "path": "summary-....md", "sha256": "..." },
    "MANIFEST-01": { "type": "manifest", "path": "artifact-manifest.json", "sha256": "..." }
  }
}
```

## Methodology wording (frozen)

> Repeating each deterministic profile N=100 demonstrates harness stability and result reproducibility rather than estimating behavioural variance.

## CI: `make paper`

Runs tests → smoke → suite → manifest → verify. Fails if any gate fails or `harness_dirty=true`. Copies verified manifest to `artifacts/`.

## Version bumps

| Version | When |
| --- | --- |
| `artifact_version` | Manifest top-level shape changes |
| `provenance_version` | Per-run `provenance` block changes |
| `run_schema_version` | Per-run JSON output changes |
| `profile_catalog_version` | `profiles/catalog.json` changes |

## Per-run `normalized_failure` (AA / TA / EAC)

Optional object on AA, TA, and EAC run JSON (library-native error strings remain). Observability Score contract is unchanged.

```json
{
  "mechanism": "AA",
  "step": "internal_authenticate",
  "iso7816_sw": "6982",
  "failure_class": "chip_sw_reject",
  "surfaced": false
}
```

| Field | Meaning |
| --- | --- |
| `mechanism` | `AA`, `TA`, or `EAC` |
| `step` | Wire/API step (`internal_authenticate`, `pso_verify_certificate`, or `n/a`) |
| `iso7816_sw` | Status word when known (omitted or empty otherwise) |
| `failure_class` | `chip_sw_reject` \| `protocol_exception` \| `peer_unsupported` |
| `surfaced` | Whether the host/middleware treated the failure as a hard stop |

Go: `internal/normfail`. Java: `NormalizedFailure`.
