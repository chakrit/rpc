#!/bin/sh

if [[ "$1" == "init" ]]; then
    go build -o $RT_RESULTS/rpc ..
fi
