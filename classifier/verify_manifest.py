#!/usr/bin/env python3
"""Verify canonical artifact-manifest.json — CI gate for make paper."""

from __future__ import annotations

import argparse
import hashlib
import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from constants import ARTIFACT_MANIFEST_VERSION, METHODOLOGY_NOTE  # noqa: E402


def sha256_file(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def verify(log_dir: Path, suite_manifest: Path | None) -> list[str]:
    errors: list[str] = []
    manifest_path = log_dir / "artifact-manifest.json"
    if not manifest_path.is_file():
        return ["missing artifact-manifest.json (canonical object)"]

    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))

    if manifest.get("artifact_version") != ARTIFACT_MANIFEST_VERSION:
        errors.append(f"artifact_version must be {ARTIFACT_MANIFEST_VERSION}")

    if manifest.get("methodology_note") != METHODOLOGY_NOTE:
        errors.append("methodology_note does not match frozen constants.py text")

    if manifest.get("harness_dirty"):
        errors.append("harness_dirty=true — paper artifacts require clean checkout")

    pipeline = manifest.get("pipeline", {})
    for gate in ("tests", "smoke"):
        if pipeline.get(gate) != "pass":
            errors.append(f"pipeline gate failed or missing: {gate}={pipeline.get(gate)!r}")
    if pipeline.get("suite_complete") not in (True, "pass"):
        errors.append(f"pipeline gate failed or missing: suite_complete={pipeline.get('suite_complete')!r}")

    entries = manifest.get("entries", {})

    if suite_manifest and suite_manifest.is_file():
        spec = json.loads(suite_manifest.read_text(encoding="utf-8"))
        expected_figs = [
            e["figure_id"]
            for e in spec.get("entries", [])
            if e.get("figure_id", "").startswith("FIG-")
        ]
        for fig in expected_figs:
            if fig not in entries:
                errors.append(f"missing manifest entry for {fig}")

    for key, entry in entries.items():
        if entry.get("type") == "finding":
            for ref in entry.get("artifact_refs", []):
                p = log_dir / ref
                if not p.is_file():
                    errors.append(f"{key}: missing run artifact {ref}")
        elif entry.get("type") in ("table", "summary", "manifest"):
            p = log_dir / entry.get("path", "")
            if not p.is_file():
                errors.append(f"{key}: missing file {entry.get('path')}")
            elif entry.get("type") == "manifest":
                m2 = dict(manifest)
                e2 = dict(entries)
                e2.pop("MANIFEST-01", None)
                m2["entries"] = e2
                expected = hashlib.sha256((json.dumps(m2, indent=2) + "\n").encode()).hexdigest()
                if expected != entry.get("sha256"):
                    errors.append(f"{key}: manifest self-hash mismatch")
            elif sha256_file(p) != entry.get("sha256"):
                errors.append(f"{key}: sha256 mismatch for {entry.get('path')}")

    if "TABLE-01" not in entries:
        errors.append("missing TABLE-01 entry")
    if "MANIFEST-01" not in entries:
        errors.append("missing MANIFEST-01 self-reference")

    return errors


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("log_dir", type=Path)
    parser.add_argument("--manifest", type=Path, default=None)
    args = parser.parse_args()

    errors = verify(args.log_dir, args.manifest)
    if errors:
        for e in errors:
            print(f"VERIFY FAIL: {e}", file=sys.stderr)
        return 1
    print(f"VERIFY OK — {args.log_dir / 'artifact-manifest.json'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
