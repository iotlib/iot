#!/bin/bash

DIR=./backend
REMOTEDIR=/root/go/src/github.com/twinone/iot
while :; do
	fswatch -1 $DIR >/dev/null 2>&1
	echo -n ".";
	sleep 1;
	echo -n ".."
	rsync -r . s:$REMOTEDIR
	echo "Synced"
done
