#!/bin/bash
set -e

WATCHDOG=$1
shift

FILENAME=$1
shift

DESTNAME=$1
shift

cont=$(docker create --cpus=1 -m 512M --rm "$@")
docker cp $FILENAME $cont:/data/$DESTNAME
docker start -a $cont &
code=$(timeout -t "${WATCHDOG}" docker wait "$cont" || true)

if [ -z "$code" ]; then
    docker kill $cont >/dev/null 2>&1
fi
