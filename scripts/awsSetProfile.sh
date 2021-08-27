#!/bin/bash

echo "Choose from the following to set the aws profile"
select profile in $(aws configure list-profiles); do
    echo 'export AWS_DEFAULT_PROFILE="$profile"';
    echo "Setting default profile to: $profile";
    echo "Profile information:"
    #aws sts get-caller-identity
    exit 0
done
