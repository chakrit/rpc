#!/bin/sh

set -xe

drone lint --trusted
drone fmt --save
drone sign --save chakrit/rpc
drone exec --trusted
