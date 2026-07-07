#!/usr/bin/env bash
set -euo pipefail
cd $HOME/emrtd-workspace/repos/_vendor/gmrtd
go build -o /tmp/gmrtd-reader ./cmd/gmrtd-reader
cd $HOME/emrtd-workspace/repos/emrtd-differential-harness
python3 -m venv .venv
.venv/bin/pip install -q --upgrade pip
.venv/bin/pip install -q -e $HOME/emrtd-workspace/repos/_vendor/pymrtd
.venv/bin/python -c "import pymrtd; print('pymrtd import OK')"
.venv/bin/python -c "import sys; sys.path.insert(0,'classifier'); from observability import RunOutcome, classify, ObservabilityScore; o=RunOutcome('gmrtd','PACE','TC-AC-01',None,True,False); assert classify(o)==ObservabilityScore.LOGGED; print('classifier smoke OK')"
/tmp/gmrtd-reader -h 2>&1 | head -3 || true
echo FINISH_OK
