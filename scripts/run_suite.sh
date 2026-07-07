#!/usr/bin/env bash
# Manifest-driven suite runner (see suites/*.json).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

MANIFEST="${MANIFEST:-suites/ac-01-wire.json}"
SUITE_N="${SUITE_N:-}"
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"

ARGS=(--manifest "$MANIFEST")
if [[ -n "$SUITE_N" ]]; then
  ARGS+=(--n "$SUITE_N")
fi

python3 classifier/run_suite.py "${ARGS[@]}"
echo "SUITE OK"
