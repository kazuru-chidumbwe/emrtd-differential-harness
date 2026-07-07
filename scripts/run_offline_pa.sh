#!/usr/bin/env bash
# TC-PA-01 / TC-PA-03 offline smoke (pymrtd tier — stratified).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

LOG_DIR="${LOG_DIR:-logs}"
export LOG_DIR

echo "==> TC-PA offline smoke (pymrtd tier)"
python3 drivers/pymrtd-offline/run_smoke.py
