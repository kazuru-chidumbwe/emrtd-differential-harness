#!/usr/bin/env python3
"""Run a single offline PA case with provenance metadata."""

from __future__ import annotations

import json
import sys
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Optional

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "classifier"))
from provenance import collect  # noqa: E402


@dataclass
class OfflineResult:
    run_id: str
    test_case: str
    library: str
    tier: str
    mechanism: str
    condition: str
    variant: str
    figure_id: Optional[str]
    fixture_ready: bool
    pymrtd_available: bool
    verify_attempted: bool
    verify_error: Optional[str]
    observability_score: int
    observability_meaning: str
    notes: str
    provenance: dict[str, Any]


def load_case(path: Path) -> dict[str, Any]:
    return json.loads(path.read_text(encoding="utf-8"))


def fixtures_present(case: dict[str, Any]) -> bool:
    keys = ("fixture", "dsc_fixture", "csca_fixture")
    return all((ROOT / case[k]).is_file() for k in keys)


def try_pymrtd_verify(case: dict[str, Any]) -> tuple[bool, Optional[str]]:
    try:
        from pymrtd.ef.sod import SOD  # type: ignore
    except ImportError:
        return False, "pymrtd not installed"

    raw_text = (ROOT / case["fixture"]).read_text(encoding="utf-8").strip()
    try:
        raw = bytes.fromhex(raw_text)
    except ValueError as exc:
        # Placeholder fixtures are non-hex comment text ("# Placeholder ...") — this branch
        # is what "fixture present but not real" looks like, distinct from a genuine parse
        # or verification failure below.
        return True, f"fixture content is not valid hex ({exc}) — placeholder, not a real fixture"
    try:
        sod = SOD.load(raw)
        # SOD.verify(si, dsc) requires the specific SignerInfo and its matching Document
        # Signer Certificate — it is not a no-argument call. A previous version of this
        # driver called sod.verify() with no arguments, which always raised a TypeError
        # about missing arguments regardless of the fixture's actual content; that was
        # misreported as a "surfaced" finding for the wrong reason (a driver bug, not a
        # pymrtd digest-policy decision). Verify against every signer present.
        errors = []
        for si in sod.signers:
            dsc = sod.getDscCertificate(si)
            try:
                sod.verify(si, dsc)
            except Exception as verify_exc:  # noqa: BLE001
                errors.append(str(verify_exc))
        if errors:
            return True, "; ".join(errors)
        return True, None
    except Exception as exc:  # noqa: BLE001
        return True, str(exc)


def classify(case: dict[str, Any], verify_err: Optional[str], attempted: bool) -> tuple[int, str]:
    if not attempted:
        # Intentional stub fixtures only — paper/CI must install pymrtd (fail closed below).
        return 1, "deferred — verify not run (non-paper scaffold path)"
    if verify_err:
        return 2, "surfaced — verify raised or returned error"
    if case.get("expect_policy_rejection"):
        # verify() succeeded with no error, but this case's fixture was constructed
        # specifically to violate a policy condition (weak digest, expired DSC, etc.) that
        # a caller relying on verify() alone would want flagged. Silent success under these
        # conditions is the same "green light with no built-in signal" pattern as a silent
        # PACE-to-BAC downgrade — Observability Score 0, not a case-by-case exception.
        # Fixed per reviewer ruling, round 3: this was previously a hardcoded
        # `case_id == "TC-PA-03"` special case, which meant a structurally identical
        # finding for TC-PA-01 fell through to a weaker "logged" score by omission, not by
        # a considered classification rule. All future PA cases should set
        # expect_policy_rejection explicitly in their case JSON rather than relying on a
        # case-ID allowlist here.
        return 0, "silent risk — verify succeeded with no policy-rejection signal"
    return 1, "logged — verify outcome requires policy comparison"


def main() -> int:
    if len(sys.argv) < 2:
        print("usage: run_case.py <fixture.json> [variant] [run_index] [suite_id] [seed] [n]", file=sys.stderr)
        return 2

    case_path = Path(sys.argv[1])
    variant = sys.argv[2] if len(sys.argv) > 2 else "baseline"
    run_index = int(sys.argv[3]) if len(sys.argv) > 3 else 1
    suite_id = sys.argv[4] if len(sys.argv) > 4 else ""
    seed = int(sys.argv[5]) if len(sys.argv) > 5 else 1
    suite_n = int(sys.argv[6]) if len(sys.argv) > 6 else 1

    case = load_case(case_path)
    ready = fixtures_present(case)
    attempted = False
    verify_err: Optional[str] = None

    if ready:
        attempted, verify_err = try_pymrtd_verify(case)
    else:
        verify_err = "fixture files pending"

    # Fail closed: missing pymrtd must not score as evidence.
    if verify_err == "pymrtd not installed":
        print("FATAL: pymrtd is required for offline PA evidence runs", file=sys.stderr)
        return 2

    obs, meaning = classify(case, verify_err, attempted and ready)
    run_id = f"{case['id']}-pymrtd-{datetime.now(timezone.utc).strftime('%Y%m%dT%H%M%S%fZ')}-{run_index:06d}"

    prov = collect(
        root=ROOT,
        profile_path=case_path,
        suite_id=suite_id,
        suite_seed=seed,
        suite_n=suite_n,
        run_index=run_index,
        driver="python/pymrtd-offline",
        variant=variant,
    )

    import os

    log_dir = Path(os.environ.get("LOG_DIR", ROOT / "logs"))
    log_dir.mkdir(parents=True, exist_ok=True)

    result = OfflineResult(
        run_id=run_id,
        test_case=case["id"],
        library="pymrtd",
        tier=case.get("tier", "offline"),
        mechanism=case.get("mechanism", "PA"),
        condition=case.get("condition", ""),
        variant=variant,
        figure_id=case.get("figure_id"),
        fixture_ready=ready,
        pymrtd_available=attempted,
        verify_attempted=attempted and ready,
        verify_error=verify_err,
        observability_score=obs,
        observability_meaning=meaning,
        notes=case.get("notes", ""),
        provenance=prov,
    )

    out = log_dir / f"{run_id}.json"
    out.write_text(json.dumps(asdict(result), indent=2), encoding="utf-8")
    print(json.dumps(asdict(result), indent=2))
    return 0 if ready else 1


if __name__ == "__main__":
    raise SystemExit(main())
