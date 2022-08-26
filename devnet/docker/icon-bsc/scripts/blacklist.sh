#!/bin/bash
set -e

source config.sh
source rpc.sh
source utils.sh

# change users to blacklist
USER_LIST=(hxd47ad924eba01ec91330e4e996cf7b8c658f4e4c
    hxfb6251ac765fd428c4f961ad649050bbbf77210d
    hxd75437c389ff4bf4cd42a6a88c34fcb7cdcbce8a
    hx4e0343b6bb01abe3deac0d5be0b578addd700b35)

# mainnet
# NET="0x1.icon"

# testnet
NET="0x2.icon"


isUserBlacklisted() {
    echo -n "${2} : "
    resp=$(goloop rpc call --to $(cat icon.addr.bts) \
        --method isUserBlackListed \
        --param _net=${1} \
        --param _address=${2} | jq -r .)
    if [ $resp == "0x0" ]; then
        echo "Not blacklisted"
    else
        echo "blacklisted"
    fi
}

getBlacklistedUsers() {
    cd $CONFIG_DIR
    echo "Blacklisted Users On Network: ${1} "
    goloop rpc call --to $(cat icon.addr.bts) \
        --method getBlackListedUsers \
        --param _net=${1} \
        --param _start=0x0 \
        --param _end=0x64
}

addBlacklistedUser() {
    cd $CONFIG_DIR

    if [ ! -f icon.addr.bts ]; then
        echo "BTS address file icon.addr.bts does not exist"
        exit
    fi

    net=$1
    shift
    echo "Add the following users to blacklist on ${net}"
    for i in ${@}; do
        echo ${i}
    done
    param="{\"params\":{\"_net\":\"${net}\",\"_addresses\":$(toJsonArray ${@})}}"

    goloop rpc sendtx call --to $(cat icon.addr.bts) \
        --method addBlacklistAddress \
        --raw $param | jq -r . >tx/addBlacklistedUser.$NET.icon
    sleep 3
    ensure_txresult tx/addBlacklistedUser.$NET.icon
    echo "Added to blacklist"
    getBlacklistedUsers ${net}
}

removeBlacklistedUser() {
    cd $CONFIG_DIR
    if [ ! -f icon.addr.bts ]; then
        echo "BTS address file icon.addr.bts does not exist"
        exit
    fi

    net=$1
    shift
    echo "Current Blacklist status of users on ${net}: "
    for i in ${@}; do
        isUserBlacklisted ${net} ${i}
    done

    echo "To remove given users on ${net}"
    param="{\"params\":{\"_net\":\"${net}\",\"_addresses\":$(toJsonArray ${@})}}"
    goloop rpc sendtx call --to $(cat icon.addr.bts) \
        --method removeBlacklistAddress \
        --raw $param | jq -r . >tx/removeBlacklistedUser.$NET.icon
    sleep 3
    ensure_txresult tx/removeBlacklistedUser.$NET.icon
    echo "Removed from blacklist"
    getBlacklistedUsers ${net}
}

if [ $# -eq 0 ]; then
    echo "No arguments supplied: Pass --show to get blacklisted users, --add to add to blacklist, --remove to remove from blacklist"
elif [ $1 == "--show" ]; then
    getBlacklistedUsers $NET
elif [ $1 == "--add" ]; then
    addBlacklistedUser $NET "${USER_LIST[@]}"
elif [ $1 == "--remove" ]; then
    removeBlacklistedUser $NET "${USER_LIST[@]}"
else
    echo "Invalid argument: Pass --show to get blacklisted users, --add to add to blacklist, --remove to remove from blacklist"
fi
