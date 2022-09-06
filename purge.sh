#!/bin/bash

if [[ $UID != 0 ]]; then
    echo "ERR: Please run this script as root"
    exit 1
fi

echo "------------------------------------ BEFORE -----------------------------------"
free

sync && echo 3 > /proc/sys/vm/drop_caches

echo "------------------------------------ AFTER ------------------------------------"
free