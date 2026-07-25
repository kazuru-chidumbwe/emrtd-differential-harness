#!/usr/bin/env bash
# AA + TA/EAC smoke (bilateral AA; JMRTD TA/EAC; gmrtd peer_unsupported).
# Requires: vendored gmrtd, compiled shade jar at drivers/jmrtd/target/jmrtd-tc-ac-01-0.1.0.jar
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

LOG_DIR="${LOG_DIR:-logs/aa-ta-eac-smoke}"
JAR="${JAR:-drivers/jmrtd/target/jmrtd-tc-ac-01-0.1.0.jar}"
mkdir -p "$LOG_DIR"

if [[ ! -f "$JAR" ]]; then
  echo "Building JMRTD driver jar..."
  (cd drivers/jmrtd && mvn -q -DskipTests package)
fi

echo "==> AA baseline + mitigated (gmrtd)"
go run ./cmd/tc-aa-01 -profile profiles/aa-internal-auth-reject.json -log-dir "$LOG_DIR" -variant baseline
go run ./cmd/tc-aa-01-mitigated -profile profiles/aa-internal-auth-reject.json -log-dir "$LOG_DIR" -variant mitigated

echo "==> AA baseline + mitigated (JMRTD)"
java -cp "$JAR" org.emrtd.harness.jmrtd.TcAa01Runner \
  -profile profiles/aa-internal-auth-reject.json -log-dir "$LOG_DIR" -variant baseline
java -cp "$JAR" org.emrtd.harness.jmrtd.TcAa01MitigatedRunner \
  -profile profiles/aa-internal-auth-reject.json -log-dir "$LOG_DIR" -variant mitigated

echo "==> TA + EAC smoke (JMRTD SW-proxy)"
java -cp "$JAR" org.emrtd.harness.jmrtd.TcTa01Runner \
  -profile profiles/ta-pso-verify-reject.json -log-dir "$LOG_DIR" -variant baseline
java -cp "$JAR" org.emrtd.harness.jmrtd.TcEac01Runner \
  -profile profiles/eac-ta-fail-dg-access.json -log-dir "$LOG_DIR" -variant baseline

echo "==> TA + EAC gmrtd peer_unsupported"
go run ./cmd/tc-ta-eac-unsupported \
  -profile profiles/ta-pso-verify-reject.json -log-dir "$LOG_DIR" -variant unsupported
go run ./cmd/tc-ta-eac-unsupported \
  -profile profiles/eac-ta-fail-dg-access.json -log-dir "$LOG_DIR" -variant unsupported

echo "==> Assert normalized_failure + scores"
python3 - <<PY
import json, glob, sys
from pathlib import Path
log = Path("$LOG_DIR")
files = list(log.glob("*.json"))
assert files, f"no JSON under {log}"
need = {
    ("TC-AA-01", "gmrtd", "baseline"): {"score": 0, "class": "chip_sw_reject", "surfaced": False},
    ("TC-AA-01", "gmrtd", "mitigated"): {"score": 2, "class": "chip_sw_reject", "surfaced": True},
    ("TC-AA-01", "jmrtd", "baseline"): {"score": 0, "class": "chip_sw_reject", "surfaced": False},
    ("TC-AA-01", "jmrtd", "mitigated"): {"score": 2, "class": "chip_sw_reject", "surfaced": True},
    ("TC-TA-01", "jmrtd", "baseline"): {"score": 0, "class": "chip_sw_reject", "surfaced": False},
    ("TC-EAC-01", "jmrtd", "baseline"): {"score": 0, "class": "chip_sw_reject", "surfaced": False},
    ("TC-TA-01", "gmrtd", "unsupported"): {"score": None, "class": "peer_unsupported", "surfaced": True},
    ("TC-EAC-01", "gmrtd", "unsupported"): {"score": None, "class": "peer_unsupported", "surfaced": True},
}
seen = set()
for f in files:
    d = json.load(open(f))
    key = (d.get("test_case"), d.get("library"), d.get("variant"))
    if key not in need:
        continue
    nf = d.get("normalized_failure")
    assert isinstance(nf, dict), f"{f}: missing normalized_failure"
    exp = need[key]
    assert nf.get("failure_class") == exp["class"], f"{key}: class {nf.get('failure_class')}"
    assert nf.get("surfaced") is exp["surfaced"], f"{key}: surfaced {nf.get('surfaced')}"
    if exp["score"] is not None:
        assert d.get("observability_score") == exp["score"], f"{key}: score {d.get('observability_score')}"
    if exp["class"] == "chip_sw_reject":
        assert nf.get("iso7816_sw"), f"{key}: missing iso7816_sw"
    seen.add(key)
missing = set(need) - seen
if missing:
    print("MISSING runs:", sorted(missing), file=sys.stderr)
    sys.exit(1)
print("normalized_failure OK for", len(seen), "runs")
PY

echo "AA/TA/EAC SMOKE OK — artifacts under $LOG_DIR/"
