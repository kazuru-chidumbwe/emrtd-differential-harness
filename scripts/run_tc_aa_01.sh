#!/usr/bin/env bash
# TC-AA-01 smoke — AA INTERNAL AUTHENTICATE failure after BAC (gmrtd).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

PROFILE="${PROFILE:-profiles/aa-internal-auth-reject.json}"
LOG_DIR="${LOG_DIR:-logs}"

mkdir -p "$LOG_DIR"

echo "==> TC-AA-01 baseline (gmrtd)"
go run ./cmd/tc-aa-01 -profile "$PROFILE" -log-dir "$LOG_DIR"

echo "==> TC-AA-01 mitigated (gmrtd)"
go run ./cmd/tc-aa-01-mitigated -profile "$PROFILE" -log-dir "$LOG_DIR"

echo ""
echo "TC-AA-01 SMOKE OK — traces under $LOG_DIR/"
