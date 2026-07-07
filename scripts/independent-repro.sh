#!/usr/bin/env bash
# Independent reproduction from a clean copy (portable paths only).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPRO="${REPRO_ROOT:-/tmp/emrtd-harness-repro}"
VENDOR_SRC="${VENDOR_ROOT:-$ROOT/../_vendor}"

rm -rf "$REPRO"
mkdir -p "$REPRO/repos"
cp -a "$ROOT" "$REPRO/repos/emrtd-differential-harness"
cp -a "$VENDOR_SRC" "$REPRO/repos/_vendor"
cd "$REPRO/repos/emrtd-differential-harness"
export GMRTD_PATH="$REPRO/repos/_vendor/gmrtd"
export JMRTD_PATH="$REPRO/repos/_vendor/JMRTD/jmrtd"
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
go mod tidy

echo "=== Independent reproduction ==="
echo "KERNEL=$(uname -r)"
echo "GO=$(go version)"
echo "JAVA=$(java -version 2>&1 | head -1)"
echo "WORKDIR=$(pwd)"
echo "UTC=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
rm -rf logs
/usr/bin/time -f ELAPSED_SEC=%e make smoke 2>&1 | tee "$REPRO/repro.log"
echo "=== run_ids ==="
ls -1t logs/*.json | head -2 | while read -r f; do
  python3 -c 'import json,sys; d=json.load(open(sys.argv[1])); print(d["run_id"], d["library"], d.get("observability_score"))' "$f"
done
echo "REPRO OK — log at $REPRO/repro.log"
