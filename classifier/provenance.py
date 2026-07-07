"""Experiment provenance helpers (shared with internal/provenance)."""

from __future__ import annotations

import hashlib
import json
import platform
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


def _git_head(root: Path) -> tuple[str, bool]:
    try:
        commit = subprocess.check_output(
            ["git", "-C", str(root), "rev-parse", "HEAD"], text=True
        ).strip()
        status = subprocess.check_output(
            ["git", "-C", str(root), "status", "--porcelain"], text=True
        ).strip()
        return commit, bool(status)
    except (subprocess.CalledProcessError, FileNotFoundError):
        return "unknown", False


def file_sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def collect(
    *,
    root: Path,
    profile_path: Path,
    suite_id: str,
    suite_seed: int,
    suite_n: int,
    run_index: int,
    driver: str,
    variant: str,
    middleware: str = "",
    java_version: str = "",
) -> dict[str, Any]:
    commit, dirty = _git_head(root)
    return {
        "harness_commit": commit,
        "harness_dirty": dirty,
        "python_version": platform.python_version(),
        "java_version": java_version or None,
        "profile_path": profile_path.as_posix(),
        "profile_sha256": file_sha256(profile_path),
        "suite_id": suite_id,
        "suite_seed": suite_seed,
        "suite_n": suite_n,
        "run_index": run_index,
        "driver": driver,
        "variant": variant,
        "middleware": middleware or None,
        "captured_at_utc": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
    }


def tool_versions(root: Path) -> dict[str, str]:
    versions: dict[str, str] = {"python": platform.python_version()}
    for name, cmd in (
        ("go", ["go", "version"]),
        ("java", ["java", "-version"]),
    ):
        try:
            proc = subprocess.run(cmd, capture_output=True, text=True, check=False)
            versions[name] = (proc.stdout or proc.stderr).strip().splitlines()[0]
        except FileNotFoundError:
            versions[name] = "not found"
    return versions
