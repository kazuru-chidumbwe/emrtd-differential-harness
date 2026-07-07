"""Frozen schema version constants — bump only with explicit migration."""

ARTIFACT_MANIFEST_VERSION = 1
PROVENANCE_VERSION = 1
RUN_ARTIFACT_VERSION = 1
PROFILE_CATALOG_VERSION = 1

METHODOLOGY_NOTE = (
    "Repeating each deterministic profile N=100 demonstrates harness stability and "
    "result reproducibility rather than estimating behavioural variance."
)

FINDING_THRESHOLD_PCT = 95.0
SUITE_N_DEFAULT = 100
