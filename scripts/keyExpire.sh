#!/bin/bash

now=$(date +%s)
validTill=$(date -d "$(ssh-keygen -L -h -f ~/.ssh/*-cert.pub | grep Valid | sed 's|.*to.\(.*\)|\1|')" +%s)
keyName=$(ssh-keygen -L -h -f ~/.ssh/*-cert.pub | grep 'Key ID' | xargs | awk '{print $NF}')

diff=$(( ((10#$validTill - 10#$now))  / 86400))

if [ $diff -lt 5 ] ; then
	error="Your ssh key [$keyName] is going to expire in $diff days! Get a new one!!"
	echo -e "\e[01;31m============================================\e[0m" >&2
	echo -e "\e[01;31m$error\e[0m" >&2
	echo -e "\e[01;31m============================================\e[0m" >&2

fi

