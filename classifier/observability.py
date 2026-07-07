"""Observability Score classifier skeleton for R1 harness."""

from dataclasses import dataclass
from enum import IntEnum
from typing import Optional


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


def classify(outcome: RunOutcome) -> ObservabilityScore:
    """Map trace + return-value pair to Observability Score (0/1/2)."""
    if outcome.returned_error:
        return ObservabilityScore.SURFACED
    if outcome.trace_has_warning or outcome.trace_has_apdu_error:
        return ObservabilityScore.LOGGED
    return ObservabilityScore.SILENT


def consistency_pct(scores: list[ObservabilityScore], target: ObservabilityScore) -> float:
    if not scores:
        return 0.0
    matches = sum(1 for s in scores if s == target)
    return 100.0 * matches / len(scores)
