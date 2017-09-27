#!/bin/bash

DIR=./backend
REMOTEDIR=/root/go/src/github.com/twinone/iot
while :; do
	rsync -r --exclude "*node_modules*" --exclude "*BuildRoot*" . s:$REMOTEDIR
	echo "Synced"

	fswatch -1 $DIR >/dev/null 2>&1
	echo -n ".";
	sleep 1;
	echo -n ".."
done
