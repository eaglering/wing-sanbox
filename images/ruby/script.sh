#!/bin/sh

compiler=$1
file=$2
output=$3
addtionalArg=$4

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