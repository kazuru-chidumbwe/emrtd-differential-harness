"""Cross-language Observability Score vectors (must match Go + Java)."""

from __future__ import annotations

import json
import unittest
from pathlib import Path

from observability import (
    TCAA01Outcome,
    TCAC01Outcome,
    TCCA01Outcome,
    classify_tc_aa_01,
    classify_tc_ac_01,
    classify_tc_ca_01,
)

ROOT = Path(__file__).resolve().parents[1]
VECTORS = ROOT / "testdata" / "observability-vectors.json"


class TestObservabilityVectors(unittest.TestCase):
    def test_shared_vectors(self) -> None:
        rows = json.loads(VECTORS.read_text(encoding="utf-8"))
        for row in rows:
            mech = row["mechanism"]
            inp = row["input"]
            want = row["expected_score"]
            with self.subTest(row["id"]):
                if mech == "TC-AC-01":
                    got = classify_tc_ac_01(
                        TCAC01Outcome(
                            inp["pace_failed"],
                            inp["bac_success"],
                            inp.get("bac_err") or "",
                            inp["pace_surfaced_to_caller"],
                        )
                    )
                elif mech == "TC-CA-01":
                    got = classify_tc_ca_01(
                        TCCA01Outcome(
                            inp["chip_auth_failed"],
                            inp["chip_auth_success"],
                            inp["session_continue_ok"],
                            inp["failure_surfaced_to_caller"],
                        )
                    )
                elif mech == "TC-AA-01":
                    got = classify_tc_aa_01(
                        TCAA01Outcome(
                            inp["active_auth_failed"],
                            inp["active_auth_success"],
                            inp["failure_surfaced_to_caller"],
                        )
                    )
                else:
                    self.fail(f"unsupported mechanism {mech}")
                self.assertEqual(int(got), want)


if __name__ == "__main__":
    unittest.main()
