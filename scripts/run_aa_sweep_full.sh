#!/bin/bash
# Full AA 20-run: 5 profiles × gmrtd/jmrtd × baseline/mitigated
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
STAMP=$(date -u +%Y%m%dT%H%M%SZ)
LOG_DIR="${LOG_DIR:-logs/aa-01-sweep-full-$STAMP}"
mkdir -p "$LOG_DIR"
JAR=drivers/jmrtd/target/jmrtd-tc-ac-01-0.1.0.jar
if [[ ! -f "$JAR" ]]; then
  (cd drivers/jmrtd && mvn -q -DskipTests package)
fi
shopt -s nullglob
profiles=(profiles/aa-sweep/TC-AA-01-sweep-*.json)
if [[ ${#profiles[@]} -eq 0 ]]; then
  python3 profiles/generate_aa_sweep.py
  profiles=(profiles/aa-sweep/TC-AA-01-sweep-*.json)
fi
echo "==> AA full sweep: ${#profiles[@]} profiles × 4 arms → $LOG_DIR"
g_base=0; g_mit=0; j_base=0; j_mit=0
idx=1
for p in "${profiles[@]}"; do
  go run ./cmd/tc-aa-01 -profile "$p" -log-dir "$LOG_DIR" -variant baseline -suite-id aa-01-sweep -run-index "$idx" >/dev/null
  g_base=$((g_base+1)); idx=$((idx+1))
  go run ./cmd/tc-aa-01-mitigated -profile "$p" -log-dir "$LOG_DIR" -variant mitigated -suite-id aa-01-sweep -run-index "$idx" >/dev/null
  g_mit=$((g_mit+1)); idx=$((idx+1))
  java -cp "$JAR" org.emrtd.harness.jmrtd.TcAa01Runner -profile "$p" -log-dir "$LOG_DIR" -variant baseline -suite-id aa-01-sweep -run-index "$idx" >/dev/null
  j_base=$((j_base+1)); idx=$((idx+1))
  java -cp "$JAR" org.emrtd.harness.jmrtd.TcAa01MitigatedRunner -profile "$p" -log-dir "$LOG_DIR" -variant mitigated -suite-id aa-01-sweep -run-index "$idx" >/dev/null
  j_mit=$((j_mit+1)); idx=$((idx+1))
done
python3 - <<PY
import json, sys
from pathlib import Path
log = Path("$LOG_DIR")
files = sorted(log.glob("*.json"))
print(f"json_files={len(files)}")
counts = {}
bad = []
for f in files:
    d = json.loads(f.read_text(encoding="utf-8"))
    lib, var, score = d.get("library"), d.get("variant"), d.get("observability_score")
    counts[(lib, var, score)] = counts.get((lib, var, score), 0) + 1
    expect = 0 if var == "baseline" else 2
    if score != expect:
        bad.append((f.name, lib, var, score, expect))
print("counts:", counts)
if bad:
    print("BAD:", bad)
    sys.exit(1)
if len(files) != 20:
    print(f"expected 20 json logs, got {len(files)}")
    sys.exit(1)
print("SCORE_GRID_OK")
PY
BUNDLE=$(find "$LOG_DIR" -name '*.json' -type f | sort | xargs cat | sha256sum | awk '{print $1}')
echo "bundle_sha256=$BUNDLE"
echo "LOG_DIR=$LOG_DIR"
echo "$BUNDLE" > "$LOG_DIR/bundle_sha256.txt"
echo "AA FULL 20-RUN OK"
