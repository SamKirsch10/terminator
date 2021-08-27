#!/bin/bash

grep -r "$1" "$2" | awk -F':' '{print $1}' | xargs $3