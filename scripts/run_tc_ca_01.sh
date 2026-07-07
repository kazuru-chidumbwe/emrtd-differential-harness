#!/usr/bin/env bash
# TC-CA-01 smoke — EAC-CA MSE failure after BAC (gmrtd wire tier).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

PROFILE="${PROFILE:-profiles/ca-v1-v2-skew.json}"
LOG_DIR="${LOG_DIR:-logs}"
VENDOR_GMRTD="${GMRTD_PATH:-$ROOT/../_vendor/gmrtd}"

if [[ ! -d "$VENDOR_GMRTD" ]]; then
  echo "error: gmrtd vendor not found at $VENDOR_GMRTD" >&2
  exit 2
fi

mkdir -p "$LOG_DIR"

echo "==> TC-CA-01 smoke (gmrtd)"
go run ./cmd/tc-ca-01 -profile "$PROFILE" -log-dir "$LOG_DIR"

echo ""
echo "TC-CA-01 SMOKE OK — trace under $LOG_DIR/"
