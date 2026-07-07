#!/usr/bin/env bash
# N-run suite for TC-AC-01 wire tier (baseline + gmrtd mitigated). Aggregates to summary JSON/MD.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

N="${SUITE_N:-100}"
PROFILE="${PROFILE:-profiles/pace-then-bac-downgrade.json}"
LOG_DIR="${LOG_DIR:-logs/suite-$(date -u +%Y%m%dT%H%M%SZ)}"
VENDOR_GMRTD="${GMRTD_PATH:-$ROOT/../_vendor/gmrtd}"
VENDOR_JMRTD="${JMRTD_PATH:-$ROOT/../_vendor/JMRTD/jmrtd}"
export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"

if [[ ! -d "$VENDOR_GMRTD" ]]; then
  echo "error: gmrtd vendor not found at $VENDOR_GMRTD (run scripts/bootstrap-vendor.sh)" >&2
  exit 2
fi
if [[ ! -f "$VENDOR_JMRTD/pom.xml" ]]; then
  echo "error: JMRTD vendor not found at $VENDOR_JMRTD" >&2
  exit 2
fi

mkdir -p "$LOG_DIR"

echo "==> TC-AC-01 suite N=$N profile=$PROFILE log_dir=$LOG_DIR"

echo "==> gmrtd baseline ($N runs)"
for _ in $(seq 1 "$N"); do
  go run ./cmd/tc-ac-01 -profile "$PROFILE" -log-dir "$LOG_DIR" -variant baseline
done

echo "==> gmrtd mitigated ($N runs)"
for _ in $(seq 1 "$N"); do
  go run ./cmd/tc-ac-01-mitigated -profile "$PROFILE" -log-dir "$LOG_DIR"
done

echo "==> jmrtd baseline ($N runs)"
if ! jar tf "$HOME/.m2/repository/org/jmrtd/jmrtd/0.5.2/jmrtd-0.5.2.jar" 2>/dev/null | grep -q 'org/jmrtd/PassportService.class'; then
  echo "    building JMRTD 0.5.2 from vendor sources..."
  bash scripts/install-jmrtd-local.sh
fi
( cd drivers/jmrtd && mvn -q -DskipTests package )
for _ in $(seq 1 "$N"); do
  java -jar drivers/jmrtd/target/jmrtd-tc-ac-01-0.1.0.jar -profile "$PROFILE" -log-dir "$LOG_DIR" -variant baseline
done

EXPECTED_RUNS=$((N * 3))
ACTUAL_RUNS=$(find "$LOG_DIR" -maxdepth 1 -name '*.json' ! -name 'summary-*' | wc -l | tr -d ' ')
if [[ "$ACTUAL_RUNS" -ne "$EXPECTED_RUNS" ]]; then
  echo "error: expected $EXPECTED_RUNS run artifacts, found $ACTUAL_RUNS" >&2
  exit 1
fi

echo "==> aggregate"
python3 classifier/aggregate.py --log-dir "$LOG_DIR"

echo ""
echo "SUITE OK — $ACTUAL_RUNS runs under $LOG_DIR/"
