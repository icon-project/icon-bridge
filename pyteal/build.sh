#!/usr/bin/env bash

mkdir -p ./build/$2
# clean
rm -f ./build/$2/*.teal
rm -f ./build/$2/contract.json

set -e # die on error

python3 ./compile.py "$1" ./build/$2/approval.teal ./build/$2/clear.teal ./build/$2