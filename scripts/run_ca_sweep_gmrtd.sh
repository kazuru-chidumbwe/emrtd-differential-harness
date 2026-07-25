#!/usr/bin/env bash
# Run gmrtd TC-CA-01 matrix (baseline + mitigated) over profiles/ca-sweep/.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

LOG_DIR="${LOG_DIR:-logs/ca-01-sweep-$(date -u +%Y%m%dT%H%M%SZ)}"
VENDOR_GMRTD="${GMRTD_PATH:-$ROOT/../_vendor/gmrtd}"
mkdir -p "$LOG_DIR"

if [[ ! -d "$VENDOR_GMRTD" ]]; then
  echo "error: gmrtd vendor not found at $VENDOR_GMRTD" >&2
  exit 2
fi

export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
export GMRTD_PATH="$VENDOR_GMRTD"

shopt -s nullglob
profiles=(profiles/ca-sweep/TC-CA-01-sweep-*.json)
if [[ ${#profiles[@]} -eq 0 ]]; then
  echo "==> generating CA sweep profiles"
  python3 profiles/generate_ca_sweep.py
  profiles=(profiles/ca-sweep/TC-CA-01-sweep-*.json)
fi

echo "==> CA sweep: ${#profiles[@]} profiles × gmrtd baseline+mitigated → $LOG_DIR"
baseline_ok=0
mitigated_ok=0
for p in "${profiles[@]}"; do
  go run ./cmd/tc-ca-01 -profile "$p" -log-dir "$LOG_DIR" -variant baseline
  baseline_ok=$((baseline_ok + 1))
  go run ./cmd/tc-ca-01-mitigated -profile "$p" -log-dir "$LOG_DIR" -variant mitigated
  mitigated_ok=$((mitigated_ok + 1))
done

echo "CA SWEEP OK — baseline=$baseline_ok mitigated=$mitigated_ok logs=$LOG_DIR"
