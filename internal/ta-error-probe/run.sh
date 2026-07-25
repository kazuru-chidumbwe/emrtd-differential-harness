#!/bin/bash
set -euo pipefail
JAR=/root/.m2/repository/org/jmrtd/jmrtd/0.8.6/jmrtd-0.8.6.jar
SCUBA=/root/.m2/repository/net/sf/scuba/scuba-smartcards/0.0.20/scuba-smartcards-0.0.20.jar
BC=$(find /root/.m2/repository/org/bouncycastle/bcprov-jdk18on -name 'bcprov-jdk18on-*.jar' | head -1)
echo "=== APDULevelEACTACapable methods ==="
javap -classpath "$JAR" -public org.jmrtd.APDULevelEACTACapable
echo "=== compile+run probe ==="
cd /probe
CP="${JAR}:${SCUBA}:${BC}"
javac -cp "$CP" TaFailureProbe.java
java -cp ".:${CP}" TaFailureProbe
