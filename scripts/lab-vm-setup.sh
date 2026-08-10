#!/usr/bin/env bash
# Lab VM bootstrap for eMRTD harness (Ubuntu 22.04/24.04)
set -euo pipefail

HARNESS_DIR="$(cd "$(dirname "$0")/.." && pwd)"
VENDOR_DIR="${VENDOR_DIR:-$(cd "$HARNESS_DIR/.." && pwd)/_vendor}"

echo "==> System packages"
export DEBIAN_FRONTEND=noninteractive
sudo apt-get update -qq
sudo apt-get install -y -qq \
  git curl ca-certificates \
  openjdk-17-jdk maven \
  golang-go \
  python3 python3-venv python3-pip \
  build-essential unzip \
  libpcsclite-dev pcscd

echo "==> Tool versions"
java -version
mvn -version | head -1
go version
python3 --version

echo "==> Directory layout"
mkdir -p "$VENDOR_DIR" "$HARNESS_DIR/logs"

echo "==> Vendor: gmrtd pin + Option A JMRTD (Maven) + pymrtd"
# Prefer paper bootstrap (gmrtd pin + go mod). JMRTD is Maven Central 0.8.6 — not E3V3A 0.5.2.
bash "$HARNESS_DIR/scripts/bootstrap-vendor.sh"
bash "$HARNESS_DIR/scripts/install-jmrtd-local.sh" || true

if [[ ! -d "$VENDOR_DIR/pymrtd/.git" ]]; then
  git clone --depth 1 --branch v0.6.6 https://github.com/ZeroPass/pymrtd.git "$VENDOR_DIR/pymrtd"
fi

echo "==> gmrtd build (smoke)"
cd "$VENDOR_DIR/gmrtd"
go build -o /tmp/gmrtd-reader ./cmd/gmrtd-reader

echo "==> pymrtd venv"
cd "$HARNESS_DIR"
python3 -m venv .venv
# shellcheck disable=SC1091
source .venv/bin/activate
pip install -q --upgrade pip
pip install -q -e "$VENDOR_DIR/pymrtd"
python -c "import pymrtd; print('pymrtd import OK')"

echo "==> Harness classifier smoke"
python3 -c "
import sys
sys.path.insert(0, '$HARNESS_DIR/classifier')
from observability import TCAC01Outcome, classify_tc_ac_01, ObservabilityScore
assert classify_tc_ac_01(TCAC01Outcome(True, True, None, False)) == ObservabilityScore.SILENT
print('classifier smoke OK')
"

echo ""
echo "SETUP COMPLETE"
echo "  Harness: $HARNESS_DIR"
echo "  Vendors: $VENDOR_DIR/{JMRTD,gmrtd,pymrtd}"
echo "Next: synthetic profile in profiles/ + JMRTD simulator wiring (TC-AC-01)"
