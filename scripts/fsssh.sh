#!/usr/bin/env bash

SCRIPT_NAME="$(basename "$(test -L "$0" && readlink "$0" || echo "$0")")"


VM="$1"

if [ -z $VM ]; then
	echo "This script takes in a VM as an argument!"
	echo "Example: $SCRIPT_NAME bastion-oslogin-rpn9"
	echo "Example: $SCRIPT_NAME \"bastion.*play\""
	exit 1
fi

find_vm() {
	select choice in $1; do
		echo "$choice"
		break
	done
}


VM_LIST=$(cat ~/.ssh/config | grep 'Host ' | grep -Ei "$VM" | awk '{print $2}' )
TARGET=""

count="$(echo "$VM_LIST" | wc -l | tr -d ' ')"
if [[ -z "$VM_LIST" ]]; then
	# not in ssh config, so pass thru
	ssh "$@"
	exit 0
elif [[ "$count" == "1" ]]; then
	TARGET="$VM_LIST"
else
	echo "Select the VM to ssh to"
	TARGET=$(find_vm "${VM_LIST}")
fi

shift

echo ssh "${TARGET}" "$@"
ssh "${TARGET}" "$@"
