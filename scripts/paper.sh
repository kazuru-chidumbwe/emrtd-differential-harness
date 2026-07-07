#!/usr/bin/env bash
# CI pipeline: tests → smoke → suite → canonical manifest → verify → publish pointer.
# Fails on any gate — paper artifacts are never generated from stale/incomplete runs.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

MANIFEST="${MANIFEST:-suites/ac-01-wire.json}"
SUITE_N="${SUITE_N:-}"
PIPELINE_STATE="$ROOT/logs/.pipeline-state.json"
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"

mkdir -p logs artifacts

write_pipeline() {
  local tests="$1" smoke="$2" suite="$3"
  python3 -c "
import json, pathlib
p = pathlib.Path('$PIPELINE_STATE')
p.write_text(json.dumps({
    'tests': '$tests',
    'smoke': '$smoke',
    'suite_complete': $suite,
}, indent=2) + '\n')
"
}

echo "==> [1/5] tests"
if make test; then
  TESTS=pass
else
  write_pipeline fail fail false
  echo "paper: ABORT — tests failed" >&2
  exit 1
fi

echo "==> [2/5] smoke"
if make smoke; then
  SMOKE=pass
else
  write_pipeline pass fail false
  echo "paper: ABORT — smoke failed" >&2
  exit 1
fi

echo "==> [3/5] suite ($MANIFEST)"
ARGS=(python3 classifier/run_suite.py --manifest "$MANIFEST")
if [[ -n "$SUITE_N" ]]; then
  ARGS+=(--n "$SUITE_N")
fi
if ! "${ARGS[@]}"; then
  write_pipeline pass pass false
  echo "paper: ABORT — suite incomplete" >&2
  exit 1
fi

SUITE_ID="$(python3 -c "import json; print(json.load(open('$MANIFEST'))['suite_id'])")"
LOG_DIR="$(ls -1dt logs/suite-${SUITE_ID}-* 2>/dev/null | head -1)"
if [[ -z "$LOG_DIR" || ! -f "$LOG_DIR/artifact-manifest.json" ]]; then
  write_pipeline pass pass false
  echo "paper: ABORT — no canonical manifest in $LOG_DIR" >&2
  exit 1
fi

write_pipeline pass pass true

echo "==> [4/5] re-aggregate with pipeline gates"
python3 classifier/aggregate.py \
  --log-dir "$LOG_DIR" \
  --manifest "$MANIFEST" \
  --suite-id "$SUITE_ID" \
  --suite-seed "$(python3 -c "import json; print(json.load(open('$MANIFEST')).get('seed',1))")" \
  --pipeline-json "$PIPELINE_STATE"

echo "==> [5/5] verify manifest"
python3 classifier/verify_manifest.py "$LOG_DIR" --manifest "$MANIFEST"

COMMIT="$(git rev-parse HEAD)"
DEST="artifacts/${SUITE_ID}-${COMMIT}-artifact-manifest.json"
cp "$LOG_DIR/artifact-manifest.json" "$DEST"
ln -sf "$(basename "$DEST")" artifacts/latest-artifact-manifest.json

echo ""
echo "PAPER OK"
echo "  canonical: $LOG_DIR/artifact-manifest.json"
echo "  published: $DEST"
echo "  commit:    $COMMIT"
