#!/bin/env bash

# Make the script fail when a command fails
set -e

# Print all commands
set -x

chmod -R 700 /testnet

goal network start -r /testnet

while true; do sleep 1; done
