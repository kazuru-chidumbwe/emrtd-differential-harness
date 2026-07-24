#!/usr/bin/env bash
# Success-path / FP control pin — run on test-server (lab host).
# Usage (from harness root):
#   bash scripts/run-success-path-control-test-server.sh
# Or background:
#   nohup bash scripts/run-success-path-control-test-server.sh > /tmp/sp-pin.out 2>&1 &
#
# Status file for remote monitoring: /tmp/sp-pin-status.txt
# Result dir: logs/suite-ac-01-success-path-control-<UTC>-test-server/

set -euo pipefail

STATUS="${SP_PIN_STATUS:-/tmp/sp-pin-status.txt}"
HARNESS_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$HARNESS_ROOT"

stamp="$(date -u +%Y%m%dT%H%M%SZ)"
LOG_DIR="logs/suite-ac-01-success-path-control-${stamp}-test-server"
mkdir -p "$LOG_DIR"

status() {
  local msg="$1"
  printf '%s\n' "$msg" | tee -a "$STATUS"
  echo "[sp-pin] $msg"
}

: > "$STATUS"
status "START stamp=$stamp host=$(hostname) cwd=$HARNESS_ROOT"
status "PHASE go-test-middleware"

if ! command -v go >/dev/null 2>&1; then
  status "FAIL no go on PATH"
  exit 2
fi

go test ./middleware/ -count=1 2>&1 | tee "$LOG_DIR/middleware-test.txt"
status "OK middleware unit tests"

status "PHASE success-path-runs"
run_one() {
  local profile="$1"
  local allow="$2"
  local tag="$3"
  status "RUN $tag profile=$profile allow=$allow"
  go run ./cmd/tc-ac-01-mitigated \
    -profile "$profile" \
    -allow-bac-fallback="$allow" \
    -success-path \
    -suite-id ac-01-success-path-control \
    -suite-seed 1 \
    -log-dir "$LOG_DIR" \
    >"$LOG_DIR/${tag}.stdout.json" 2>"$LOG_DIR/${tag}.stderr.txt"
  status "OK $tag"
}

run_one profiles/success-path/bac-only.json false bac-only-allow-false
run_one profiles/success-path/bac-only.json true bac-only-allow-true
run_one profiles/success-path/pace-fail-allow.json true pace-fail-allow-true

status "PHASE manifest"
(
  cd "$LOG_DIR"
  find . -type f | sort | xargs sha256sum
) | tee "$LOG_DIR/file-sha256sums.txt"

# Directory fingerprint (sorted file hashes)
MANIFEST_SHA="$(sha256sum "$LOG_DIR/file-sha256sums.txt" | awk '{print $1}')"
echo "$MANIFEST_SHA" | tee "$LOG_DIR/MANIFEST.sha256"
status "MANIFEST_SHA256=$MANIFEST_SHA"
status "LOG_DIR=$HARNESS_ROOT/$LOG_DIR"

if command -v python3 >/dev/null 2>&1 && [[ -f drivers/pymrtd-offline/run_smoke.py ]]; then
  status "PHASE offline-pa-smoke (best-effort)"
  if LOG_DIR="$LOG_DIR/offline-pa" python3 drivers/pymrtd-offline/run_smoke.py \
      >"$LOG_DIR/offline-pa-smoke.stdout.txt" 2>"$LOG_DIR/offline-pa-smoke.stderr.txt"; then
    status "OK offline-pa-smoke"
  else
    status "WARN offline-pa-smoke exit=$? (pymrtd may be missing; wire pin still valid)"
  fi
else
  status "SKIP offline-pa-smoke"
fi

status "DONE OK"
exit 0
