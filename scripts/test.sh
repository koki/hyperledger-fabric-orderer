#!/usr/bin/env bash

set -ax
set -e

SCRIPTS_DIR=$(dirname $0)

cd $SCRIPTS_DIR/..

go test ./...
