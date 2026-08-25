#!/usr/bin/env python3
"""Fail if a locked-run staging tree contains banned provenance strings.

Wire this into any script that builds a Release / Zenodo zip:

  python3 scripts/preflight_banned_terms.py /path/to/staging_dir

Exit 1 on any hit. Git cleanliness does not imply Release-asset cleanliness.
"""
from __future__ import annotations

import argparse
import sys
from pathlib import Path

BANNED = (
    b"Project-Atlas",
    b"/home/boma",
    b"boma.gov",
    b"Harsh",
    b"boma",
)


def scan(root: Path) -> list[tuple[str, str, int]]:
    hits: list[tuple[str, str, int]] = []
    self_path = Path(__file__).resolve()
    for path in root.rglob("*"):
        if not path.is_file():
            continue
        # Skip this guard script itself (embeds banned literals by design).
        try:
            if path.resolve() == self_path:
                continue
        except OSError:
            pass
        if path.name == "preflight_banned_terms.py":
            continue
        raw = path.read_bytes()
        for needle in BANNED:
            n = raw.count(needle)
            if n:
                hits.append((str(path.relative_to(root)).replace("\\", "/"), needle.decode(), n))
    return hits


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("staging_dir", type=Path, help="Directory that will be zipped for deposit")
    args = ap.parse_args()
    root = args.staging_dir
    if not root.is_dir():
        print(f"ERROR: not a directory: {root}", file=sys.stderr)
        return 2
    hits = scan(root)
    if hits:
        print(f"BANNED-TERM PREFLIGHT FAILED: {len(hits)} file/needle hits under {root}", file=sys.stderr)
        for path, needle, n in hits[:50]:
            print(f"  {n:4d}× {needle!r}  in  {path}", file=sys.stderr)
        if len(hits) > 50:
            print(f"  … {len(hits) - 50} more", file=sys.stderr)
        return 1
    print(f"OK: banned-term preflight clean under {root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
