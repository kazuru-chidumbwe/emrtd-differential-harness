#!/usr/bin/env python3
"""TC-PA-01 / TC-PA-03 / TC-PA-04* offline smoke — pymrtd tier (stratified, not pooled with wire).

Emits the same required run-artifact fields as run_case.py (variant + provenance)
so smoke outputs validate against schemas/run-artifact-v1.json.
"""

from __future__ import annotations

import json
import os
import sys
from dataclasses import asdict, dataclass
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Optional

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "classifier"))
from provenance import collect  # noqa: E402

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
    variant: str
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
        return True, f"fixture content is not valid hex ({exc}) — placeholder, not a real fixture"
    try:
        sod = SOD.load(raw)
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
        return 1, "fixture scaffold — pymrtd verify not run (install pymrtd to execute)"
    if verify_err:
        return 2, "surfaced — verify raised or returned error"
    if case.get("expect_policy_rejection"):
        return 0, "silent risk — verify succeeded with no policy-rejection signal"
    return 1, "logged — verify outcome requires policy comparison"


def main() -> int:
    log_dir = Path(os.environ.get("LOG_DIR", ROOT / "logs"))
    log_dir.mkdir(parents=True, exist_ok=True)
    variant = os.environ.get("VARIANT", "baseline")
    suite_id = os.environ.get("SUITE_ID", "offline-pa-smoke")
    suite_seed = int(os.environ.get("SUITE_SEED", "1"))
    exit_code = 0
    run_index = 0

    for case_path in CASES:
        run_index += 1
        case = load_case(case_path)
        ready = fixtures_present(case)
        attempted = False
        verify_err: Optional[str] = None

        if ready:
            attempted, verify_err = try_pymrtd_verify(case)
        else:
            verify_err = "fixture files pending code-signing lab pass"

        obs, meaning = classify(case, verify_err, attempted and ready)
        run_id = (
            f"{case['id']}-pymrtd-{datetime.now(timezone.utc).strftime('%Y%m%dT%H%M%S%fZ')}"
            f"-{run_index:06d}"
        )

        prov = collect(
            root=ROOT,
            profile_path=case_path,
            suite_id=suite_id,
            suite_seed=suite_seed,
            suite_n=len(CASES),
            run_index=run_index,
            driver="python/pymrtd-offline-smoke",
            variant=variant,
        )

        result = OfflineResult(
            run_id=run_id,
            test_case=case["id"],
            library="pymrtd",
            tier=case.get("tier", "offline"),
            mechanism=case.get("mechanism", "PA"),
            condition=case.get("condition", ""),
            variant=variant,
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
        payload = asdict(result)
        out.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
        print(json.dumps(payload, indent=2))

        if not ready:
            exit_code = 1

    if exit_code == 0:
        print("OFFLINE PA SMOKE OK — traces under logs/")
    else:
        print(
            "OFFLINE PA SCAFFOLD — fixture PEM/SOD hex pending; metadata and runner ready",
            file=sys.stderr,
        )
    return exit_code


if __name__ == "__main__":
    raise SystemExit(main())
