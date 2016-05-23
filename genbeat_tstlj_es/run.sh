#!/usr/bin/env bash

DIR=$(dirname $0)
source $DIR/../env.sh

# start lumberjack test server
$TST_LJ -es -bind localhost:9200 >/dev/null 2>&1 &

# start generatorbeat sending events to test server
$GENERATORBEAT -httpprof localhost:6060 -c $DIR/genbeat.yml &

# start metrics collection
$COLLECTBEAT -c $DIR/collect.yml > $DIR/stats.json &


watch http://localhost:6060 &

wait_shut $1
