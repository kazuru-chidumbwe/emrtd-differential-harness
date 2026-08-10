"""Unit tests for frozen v1 JSON Schemas + validator wiring."""

from __future__ import annotations

import json
import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(Path(__file__).resolve().parent))

from schema_validate import (  # noqa: E402
    SchemaUnavailableError,
    validate_artifact_manifest,
    validate_profile_catalog,
    validate_run_artifact,
)


class TestFrozenSchemas(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        try:
            validate_artifact_manifest({"artifact_version": 1})
        except SchemaUnavailableError as e:
            raise unittest.SkipTest(str(e)) from e

    def test_schema_files_exist(self) -> None:
        for name in (
            "artifact-manifest-v1.json",
            "run-artifact-v1.json",
            "provenance-v1.json",
            "profile-catalog-v1.json",
        ):
            self.assertTrue((ROOT / "schemas" / name).is_file(), name)

    def test_minimal_run_valid(self) -> None:
        run = {
            "run_id": "TC-AC-01-example-gmrtd-000001",
            "test_case": "TC-AC-01",
            "library": "gmrtd",
            "mechanism": "PACE",
            "condition": "pace_fail_then_bac",
            "tier": "wire",
            "variant": "baseline",
            "observability_score": 0,
            "provenance": {
                "harness_commit": "deadbeef",
                "harness_dirty": False,
                "profile_path": "profiles/pace-then-bac-downgrade.json",
                "profile_sha256": "a" * 64,
                "suite_id": "ac-01-wire",
                "suite_seed": 1,
                "suite_n": 1,
                "run_index": 1,
                "driver": "go/tc-ac-01",
                "variant": "baseline",
                "captured_at_utc": "2026-08-10T00:00:00Z",
            },
        }
        self.assertEqual(validate_run_artifact(run), [])

    def test_run_rejects_bad_score(self) -> None:
        run = {
            "run_id": "x",
            "test_case": "TC-AC-01",
            "library": "gmrtd",
            "mechanism": "PACE",
            "condition": "pace_fail_then_bac",
            "tier": "wire",
            "variant": "baseline",
            "observability_score": 9,
            "provenance": {
                "harness_commit": "deadbeef",
                "harness_dirty": False,
                "profile_path": "profiles/x.json",
                "profile_sha256": "b" * 64,
                "suite_id": "ac-01-wire",
                "suite_seed": 1,
                "run_index": 1,
                "driver": "go/tc-ac-01",
                "variant": "baseline",
            },
        }
        self.assertTrue(validate_run_artifact(run))

    def test_minimal_manifest_valid(self) -> None:
        manifest = {
            "artifact_version": 1,
            "provenance_version": 1,
            "run_schema_version": 1,
            "profile_catalog_version": 1,
            "suite": "ac-01-wire",
            "commit": "deadbeef",
            "harness_dirty": False,
            "methodology_note": "note",
            "pipeline": {"tests": "pass"},
            "entries": {
                "FIG-01": {"type": "finding", "artifact_refs": []},
                "TABLE-01": {"type": "table", "path": "summary.md", "sha256": "c" * 64},
                "MANIFEST-01": {"type": "manifest", "path": "artifact-manifest.json", "sha256": "d" * 64},
            },
        }
        self.assertEqual(validate_artifact_manifest(manifest), [])

    def test_catalog_on_disk(self) -> None:
        catalog = json.loads((ROOT / "profiles" / "catalog.json").read_text(encoding="utf-8"))
        self.assertEqual(validate_profile_catalog(catalog), [])


if __name__ == "__main__":
    unittest.main()
