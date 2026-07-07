#!/usr/bin/env python3
"""Aggregate N-run suite logs into provenance-linked summary + artifact manifest."""

from __future__ import annotations

import argparse
import hashlib
import json
import subprocess
import sys
from collections import defaultdict
from datetime import datetime, timezone
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from observability import FINDING_THRESHOLD_PCT, ObservabilityScore, consistency_pct
from provenance import tool_versions


def sha256_file(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def load_runs(log_dir: Path) -> list[dict]:
    runs: list[dict] = []
    for path in sorted(log_dir.glob("*.json")):
        if path.name.startswith(("summary-", "artifact-", "environment")):
            continue
        with path.open(encoding="utf-8") as fh:
            runs.append(json.load(fh))
    return runs


def group_key(run: dict) -> tuple:
    return (
        run.get("test_case", ""),
        run.get("library", ""),
        run.get("variant", "baseline"),
    )


def build_artifact_manifest(log_dir: Path, runs: list[dict]) -> dict:
    files = []
    for run in runs:
        run_id = run["run_id"]
        path = log_dir / f"{run_id}.json"
        files.append(
            {
                "run_id": run_id,
                "path": path.name,
                "sha256": sha256_file(path),
                "observability_score": run.get("observability_score"),
                "figure_id": run.get("figure_id"),
                "provenance": run.get("provenance"),
            }
        )
    return {"artifact_count": len(files), "artifacts": files}


def aggregate(
    runs: list[dict],
    *,
    suite_id: str,
    suite_seed: int,
    manifest_path: Path | None,
    log_dir: Path,
) -> dict:
    groups: dict[tuple, list[dict]] = defaultdict(list)
    for run in runs:
        groups[group_key(run)].append(run)

    manifest_entries = {}
    if manifest_path and manifest_path.is_file():
        manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
        for entry in manifest.get("entries", []):
            key = (
                entry.get("test_case", ""),
                entry.get("library", ""),
                entry.get("variant", "baseline"),
            )
            manifest_entries[key] = entry

    tuples = []
    figures = []
    for key, group_runs in sorted(groups.items()):
        test_case, library, variant = key
        sample = group_runs[0]
        mechanism = sample.get("mechanism", "")
        condition = sample.get("condition", "")
        scores = [int(r["observability_score"]) for r in group_runs]
        n = len(scores)
        silent_pct = consistency_pct(scores, ObservabilityScore.SILENT)
        logged_pct = consistency_pct(scores, ObservabilityScore.LOGGED)
        surfaced_pct = consistency_pct(scores, ObservabilityScore.SURFACED)
        dominant = max(
            (silent_pct, "silent"),
            (logged_pct, "logged"),
            (surfaced_pct, "surfaced"),
            key=lambda x: x[0],
        )
        entry = manifest_entries.get(key, {})
        figure_id = entry.get("figure_id", f"fig-{test_case}-{library}-{variant}".lower())
        artifact_refs = [r["run_id"] + ".json" for r in group_runs]
        row = {
            "figure_id": figure_id,
            "test_case": test_case,
            "library": library,
            "mechanism": mechanism,
            "condition": condition,
            "variant": variant,
            "middleware": entry.get("middleware"),
            "n": n,
            "scores": scores,
            "silent_pct": round(silent_pct, 2),
            "logged_pct": round(logged_pct, 2),
            "surfaced_pct": round(surfaced_pct, 2),
            "dominant_outcome": dominant[1],
            "finding_threshold_pct": FINDING_THRESHOLD_PCT,
            "meets_finding_threshold": dominant[0] >= FINDING_THRESHOLD_PCT,
            "artifact_refs": artifact_refs,
        }
        tuples.append(row)
        figures.append(
            {
                "figure_id": figure_id,
                "tuple": {
                    "test_case": test_case,
                    "library": library,
                    "variant": variant,
                },
                "published_values": {
                    "silent_pct": row["silent_pct"],
                    "logged_pct": row["logged_pct"],
                    "surfaced_pct": row["surfaced_pct"],
                    "n": n,
                },
                "artifact_refs": artifact_refs,
            }
        )

    prov = runs[0].get("provenance", {}) if runs else {}
    return {
        "generated_at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "suite_id": suite_id,
        "suite_seed": suite_seed,
        "harness_commit": prov.get("harness_commit"),
        "profile_sha256": prov.get("profile_sha256"),
        "finding_threshold_pct": FINDING_THRESHOLD_PCT,
        "methodology_note": (
            "Fixed synthetic profile — N runs prove harness reproducibility, "
            "not input-population variance. Every published percentage links to artifact_refs."
        ),
        "environment": tool_versions(log_dir.parent.parent),
        "tuples": tuples,
        "figures": figures,
        "artifact_manifest": build_artifact_manifest(log_dir, runs),
    }


def markdown_table(summary: dict) -> str:
    lines = [
        f"<!-- suite_id={summary['suite_id']} seed={summary['suite_seed']} commit={summary.get('harness_commit')} -->",
        "| figure_id | library | mechanism | variant | N | silent % | logged % | surfaced % | ≥95% | artifact_refs |",
        "| --- | --- | --- | --- | ---: | ---: | ---: | ---: | --- | --- |",
    ]
    for row in summary["tuples"]:
        finding = "yes" if row["meets_finding_threshold"] else "no"
        refs = ", ".join(row["artifact_refs"][:3])
        if len(row["artifact_refs"]) > 3:
            refs += f", … (+{len(row['artifact_refs']) - 3})"
        lines.append(
            f"| {row['figure_id']} | {row['library']} | {row['mechanism']} | {row['variant']} "
            f"| {row['n']} | {row['silent_pct']} | {row['logged_pct']} | {row['surfaced_pct']} | {finding} | {refs} |"
        )
    return "\n".join(lines) + "\n"


def main() -> int:
    parser = argparse.ArgumentParser(description="Aggregate harness run logs")
    parser.add_argument("--log-dir", type=Path, default=Path("logs"))
    parser.add_argument("--out-prefix", type=Path, default=None)
    parser.add_argument("--manifest", type=Path, default=None)
    parser.add_argument("--suite-id", type=str, default="unknown")
    parser.add_argument("--suite-seed", type=int, default=1)
    args = parser.parse_args()

    runs = load_runs(args.log_dir)
    if not runs:
        print(f"error: no run JSON under {args.log_dir}", file=sys.stderr)
        return 1

    summary = aggregate(
        runs,
        suite_id=args.suite_id,
        suite_seed=args.suite_seed,
        manifest_path=args.manifest,
        log_dir=args.log_dir,
    )
    ts = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    prefix = args.out_prefix or (args.log_dir / f"summary-{ts}")

    json_path = prefix.with_suffix(".json") if prefix.suffix != ".json" else prefix
    md_path = json_path.with_suffix(".md")
    manifest_path = args.log_dir / f"artifact-manifest-{ts}.json"

    json_path.write_text(json.dumps(summary, indent=2) + "\n", encoding="utf-8")
    md_path.write_text(markdown_table(summary), encoding="utf-8")
    manifest_path.write_text(json.dumps(summary["artifact_manifest"], indent=2) + "\n", encoding="utf-8")

    print(json.dumps(summary, indent=2))
    print(f"\nWrote {json_path}, {md_path}, {manifest_path}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
