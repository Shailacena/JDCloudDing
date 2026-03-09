#!/usr/bin/env bash
export BINPATH=$(cd `dirname $0`; pwd)
export JCQ_TEST_BASEDIR=$(dirname "$BINPATH")

# JAVA
#export JAVA_HOME="/export/share/java"
export CLASSPATH="$JCQ_TEST_BASEDIR/config:$JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar:$JCQ_TEST_BASEDIR/lib/*:."

# startup entry
export STARTUP_ENTRY="demo.JCQHttpProcessor"

JAVA_OPTS="${JAVA_OPTS} -XX:+HeapDumpOnOutOfMemoryError -XX:+DisableExplicitGC -verbose:gc -Xloggc:./jcq_gc_%p.log"

$JAVA_HOME/bin/java ${JAVA_OPTS} -classpath "$CLASSPATH" -Dfile.encoding="UTF-8" ${STARTUP_ENTRY} $*