#!/usr/bin/env bash

mkdir -p ./teal/$2
# clean
rm -f ./teal/$2/*.teal
rm -f ./teal/$2/contract.json

set -e # die on error

python3 ./compile.py "$1" ./teal/$2/approval.teal ./teal/$2/clear.teal ./teal/$2