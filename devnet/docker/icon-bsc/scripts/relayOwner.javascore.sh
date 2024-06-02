#!/bin/bash

set -e 
source config.sh # config file
source rpc.sh
source keystore.sh

function bmcExists() {
    if [ ! -f icon.addr.bmc ]; then 
        echo "icon.addr.bmc does not exist"
        exit 0
    fi
}

function addRelayIcon() {
    cd $CONFIG_DIR
    bmcExists
    getRelaysIcon
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
        --method addRelay \
        --param _link=$(cat bsc.addr.bmcbtp)\
        --param _addr=${1} | jq -r . >tx/addRelay.icon
    sleep 3
    ensure_txresult tx/addRelay.icon

    getRelaysIcon
    echo "Relay added to BMC javascore!"
}

function getRelaysIcon() {
    cd $CONFIG_DIR
    echo
    echo "Existing relays"
    goloop rpc call --to $(cat icon.addr.bmc) \
        --method getRelays \
        --param _link=$(cat bsc.addr.bmcbtp) 
    echo 
}

function removeRelayIcon() {
    cd $CONFIG_DIR

    bmcExists
    getRelaysIcon

    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
        --method removeRelay \
        --param _link=$(cat bsc.addr.bmcbtp) \
        --param _addr=${1} | jq -r . >tx/removeRelay.icon
    sleep 3
    ensure_txresult tx/addRelay.icon

    getRelaysIcon
    echo "Relay removed from BMC javascore!"
}

if [ $# -eq 0 ]; then
    echo "No arguments supplied: Pass --help or -h"
elif [ $1 == "--remove" ]; then
    echo "Removing relay " $2
    removeRelayIcon $2
elif [ $1 == "--add" ]; then
    echo "Adding relay " $2
    addRelayIcon $2
elif [ $1 == "--get" ]; then
    getRelaysIcon
else
    echo "Invalid argument "
    echo "Ensure config.sh is for relevant configuration"
    echo 
    echo "Usage:  "
    echo "      --add  addr       : Add Relayer"
    echo "      --remove  addr    : Remove Relayer"
    echo "      --get             : Get Relayers"
fi
