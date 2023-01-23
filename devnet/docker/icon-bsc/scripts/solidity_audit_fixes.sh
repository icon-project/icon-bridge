#!/bin/bash
set -e

source config.sh
source keystore.sh
source rpc.sh
source ownerManager.solidity.sh
source upgrade.solidity.bts.sh

echo "Loaded Configuration Files For " $TAG

read -p "Confirm? [y/N]: " proceed

case $TAG in
"ICON BSC TESTNET")
    f=bsc.audit
    ;;
"ICON BSC MAINNET")
    f=bsc.audit
    ;;
"ICON SNOW TESTNET")
    f=snow.audit
    ;;
"ICON SNOW MAINNET")
    f=snow.audit
    ;;
*)
    f=None
    ;;
esac

if [[ $proceed == "y" ]]; then
    if [ ! -f $ICONBRIDGE_CONFIG_DIR/$f ]; then
        echo "Proceeding after 10 seconds. You can still press Ctrl + C to exit ..."
        sleep 10

        echo "========= BTS OWNER MANAGER DEPLOYMENT =========="
        sleep 3
        deploy_solidity_owner_manager

        # add to json
        ownerManagerAddr=$(cat $ICONBRIDGE_CONFIG_DIR/bsc.addr.btsownermanager)
        cat $ICONBRIDGE_CONFIG_DIR/addresses.json | jq -r ".solidity += {\"BTSOwnerManager\":\"${ownerManagerAddr}\"}" >>copy_temp.json

        mv copy_temp.json $ICONBRIDGE_CONFIG_DIR/addresses.json

        echo "========= Upgrading BTSCore to BTSCoreV3 =========="
        sleep 3
        migrate_and_upgrade_bts_core BTSCore BTSCoreV3

        echo "========= Add BTSOwnerManager address to BTSCore =========="
        sleep 3
        set_owner_manager_in_bts_core

        echo "========= Update CoinDB =========="
        sleep 3
        update_coin_db

        echo "========= Upgrading BTSPeriphery to BTSPeripheryV2 =========="
        sleep 3
        migrate_and_upgrade_bts_periphery BTSPeriphery BTSPeripheryV2

        echo "Audit fixes for $TAG deployed" >$ICONBRIDGE_CONFIG_DIR/$f
        echo "=============================================="
        echo "Audit fixes for $TAG successfully deployed ..."
        echo "=============================================="
    else
        echo "================================"
        echo "$TAG Audit Fixes already deployed"
        echo "================================"
    fi
else
    echo "Exit."
fi
