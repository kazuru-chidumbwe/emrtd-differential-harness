#!/usr/bin/env python3
"""Manifest-driven suite runner — produces reproducible per-run artifacts + summary."""

from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "classifier"))
from provenance import tool_versions  # noqa: E402


def load_manifest(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


def ensure_jmrtd_jar(root: Path) -> Path:
    jar = root / "drivers/jmrtd/target/jmrtd-tc-ac-01-0.1.0.jar"
    if jar.is_file():
        return jar
    subprocess.check_call(["bash", str(root / "scripts/install-jmrtd-local.sh")], cwd=root)
    subprocess.check_call(["mvn", "-q", "-DskipTests", "package"], cwd=root / "drivers/jmrtd")
    return jar


def build_go_bins(root: Path, commands: set[str]) -> dict[str, Path]:
    bindir = root / "bin"
    bindir.mkdir(exist_ok=True)
    out: dict[str, Path] = {}
    for cmd in sorted(commands):
        dest = bindir / cmd
        subprocess.check_call(["go", "build", "-o", str(dest), f"./cmd/{cmd}"], cwd=root)
        out[cmd] = dest
    return out


def common_flags(entry: dict, manifest: dict, profile: Path, i: int, n: int) -> list[str]:
    return [
        "-profile",
        str(profile),
        "-suite-id",
        manifest["suite_id"],
        "-suite-seed",
        str(manifest.get("seed", 1)),
        "-suite-n",
        str(n),
        "-run-index",
        str(i),
        "-variant",
        entry.get("variant", "baseline"),
        "-figure-id",
        entry.get("figure_id", ""),
    ]


def run_go(bin_path: Path, log_dir: Path, flags: list[str]) -> None:
    subprocess.check_call([str(bin_path), "-log-dir", str(log_dir), *flags], cwd=ROOT)


def run_java(jar: Path, main_class: str, log_dir: Path, flags: list[str]) -> None:
    subprocess.check_call(
        ["java", "-cp", str(jar), main_class, "-log-dir", str(log_dir), *flags],
        cwd=ROOT,
    )


def run_python(script: Path, fixture: Path, log_dir: Path, entry: dict, manifest: dict, i: int, n: int) -> None:
    env = os.environ.copy()
    env["LOG_DIR"] = str(log_dir)
    subprocess.check_call(
        [
            sys.executable,
            str(script),
            str(fixture),
            entry.get("variant", "baseline"),
            str(i),
            manifest["suite_id"],
            str(manifest.get("seed", 1)),
            str(n),
        ],
        cwd=ROOT,
        env=env,
    )


def main() -> int:
    parser = argparse.ArgumentParser(description="Run harness suite from JSON manifest")
    parser.add_argument("--manifest", type=Path, required=True)
    parser.add_argument("--log-dir", type=Path, default=None)
    parser.add_argument("--n", type=int, default=None, help="override manifest n")
    args = parser.parse_args()

    manifest = load_manifest(args.manifest)
    suite_id = manifest["suite_id"]
    default_n = int(manifest.get("n", 100))
    n_override = args.n

    ts = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    log_dir = args.log_dir or (ROOT / "logs" / f"suite-{suite_id}-{ts}")
    log_dir.mkdir(parents=True, exist_ok=True)

    go_cmds = {e["command"] for e in manifest["entries"] if e.get("driver") == "go"}
    go_bins = build_go_bins(ROOT, go_cmds) if go_cmds else {}
    jmrtd_jar = ensure_jmrtd_jar(ROOT) if any(e.get("driver") == "java" for e in manifest["entries"]) else None

    env_path = log_dir / "environment.json"
    env_path.write_text(json.dumps(tool_versions(ROOT), indent=2) + "\n", encoding="utf-8")

    expected = 0
    for entry in manifest["entries"]:
        n = int(entry.get("n", n_override if n_override is not None else default_n))
        profile = ROOT / entry.get("profile", manifest.get("profile", "profiles/pace-then-bac-downgrade.json"))
        driver = entry["driver"]

        for i in range(1, n + 1):
            expected += 1
            flags = common_flags(entry, manifest, profile, i, n)
            if driver == "go":
                run_go(go_bins[entry["command"]], log_dir, flags)
            elif driver == "java":
                assert jmrtd_jar is not None
                run_java(jmrtd_jar, entry["main_class"], log_dir, flags)
            elif driver == "python":
                run_python(
                    ROOT / entry["script"],
                    ROOT / entry["fixture"],
                    log_dir,
                    entry,
                    manifest,
                    i,
                    n,
                )
            else:
                raise SystemExit(f"unknown driver: {driver}")

    skip_names = {"artifact-manifest.json", "environment.json"}
    actual = len([
        p for p in log_dir.glob("*.json")
        if p.name not in skip_names and not p.name.startswith("summary-")
    ])
    if actual != expected:
        print(f"error: expected {expected} artifacts, found {actual}", file=sys.stderr)
        return 1

    agg = subprocess.run(
        [
            sys.executable,
            str(ROOT / "classifier/aggregate.py"),
            "--log-dir",
            str(log_dir),
            "--manifest",
            str(args.manifest),
            "--suite-id",
            suite_id,
            "--suite-seed",
            str(manifest.get("seed", 1)),
        ],
        cwd=ROOT,
        check=False,
    )
    if agg.returncode != 0:
        return agg.returncode

    print(f"SUITE OK — {actual} runs under {log_dir}/")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
