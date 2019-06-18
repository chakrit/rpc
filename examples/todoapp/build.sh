#!/bin/sh

set -e

go build -o ./bin/rpc ../.. # rpc itself

./bin/rpc -gen go  -out ./goapi/api/ ./todo.rpc
./bin/rpc -gen elm -out ./web/src/ ./todo.rpc

cd goapi && go build -o ../bin/goapi . \
  && cd ..

cd web && elm make --optimize src/Main.elm \
  && cd ..

