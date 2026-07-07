#!/usr/bin/env bash
# Clone upstream libraries for local harness runs (not vendored in git).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VENDOR_DIR="${VENDOR_DIR:-$ROOT/../_vendor}"

mkdir -p "$VENDOR_DIR"

if [[ ! -d "$VENDOR_DIR/gmrtd/.git" ]]; then
  echo "==> clone gmrtd"
  git clone --depth 1 https://github.com/gmrtd/gmrtd.git "$VENDOR_DIR/gmrtd"
fi

if [[ ! -d "$VENDOR_DIR/JMRTD/.git" ]]; then
  echo "==> clone JMRTD 0.5.2"
  git clone --depth 1 --branch 0.5.2 https://github.com/E3V3A/JMRTD.git "$VENDOR_DIR/JMRTD"
fi

echo "Vendor ready under $VENDOR_DIR"
echo "  GMRTD_PATH=$VENDOR_DIR/gmrtd"
echo "  JMRTD_PATH=$VENDOR_DIR/JMRTD/jmrtd"

export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
export GMRTD_PATH="$VENDOR_DIR/gmrtd"
cd "$ROOT"
go mod tidy
echo "go mod tidy OK"
