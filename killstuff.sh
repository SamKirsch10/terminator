#!/bin/bash


killorphaned() {
	if [[ ! -z $SSH_AGENT_PID ]]; then
		kill $SSH_AGENT_PID
	fi
	## Can add more background shit later
}


# Let's get our shell PID so we know if it ever closes
SHELL_PID=$(echo $PPID)

while :
do
	# if our terminal/shell is now closed...
	exited=$(ps aux | grep $SHELL_PID | grep -vq grep)
	if [[ ! -z $exited ]]; then
		killorphaned
		break
	fi

	sleep 10
done
