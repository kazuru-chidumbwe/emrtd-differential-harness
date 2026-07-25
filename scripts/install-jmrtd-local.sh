#!/usr/bin/env bash
# Fetch JMRTD 0.8.6 from Maven Central (Option A — live maintainer line).
# Does not clone E3V3A. Requires network + Maven.
set -euo pipefail

VERSION="${JMRTD_VERSION:-0.8.6}"
# Paper / reproducibility path MUST use the default pin below.
# Override JMRTD_JAR_SHA256 only for intentional version bumps (never in CI or paper reproduction).
EXPECTED_SHA256="${JMRTD_JAR_SHA256:-5C303D7BA0DB892411E739A9920B3E0FB3C62416344CD7F220F359BDD91C0C5B}"

echo "==> Resolving org.jmrtd:jmrtd:${VERSION} from Maven Central"
mvn -q dependency:get -Dartifact="org.jmrtd:jmrtd:${VERSION}"

JAR="$HOME/.m2/repository/org/jmrtd/jmrtd/${VERSION}/jmrtd-${VERSION}.jar"
if [[ ! -f "$JAR" ]]; then
  echo "error: jar not found at $JAR" >&2
  exit 2
fi

# SHA-256 check (Linux sha256sum or macOS/Windows-compatible python)
ACTUAL="$(python3 - <<PY
import hashlib, pathlib
p = pathlib.Path(r"""$JAR""")
print(hashlib.sha256(p.read_bytes()).hexdigest().upper())
PY
)"
if [[ "$ACTUAL" != "$EXPECTED_SHA256" ]]; then
  echo "error: JMRTD ${VERSION} jar SHA-256 mismatch" >&2
  echo "  expected $EXPECTED_SHA256" >&2
  echo "  actual   $ACTUAL" >&2
  exit 3
fi

echo "Installed $JAR"
echo "SHA-256 OK ($ACTUAL)"
