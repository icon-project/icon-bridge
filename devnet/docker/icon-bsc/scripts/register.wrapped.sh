#! /bin/bash
set -e

source config.sh
source token.javascore.sh
source token.solidity.sh
source utils.sh

register_icon_wrapped_coin() {
    echo "register_icon_wrapped_coin. Checking..."
    cd $CONFIG_DIR
    for i in "${!ICON_WRAPPED_COIN_SYM[@]}" 
    do 
        echo ${ICON_WRAPPED_COIN_NAME[$i]}
        local coinIDRes=$(goloop rpc call --to $(cat icon.addr.bts) --method coinId --param _coinName="${ICON_WRAPPED_COIN_NAME[$i]}" | jq -r .)
        if [ "$coinIDRes" == "null" ]; then
            configure_javascore_register_wrapped_coin "${ICON_WRAPPED_COIN_NAME[$i]}" "${ICON_WRAPPED_COIN_SYM[$i]}" "${ICON_WRAPPED_COIN_FIXED_FEE[$i]}" "${ICON_WRAPPED_COIN_FEE_NUMERATOR[$i]}" "${ICON_WRAPPED_COIN_DECIMALS[$i]}"
            get_btp_icon_coinId "${ICON_WRAPPED_COIN_NAME[$i]}" "${ICON_WRAPPED_COIN_SYM[$i]}"
        fi
    done
}

register_eth_wrapped_coin() {
    echo "register_eth_wrapped_coin. Checking..."
    cd $CONTRACTS_DIR/solidity/bts
    export PRIVATE_KEY="[\""$(cat $BSC_KEY_STORE.priv)"\"]"
    for i in "${!BSC_WRAPPED_COIN_SYM[@]}" 
    do
        echo ${BSC_WRAPPED_COIN_NAME[$i]}
        tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js --method coinId --coinName "${BSC_WRAPPED_COIN_NAME[$i]}")
        coinId=$(echo "$tx" | grep "coinId:" | sed -e "s/^coinId: //")
        exists=$(echo $coinId | wc -l | awk '{$1=$1;print}')
        if [ "$exists" != "1" ]; then
            sleep 3
            bsc_register_wrapped_coin "${BSC_WRAPPED_COIN_NAME[$i]}" "${BSC_WRAPPED_COIN_SYM[$i]}" "${BSC_WRAPPED_COIN_FIXED_FEE[$i]}" "${BSC_WRAPPED_COIN_FEE_NUMERATOR[$i]}" "${BSC_WRAPPED_COIN_DECIMALS[$i]}"
            get_coinID "${BSC_WRAPPED_COIN_NAME[$i]}" "${BSC_WRAPPED_COIN_SYM[$i]}"
        fi
        sleep 3
    done
}

register_icon_wrapped_coin
register_eth_wrapped_coin
generate_addresses_json >$CONFIG_DIR/addresses.json