#!/bin/bash

if [[ ! -z "$1" ]]; then
    kubectl config use-context $1
    exit 0
fi

echo "Choose from the following to set the k8s cluster"
select cluster in $(kubectl config get-contexts | tail -n +2 | awk '{print $2}' | sort -n); do
    kubectl config use-context $cluster
    exit 0
done