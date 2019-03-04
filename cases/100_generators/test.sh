#!/bin/sh

# SUMMARY: Test generator works and produce correct and stable output.
# AUTHOR: Chakrit Wichian <service@chakrit.net>

set -e

RPC=$RT_RESULTS/rpc
GO_OUT=$RT_RESULTS/$RT_TEST_NAME/golang
RUBY_OUT=$RT_RESULTS/$RT_TEST_NAME/ruby
ELM_OUT=$RT_RESULTS/$RT_TEST_NAME/elm

echo [info] Generating...
$RPC -gen go -out $GO_OUT todo.rpc
$RPC -gen ruby -out $RUBY_OUT todo.rpc
$RPC -gen elm -out $ELM_OUT todo.rpc

echo [info] Comparing...
diff -r $GO_OUT ./golang
diff -r $RUBY_OUT ./ruby
diff -r $ELM_OUT ./elm

echo [info] Building...
cp ./go.mod $GO_OUT/go.mod
cd $GO_OUT && go build ./...
# [TODO]: elm make && rake test?

exit 0
