#!/usr/bin/env bash
# Clone upstream libraries for local harness runs (not vendored in git).
# Option A (locked): gmrtd from GitHub; JMRTD from Maven Central 0.8.6 (see install-jmrtd-local.sh).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VENDOR_DIR="${VENDOR_DIR:-$ROOT/../_vendor}"

mkdir -p "$VENDOR_DIR"

if [[ ! -d "$VENDOR_DIR/gmrtd/.git" ]]; then
  echo "==> clone gmrtd"
  git clone --depth 1 https://github.com/gmrtd/gmrtd.git "$VENDOR_DIR/gmrtd"
fi

echo "Vendor ready under $VENDOR_DIR"
echo "  GMRTD_PATH=$VENDOR_DIR/gmrtd"
echo "  JMRTD: run bash scripts/install-jmrtd-local.sh (Maven Central org.jmrtd:jmrtd:0.8.6)"

export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
export GMRTD_PATH="$VENDOR_DIR/gmrtd"
cd "$ROOT"
go mod tidy
echo "go mod tidy OK"
