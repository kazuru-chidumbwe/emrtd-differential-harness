#!/usr/bin/env python3
"""TC-PA-01 / TC-PA-03 offline smoke — pymrtd tier (stratified, not pooled with wire)."""

from __future__ import annotations

import json
import os
import sys
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Optional

ROOT = Path(__file__).resolve().parents[2]
CASES = [
    ROOT / "testdata/sod/tc-pa-01-weak-digest.json",
    ROOT / "testdata/sod/tc-pa-03-expired-dsc.json",
    ROOT / "testdata/sod/tc-pa-04a-chained-fresh.json",
    ROOT / "testdata/sod/tc-pa-04b-chained-expired.json",
]


@dataclass
class OfflineResult:
    run_id: str
    test_case: str
    library: str
    tier: str
    mechanism: str
    condition: str
    fixture_ready: bool
    pymrtd_available: bool
    verify_attempted: bool
    verify_error: Optional[str]
    observability_score: int
    observability_meaning: str
    notes: str


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

    sod_path = ROOT / case["fixture"]
    raw = sod_path.read_bytes()
    try:
        sod = SOD.load(raw)
        sod.verify()
        return True, None
    except Exception as exc:  # noqa: BLE001 — harness captures library surface
        return True, str(exc)


def classify(case_id: str, verify_err: Optional[str], attempted: bool) -> tuple[int, str]:
    if not attempted:
        return 1, "fixture scaffold — pymrtd verify not run (install pymrtd to execute)"
    if verify_err:
        return 2, "surfaced — verify raised or returned error"
    if case_id in ("TC-PA-03", "TC-PA-04b"):
        return 0, "silent risk — verify succeeded without inspection-date / chain-validity policy in naive caller"
    if case_id == "TC-PA-04a":
        return 0, "chained fresh DSC — naive verify succeeded (CMS OK; trust-policy not applied)"
    return 1, "logged — verify outcome requires policy comparison"


def main() -> int:
    log_dir = Path(os.environ.get("LOG_DIR", ROOT / "logs"))
    log_dir.mkdir(parents=True, exist_ok=True)
    exit_code = 0

    for case_path in CASES:
        case = load_case(case_path)
        ready = fixtures_present(case)
        attempted = False
        verify_err: Optional[str] = None

        if ready:
            attempted, verify_err = try_pymrtd_verify(case)
        else:
            verify_err = "fixture files pending code-signing lab pass"

        obs, meaning = classify(case["id"], verify_err, attempted and ready)
        run_id = f"{case['id']}-pymrtd-{datetime.now(timezone.utc).strftime('%Y%m%dT%H%M%SZ')}"

        result = OfflineResult(
            run_id=run_id,
            test_case=case["id"],
            library="pymrtd",
            tier=case.get("tier", "offline"),
            mechanism=case.get("mechanism", "PA"),
            condition=case.get("condition", ""),
            fixture_ready=ready,
            pymrtd_available=attempted,
            verify_attempted=attempted and ready,
            verify_error=verify_err,
            observability_score=obs,
            observability_meaning=meaning,
            notes=case.get("notes", ""),
        )

        out = log_dir / f"{run_id}.json"
        out.write_text(json.dumps(asdict(result), indent=2), encoding="utf-8")
        print(json.dumps(asdict(result), indent=2))

        if not ready:
            exit_code = 1

    if exit_code == 0:
        print("OFFLINE PA SMOKE OK — traces under logs/")
    else:
        print("OFFLINE PA SCAFFOLD — fixture PEM/SOD hex pending; metadata and runner ready", file=sys.stderr)
    return exit_code


if __name__ == "__main__":
    sys.exit(main())
