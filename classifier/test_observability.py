"""Unit tests for Observability Score classifiers (must match Go + Java contracts)."""

from __future__ import annotations

import unittest

from observability import (
    ObservabilityScore,
    TCAA01Outcome,
    TCAC01Outcome,
    TCCA01Outcome,
    TCEAC01Outcome,
    TCTA01Outcome,
    classify_tc_aa_01,
    classify_tc_ac_01,
    classify_tc_ca_01,
    classify_tc_eac_01,
    classify_tc_ta_01,
)


class TestClassifyTCAC01(unittest.TestCase):
    def test_silent(self) -> None:
        self.assertEqual(
            classify_tc_ac_01(
                TCAC01Outcome(True, True, "", False)
            ),
            ObservabilityScore.SILENT,
        )

    def test_logged(self) -> None:
        self.assertEqual(
            classify_tc_ac_01(
                TCAC01Outcome(True, False, "bac failed", False)
            ),
            ObservabilityScore.LOGGED,
        )

    def test_surfaced(self) -> None:
        self.assertEqual(
            classify_tc_ac_01(
                TCAC01Outcome(True, False, "", True)
            ),
            ObservabilityScore.SURFACED,
        )


class TestClassifyTCCA01(unittest.TestCase):
    def test_silent(self) -> None:
        self.assertEqual(
            classify_tc_ca_01(TCCA01Outcome(True, False, True, False)),
            ObservabilityScore.SILENT,
        )

    def test_logged(self) -> None:
        self.assertEqual(
            classify_tc_ca_01(TCCA01Outcome(True, False, False, False)),
            ObservabilityScore.LOGGED,
        )

    def test_surfaced(self) -> None:
        self.assertEqual(
            classify_tc_ca_01(TCCA01Outcome(True, False, True, True)),
            ObservabilityScore.SURFACED,
        )


class TestClassifyTCAA01(unittest.TestCase):
    def test_silent(self) -> None:
        self.assertEqual(
            classify_tc_aa_01(TCAA01Outcome(True, False, False)),
            ObservabilityScore.SILENT,
        )

    def test_logged(self) -> None:
        self.assertEqual(
            classify_tc_aa_01(TCAA01Outcome(True, True, False)),
            ObservabilityScore.LOGGED,
        )

    def test_surfaced(self) -> None:
        self.assertEqual(
            classify_tc_aa_01(TCAA01Outcome(True, False, True)),
            ObservabilityScore.SURFACED,
        )


class TestClassifyTCTA01(unittest.TestCase):
    def test_peer_unsupported(self) -> None:
        self.assertEqual(
            classify_tc_ta_01(TCTA01Outcome(False, False, False, True)),
            ObservabilityScore.SURFACED,
        )

    def test_silent(self) -> None:
        self.assertEqual(
            classify_tc_ta_01(TCTA01Outcome(True, False, False, False)),
            ObservabilityScore.SILENT,
        )

    def test_logged(self) -> None:
        self.assertEqual(
            classify_tc_ta_01(TCTA01Outcome(True, True, False, False)),
            ObservabilityScore.LOGGED,
        )

    def test_surfaced(self) -> None:
        self.assertEqual(
            classify_tc_ta_01(TCTA01Outcome(True, False, True, False)),
            ObservabilityScore.SURFACED,
        )


class TestClassifyTCEAC01(unittest.TestCase):
    def test_peer_unsupported(self) -> None:
        self.assertEqual(
            classify_tc_eac_01(TCEAC01Outcome(False, False, False, False, True)),
            ObservabilityScore.SURFACED,
        )

    def test_silent_protected_dg(self) -> None:
        self.assertEqual(
            classify_tc_eac_01(TCEAC01Outcome(True, False, True, False, False)),
            ObservabilityScore.SILENT,
        )

    def test_logged_no_dg(self) -> None:
        self.assertEqual(
            classify_tc_eac_01(TCEAC01Outcome(True, False, False, False, False)),
            ObservabilityScore.LOGGED,
        )

    def test_surfaced(self) -> None:
        self.assertEqual(
            classify_tc_eac_01(TCEAC01Outcome(True, False, True, True, False)),
            ObservabilityScore.SURFACED,
        )


if __name__ == "__main__":
    unittest.main()
