#!/bin/bash
set -euo pipefail

echo "######## GMRTD PROBE ########"
cd /probe/gmrtd-probe
go mod tidy
go run .

echo
echo "######## JMRTD PROBE ########"
JAR=/root/.m2/repository/org/jmrtd/jmrtd/0.8.6/jmrtd-0.8.6.jar
SCUBA=/root/.m2/repository/net/sf/scuba/scuba-smartcards/0.0.20/scuba-smartcards-0.0.20.jar
BC=$(find /root/.m2/repository/org/bouncycastle/bcprov-jdk18on -name 'bcprov-jdk18on-*.jar' 2>/dev/null | head -1)
if [ -z "$BC" ]; then
  BC=/workspace/emrtd-differential-harness/drivers/jmrtd/.deps/bcprov-jdk18on-1.80.jar
fi
echo "SCUBA=$SCUBA"
echo "BC=$BC"
cd /probe/jmrtd-probe
CP="${JAR}:${SCUBA}:${BC}"
javac -cp "$CP" AaFailureProbe.java
java -cp ".:${CP}" AaFailureProbe
echo
echo "######## JMRTD PassportService.doAA (bytecode excerpt) ########"
javap -classpath "$JAR" -c org.jmrtd.PassportService 2>/dev/null | awk '/doAA\(/,/^  public |^}/' | head -40
echo
echo "######## JMRTD AA-related types ########"
jar tf "$JAR" | grep -iE 'AAProtocol|AAResult|ActiveAuth' | head -20
