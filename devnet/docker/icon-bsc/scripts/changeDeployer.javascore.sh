#!/bin/bash

set -e

source config.sh
source rpc.sh

NEW_BTS_OWNER=hxdb23ace5d4cb14682af9fd85feb499f76edaea6a
NEW_BMC_OWNER=hxdb23ace5d4cb14682af9fd85feb499f76edaea6a

getOwner() {
    goloop rpc call --to cx0000000000000000000000000000000000000000 \
        --method getScoreOwner \
        --param score=${1}
}

changeBtsOwner() {
    cd $CONFIG_DIR

    if [ ! -f icon.addr.bts ]; then
        echo "BTS address file icon.addr.bts does not exist"
        exit
    fi

    echo "The current contract owner of BTS is :"

    getOwner $(cat icon.addr.bts)

    echo "Changing score owner for BTS: "

    goloop rpc sendtx call --to cx0000000000000000000000000000000000000000 \
        --method setScoreOwner \
        --param score=$(cat icon.addr.bts) \
        --param owner=${1} | jq -r . >tx/changeJavascoreBtsContractOwner.icon
    sleep 3
    ensure_txresult tx/changeJavascoreBtsContractOwner.icon
    echo "Contract Owner of BTS Contract updated!"

    echo "The new contract owner of BTS is: "
    getOwner $(cat icon.addr.bts)
}

changeBmcOwner() {
    cd $CONFIG_DIR

    if [ ! -f icon.addr.bmc ]; then
        echo "BMC address file icon.addr.bmc does not exist"
        exit
    fi

    echo "The current contract owner of BMC is :"

    getOwner $(cat icon.addr.bmc)

    echo "Changing score owner for BMC: "

    goloop rpc sendtx call --to cx0000000000000000000000000000000000000000 \
        --method setScoreOwner \
        --param score=$(cat icon.addr.bmc) \
        --param owner=${1} | jq -r . >tx/changeJavascoreBmcContractOwner.icon
    sleep 3
    ensure_txresult tx/changeJavascoreBmcContractOwner.icon
    echo "Contract Owner of BMC Contract updated!"

    echo "The new contract owner of BMC is: "
    getOwner $(cat icon.addr.bmc)
}

listOwners() {
    cd $CONFIG_DIR

    if [ ! -f icon.addr.bmc ]; then
        echo "BMC address file icon.addr.bmc does not exist"
        exit
    fi

    echo "The current contract owner of BMC is :"

    getOwner $(cat icon.addr.bmc)

    if [ ! -f icon.addr.bts ]; then
        echo "BTS address file icon.addr.bmc does not exist"
        exit
    fi

    echo "The current contract owner of BTS is :"

    getOwner $(cat icon.addr.bts)

}

if [ $# -eq 0 ]; then
    echo "No arguments supplied: Pass --show to get score owners, --bts to update BTS score owner, --bmc to update BMC score owner"
elif [ $1 == "--show" ]; then
    listOwners
elif [ $1 == "--bts" ]; then
    changeBtsOwner ${NEW_BTS_OWNER}
elif [ $1 == "--bmc" ]; then
    changeBmcOwner ${NEW_BMC_OWNER}
else
    echo "Invalid argument: Pass --show to get score owners, --bts to update BTS score owner, --bmc to update BMC score owner "
fi
