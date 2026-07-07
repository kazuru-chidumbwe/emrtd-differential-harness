#!/usr/bin/env python3
"""Aggregate N-run suite logs into summary JSON + markdown table."""

from __future__ import annotations

import argparse
import json
import sys
from collections import defaultdict
from datetime import datetime, timezone
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from observability import FINDING_THRESHOLD_PCT, ObservabilityScore, consistency_pct


def load_runs(log_dir: Path) -> list[dict]:
    runs: list[dict] = []
    for path in sorted(log_dir.glob("*.json")):
        if path.name.startswith("summary-"):
            continue
        with path.open(encoding="utf-8") as fh:
            runs.append(json.load(fh))
    return runs


def group_key(run: dict) -> tuple:
    return (
        run.get("library", ""),
        run.get("mechanism", ""),
        run.get("condition", ""),
        run.get("variant", "baseline"),
    )


def aggregate(runs: list[dict]) -> dict:
    groups: dict[tuple, list[int]] = defaultdict(list)
    for run in runs:
        groups[group_key(run)].append(int(run["observability_score"]))

    tuples = []
    for key, scores in sorted(groups.items()):
        library, mechanism, condition, variant = key
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
        tuples.append(
            {
                "library": library,
                "mechanism": mechanism,
                "condition": condition,
                "variant": variant,
                "n": n,
                "scores": scores,
                "silent_pct": round(silent_pct, 2),
                "logged_pct": round(logged_pct, 2),
                "surfaced_pct": round(surfaced_pct, 2),
                "dominant_outcome": dominant[1],
                "finding_threshold_pct": FINDING_THRESHOLD_PCT,
                "meets_finding_threshold": dominant[0] >= FINDING_THRESHOLD_PCT,
            }
        )

    return {
        "generated_at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "finding_threshold_pct": FINDING_THRESHOLD_PCT,
        "note": (
            "Fixed synthetic profile — N runs prove harness reproducibility, "
            "not input variance. See README limitations."
        ),
        "tuples": tuples,
    }


def markdown_table(summary: dict) -> str:
    lines = [
        "| library | mechanism | condition | variant | N | silent % | logged % | surfaced % | ≥95% finding |",
        "| --- | --- | --- | --- | ---: | ---: | ---: | ---: | --- |",
    ]
    for row in summary["tuples"]:
        finding = "yes" if row["meets_finding_threshold"] else "no"
        lines.append(
            f"| {row['library']} | {row['mechanism']} | {row['condition']} | {row['variant']} "
            f"| {row['n']} | {row['silent_pct']} | {row['logged_pct']} | {row['surfaced_pct']} | {finding} |"
        )
    return "\n".join(lines) + "\n"


def main() -> int:
    parser = argparse.ArgumentParser(description="Aggregate harness run logs")
    parser.add_argument("--log-dir", type=Path, default=Path("logs"), help="directory of per-run JSON")
    parser.add_argument("--out-prefix", type=Path, default=None, help="summary path prefix (default: logs/summary-<ts>)")
    args = parser.parse_args()

    runs = load_runs(args.log_dir)
    if not runs:
        print(f"error: no run JSON under {args.log_dir}", file=sys.stderr)
        return 1

    summary = aggregate(runs)
    ts = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    prefix = args.out_prefix or (args.log_dir / f"summary-{ts}")

    json_path = prefix.with_suffix(".json") if prefix.suffix != ".json" else prefix
    md_path = json_path.with_suffix(".md")

    json_path.parent.mkdir(parents=True, exist_ok=True)
    json_path.write_text(json.dumps(summary, indent=2) + "\n", encoding="utf-8")
    md_path.write_text(markdown_table(summary), encoding="utf-8")

    print(json.dumps(summary, indent=2))
    print(f"\nWrote {json_path} and {md_path}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
