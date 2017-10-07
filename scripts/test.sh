#!/bin/bash

set -e

function cleanup {
  rm -f main.exe
  rm -f tell-me-a-joke-plugin
}
trap cleanup EXIT

SCRIPT_DIR=`dirname $0`
cd ${SCRIPT_DIR}/..

echo "Go formatting..."
go fmt ./...

echo "Go vetting..."
go vet ./...

echo "Go imports..."
goimports .

echo "Recursive ginkgo... ${*:+(with parameter(s) }$*${*:+)}"
ginkgo -r --race --randomizeAllSpecs --failOnPending -cover $*
