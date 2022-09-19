#!/bin/bash

set -e

source config.sh
source rpc.sh

# Address of owner to add/remove
ADDR=hxdb23ace5d4cb14682af9fd85feb499f76edaea6b

getOwnerBTS() {
    goloop rpc call --to $(cat icon.addr.bts) \
        --method getOwners
}

getOwnerBMC() {
    goloop rpc call --to $(cat icon.addr.bmc) \
        --method getOwners
}

isOwnerBMC() {
    resp=$(
        goloop rpc call --to $(cat icon.addr.bmc) \
            --method isOwner \
            --param _addr=${1} | jq -r .
    )
    if [ $resp == "0x0" ]; then
        echo "${1} is currently not an owner"
    else
        echo "${1} is currently an owner"
    fi
}

isOwnerBTS() {
    resp=$(
        goloop rpc call --to $(cat icon.addr.bts) \
            --method isOwner \
            --param _addr=${1} | jq -r .
    )
    if [ $resp == "0x0" ]; then
        echo "${1} is currently not an owner"
    else
        echo "${1} is currently an owner"
    fi
}

addOwnerBTS() {
    cd $CONFIG_DIR

    checkBTSExists
    echo "Current Owners of BTS are: "
    getOwnerBTS

    echo "Adding owner for BTS: "

    goloop rpc sendtx call --to $(cat icon.addr.bts) \
        --method addOwner \
        --param _addr=${1} | jq -r . >tx/addJavascoreBtsOwner.icon
    sleep 3
    ensure_txresult tx/addJavascoreBtsOwner.icon
    echo "Owner added to BTS javascore!"

    echo "The new owners of BTS are: "
    getOwnerBTS
}

removeOwnerBTS() {
    cd $CONFIG_DIR

    checkBTSExists
    echo "Current Owners of BTS are: "
    getOwnerBTS

    isOwnerBTS ${1}

    echo "Removing owner for BTS: "

    goloop rpc sendtx call --to $(cat icon.addr.bts) \
        --method removeOwner \
        --param _addr=${1} | jq -r . >tx/removeJavascoreBtsOwner.icon
    sleep 3
    ensure_txresult tx/removeJavascoreBtsOwner.icon
    echo "Owner removed from BTS javascore!"

    echo "The new owners of BTS are: "
    getOwnerBTS

    isOwnerBTS ${1}
}

addOwnerBMC() {
    cd $CONFIG_DIR

    checkBMCExists
    echo "Current Owners of BMC are: "
    getOwnerBMC

    echo "Adding owner for BMC: "

    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
        --method addOwner \
        --param _addr=${1} | jq -r . >tx/addJavascoreBmcOwner.icon
    sleep 3
    ensure_txresult tx/addJavascoreBmcOwner.icon
    echo "Owner added to BMC javascore!"

    echo "The new owners of BMC are: "
    getOwnerBMC
}

checkBMCExists() {
    if [ ! -f icon.addr.bmc ]; then
        echo "BMC address file icon.addr.bmc does not exist"
        exit
    fi
}

checkBTSExists() {
    if [ ! -f icon.addr.bts ]; then
        echo "BTS address file icon.addr.bts does not exist"
        exit
    fi
}

removeOwnerBMC() {
    cd $CONFIG_DIR

    checkBMCExists

    echo "Current Owners of BMC are: "
    getOwnerBMC

    isOwnerBMC ${1}

    echo "Removing owner for BMC: "

    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
        --method removeOwner \
        --param _addr=${1} | jq -r . >tx/removeJavascoreBmcOwner.icon
    sleep 3
    ensure_txresult tx/removeJavascoreBmcOwner.icon
    echo "Owner removed from BMC javascore!"

    echo "The new owners of BMC are: "
    getOwnerBMC

    isOwnerBMC ${1}
}

if [ $# -eq 0 ]; then
    echo "No arguments supplied: Pass --help for details"
elif [ $1 == "--show-bts" ]; then
    cd $CONFIG_DIR
    checkBTSExists
    echo "Current Owners of BTS are: "
    getOwnerBTS
elif [ $1 == "--show-bmc" ]; then
    cd $CONFIG_DIR
    checkBMCExists
    echo "Current Owners of BMC are: "
    getOwnerBMC
elif [ $1 == "--add-bts" ]; then
    addOwnerBTS ${ADDR}
elif [ $1 == "--remove-bts" ]; then
    removeOwnerBTS ${ADDR}
elif [ $1 == "--add-bmc" ]; then
    addOwnerBMC ${ADDR}
elif [ $1 == "--remove-bmc" ]; then
    removeOwnerBMC ${ADDR}
else
    echo "Invalid argument: "
    echo "Valid arguments: "
    echo "--show-bmc: Show BMC Owners"
    echo "--show-bts: Show BTS Owners"
    echo "--add-bmc: Add BMC Owner"
    echo "--add-bts: Add BTS Owner"
    echo "--remove-bmc: Remove BMC Owner"
    echo "--remove-bts: Remove BTS Owner"
fi
