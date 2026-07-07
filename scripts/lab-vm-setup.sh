#!/usr/bin/env bash
# Lab VM bootstrap for R1 eMRTD harness (Ubuntu 22.04/24.04)
set -euo pipefail

WORKSPACE_ROOT="${WORKSPACE_ROOT:-$HOME/emrtd-workspace}"
VENDOR_DIR="$WORKSPACE_ROOT/repos/_vendor"
HARNESS_DIR="$WORKSPACE_ROOT/repos/emrtd-differential-harness"

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
mkdir -p "$WORKSPACE_ROOT/repos" "$VENDOR_DIR" "$HARNESS_DIR/logs"

echo "==> Vendor clones (pinned refs for G1 static audit)"
if [[ ! -d "$VENDOR_DIR/JMRTD/.git" ]]; then
  git clone --depth 1 --branch 0.5.2 https://github.com/E3V3A/JMRTD.git "$VENDOR_DIR/JMRTD"
fi
if [[ ! -d "$VENDOR_DIR/gmrtd/.git" ]]; then
  git clone --depth 1 https://github.com/gmrtd/gmrtd.git "$VENDOR_DIR/gmrtd"
fi
if [[ ! -d "$VENDOR_DIR/pymrtd/.git" ]]; then
  git clone --depth 1 --branch v0.6.6 https://github.com/ZeroPass/pymrtd.git "$VENDOR_DIR/pymrtd"
fi

echo "==> JMRTD build (smoke)"
cd "$VENDOR_DIR/JMRTD/jmrtd"
mvn -q -DskipTests package

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
from observability import RunOutcome, classify, ObservabilityScore
o = RunOutcome('gmrtd', 'PACE', 'TC-AC-01', None, True, False)
assert classify(o) == ObservabilityScore.LOGGED
print('classifier smoke OK')
"

echo ""
echo "SETUP COMPLETE"
echo "  WORKSPACE_ROOT=$WORKSPACE_ROOT"
echo "  Harness: $HARNESS_DIR"
echo "  Vendors: $VENDOR_DIR/{JMRTD,gmrtd,pymrtd}"
echo "Next: synthetic profile in profiles/ + JMRTD simulator wiring (TC-AC-01)"
