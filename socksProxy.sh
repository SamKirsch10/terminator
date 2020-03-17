#!/bin/bash
# Assumes you've created a firefox profile 'SocksProfile' that is setup to use localhost:1337 socks v5.
# to do so, you can start firefox with firefox -ProfileManager to make a new profile
echo 'Starting proxy via ssh -D 1337 -N bastionNode'
ssh -D 1337 -N -f -S /tmp/socksProxy-%h-%p-%r.sock bastionNode
firefox -P "SocksProfile" &
FFpid=$!
wait $FFpid
ssh -S /tmp/socksProxy-%h-%p-%r.sock -O exit bastionNode
