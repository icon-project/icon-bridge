#!/bin/bash
set -e
source config.sh
source rpc.sh

export PRIVATE_KEY="[\""$(cat $BSC_KEY_STORE.priv)"\"]"

function bmcExists() {
    if [ ! -f bsc.addr.bmcmanagement ]; then
        echo "bsc.addr.bmcmanagement does not exist"
        exit 0
    fi
}

function getRelaysSolidity() {
    cd $CONFIG_DIR
    bmcExists

    cd $CONTRACTS_DIR/solidity/bmc
    truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
        --method "getRelays" --link "$(cat $CONFIG_DIR/icon.addr.bmcbtp)"
}

function removeRelaySolidity() {
    cd $CONFIG_DIR

    getRelaysSolidity
    echo
    cd $CONTRACTS_DIR/solidity/bmc
    tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
        --method removeRelay --link "$(cat $CONFIG_DIR/icon.addr.bmcbtp)" --addr "${1}")
    echo "$tx" >$CONFIG_DIR/tx/removeRelay.bmc.bsc
    echo
    getRelaysSolidity
}

function addRelaySolidity() {
    # USAGE 
    # For single relay
    # addRelaySolidity 0x4300148436d51f7f270cb6e76cbc82fa0ce1b359
    # For multiple relays
    # addRelaySolidity 0x4300148436d51f7f270cb6e76cbc82fa0ce1b359,0x70e789d2f5d469ea30e0525dbfdd5515d6ead30d
    #
    cd $CONFIG_DIR

    getRelaysSolidity
    echo
    echo

    echo "Adding relay to BMC Management"
    echo "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
    echo -e "\033[0;31m DANGER:: This method removes all existing previous relays, and adds $1 as new relay addresses \033[0m"
    echo "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
    sleep 1
    read -p "Are you sure you want to proceed [y/N]: " proceed
    if [[ $proceed == "y" ]]; then

        cd $CONTRACTS_DIR/solidity/bmc
        tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
            --method addRelay --link "$(cat $CONFIG_DIR/icon.addr.bmcbtp)" --addr "${1}")
        echo "$tx" >$CONFIG_DIR/tx/addRelay.bmc.bsc
        getRelaysSolidity
    fi

}

if [ $# -eq 0 ]; then
    echo "No arguments supplied: Pass --help or -h"
elif [ $1 == "--remove" ]; then
    echo "Removing relay " $2
    removeRelaySolidity $2
elif [ $1 == "--add" ]; then
    echo "Adding relay " $2
    addRelaySolidity $2
elif [ $1 == "--get" ]; then
    getRelaysSolidity
else
    echo "Invalid argument "
    echo "Ensure config.sh is for relevant configuration"
    echo
    echo "Usage:  "
    echo "      --add addr    : Add Relayer"
    echo "      --remove addr : Remove Relayer"
    echo "      --get         : Get Relayers"
fi
