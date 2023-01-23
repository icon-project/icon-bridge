#!/bin/bash

set -e
source utils.sh
source config.sh

ROOT_DIR=$(echo "$(
    cd "$(dirname "../../../../../")"
    pwd
)")

export PRIVATE_KEY="[\""$(cat $BSC_KEY_STORE.priv)"\"]"

deploy_solidity_owner_manager() {
    cd $ROOT_DIR/solidity

    if [ ! -f $ROOT_DIR/solidity/bts/contracts/interfaces/IBTSOwnerManager.sol ]; then
        echo "Contract Interface IBTSOwnerManager does not exist"
        exit 0
    fi

    if [ ! -f $ROOT_DIR/solidity/bts/contracts/BTSOwnerManager.sol ]; then
        echo "Contract BTSOwnerManager does not exist"
        exit 0
    fi

    if [ ! -f $ROOT_DIR/solidity/bts/migrations/5_deploy_owner_manager.js ]; then
        echo "Owner Manager deployment migration file does not exist"
        exit 0
    fi

    cp bts/contracts/BTSOwnerManager.sol $CONTRACTS_DIR/solidity/bts/contracts/BTSOwnerManager.sol
    cp bts/contracts/interfaces/IBTSOwnerManager.sol $CONTRACTS_DIR/solidity/bts/contracts/interfaces/IBTSOwnerManager.sol
    cp bts/migrations/5_deploy_owner_manager.js $CONTRACTS_DIR/solidity/bts/migrations/5_deploy_owner_manager.js

    sleep 2

    cd $CONTRACTS_DIR/solidity/bts

    if [ ! -f $CONFIG_DIR/bsc.deploy.btsownermanager ]; then
        echo "Deploying solidity owner manager"
        rm -rf contracts/test
        truffle compile --all
        set +e
        local status="retry"
        for i in $(seq 1 20); do
            truffle migrate --compile-none --network bsc --f 5 --to 5
            if [ $? == 0 ]; then
                status="ok"
                break
            fi
            echo "Retry: "$i
        done
        set -e
        if [ "$status" == "retry" ]; then
            exit 1
        fi
        jq -r '.networks[] | .address' build/contracts/BTSOwnerManager.json >$CONFIG_DIR/bsc.addr.btsownermanager

        echo -n "bts owner manager" >$CONFIG_DIR/bsc.deploy.btsownermanager

        echo "BTS OWNER MANAGER SUCCESSFULLY DEPLOYED"
    else
        echo "BTS Owner Manager deployed already"
    fi
}

set_owner_manager_in_bts_core() {
    echo "Setting owner manager in BTS Core"

    cd $CONTRACTS_DIR/solidity/bts

    if [ ! -f $CONFIG_DIR/bsc.addr.btsownermanager ]; then
        echo "bsc.addr.btsownermanager file does not exist"
        exit 0
    fi

    ownerManagerAddr=$(cat $CONFIG_DIR/bsc.addr.btsownermanager)

    if [ ! -f $CONFIG_DIR/bsc.configure.addOwnerManager ]; then
        tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
            --method addBTSOwnerManager --addr "$ownerManagerAddr")
        echo "$tx" >$CONFIG_DIR/tx/addOwnerManager.bts.bsc
        isTrue=$(echo "$tx" | grep "status: true" | wc -l | awk '{$1=$1;print}')
        if [ "$isTrue" == "1" ]; then
            echo "OwnerManagerAdded" >$CONFIG_DIR/bsc.configure.addOwnerManager
        else
            echo "Error Addding Owner Manager"
            return 1
        fi
    else
        echo "Owner Manager already updated"
    fi

}

update_coin_db() {
    echo "Setting coinDB..."
    cd $CONTRACTS_DIR/solidity/bts

    if [ ! -f $CONFIG_DIR/bsc.configure.updateCoinDb ]; then
        tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
            --method updateCoinDb)
        echo "$tx" >$CONFIG_DIR/tx/updatecoindb.bts.bsc
        isTrue=$(echo "$tx" | grep "status: true" | wc -l | awk '{$1=$1;print}')
        if [ "$isTrue" == "1" ]; then
            echo "CoinDBUpdated" >$CONFIG_DIR/bsc.configure.updateCoinDb
            echo "CoinDB Successfully Updated"
        else
            echo "Error Addding Owner Manager"
            return 1
        fi
    else
        echo "coinDB already updated"
    fi

}
