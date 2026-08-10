"""Frozen schema version constants — bump only with explicit migration."""

ARTIFACT_MANIFEST_VERSION = 1
PROVENANCE_VERSION = 1
RUN_ARTIFACT_VERSION = 1
PROFILE_CATALOG_VERSION = 1

METHODOLOGY_NOTE = (
    "In-process simulators are deterministic (no RNG, no concurrency). "
    "Each profile is therefore run at N=1; multi-profile cells (e.g. the 50-profile "
    "AC-01 sweep) establish consistency across profiles, not repeated identical trials. "
    "The ≥95% finding threshold applies to multi-profile cells."
)

FINDING_THRESHOLD_PCT = 95.0
SUITE_N_DEFAULT = 1
