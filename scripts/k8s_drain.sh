#! /usr/bin/env bash


if [[ -z $1 ]]; then
    echo "No argument supplied! Argument = node to drain"
    exit 1
fi

echo "kubectl drain $1 --force=true --grace-period=0  --ignore-daemonsets --disable-eviction=true"
kubectl drain $1 --force=true --grace-period=0 --ignore-daemonsets --disable-eviction=true
