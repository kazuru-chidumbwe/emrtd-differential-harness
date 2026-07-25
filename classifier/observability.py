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


@dataclass
class TCCA01Outcome:
    chip_auth_failed: bool
    chip_auth_success: bool
    failure_surfaced_to_caller: bool


def classify_tc_ca_01(outcome: TCCA01Outcome) -> ObservabilityScore:
    """TC-CA-01 wire-tier classifier (must match classifier/observability.go)."""
    if (
        outcome.chip_auth_failed
        and not outcome.chip_auth_success
        and not outcome.failure_surfaced_to_caller
    ):
        return ObservabilityScore.SILENT
    if outcome.chip_auth_failed and not outcome.failure_surfaced_to_caller:
        return ObservabilityScore.LOGGED
    return ObservabilityScore.SURFACED


@dataclass
class TCAA01Outcome:
    active_auth_failed: bool
    active_auth_success: bool
    failure_surfaced_to_caller: bool


def classify_tc_aa_01(outcome: TCAA01Outcome) -> ObservabilityScore:
    """TC-AA-01 wire-tier classifier (must match classifier/observability.go)."""
    if (
        outcome.active_auth_failed
        and not outcome.active_auth_success
        and not outcome.failure_surfaced_to_caller
    ):
        return ObservabilityScore.SILENT
    if outcome.active_auth_failed and not outcome.failure_surfaced_to_caller:
        return ObservabilityScore.LOGGED
    return ObservabilityScore.SURFACED


@dataclass
class TCTA01Outcome:
    terminal_auth_failed: bool
    terminal_auth_success: bool
    failure_surfaced_to_caller: bool
    peer_unsupported: bool = False


def classify_tc_ta_01(outcome: TCTA01Outcome) -> ObservabilityScore:
    """TC-TA-01 classifier (JMRTD-asymmetric; peer_unsupported => surfaced/unsupported)."""
    if outcome.peer_unsupported:
        return ObservabilityScore.SURFACED
    if (
        outcome.terminal_auth_failed
        and not outcome.terminal_auth_success
        and not outcome.failure_surfaced_to_caller
    ):
        return ObservabilityScore.SILENT
    if outcome.terminal_auth_failed and not outcome.failure_surfaced_to_caller:
        return ObservabilityScore.LOGGED
    return ObservabilityScore.SURFACED


@dataclass
class TCEAC01Outcome:
    eac_failed: bool
    eac_success: bool
    protected_dg_accessible: bool
    failure_surfaced_to_caller: bool
    peer_unsupported: bool = False


def classify_tc_eac_01(outcome: TCEAC01Outcome) -> ObservabilityScore:
    """TC-EAC-01 classifier (must match classifier/observability.go)."""
    if outcome.peer_unsupported:
        return ObservabilityScore.SURFACED
    if (
        outcome.eac_failed
        and not outcome.eac_success
        and outcome.protected_dg_accessible
        and not outcome.failure_surfaced_to_caller
    ):
        return ObservabilityScore.SILENT
    if outcome.eac_failed and not outcome.failure_surfaced_to_caller:
        return ObservabilityScore.LOGGED
    return ObservabilityScore.SURFACED


def consistency_pct(scores: list[ObservabilityScore | int], target: ObservabilityScore) -> float:
    if not scores:
        return 0.0
    matches = sum(1 for s in scores if int(s) == int(target))
    return 100.0 * matches / len(scores)

