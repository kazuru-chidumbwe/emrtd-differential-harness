#!/bin/bash
# Offline PA combinatorial grid (≥24 fixtures × n=1).
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
STAMP=$(date -u +%Y%m%dT%H%M%SZ)
LOG_DIR="${LOG_DIR:-logs/pa-sweep-full-$STAMP}"
mkdir -p "$LOG_DIR"
if [[ ! -f testdata/sod/pa-sweep/index.json ]]; then
  python3 profiles/generate_pa_sweep.py
fi
shopt -s nullglob
cases=(testdata/sod/pa-sweep/TC-PA-sweep-*.json)
echo "==> PA sweep: ${#cases[@]} cases → $LOG_DIR"
export LOG_DIR
ok=0
for c in "${cases[@]}"; do
  python3 drivers/pymrtd-offline/run_case.py "$c" baseline 1 pa-01-sweep 1 1 >/dev/null
  ok=$((ok+1))
done
python3 - <<PY
import json, sys
from pathlib import Path
log = Path("$LOG_DIR")
files = sorted(log.glob("*.json"))
print(f"json_files={len(files)}")
silent = 0
surfaced = 0
control_ok = 0
bad = []
for f in files:
    d = json.loads(f.read_text(encoding="utf-8"))
    score = d.get("observability_score")
    cond = d.get("condition", "")
    is_control = "sha256_fresh" in cond
    if is_control:
        # control: verify should succeed without policy-gap scoring (0 or 1)
        if score not in (0, 1):
            bad.append((f.name, cond, score, "control"))
        else:
            control_ok += 1
    else:
        # Policy cells: score 0 = silent verify success (gap); score 2 = verify raised
        # (library already rejects). Both are valid combinatorial outcomes.
        if score == 0:
            silent += 1
        elif score == 2:
            surfaced += 1
        else:
            bad.append((f.name, cond, score, "policy"))
print(f"silent_policy_cells={silent} surfaced_reject_cells={surfaced} control_cells={control_ok}")
if bad:
    print("BAD:", bad[:20], "...", len(bad))
    sys.exit(1)
if len(files) < 24:
    print(f"expected ≥24 json logs, got {len(files)}")
    sys.exit(1)
if silent + surfaced + control_ok != len(files):
    print("internal count mismatch")
    sys.exit(1)
print("PA_SWEEP_OK")
PY
BUNDLE=$(find "$LOG_DIR" -name '*.json' -type f | sort | xargs cat | sha256sum | awk '{print $1}')
echo "bundle_sha256=$BUNDLE"
echo "$BUNDLE" > "$LOG_DIR/bundle_sha256.txt"
echo "LOG_DIR=$LOG_DIR"
echo "PA SWEEP OK ($ok cases)"
