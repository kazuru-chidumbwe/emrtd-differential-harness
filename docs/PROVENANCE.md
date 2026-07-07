# Experiment provenance and canonical artifacts

## Canonical object

**`artifact-manifest.json`** is the primary published object. Every `FIG-xx` entry is dispute-resistant: scenario, execution, observation, finding, `artifact_refs`, and `bundle_sha256`.

Derived files (`summary-*.md`, `summary-*.json`) are views — cite the manifest.

## Reproducing a figure from clean checkout

```bash
git checkout <commit from manifest>
bash scripts/bootstrap-vendor.sh
make paper    # tests → smoke → suite → verify (fails on stale/incomplete)
```

The verified manifest is copied to `artifacts/<suite>-<commit>-artifact-manifest.json`.

## Figure IDs

Use repository IDs (`FIG-01`, `TABLE-01`). Map to manuscript numbering externally:

```
FIG-01  →  Figure 5 (manuscript)
```

## Per-run provenance (`provenance_version`: 1)

Each run JSON embeds: `harness_commit`, `profile_sha256`, `suite_id`, `suite_seed`, `driver`, `variant`, `middleware`, `run_index`.

## Methodology (frozen wording)

Repeating each deterministic profile N=100 demonstrates harness stability and result reproducibility rather than estimating behavioural variance.

## Pipeline layers

See [SCHEMA.md](SCHEMA.md): Scenario → Execution → Observation → Finding → Artifact.
