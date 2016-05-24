#!/usr/bin/env bash

source $(dirname $0)/../env.sh

# start lumberjack test server
$TST_LJ -es -bind localhost:9200 >/dev/null 2>&1 &

# start generatorbeat sending events to test server
$GENERATORBEAT -httpprof localhost:6060 -c $DIR/genbeat.yml &

watch http://localhost:6060 &

# start metrics collection
collect &

wait_shut $1
