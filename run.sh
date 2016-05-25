#!/usr/bin/env bash

trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

export DURATION=45
export WAIT_COLLECT=10
export ENABLE_WATCH=false

while true; do
    ./genbeat_logstash/compressed/run.sh
    ./genbeat_logstash/compressed_pipelined/run.sh
    #./genbeat_logstash/uncompressed_pipelined/run.sh

    ./genbeat_tstlj_v2/compressed/run.sh
    ./genbeat_tstlj_v2/compressed_pipelined/run.sh

    ./genbeat_tstlj_es/run.sh
done
