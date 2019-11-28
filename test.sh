#!/bin/sh

set -e

# go get -v -u github.com/chakrit/smoke
cd smoketests && \
    "$(go env GOPATH)/bin/smoke" "$@" ./smoketests.yml && \
    cd ..
