#!/bin/bash

echo "waiting for ssh to come back up"
until ssh $1; do sleep 2; done
