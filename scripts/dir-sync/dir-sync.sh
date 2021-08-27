#!/bin/bash

# This simple script is used to automatically sync changes made in
# directory 1 to directory 2. It uses the inotifywait command to 
# detect changes in dir1 and then rsync's them over to dir2
# To use this command in Ubuntu: `sudo apt install inotify-tools`
# Note: This will delete files if deleted in dir1!!

echo "Syncing $1 -> $2"
echo "Watching for changes in directory $1"

while true; do 
	inotifywait -r -e modify,create,delete,move $1
	rsync -avz --delete $1 $2
done
test   