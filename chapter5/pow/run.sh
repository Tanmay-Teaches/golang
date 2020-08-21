#!/usr/bin/env bash
#Run client and server for POW
#run worker on port 8080 and 8081
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"


function cleanup() {
    kill $(jobs -p)
    rm -f ${DIR}/worker/worker
    rm -f ${DIR}/master/master
}
trap cleanup EXIT

cd "${DIR}/master"
go build .
./master -port 8079 -host "$(hostname)" &
master_host="http://$(hostname):8079"
echo ${master_host}
cd "${DIR}/worker"
go build .
./worker -port 8080 -master ${master_host} &
./worker -port 8081 -master ${master_host}


