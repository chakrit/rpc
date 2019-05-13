#!/bin/sh

rpc -gen go  -out ./goapi/api/ ./todo.rpc
rpc -gen elm -out ./web/src/ ./todo.rpc
