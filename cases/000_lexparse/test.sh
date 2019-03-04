#!/bin/sh

# SUMMARY: Test that basic lexing and parsing works.
# AUTHOR: Chakrit Wichian <service@chakrit.net>

set -e

RPC=$RT_RESULTS/rpc
LEX_OUT=$RT_RESULTS/$RT_TEST_NAME/lex.txt
PARSE_OUT=$RT_RESULTS/$RT_TEST_NAME/parse.json

mkdir -p $RT_RESULTS/$RT_TEST_NAME/

echo [info] Lexing...
$RPC -lex complex.rpc > $LEX_OUT
diff $LEX_OUT ./lex.txt 1>&2

echo [info] Parsing...
$RPC -parse complex.rpc > $PARSE_OUT
diff $PARSE_OUT ./parse.json 1>&2

exit 0
