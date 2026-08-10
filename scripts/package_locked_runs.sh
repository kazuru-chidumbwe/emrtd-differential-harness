#!/usr/bin/env bash
# Package a locked-run staging directory into a SemVer-named deposit zip.
# Hard-fails on banned-term / abs-path preflight (git-clean ≠ asset-clean).
#
# Usage:
#   make package-locked-runs STAGING=path/to/staging VERSION=1.0.6
#   # or:
#   bash scripts/package_locked_runs.sh path/to/staging 1.0.6 [outfile.zip]
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
STAGING="${1:?staging directory required}"
VERSION="${2:?SemVer without v required, e.g. 1.0.6}"
OUT="${3:-$ROOT/artifacts/emrtd-locked-runs-v${VERSION}.zip}"

mkdir -p "$(dirname "$OUT")"
python3 "$ROOT/scripts/preflight_banned_terms.py" "$STAGING"

# Also refuse host absolute harness roots that are not on the banned list.
if grep -RIl --exclude='*.png' --exclude='*.pdf' \
    -e '/opt/atlas/repos/emrtd-differential-harness/' \
    -e '/home/.*/emrtd-differential-harness/' \
    "$STAGING" >/tmp/abs-path-hits.txt 2>/dev/null; then
  if [[ -s /tmp/abs-path-hits.txt ]]; then
    echo "ABS-PATH PREFLIGHT FAILED:" >&2
    head -50 /tmp/abs-path-hits.txt >&2
    exit 1
  fi
fi

# Schema-validate every run artifact under staging
python3 - <<PY
import json, sys
from pathlib import Path
sys.path.insert(0, "$ROOT/classifier")
from schema_validate import validate_run_artifact

root = Path("$STAGING")
total = invalid = 0
for path in root.rglob("*.json"):
    try:
        obj = json.loads(path.read_text(encoding="utf-8"))
    except Exception:
        continue
    if not isinstance(obj, dict) or "observability_score" not in obj:
        continue
    total += 1
    errs = validate_run_artifact(obj)
    if errs:
        invalid += 1
        if invalid <= 10:
            print(f"SCHEMA FAIL {path.relative_to(root)}: {errs[:2]}", file=sys.stderr)
if invalid:
    print(f"SCHEMA PREFLIGHT FAILED: {invalid}/{total} run artifacts invalid", file=sys.stderr)
    raise SystemExit(1)
print(f"SCHEMA PREFLIGHT OK — {total} run artifacts")
PY

rm -f "$OUT"
(
  cd "$STAGING"
  if command -v zip >/dev/null 2>&1; then
    zip -r -q "$OUT" .
  else
    python3 - <<PY
import zipfile
from pathlib import Path
root = Path(".")
out = Path(r"$OUT")
with zipfile.ZipFile(out, "w", compression=zipfile.ZIP_DEFLATED) as z:
    for p in sorted(root.rglob("*")):
        if p.is_file():
            z.write(p, p.as_posix())
print("wrote", out)
PY
  fi
)

python3 - <<PY
import hashlib
from pathlib import Path
p = Path(r"$OUT")
print("ZIP", p.name, "SHA256", hashlib.sha256(p.read_bytes()).hexdigest(), "bytes", p.stat().st_size)
PY
echo "PACKAGE_OK $OUT"
