#!/bin/sh

cd /data

for classfile in *.class; do
    classname=${classfile%.*}

    if javap -public $classname | fgrep -q 'public static void main(java.lang.String[])'; then
        java $classname "$@"
        exit 0;
    fi
done
