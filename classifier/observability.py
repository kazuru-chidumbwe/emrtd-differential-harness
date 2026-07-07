"""Observability Score classifier for R1 harness (shared contract with classifier/observability.go)."""

from dataclasses import dataclass
from enum import IntEnum
from typing import Optional

from constants import FINDING_THRESHOLD_PCT, SUITE_N_DEFAULT


class ObservabilityScore(IntEnum):
    SILENT = 0
    LOGGED = 1
    SURFACED = 2


@dataclass
class RunOutcome:
    library: str
    mechanism: str
    condition: str
    returned_error: Optional[str]
    trace_has_warning: bool
    trace_has_apdu_error: bool


@dataclass
class TCAC01Outcome:
    pace_failed: bool
    bac_success: bool
    bac_err: Optional[str]
    pace_surfaced_to_caller: bool


def classify(outcome: RunOutcome) -> ObservabilityScore:
    """Generic trace + return-value classifier for non-TC-AC-01 tiers."""
    if outcome.returned_error:
        return ObservabilityScore.SURFACED
    if outcome.trace_has_warning or outcome.trace_has_apdu_error:
        return ObservabilityScore.LOGGED
    return ObservabilityScore.SILENT


def classify_tc_ac_01(outcome: TCAC01Outcome) -> ObservabilityScore:
    """TC-AC-01 wire-tier classifier (must match classifier/observability.go)."""
    bac_err = outcome.bac_err or ""
    if (
        outcome.pace_failed
        and outcome.bac_success
        and bac_err == ""
        and not outcome.pace_surfaced_to_caller
    ):
        return ObservabilityScore.SILENT
    if outcome.pace_failed and not outcome.pace_surfaced_to_caller:
        return ObservabilityScore.LOGGED
    return ObservabilityScore.SURFACED


def consistency_pct(scores: list[ObservabilityScore | int], target: ObservabilityScore) -> float:
    if not scores:
        return 0.0
    matches = sum(1 for s in scores if int(s) == int(target))
    return 100.0 * matches / len(scores)

