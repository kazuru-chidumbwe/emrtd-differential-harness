#!/usr/bin/env bash
# Clone/pin upstream libraries for local harness runs (not vendored in git).
# JMRTD Option A (locked): Maven Central org.jmrtd:jmrtd:0.8.6 (see install-jmrtd-local.sh / docs/JMRTD-PIN.md).
# gmrtd: pinned commit (see docs/GMRTD-PIN.md) — not floating HEAD.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VENDOR_DIR="${VENDOR_DIR:-$ROOT/../_vendor}"
# shellcheck source=/dev/null
source "$ROOT/scripts/gmrtd-pin.sh"

mkdir -p "$VENDOR_DIR"

if [[ ! -d "$VENDOR_DIR/gmrtd/.git" ]]; then
  echo "==> clone gmrtd"
  git clone https://github.com/gmrtd/gmrtd.git "$VENDOR_DIR/gmrtd"
fi

echo "==> pin gmrtd to ${GMRTD_COMMIT}"
git -C "$VENDOR_DIR/gmrtd" fetch --depth 1 origin "$GMRTD_COMMIT" 2>/dev/null \
  || git -C "$VENDOR_DIR/gmrtd" fetch origin "$GMRTD_COMMIT"
git -C "$VENDOR_DIR/gmrtd" checkout --detach "$GMRTD_COMMIT"
ACTUAL="$(git -C "$VENDOR_DIR/gmrtd" rev-parse HEAD)"
if [[ "$ACTUAL" != "$GMRTD_COMMIT" ]]; then
  echo "gmrtd pin mismatch: want $GMRTD_COMMIT got $ACTUAL" >&2
  exit 1
fi
echo "  gmrtd HEAD=$(git -C "$VENDOR_DIR/gmrtd" rev-parse --short HEAD)"

echo "Vendor ready under $VENDOR_DIR"
echo "  GMRTD_PATH=$VENDOR_DIR/gmrtd"
echo "  JMRTD: run bash scripts/install-jmrtd-local.sh (Maven Central org.jmrtd:jmrtd:0.8.6)"

export GOTOOLCHAIN="${GOTOOLCHAIN:-auto}"
export GMRTD_PATH="$VENDOR_DIR/gmrtd"
cd "$ROOT"
go mod tidy
echo "go mod tidy OK"
