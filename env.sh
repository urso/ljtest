# make sure all child processes are killed
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

# get testing home directory from finding absolute path of env.sh
TEST_HOME=$(cd "$(dirname $BASH_SOURCE)"; pwd)

# test config variables
: ${DURATION:=60}
: ${ENABLE_WATCH:=true}

# test executables
: ${GENERATORBEAT:=$GOPATH/src/github.com/urso/generatorbeat/generatorbeat}

: ${TST_LJ:=$GOPATH/src/github.com/urso/go-lumber/tst-lj}

: ${COLLECTBEAT:=$GOPATH/src/github.com/urso/collectbeat/collectbeat}

: ${EXPVAR_RATES:=$TEST_HOME/scripts/expvar_rates.py}

: ${LOGSTASH:=logstash}

watch() {
    if $ENABLE_WATCH && [ -f "$EXPVAR_RATES" ]; then
        python $EXPVAR_RATES $1/debug/vars
    fi
}

wait_shut() {
    DURATION=${1:-$DURATION}
    DURATION=${DURATION:-60}

    if [ "$DURATION" -gt 0 ]; then
        sleep $DURATION
    else
        read -s -n 1
    fi
}
