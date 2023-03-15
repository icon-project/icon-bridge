#!/bin/bash

# Make the whole script fail when any of the commands fails
set -e

INITIAL_DIR=$PWD
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

for program in deploy-contract kmd-extract-private-key get-app-id init-and-register-dbsh init-and-register-wtt dbsh-call-send-service-message get-global-state-by-key deploy-asset algorand-deposit-token get-public-key-hex get-asset-holding-amount init-and-register-i2a algorand-burn-token algorand-send-asset get-algorand-address get-last-round
do
    cd $SCRIPT_DIR/$program
    go install
    echo "Installed $program"
done

cd $INITIAL_DIR