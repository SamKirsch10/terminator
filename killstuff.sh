#!/bin/bash


killorphaned() {

	# Takes in argument:
	# $1 = PID of parent shell
	# $2 = SSH_AGENT_PID

	echo $1 $2 > /tmp/samtest

	while :
	do
		# if our terminal/shell is now closed...
		ps aux | grep -v grep | grep -q "$1"
		exited=$?
		if [[ "$exited" != "0" ]]; then
			kill $2
			exit 0
		fi

		sleep 5
	done
	
}


# Let's get our shell PID so we know if it ever closes
SHELL_PID=$(echo $PPID)
killorphaned $SHELL_PID $SSH_AGENT_PID
