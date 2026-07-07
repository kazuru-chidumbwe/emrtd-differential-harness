#!/usr/bin/env bash
# R1 smoke gate — TC-AC-01 (synthetic profile, no physical passport).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

PROFILE="${PROFILE:-profiles/pace-then-bac-downgrade.json}"
LOG_DIR="${LOG_DIR:-logs}"
VENDOR_GMRTD="${GMRTD_PATH:-$ROOT/../_vendor/gmrtd}"
VENDOR_JMRTD="${JMRTD_PATH:-$ROOT/../_vendor/JMRTD/jmrtd}"

if [[ ! -d "$VENDOR_GMRTD" ]]; then
  echo "error: gmrtd vendor not found at $VENDOR_GMRTD" >&2
  exit 2
fi
if [[ ! -f "$VENDOR_JMRTD/pom.xml" ]]; then
  echo "error: JMRTD vendor not found at $VENDOR_JMRTD" >&2
  exit 2
fi

mkdir -p "$LOG_DIR"

echo "==> TC-AC-01 smoke (gmrtd)"
go run ./cmd/tc-ac-01 -profile "$PROFILE" -log-dir "$LOG_DIR"

echo ""
echo "==> TC-AC-01 smoke (jmrtd)"
if ! jar tf "$HOME/.m2/repository/org/jmrtd/jmrtd/0.5.2/jmrtd-0.5.2.jar" 2>/dev/null | grep -q 'org/jmrtd/PassportService.class'; then
  echo "    building JMRTD 0.5.2 from vendor sources..."
  bash scripts/install-jmrtd-local.sh
fi
( cd drivers/jmrtd && mvn -q -DskipTests package )
java -jar drivers/jmrtd/target/jmrtd-tc-ac-01-0.1.0.jar -profile "$PROFILE" -log-dir "$LOG_DIR"

echo ""
echo "SMOKE OK — traces written under $LOG_DIR/"
