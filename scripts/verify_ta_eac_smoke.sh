#!/bin/bash
set -euo pipefail
JAR=drivers/jmrtd/target/jmrtd-tc-ac-01-0.1.0.jar
java -cp "$JAR" org.emrtd.harness.jmrtd.TcTa01Runner -profile profiles/ta-pso-verify-reject.json -log-dir /tmp/x -variant baseline > /tmp/ta.json 2>/tmp/ta.err
python3 - <<'PY'
import json
d=json.load(open("/tmp/ta.json"))
print("TA score", d["observability_score"], "success", d["terminal_auth_success"])
assert d["observability_score"]==0
PY
java -cp "$JAR" org.emrtd.harness.jmrtd.TcEac01Runner -profile profiles/eac-ta-fail-dg-access.json -log-dir /tmp/x -variant baseline > /tmp/eac.json 2>/tmp/eac.err
python3 - <<'PY'
import json
d=json.load(open("/tmp/eac.json"))
print("EAC score", d["observability_score"], "dg", d["protected_dg_accessible"], "eac_ok", d["eac_success"])
assert d["observability_score"]==0 and d["protected_dg_accessible"] and not d["eac_success"]
PY
echo ALL_OK
