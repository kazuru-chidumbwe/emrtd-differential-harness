# Experiment provenance

Every run artifact includes a `provenance` block. Aggregated summaries link published figures to those artifacts.

## Per-run fields

| Field | Meaning |
| --- | --- |
| `harness_commit` | Git commit of the harness at run time |
| `harness_dirty` | Whether the working tree had uncommitted changes |
| `profile_path` | Relative path to the synthetic profile JSON |
| `profile_sha256` | SHA-256 of the profile file |
| `suite_id` | Manifest identifier (e.g. `ac-01-wire`) |
| `suite_seed` | PRNG / catalog seed (metadata; fixed profiles are deterministic) |
| `suite_n` | Repetitions per matrix cell |
| `run_index` | 1-based index within the cell |
| `driver` | Driver binary (e.g. `go/tc-ac-01-mitigated`) |
| `variant` | `baseline` or `mitigated` |
| `middleware` | Middleware policy when applicable |
| `captured_at_utc` | Run timestamp |

## Published figures

`classifier/aggregate.py` emits:

- `summary-{timestamp}.json` — full tuples + `figures[]` with `artifact_refs`
- `summary-{timestamp}.md` — blog-ready table with `figure_id` column
- `artifact-manifest-{timestamp}.json` — SHA-256 per run file

**Rule:** cite `figure_id` and `harness_commit` from the summary; do not hand-copy percentages without the matching summary artifact.

## Example

```bash
make suite                    # suites/ac-01-wire.json, N=100
make suite-paper              # full paper matrix
SUITE_N=10 make suite         # quick check
```

Output: `logs/suite-{suite_id}-{timestamp}/`
