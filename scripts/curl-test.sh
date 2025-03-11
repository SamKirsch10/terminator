#!/usr/bin/env bash -e


TMP=$(mktemp -d)
cd $TMP

# no matter what, make sure the tmp dir is cleaned up
trap "rm -rf $TMP" EXIT


cat <<EOF > curl-format.txt
     time_namelookup:  %{time_namelookup}s\n
        time_connect:  %{time_connect}s\n
     time_appconnect:  %{time_appconnect}s\n
    time_pretransfer:  %{time_pretransfer}s\n
       time_redirect:  %{time_redirect}s\n
  time_starttransfer:  %{time_starttransfer}s\n
                     ----------\n
          time_total:  %{time_total}s\n
EOF

if [[ -z "$@" ]]; then
    echo "No args detected!"
    echo "Pass any args you would to curl to this script!"
    exit 1
fi

curl -w "@curl-format.txt" $@

