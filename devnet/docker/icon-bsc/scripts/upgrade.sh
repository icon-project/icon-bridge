#!/bin/bash
set -e

source config.sh
source keystore.sh
source rpc.sh
source utils.sh

ROOT_DIR=$(echo "$(cd "$(dirname "../../../../../")"; pwd)")

build_bts() {
    echo "building bts jar"
    cd $ROOT_DIR/javascore
    gradle clean
    gradle bts:optimizedJar
    cp bts/build/libs/bts-optimized.jar $CONTRACTS_DIR/javascore/bts.jar
    cp lib/irc2Tradeable-0.1.0-optimized.jar  $CONTRACTS_DIR/javascore/irc2Tradeable.jar
    echo "build bts complete"
}

upgrade_javascore_bts() {
    echo "upgrading javascore bts"
    cd $CONFIG_DIR
    if [ ! -f icon.addr.bts ]; then
        echo "BTS address file icon.addr.bts does not exist"
        exit
    fi
    if [ ! -f icon.addr.bmc ]; then
        echo "BMC address file icon.addr.bmc does not exist"
        exit
    fi
    if [ ! -f icon.addr.bts.upgrade ]; then
        goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/bts.jar \
        --content_type application/java \
        --to $(cat icon.addr.bts) \
        --param _name="${ICON_NATIVE_COIN_NAME[0]}" \
        --param _bmc=$(cat icon.addr.bmc) \
        --param _decimals=$(decimal2Hex $3) \
        --param _feeNumerator=$(decimal2Hex $2) \
        --param _fixedFee=$(decimal2Hex $1) \
        --param _serializedIrc2=$(xxd -p $CONTRACTS_DIR/javascore/irc2Tradeable.jar | tr -d '\n') | jq -r . > tx/tx.icon.bts.upgrade
        sleep 5
        extract_scoreAddress tx/tx.icon.bts.upgrade icon.addr.bts.upgrade
        echo "Upgraded Address: "
        cat icon.addr.bts.upgrade
    fi
}
echo "Start Upgrade "
build_bts
upgrade_javascore_bts "${ICON_NATIVE_COIN_FIXED_FEE[0]}" "${ICON_NATIVE_COIN_FEE_NUMERATOR[0]}" "${ICON_NATIVE_COIN_DECIMALS[0]}"
echo "Done"
