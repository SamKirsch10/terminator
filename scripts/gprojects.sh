#!/bin/bash

if [[ ! -z "$1" ]]; then
	project_id=$1
    gcloud config set project $project_id
else
	echo "Choose from the following to set the GCP project"
	select project_id in $(gcloud projects list | tail -n +2 | awk '{print $1}' | sort); do
		gcloud config set project $project_id
		break
	done
fi

echo "Switching to GCP project $project_id"

# if [[ "$project_id" == *"pr-"* ]] || [[ "$project_id" == *"rc-"* ]]; then
# 	echo -e "\033]50;SetProfile=Prod\a"
# else
# 	echo -e "\033]50;SetProfile=Default\a"
# fi