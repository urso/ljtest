# make sure all child processes are killed
trap "pkill -9 -P $$" SIGINT SIGTERM EXIT

# get testing home directory from finding absolute path of env.sh
DIR=$(dirname $0)
TEST_HOME=$(cd "$(dirname $BASH_SOURCE)"; pwd)

# test config variables
: ${DURATION:=60}
: ${ENABLE_WATCH:=true}

: ${WORKDIR:=$GOPATH/src}

# test executables
: ${GENERATORBEAT:=$WORKDIR/github.com/urso/generatorbeat/generatorbeat}

: ${JAVA_LUMBER:=$WORKDIR/github.com/ph/java-lumber}

: ${TST_LJ:=$WORKDIR/github.com/urso/go-lumber/tst-lj}

: ${COLLECTBEAT:=$WORKDIR/github.com/urso/collectbeat/collectbeat}

: ${EXPVAR_RATES:=$TEST_HOME/scripts/expvar_rates.py}

: ${LOGSTASH:=logstash}

: ${LOGSTASH2=$WORKDIR/logstash-2.3.2/bin/logstash}

: ${LOGSTASH5=$WORKDIR/logstash-5.0.0-alpha2/bin/logstash}

: ${OUTDIR:=$DIR}

: ${OUTCONFIG:=$TEST_HOME/collect_es.yml}

: ${WAIT_COLLECT:=0}

# run script utilities

gradle_run() {
    cd $1
    gradle run &
    cd -
}

collect() {
    sleep $WAIT_COLLECT
    if [ -f "$OUTCONFIG" ]; then
        $COLLECTBEAT -c <(cat "$DIR/collect.yml" "$OUTCONFIG") &
    else
        $COLLECTBEAT -c <(cat $DIR/collect.yml <(echo "output.console:")) > $OUTDIR/stats.json &
    fi
}

watch() {
    if $ENABLE_WATCH && [ -f "$EXPVAR_RATES" ]; then
        python $EXPVAR_RATES $1/debug/vars &
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
