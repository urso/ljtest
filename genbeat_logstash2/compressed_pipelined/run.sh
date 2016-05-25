#!/usr/bin/env bash

source $(dirname $0)/../../env.sh

## start logstash test server
$LOGSTASH2 -f $DIR/logstash.conf >/dev/null 2>&1 &

# start generatorbeat sending events to test server
$GENERATORBEAT -httpprof localhost:6060 -c $DIR/genbeat.yml &

watch http://localhost:6060

sleep 5
collect

wait_shut $1
