#!/bin/bash

lsof -Fnst | awk '
    { field = substr($0,1,1); sub(/^./,""); }
    field == "p" { pid = $0; }
    field == "t" { if ($0 == "REG") size = 0; else next; }
    field == "s" { size = $0; }
    field == "n" && size != 0 { print size, $0; }
' | sort -k1n -u | tail -n42 | sed 's/^[0-9]* //'
