#!/bin/bash

set -e

source config.sh
source rpc.sh

# Address of new owner of ETH token on ICON
ADDR=hx80e312d8e68ee2db8fb95e6e7daae2e770d0e368

transferEthOwnership() {
    cd $CONFIG_DIR

    if [ ! -f icon.addr.bts ]; then
        echo "BTS address file icon.addr.bts does not exist"
        exit
    fi

    goloop rpc sendtx call --to $(cat icon.addr.bts) \
        --method transferOwnership \
        --param _name=btp-0x38.bsc-ETH \
        --param to=${ADDR} | jq -r . >tx/transferTokenOwnership.icon
    sleep 3
    ensure_txresult tx/transferTokenOwnership.icon
    echo "Ownership of ETH transferred to ${ADDR} " 
}

transferEthOwnership
