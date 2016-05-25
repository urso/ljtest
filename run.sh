#!/usr/bin/env bash

trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

export DURATION=$((60 * 2))
#export DURATION=20
export WAIT_COLLECT=20
export ENABLE_WATCH=false

while true; do
    #./genbeat_logstash/compressed/run.sh
    #./genbeat_logstash/compressed_pipelined/run.sh
    #./genbeat_logstash/uncompressed_pipelined/run.sh

    ./genbeat_logstash5/compressed_pipelined/run.sh; sleep 5
    ./genbeat_logstash5/compressed/run.sh; sleep 5
    ./genbeat_logstash2/compressed/run.sh; sleep 5

    #./genbeat_jlumber_v2/compressed/run.sh
    #./genbeat_tstlj_v2/compressed/run.sh
    #./genbeat_tstlj_v2/compressed_pipelined/run.sh
    #./genbeat_tstlj_es/run.sh
    #./genbeat_jlumber_v2/compressed_pipelined/run.sh
done
