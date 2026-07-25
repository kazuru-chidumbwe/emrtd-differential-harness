#!/usr/bin/env bash
set -euo pipefail
HARNESS_DIR="$(cd "$(dirname "$0")/.." && pwd)"
VENDOR_DIR="${VENDOR_DIR:-$(cd "$HARNESS_DIR/.." && pwd)/_vendor}"
cd "$VENDOR_DIR/gmrtd"
go build -o /tmp/gmrtd-reader ./cmd/gmrtd-reader
cd "$HARNESS_DIR"
python3 -m venv .venv
.venv/bin/pip install -q --upgrade pip
.venv/bin/pip install -q -e "$VENDOR_DIR/pymrtd"
.venv/bin/python -c "import pymrtd; print('pymrtd import OK')"
.venv/bin/python -c "import sys; sys.path.insert(0,'classifier'); from observability import RunOutcome, classify, ObservabilityScore; o=RunOutcome('gmrtd','PACE','TC-AC-01',None,True,False); assert classify(o)==ObservabilityScore.LOGGED; print('classifier smoke OK')"
/tmp/gmrtd-reader -h 2>&1 | head -3 || true
echo FINISH_OK
