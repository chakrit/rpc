#!/bin/bash

# !! NOTE !! : that this one uses /bin/bash instead of /bin/sh

# SUMMARY: Test that basic client-server communication work.
# AUTHOR: Chakrit Wichian <service@chakrit.net>

set -e

RPC=$RT_RESULTS/rpc
GO_OUT=$RT_RESULTS/$RT_TEST_NAME/golang

echo [info] Generating...
$RPC -gen go -out $GO_OUT/api todo.rpc
cp -r todo/* $GO_OUT/

echo [info] Building...
cd $GO_OUT && go build -v -o ./todotest .

echo [info] Running...
cd $GO_OUT && ./todotest 2>&1 | tee client.log

echo [Info] Comparing...
diff client.log client.expected.log 1>&2

exit 0
