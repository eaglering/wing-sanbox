#!/bin/sh

compiler=$1
file=$2
output=$3
addtionalArg=$4

mkdir -p /usercode
# Java
CLASSNAME=$(egrep "public\s+class\s+\w+" /sandbox/file.java | awk '{print $3}')
if [ -n "${CLASSNAME}" ]; then
	cp /sandbox/$file /usercode/$CLASSNAME.java
	file=$CLASSNAME.java
else
	cp /sandbox/$file /usercode/$file
fi

START=$(date +%s.%2N)
if [ "$output" = "" ]; then
    $compiler /usercode/$file $addtionalArg
else
    $compiler /usercode/$file $addtionalArg
	if [ $? -eq 0 ]; then
		$output
	fi
fi

END=$(date +%s.%2N)
runtime=$(echo "$END - $START" | bc)

echo -n "*-SANDBOX::ENDOFOUTPUT-*"${runtime}