#!/bin/bash
# Full CA 40-run: 10 profiles × gmrtd/jmrtd × baseline/mitigated
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
STAMP=$(date -u +%Y%m%dT%H%M%SZ)
LOG_DIR="${LOG_DIR:-logs/ca-01-sweep-full-$STAMP}"
mkdir -p "$LOG_DIR"
JAR=drivers/jmrtd/target/jmrtd-tc-ac-01-0.1.0.jar
if [[ ! -f "$JAR" ]]; then
  (cd drivers/jmrtd && mvn -q -DskipTests package)
fi
shopt -s nullglob
profiles=(profiles/ca-sweep/TC-CA-01-sweep-*.json)
if [[ ${#profiles[@]} -eq 0 ]]; then
  python3 profiles/generate_ca_sweep.py
  profiles=(profiles/ca-sweep/TC-CA-01-sweep-*.json)
fi
echo "==> CA full sweep: ${#profiles[@]} profiles × 4 arms → $LOG_DIR"
g_base=0; g_mit=0; j_base=0; j_mit=0
for p in "${profiles[@]}"; do
  go run ./cmd/tc-ca-01 -profile "$p" -log-dir "$LOG_DIR" -variant baseline -suite-id ca-01-sweep >/dev/null
  g_base=$((g_base+1))
  go run ./cmd/tc-ca-01-mitigated -profile "$p" -log-dir "$LOG_DIR" -variant mitigated -suite-id ca-01-sweep >/dev/null
  g_mit=$((g_mit+1))
  java -cp "$JAR" org.emrtd.harness.jmrtd.TcCa01Runner -profile "$p" -log-dir "$LOG_DIR" -variant baseline -suite-id ca-01-sweep >/dev/null
  j_base=$((j_base+1))
  java -cp "$JAR" org.emrtd.harness.jmrtd.TcCa01MitigatedRunner -profile "$p" -log-dir "$LOG_DIR" -variant mitigated -suite-id ca-01-sweep >/dev/null
  j_mit=$((j_mit+1))
done
# Score check
python3 - <<PY
import json, glob, sys
from pathlib import Path
log = Path("$LOG_DIR")
files = sorted(log.glob("*.json"))
print(f"json_files={len(files)}")
counts = {}
bad = []
for f in files:
    d = json.loads(f.read_text(encoding="utf-8"))
    lib = d.get("library")
    var = d.get("variant")
    score = d.get("observability_score")
    key = (lib, var, score)
    counts[key] = counts.get(key, 0) + 1
    expect = 0 if var == "baseline" else 2
    if score != expect:
        bad.append((f.name, lib, var, score, expect))
print("counts:", counts)
if bad:
    print("BAD:", bad)
    sys.exit(1)
if len(files) != 40:
    print(f"expected 40 json logs, got {len(files)}")
    sys.exit(1)
print("SCORE_GRID_OK")
PY
# Bundle hash over sorted JSON contents
BUNDLE=$(find "$LOG_DIR" -name '*.json' -type f | sort | xargs cat | sha256sum | awk '{print $1}')
echo "bundle_sha256=$BUNDLE"
echo "LOG_DIR=$LOG_DIR"
echo "gmrtd_baseline=$g_base gmrtd_mitigated=$g_mit jmrtd_baseline=$j_base jmrtd_mitigated=$j_mit"
echo "$BUNDLE" > "$LOG_DIR/bundle_sha256.txt"
echo "CA FULL 40-RUN OK"
