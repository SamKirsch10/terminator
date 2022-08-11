#!/usr/bin/env bash

OLDIFS=$IFS
IFS=$'\n'

for project_id in $(gcloud projects list | tail -n +2 | awk '{print $1}'); do
    for cluster_info in $(gcloud container clusters list --project ${project_id} | tail -n +2); do
        cluster=$(echo $cluster_info | awk '{print $1}')
        region=$(echo $cluster_info | awk '{print $2}')
        gcloud container clusters get-credentials ${cluster} --region ${region} --project ${project_id}
    done
done

IFS=$OLDIFS

