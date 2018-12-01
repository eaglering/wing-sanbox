#!/bin/bash

WATCHDOG=$1
shift
FILENAME=$1
shift
DESTNAME=$1
shift

cont=$(docker create --cpus=1 -m 512M --rm "$@" 2>/dev/null)
if [ -z "$cont" ]; then
    echo "获取编译器失败"
    exit 0
fi

docker cp $FILENAME $cont:/data/$DESTNAME >/dev/null 2>&1
docker start -a $cont &
code=$(timeout -t "${WATCHDOG}" docker wait "$cont" || true)

if [ -z "$code" ]; then
    docker kill $cont >/dev/null 2>&1
    echo "脚本执行超时，允许最大执行时间${WATCHDOG}s"
    exit 0
fi
