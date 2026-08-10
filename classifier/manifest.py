#!/usr/bin/env python3
"""Build the canonical artifact manifest (v1) — primary published object."""

from __future__ import annotations

import hashlib
import json
from collections import defaultdict
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

from constants import (
    ARTIFACT_MANIFEST_VERSION,
    FINDING_THRESHOLD_PCT,
    METHODOLOGY_NOTE,
    PROFILE_CATALOG_VERSION,
    PROVENANCE_VERSION,
    RUN_ARTIFACT_VERSION,
)
from observability import ObservabilityScore, consistency_pct


def sha256_file(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def sha256_text(text: str) -> str:
    return hashlib.sha256(text.encode("utf-8")).hexdigest()


def bundle_sha256(hashes: list[str]) -> str:
    return hashlib.sha256("".join(sorted(hashes)).encode("utf-8")).hexdigest()


def group_key(run: dict) -> tuple[str, str, str]:
    return (
        run.get("test_case", ""),
        run.get("library", ""),
        run.get("variant", "baseline"),
    )


def build_manifest(
    *,
    log_dir: Path,
    runs: list[dict],
    suite_id: str,
    suite_seed: int,
    suite_manifest: dict | None,
    pipeline: dict[str, Any] | None = None,
) -> dict[str, Any]:
    groups: dict[tuple[str, str, str], list[dict]] = defaultdict(list)
    for run in runs:
        groups[group_key(run)].append(run)

    entry_by_key: dict[tuple[str, str, str], dict] = {}
    if suite_manifest:
        for entry in suite_manifest.get("entries", []):
            entry_by_key[
                (
                    entry.get("test_case", ""),
                    entry.get("library", ""),
                    entry.get("variant", "baseline"),
                )
            ] = entry

    prov = runs[0].get("provenance", {}) if runs else {}
    commit = prov.get("harness_commit", "unknown")

    entries: dict[str, Any] = {}
    per_run_records: list[dict] = []

    for key, group_runs in sorted(groups.items()):
        test_case, library, variant = key
        sample = group_runs[0]
        manifest_entry = entry_by_key.get(key, {})
        fig_id = manifest_entry.get("figure_id") or sample.get("figure_id") or f"FIG-UNMAPPED-{library}-{variant}"
        scenario_id = manifest_entry.get("scenario_id", f"SCN-{test_case}-{library}-{variant}")

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

        artifact_refs: list[str] = []
        run_hashes: list[str] = []
        for run in group_runs:
            fname = f"{run['run_id']}.json"
            fpath = log_dir / fname
            fh = sha256_file(fpath)
            run_hashes.append(fh)
            artifact_refs.append(fname)
            per_run_records.append(
                {
                    "run_id": run["run_id"],
                    "path": fname,
                    "sha256": fh,
                    "figure_id": fig_id,
                    "scenario_id": scenario_id,
                    "observability_score": run.get("observability_score"),
                }
            )

        driver = sample.get("provenance", {}).get("driver", manifest_entry.get("command", ""))
        middleware = manifest_entry.get("middleware") or sample.get("provenance", {}).get("middleware")

        entries[fig_id] = {
            "type": "finding",
            "scenario": {
                "id": scenario_id,
                "test_case": test_case,
                "library": library,
                "mechanism": sample.get("mechanism", ""),
                "condition": sample.get("condition", ""),
                "variant": variant,
                "tier": sample.get("tier", ""),
                "profile_path": prov.get("profile_path"),
                "profile_sha256": prov.get("profile_sha256"),
            },
            "execution": {
                "driver": driver,
                "middleware": middleware,
                "suite_n": n,
            },
            "observation": {
                "n": n,
                "silent_pct": round(silent_pct, 2),
                "logged_pct": round(logged_pct, 2),
                "surfaced_pct": round(surfaced_pct, 2),
                "dominant_outcome": dominant[1],
            },
            "finding": {
                "threshold_pct": FINDING_THRESHOLD_PCT,
                "met": dominant[0] >= FINDING_THRESHOLD_PCT,
            },
            "artifact_refs": artifact_refs,
            "bundle_sha256": bundle_sha256(run_hashes),
        }

    return {
        "artifact_version": ARTIFACT_MANIFEST_VERSION,
        "provenance_version": PROVENANCE_VERSION,
        "run_schema_version": RUN_ARTIFACT_VERSION,
        "profile_catalog_version": PROFILE_CATALOG_VERSION,
        "suite": suite_id,
        "commit": commit,
        "harness_dirty": prov.get("harness_dirty", False),
        "suite_seed": suite_seed,
        "generated_at": _utc_now_rfc3339(),
        "methodology_note": METHODOLOGY_NOTE,
        "pipeline": pipeline or {},
        "per_run": per_run_records,
        "entries": entries,
    }


def _utc_now_rfc3339() -> str:
    """Honor SOURCE_DATE_EPOCH when set (reproducible generated_at)."""
    import os

    raw = (os.environ.get("SOURCE_DATE_EPOCH") or "").strip()
    if raw:
        try:
            return datetime.fromtimestamp(int(raw), tz=timezone.utc).strftime(
                "%Y-%m-%dT%H:%M:%SZ"
            )
        except ValueError:
            pass
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def attach_derived_artifacts(
    manifest: dict[str, Any],
    *,
    log_dir: Path,
    summary_json: Path,
    summary_md: Path,
    table_id: str = "TABLE-01",
) -> dict[str, Any]:
    """Attach TABLE/SUMMARY entries only; MANIFEST-01 is hashed in write_canonical_manifest."""
    del log_dir  # reserved for callers that pass staging dir
    manifest = dict(manifest)
    entries = dict(manifest.get("entries", {}))

    entries[table_id] = {
        "type": "table",
        "path": summary_md.name,
        "sha256": sha256_file(summary_md),
        "suite": manifest["suite"],
        "commit": manifest["commit"],
    }
    entries["SUMMARY-01"] = {
        "type": "summary",
        "path": summary_json.name,
        "sha256": sha256_file(summary_json),
        "suite": manifest["suite"],
        "commit": manifest["commit"],
    }
    manifest["entries"] = entries
    return manifest


def write_canonical_manifest(manifest: dict[str, Any], log_dir: Path) -> Path:
    """Write artifact-manifest.json; single MANIFEST-01 self-hash (no sort_keys)."""
    out = log_dir / "artifact-manifest.json"
    entries = dict(manifest.get("entries", {}))
    entries.pop("MANIFEST-01", None)
    manifest = dict(manifest)
    manifest["entries"] = entries
    body = json.dumps(manifest, indent=2) + "\n"
    manifest["entries"]["MANIFEST-01"] = {
        "type": "manifest",
        "path": "artifact-manifest.json",
        "sha256": sha256_text(body),
        "suite": manifest["suite"],
        "commit": manifest["commit"],
    }
    final = json.dumps(manifest, indent=2) + "\n"
    out.write_text(final, encoding="utf-8")
    return out


def markdown_table_from_manifest(manifest: dict[str, Any]) -> str:
    lines = [
        f"<!-- artifact_version={manifest['artifact_version']} suite={manifest['suite']} commit={manifest['commit']} -->",
        "| FIG | library | variant | N | silent % | logged % | surfaced % | ≥95% | bundle_sha256 (prefix) |",
        "| --- | --- | --- | ---: | ---: | ---: | ---: | --- | --- |",
    ]
    for fig_id, entry in sorted(manifest.get("entries", {}).items()):
        if entry.get("type") != "finding":
            continue
        obs = entry["observation"]
        finding = "yes" if entry["finding"]["met"] else "no"
        bundle = entry.get("bundle_sha256", "")[:12]
        lines.append(
            f"| {fig_id} | {entry['scenario']['library']} | {entry['scenario']['variant']} "
            f"| {obs['n']} | {obs['silent_pct']} | {obs['logged_pct']} | {obs['surfaced_pct']} "
            f"| {finding} | `{bundle}…` |"
        )
    lines.append("")
    lines.append(f"Canonical manifest: `artifact-manifest.json` · commit `{manifest['commit']}`")
    lines.append("")
    lines.append(manifest.get("methodology_note", ""))
    return "\n".join(lines) + "\n"
