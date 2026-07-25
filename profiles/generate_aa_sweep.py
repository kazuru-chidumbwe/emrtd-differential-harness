#!/usr/bin/env python3
"""Generate a bounded TC-AA-01 SW matrix (5 status words).

Axes consumed by simulator/tc_aa_01_transceiver.go and TcAa01CardService.java:
  - aa_sw — returned on first INTERNAL AUTHENTICATE (INS 0x88)

Does not overwrite profiles/aa-internal-auth-reject.json (smoke anchor).
Generated files live under profiles/aa-sweep/.
"""

from __future__ import annotations

import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
OUT_DIR = ROOT / "profiles" / "aa-sweep"
DG15 = "testdata/dg15/aa-rsa-sample.hex"

MRZ = {
    "document_number": "L898902C",
    "date_of_birth": "690806",
    "date_of_expiry": "940623",
}

# Status words that cause JMRTD doAA / gmrtd DoActiveAuth to hard-fail.
# Note: 6FFF is soft-handled by JMRTD AAProtocol (warning, no throw) — excluded from AA grid.
STATUS_WORDS = ["6300", "6982", "6A80", "6A88", "6985"]


def build(sw: str) -> dict:
    pid = f"TC-AA-01-sweep-{sw.lower()}"
    return {
        "id": pid,
        "name": f"AA sweep: aa_sw={sw}",
        "mechanism": "AA",
        "condition": "aa_internal_authenticate_reject",
        "tier": "wire",
        "seed": 1,
        "mrz": MRZ,
        "card_access_hex": "",
        "dg15_hex_path": DG15,
        "aa_injection": {
            "aa_fail_on": "first_internal_authenticate",
            "aa_sw": sw,
            "notes": (
                f"Synthetic — BAC then AA fails at INTERNAL AUTHENTICATE (0x88) with SW {sw}. "
                "Naive caller that ignores ActiveAuthErr scores 0; mitigated surfaces."
            ),
        },
        "sweep_axes": {"aa_sw": sw},
        "expected_gmrtd": {
            "active_auth_err": "non_null",
            "observability_if_caller_checks_step_only": 0,
        },
    }


def main() -> None:
    OUT_DIR.mkdir(parents=True, exist_ok=True)
    entries = []
    for sw in STATUS_WORDS:
        profile = build(sw)
        path = OUT_DIR / f"{profile['id']}.json"
        path.write_text(json.dumps(profile, indent=2) + "\n", encoding="utf-8")
        entries.append(
            {
                "id": profile["id"],
                "path": str(path.relative_to(ROOT)).replace("\\", "/"),
                "aa_sw": sw,
            }
        )
    index = {
        "generator": "profiles/generate_aa_sweep.py",
        "n_profiles": len(entries),
        "axes": {"aa_sw": len(STATUS_WORDS)},
        "factorial_note": (
            "5 profiles × 2 libraries × 2 variants = 20 runs when fully executed. "
            "Scores follow the naive-host / mitigated model (no emergent continue-check)."
        ),
        "profiles": entries,
    }
    (OUT_DIR / "index.json").write_text(json.dumps(index, indent=2) + "\n", encoding="utf-8")
    print(f"wrote {len(entries)} profiles under {OUT_DIR}")


if __name__ == "__main__":
    main()
