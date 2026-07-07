#!/usr/bin/env python3
"""Aggregate runs into canonical artifact-manifest + derived summary views."""

from __future__ import annotations

import argparse
import json
import sys
from datetime import datetime, timezone
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from manifest import (  # noqa: E402
    attach_derived_artifacts,
    build_manifest,
    markdown_table_from_manifest,
    write_canonical_manifest,
)
from provenance import tool_versions


def load_runs(log_dir: Path) -> list[dict]:
    skip_prefixes = ("summary-", "environment")
    skip_names = {"artifact-manifest.json"}
    runs: list[dict] = []
    for path in sorted(log_dir.glob("*.json")):
        if path.name in skip_names or any(path.name.startswith(p) for p in skip_prefixes):
            continue
        with path.open(encoding="utf-8") as fh:
            runs.append(json.load(fh))
    return runs


def main() -> int:
    parser = argparse.ArgumentParser(description="Build canonical artifact manifest from run logs")
    parser.add_argument("--log-dir", type=Path, required=True)
    parser.add_argument("--manifest", type=Path, default=None, help="suite manifest JSON")
    parser.add_argument("--suite-id", type=str, default="unknown")
    parser.add_argument("--suite-seed", type=int, default=1)
    parser.add_argument("--pipeline-json", type=Path, default=None, help="pipeline gate results")
    args = parser.parse_args()

    runs = load_runs(args.log_dir)
    if not runs:
        print(f"error: no run JSON under {args.log_dir}", file=sys.stderr)
        return 1

    suite_manifest = None
    if args.manifest and args.manifest.is_file():
        suite_manifest = json.loads(args.manifest.read_text(encoding="utf-8"))

    pipeline = {}
    if args.pipeline_json and args.pipeline_json.is_file():
        pipeline = json.loads(args.pipeline_json.read_text(encoding="utf-8"))

    manifest = build_manifest(
        log_dir=args.log_dir,
        runs=runs,
        suite_id=args.suite_id,
        suite_seed=args.suite_seed,
        suite_manifest=suite_manifest,
        pipeline=pipeline,
    )
    manifest["environment"] = tool_versions(args.log_dir.parent.parent)

    ts = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    summary_json = args.log_dir / f"summary-{ts}.json"
    summary_md = args.log_dir / f"summary-{ts}.md"

    # Derived views (not canonical)
    summary_json.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")
    summary_md.write_text(markdown_table_from_manifest(manifest), encoding="utf-8")

    manifest = attach_derived_artifacts(
        manifest,
        log_dir=args.log_dir,
        summary_json=summary_json,
        summary_md=summary_md,
        table_id="TABLE-01",
    )
    canonical = write_canonical_manifest(manifest, args.log_dir)

    print(json.dumps(manifest, indent=2))
    print(f"\nCanonical manifest: {canonical}", file=sys.stderr)
    print(f"Derived: {summary_json}, {summary_md}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
