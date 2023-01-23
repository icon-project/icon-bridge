#!/bin/bash

set -e
source utils.sh
source config.sh

ROOT_DIR=$(echo "$(
    cd "$(dirname "../../../../../")"
    pwd
)")

export PRIVATE_KEY="[\""$(cat $BSC_KEY_STORE.priv)"\"]"

copy_bts_core_migrations() {
    echo "copying $1 migration start"

    if [ ! -f $ROOT_DIR/solidity/bts/contracts/${2}.sol ]; then
        echo "Contract ${1} to upgrade to: ${2} does not exist"
        exit 0
    fi

    if [ ! -f $ROOT_DIR/solidity/bts/migrations/4_upgrade_bts.js ]; then
        echo "Migration script does not exist"
        exit 0
    fi

    cd $ROOT_DIR/solidity
    cp bts/contracts/${2}.sol $CONTRACTS_DIR/solidity/bts/contracts/${2}.sol
    cp bts/contracts/interfaces/IBTSCoreV2.sol $CONTRACTS_DIR/solidity/bts/contracts/interfaces/IBTSCoreV2.sol
    cp bts/migrations/4_upgrade_bts.js $CONTRACTS_DIR/solidity/bts/migrations/4_upgrade_bts.js
    echo "$1 Migration copied"
}

upgrade_solidity_bts_core() {
    echo "Upgrading solidity btsCore"
    cd $CONTRACTS_DIR/solidity/bts
    if [ ! -f $CONFIG_DIR/bsc.addr.btscore ]; then
        echo "BTSCore address file bsc.addr.btscore does not exist"
        exit
    fi

    if [ ! -f $CONFIG_DIR/bsc.btscore.upgrade ]; then

        truffle compile -all

        echo "Check if ${2} contract exists: "
        if [ ! -f $CONTRACTS_DIR/solidity/bts/build/contracts/${2}.json ]; then
            echo "Contract BTSCore to upgrade to: ${2} not compiled"
            exit 0
        fi
        echo "${2} exists"

        proxyBTSCore=$(jq -r '.networks[] | .address' $CONTRACTS_DIR/solidity/bts/build/contracts/${1}.json)
        deployedBTSCore=$(cat $CONFIG_DIR/bsc.addr.btscore)

        if [ "$proxyBTSCore" != "$deployedBTSCore" ]; then
            echo "Address not verified"
            exit 0
        fi

        set +e
        local status="retry"

        for i in $(seq 1 20); do
            truffle migrate --compile-all --network bsc --f 4 --to 4 --contract ${1} --upgradeTo ${2}
            if [ $? == 0 ]; then
                status="ok"
                break
            fi
            echo "Retry: "$i
        done
        set -e
        if [ "$status" == "retry" ]; then
            echo "BTSCore Upgrade Failed after retry"
            exit 1
        fi
        echo 'BTSCore Proxy Address after upgrade'
        jq -r '.networks[] | .address' build/contracts/BTSCore.json
        echo -n "btscoreupgraded" >$CONFIG_DIR/bsc.btscore.upgrade

    fi
}

copy_bts_periphery_migrations() {
    echo "copying $1 migration start"

    if [ ! -f $ROOT_DIR/solidity/bts/contracts/${1}.sol ]; then
        echo "Contract ${1} to upgrade to: ${2} does not exist"
        exit 0
    fi

    if [ ! -f $ROOT_DIR/solidity/bts/migrations/4_upgrade_bts.js ]; then
        echo "Migration script does not exist"
        exit 0
    fi

    cd $ROOT_DIR/solidity
    cp bts/contracts/${2}.sol $CONTRACTS_DIR/solidity/bts/contracts/${2}.sol
    cp bts/migrations/4_upgrade_bts.js $CONTRACTS_DIR/solidity/bts/migrations/4_upgrade_bts.js
    echo "$1 Migration copied"
}

upgrade_solidity_bts_periphery() {
    echo "Upgrading solidity bts periphery"
    cd $CONTRACTS_DIR/solidity/bts
    if [ ! -f $CONFIG_DIR/bsc.addr.btsperiphery ]; then
        echo "BTSPeriphery address file bsc.addr.btsperiphery does not exist"
        exit
    fi

    if [ ! -f $CONFIG_DIR/bsc.btsperiphery.upgrade ]; then

        truffle compile -all

        echo "Check if ${2} contract exists: "
        if [ ! -f $CONTRACTS_DIR/solidity/bts/build/contracts/${2}.json ]; then
            echo "Contract BTSPeriphery to upgrade to: ${2} not compiled"
            exit 0
        fi
        echo "${2} exists"

        proxyBTSPeriphery=$(jq -r '.networks[] | .address' $CONTRACTS_DIR/solidity/bts/build/contracts/${1}.json)
        deployedBTSPeriphery=$(cat $CONFIG_DIR/bsc.addr.btsperiphery)

        if [ "$proxyBTSPeriphery" != "$deployedBTSPeriphery" ]; then
            echo "Address not verified"
            exit 0
        fi

        set +e
        local status="retry"

        for i in $(seq 1 20); do
            truffle migrate --compile-all --network bsc --f 4 --to 4 --contract ${1} --upgradeTo ${2}
            if [ $? == 0 ]; then
                status="ok"
                break
            fi
            echo "Retry: "$i
        done
        set -e
        if [ "$status" == "retry" ]; then
            echo "BTSPeriphery Upgrade Failed after retry"
            exit 1
        fi
        echo 'BTSPeriphery Proxy Address after upgrade'
        jq -r '.networks[] | .address' build/contracts/BTSPeriphery.json
        echo -n "btsPeripheryUpgraded" >$CONFIG_DIR/bsc.btsperiphery.upgrade
    fi
}

migrate_and_upgrade_bts_core() {
    copy_bts_core_migrations $1 $2
    upgrade_solidity_bts_core $1 $2
}

migrate_and_upgrade_bts_periphery() {
    copy_bts_periphery_migrations $1 $2
    upgrade_solidity_bts_periphery $1 $2
}

if [ $# -eq 0 ]; then
    echo "No arguments supplied: Pass --help for details."
elif [[ $1 == "--contract" && $3 == "--upgradeTo" ]]; then
    echo "Upgrade $2 to $4 "
    if [ $2 == "BTSCore" ]; then
        migrate_and_upgrade_bts_core $2 $4
        echo "Done"
    elif [ $2 == "BTSPeriphery" ]; then
        migrate_and_upgrade_bts_periphery $2 $4
        echo "Done"
    else
        echo "Invalid contract"
    fi
else
    echo "Invalid argument: Pass --contract BTSCore --upgradeTo BTSCoreV2 to upgrade BTSCore to BTSCoreV2"
    echo "Example: ./upgrade.bts.solidity.sh --contract BTSCore --upgradeTo BTSCoreV2"
fi
