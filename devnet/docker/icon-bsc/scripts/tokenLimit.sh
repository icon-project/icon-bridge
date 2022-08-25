#!/bin/bash

set -e

source config.sh
source rpc.sh
source utils.sh

# mainnet
# BSC_BMC_NET="0x38.bsc"
# ICON_BMC_NET="0x1.icon"

# testnet
BSC_BMC_NET="0x61.bsc"
ICON_BMC_NET="0x2.icon"

# array size should be equal
COINNAMES=(btp-$ICON_BMC_NET-ICX "btp-$ICON_BMC_NET-sICX" "btp-$ICON_BMC_NET-bnUSD" "btp-$BSC_BMC_NET-BUSD" "btp-$BSC_BMC_NET-USDC" "btp-$BSC_BMC_NET-BUSDT" )
TOKENLIMITS=(36500000000000000000000 30000000000000000000000 10000000000000000000000 10000000000000000000000 10000000000000000000000 10000000000000000000000 )

getTokenLimit() {
    goloop rpc call --to $(cat icon.addr.bts) \
        --method getTokenLimit \
        --param _name=${1}
}

changeTokenLimit() {
    cd $CONFIG_DIR

    if [ ! -f icon.addr.bts ]; then
        echo "BTS address file icon.addr.bts does not exist"
        exit
    fi

    x=(${@})
    coinNames=(${x[@]::$((${#x[@]} / 2))})
    tokenLimits=(${x[@]:$((${#x[@]} / 2))})

    echo "Getting current token limit for tokens: "
    for i in "${coinNames[@]}"; do
        echo "Current token limit For: ${i}"
        getTokenLimit ${i}
    done

    echo "Changing token limit"
    param="{\"params\":{\"_coinNames\":$(toJsonArray ${coinNames[@]}),\"_tokenLimits\":$(toJsonArray ${tokenLimits[@]})}}"

    goloop rpc sendtx call --to $(cat icon.addr.bts) \
        --method setTokenLimit \
        --raw $param | jq -r . >tx/setTokenLimit.icon
    sleep 3
    ensure_txresult tx/setTokenLimit.icon
    echo "Token Limit Changed"

    echo "Updated token limit for tokens: "
    for i in "${coinNames[@]}"; do
        echo "New token limit For: ${i}"
        getTokenLimit ${i}
    done
}

getTokenLimits() {
    cd $CONFIG_DIR
    if [ ! -f icon.addr.bts ]; then
        echo "BTS address file icon.addr.bts does not exist"
        exit
    fi

    echo "Getting current token limit for tokens: "
    coinNames=(${@})
    for i in ${coinNames[@]}; do
        echo "Token limit For: ${i}"
        getTokenLimit ${i}
    done

}

if [ $# -eq 0 ]; then
    echo "No arguments supplied: Pass --show to get limit, --update to update token limit "
elif [ $1 == "--show" ]; then
    getTokenLimits ${COINNAMES[@]}
elif [ $1 == "--update" ]; then
    changeTokenLimit ${COINNAMES[@]} ${TOKENLIMITS[@]}
else
    echo "Invalid argument: Pass --show to get limit, --update to update token limit "
fi
