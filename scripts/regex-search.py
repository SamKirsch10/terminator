#!/usr/bin/env python3

import re
import sys

pattern   = sys.argv[1]
inputFile = sys.argv[2]

with open(inputFile, 'r') as fh:
	content = fh.read()
	match = re.findall(pattern, content)
	if match:
		print("\n".join(match))