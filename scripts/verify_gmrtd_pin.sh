#!/usr/bin/env bash
# Fail if vendored gmrtd is not at the paper pin (docs/GMRTD-PIN.md).
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=/dev/null
source "$ROOT/scripts/gmrtd-pin.sh"
VENDOR_DIR="${VENDOR_DIR:-$ROOT/../_vendor}"
GMRTD_PATH="${GMRTD_PATH:-$VENDOR_DIR/gmrtd}"
if [[ ! -d "$GMRTD_PATH/.git" ]]; then
  echo "verify_gmrtd_pin: missing $GMRTD_PATH — run bash scripts/bootstrap-vendor.sh" >&2
  exit 1
fi
ACTUAL="$(git -C "$GMRTD_PATH" rev-parse HEAD)"
if [[ "$ACTUAL" != "$GMRTD_COMMIT" ]]; then
  echo "verify_gmrtd_pin: want $GMRTD_COMMIT got $ACTUAL" >&2
  exit 1
fi
echo "gmrtd pin OK ($GMRTD_COMMIT_SHORT)"
