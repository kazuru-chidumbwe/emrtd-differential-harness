#!/usr/bin/env python3
"""Generate a bounded TC-CA-01 profile matrix.

Axes the harness actually consumes today (verified against
simulator/tc_ca_01_transceiver.go and TcCa01CardService.java):

  - ca_sw — returned on first CA MSE:Set AT (INS 0x22, P1 0x41, P2 0xA4)
  - skew_direction — catalog condition label for v2-chip/v1-terminal vs
    v1-chip/v2-terminal. Both use the same synthetic injection point
    (first_ca_mse_set_at). DG14 fixture advertises CA; true CA version
    negotiation is not simulated — the controlled variable is the status
    word + direction label for cross-cell comparison.

Does not overwrite profiles/ca-v1-v2-skew.json or ca-v2-terminal-v1.json
(smoke anchors). Generated files live under profiles/ca-sweep/.
"""

from __future__ import annotations

import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
OUT_DIR = ROOT / "profiles" / "ca-sweep"
DG14 = "testdata/dg14/ca-v2-sample.hex"

MRZ = {
    "document_number": "L898902C",
    "date_of_birth": "690806",
    "date_of_expiry": "940623",
}

# Same SW set as AC-01 sweep (minus proprietary-only continuity extras already covered).
STATUS_WORDS = ["6FFF", "6300", "6982", "6A88", "6A80"]

DIRECTIONS = [
    {
        "id": "v2chip_v1term",
        "condition": "ca_v2_chip_v1_terminal_skew",
        "notes": (
            "Proxy for CA v2 chip / v1 terminal skew: after BAC, first CA MSE:Set AT "
            "returns ca_sw. DG14 advertises CA; full CA version negotiation is not simulated."
        ),
    },
    {
        "id": "v1chip_v2term",
        "condition": "ca_v1_chip_v2_terminal_skew",
        "notes": (
            "Proxy for CA v1 chip / v2 terminal skew (reverse label). Same first-MSE "
            "injection as the forward skew; direction is the catalog condition for matrix cells."
        ),
    },
]


def build(sw: str, direction: dict) -> dict:
    pid = f"TC-CA-01-sweep-{direction['id']}-{sw.lower()}"
    return {
        "id": pid,
        "name": f"CA sweep: direction={direction['id']} sw={sw}",
        "mechanism": "EAC-CA",
        "condition": direction["condition"],
        "tier": "wire",
        "seed": 1,
        "mrz": MRZ,
        "card_access_hex": "",
        "dg14_hex_path": DG14,
        "ca_injection": {
            "ca_fail_on": "first_ca_mse_set_at",
            "ca_sw": sw,
            "notes": direction["notes"] + f" Status word {sw}.",
        },
        "sweep_axes": {
            "skew_direction": direction["id"],
            "ca_sw": sw,
        },
        "expected_gmrtd": {
            "chip_auth_err": "non_null",
            "perform_chip_auth_step_err": "null",
            "observability_if_caller_checks_step_only": 0,
        },
    }


def main() -> None:
    OUT_DIR.mkdir(parents=True, exist_ok=True)
    entries = []
    for direction in DIRECTIONS:
        for sw in STATUS_WORDS:
            profile = build(sw, direction)
            path = OUT_DIR / f"{profile['id']}.json"
            path.write_text(json.dumps(profile, indent=2) + "\n", encoding="utf-8")
            entries.append(
                {
                    "id": profile["id"],
                    "path": str(path.relative_to(ROOT)).replace("\\", "/"),
                    "condition": profile["condition"],
                    "ca_sw": sw,
                    "skew_direction": direction["id"],
                }
            )
    index = {
        "generator": "profiles/generate_ca_sweep.py",
        "n_profiles": len(entries),
        "axes": {"skew_direction": 2, "ca_sw": len(STATUS_WORDS)},
        "factorial_note": (
            "10 profiles × 2 libraries × 2 variants = 40 runs when fully executed. "
            "Injection is first CA MSE:Set AT only; ca_fail_on is not branched in simulators."
        ),
        "profiles": entries,
    }
    (OUT_DIR / "index.json").write_text(json.dumps(index, indent=2) + "\n", encoding="utf-8")
    print(f"wrote {len(entries)} profiles under {OUT_DIR}")


if __name__ == "__main__":
    main()
