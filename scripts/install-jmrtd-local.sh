#!/usr/bin/env bash
# Build JMRTD 0.5.2 from vendor src/ (non-standard layout) and install to ~/.m2
set -euo pipefail

JMRTD_DIR="${JMRTD_DIR:-$(cd "$(dirname "$0")/../../_vendor/JMRTD/jmrtd" && pwd)}"
VERSION=0.5.2
OUT_JAR="$JMRTD_DIR/target/jmrtd-${VERSION}-harness.jar"

if [[ ! -d "$JMRTD_DIR/src/org/jmrtd" ]]; then
  echo "JMRTD sources not found at $JMRTD_DIR" >&2
  exit 2
fi

echo "==> JMRTD local build from $JMRTD_DIR"
mkdir -p "$JMRTD_DIR/target" "$JMRTD_DIR/.harness-build"

cd "$JMRTD_DIR"
mvn -q dependency:copy-dependencies -DincludeScope=compile -DoutputDirectory=.harness-build/lib
mvn -q dependency:copy -Dartifact=org.ejbca.cvc:cert-cvc:1.4.13 -DoutputDirectory=.harness-build/lib

mapfile -t SOURCES < <(find src -name '*.java' | sort)
javac -encoding UTF-8 -cp ".harness-build/lib/*" -d .harness-build/classes "${SOURCES[@]}"

jar cf "$OUT_JAR" -C .harness-build/classes .

mvn -q install:install-file \
  -Dfile="$OUT_JAR" \
  -DgroupId=org.jmrtd \
  -DartifactId=jmrtd \
  -Dversion="$VERSION" \
  -Dpackaging=jar \
  -DpomFile=pom.xml

echo "Installed $(wc -c < "$OUT_JAR") byte jar to local Maven repo"
