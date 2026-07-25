#!/bin/bash
# Thin wrapper; prefer scripts/verify_aa_ta_eac_smoke.sh for the full gate.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
LOG_DIR="${LOG_DIR:-/tmp/ta-eac-smoke}" bash scripts/verify_aa_ta_eac_smoke.sh
echo ALL_OK
