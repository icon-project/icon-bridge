#!/bin/env bash

# Make the script fail when a command fails
set -e

# Print all commands
set -x

chmod -R 700 /tmp/algotestnet

goal network start -r /tmp/algotestnet
goal kmd start -d /tmp/algotestnet/Node

while true; do sleep 1; done