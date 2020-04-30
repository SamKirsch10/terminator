#!/bin/bash

echo "Choose from the following to set the k8s cluster"
select cluster in $(kubectl config get-contexts | tail -n +2 | awk '{print $2}' | sort -n); do
    kubectl config use-context $cluster
    exit 0
done