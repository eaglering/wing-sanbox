#!/bin/sh

compiler=$1
file=$2
output=$3
addtionalArg=$4

CLASSNAME=$(egrep "public\s+class\s+\w+" /data/file.java | awk '{print $3}')
if [ -n "${CLASSNAME}" ]; then
    mv /data/$file /data/$CLASSNAME.java
    file=$CLASSNAME.java
fi

if [ "$output" = "" ]; then
    START=$(date +%s.%N)
    $compiler /data/$file $addtionalArg
else
    $compiler /data/$file $addtionalArg
    START=$(date +%s.%N)
    if [ $? -eq 0 ]; then
        $output
    fi
fi

END=$(date +%s.%N)
runtime=$(echo "$END - $START" | bc)

echo -n "*-SANDBOX::ENDOFOUTPUT-*"${runtime}