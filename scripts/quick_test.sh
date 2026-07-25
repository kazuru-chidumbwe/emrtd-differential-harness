#!/usr/bin/env bash
# R1 smoke gate — TC-AC-01 (synthetic profile, no physical passport).
# Option A: JMRTD 0.8.6 from Maven Central (~/.m2), not E3V3A vendor tree.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

PROFILE="${PROFILE:-profiles/pace-then-bac-downgrade.json}"
LOG_DIR="${LOG_DIR:-logs}"
VENDOR_GMRTD="${GMRTD_PATH:-$ROOT/../_vendor/gmrtd}"
JMRTD_VERSION="${JMRTD_VERSION:-0.8.6}"
JMRTD_JAR="${HOME}/.m2/repository/org/jmrtd/jmrtd/${JMRTD_VERSION}/jmrtd-${JMRTD_VERSION}.jar"

if [[ ! -d "$VENDOR_GMRTD" ]]; then
  echo "error: gmrtd vendor not found at $VENDOR_GMRTD" >&2
  exit 2
fi

mkdir -p "$LOG_DIR"

echo "==> TC-AC-01 smoke (gmrtd)"
go run ./cmd/tc-ac-01 -profile "$PROFILE" -log-dir "$LOG_DIR"

echo ""
echo "==> TC-AC-01 smoke (jmrtd ${JMRTD_VERSION})"
if [[ ! -f "$JMRTD_JAR" ]]; then
  echo "    resolving JMRTD ${JMRTD_VERSION} from Maven Central..."
  bash scripts/install-jmrtd-local.sh
fi
( cd drivers/jmrtd && mvn -q -DskipTests package )
java -jar drivers/jmrtd/target/jmrtd-tc-ac-01-0.1.0.jar -profile "$PROFILE" -log-dir "$LOG_DIR"

echo ""
echo "SMOKE OK — traces written under $LOG_DIR/"
