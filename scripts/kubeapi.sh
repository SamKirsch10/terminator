#! /usr/bin/env bash


kubectl proxy &
PID=$!

# trap ctrl-c and call ctrl_c()
trap ctrl_c INT

function ctrl_c() {
    kill $PID
    exit 0
}

http_proxy='' curl http://localhost:8001/logs/kube-apiserver.log