"""Experiment provenance helpers (shared with internal/provenance)."""

from __future__ import annotations

import hashlib
import json
import os
import platform
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


def _git_head(root: Path) -> tuple[str, bool]:
    env = (os.environ.get("EMRTD_HARNESS_COMMIT") or "").strip()
    if env:
        return env, False
    try:
        commit = subprocess.check_output(
            ["git", "-C", str(root), "rev-parse", "HEAD"], text=True
        ).strip()
        if not commit:
            return "unknown", False
        status = subprocess.check_output(
            ["git", "-C", str(root), "status", "--porcelain"], text=True
        ).strip()
        return commit, bool(status)
    except (subprocess.CalledProcessError, FileNotFoundError):
        return "unknown", False


def file_sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def _captured_at_utc() -> str:
    raw = (os.environ.get("SOURCE_DATE_EPOCH") or "").strip()
    if raw:
        try:
            return datetime.fromtimestamp(int(raw), tz=timezone.utc).strftime(
                "%Y-%m-%dT%H:%M:%SZ"
            )
        except ValueError:
            pass
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def _gmrtd_commit(root: Path) -> str | None:
    env = (os.environ.get("GMRTD_COMMIT") or "").strip()
    if env:
        return env
    for candidate in (
        (os.environ.get("GMRTD_PATH") or "").strip(),
        str(root.parent / "_vendor" / "gmrtd"),
        str(root / "_vendor" / "gmrtd"),
    ):
        if not candidate:
            continue
        try:
            return subprocess.check_output(
                ["git", "-C", candidate, "rev-parse", "HEAD"], text=True
            ).strip()
        except (subprocess.CalledProcessError, FileNotFoundError):
            continue
    return None


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
    if run_index < 1:
        run_index = 1
    try:
        rel = profile_path.resolve().relative_to(root.resolve())
        profile_stored = rel.as_posix()
    except ValueError:
        profile_stored = profile_path.as_posix()
        marker = "/emrtd-differential-harness/"
        if marker in profile_stored:
            profile_stored = profile_stored.split(marker, 1)[1]
    return {
        "harness_commit": commit or "unknown",
        "harness_dirty": dirty,
        "gmrtd_commit": _gmrtd_commit(root),
        "python_version": platform.python_version(),
        "java_version": java_version or None,
        "profile_path": profile_stored,
        "profile_sha256": file_sha256(profile_path),
        "suite_id": suite_id,
        "suite_seed": suite_seed,
        "suite_n": suite_n,
        "run_index": run_index,
        "driver": driver,
        "variant": variant,
        "middleware": middleware or None,
        "captured_at_utc": _captured_at_utc(),
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
