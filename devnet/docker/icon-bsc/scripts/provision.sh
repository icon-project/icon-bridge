#!/bin/sh
set -e

# Parts of this code is adapted from https://github.com/icon-project/btp/blob/goloop2moonbeam/testnet/goloop2moonbeam/scripts
source env.variables.sh

source rpc.sh

eth_blocknumber() {
  curl -s -X POST $BSC_RPC_URI --header 'Content-Type: application/json' \
    --data-raw '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[], "id": 1}' | jq -r .result | xargs printf "%d\n"
}

goloop_lastblock() {
  goloop rpc lastblock
}

provision() {

  cp -r $BTPSIMPLE_BASE_DIR/keys/* $BTPSIMPLE_CONFIG_DIR
  cp $BTPSIMPLE_CONFIG_DIR/env $BTPSIMPLE_BIN_DIR/env

  if [ ! -f $BTPSIMPLE_CONFIG_DIR/provision ]; then
    echo "start provisioning..."

    # shellcheck disable=SC2059

    echo "$GOLOOP_RPC_NID.icon" >net.btp.icon
    mkdir -p $BTPSIMPLE_CONFIG_DIR/tx
    eth_blocknumber >/btpsimple/config/offset.bsc

    source token.javascore.sh
    source token.solidity.sh

    deploy_javascore_bmc
    # deploy_javascore_bsh
    deploy_javascore_irc2

    deploy_solidity_bmc

    #bmc_javascore_addService
    # bsh_javascore_register

    source nativeCoin.javascore.sh
    deploy_javascore_nativeCoin_BSH
    bmc_javascore_addNativeService
    nativeBSH_javascore_register
    nativeBSH_javascore_register_token
    nativeBSH_javascore_setFeeRatio

    #deploy_solidity_tokenBSH_BEP20
    source nativeCoin.solidity.sh
    deploy_solidity_nativeCoin_BSH

    generate_addresses_json >$BTPSIMPLE_CONFIG_DIR/addresses.json
    cp $BTPSIMPLE_CONFIG_DIR/addresses.json $SCRIPTS_DIR/

    #bsc_addService
    bsc_registerToken

    bmc_solidity_addNativeService
    nativeBSH_solidity_register
    
    token_bsc_fundBSH
    token_icon_fundBSH

    deploy_javascore_restrictor
    configure_javascore_TokenBSH_restrictor
    configure_javascore_NativeBSH_restrictor

    bmc_javascore_addLink
    bmc_javascore_addRelay
    bmc_javascore_setFeeAggregator

    add_icon_link
    add_icon_relay

    generate_relay_config >$BTPSIMPLE_CONFIG_DIR/bmr.config.json
    wait_for_file $BTPSIMPLE_CONFIG_DIR/bmr.config.json

    cp $BTPSIMPLE_CONFIG_DIR/addresses.json $BTPSIMPLE_CONTRACTS_DIR/solidity/bsh/
    cp $BTPSIMPLE_CONFIG_DIR/addresses.json $BTPSIMPLE_CONTRACTS_DIR/solidity/TokenBSH/

    touch $BTPSIMPLE_CONFIG_DIR/provision
    echo "provision is now complete"
  else
    prepare_solidity_env
  fi
}

prepare_solidity_env() {

  cp $BTPSIMPLE_CONFIG_DIR/env $BTPSIMPLE_CONTRACTS_DIR/solidity/bmc/.env
  cp $BTPSIMPLE_CONFIG_DIR/env $BTPSIMPLE_CONTRACTS_DIR/solidity/bsh/.env  

  cp $BTPSIMPLE_CONFIG_DIR/addresses.json $SCRIPTS_DIR/

  if [ ! -f $BTPSIMPLE_CONTRACTS_DIR/solidity/bsh/build/contracts/BSHCore.json ]; then
    cd $BTPSIMPLE_CONTRACTS_DIR/solidity/bsh/
    rm -rf contracts/test
    truffle compile --network bsc
  fi
}

wait_for_file() {
  FILE_NAME=$1
  timeout=10
  while [ ! -f "$FILE_NAME" ]; do
    if [ "$timeout" == 0 ]; then
      echo "ERROR: Timeout while waiting for the file $FILE_NAME."
      exit 1
    fi
    sleep 1
    timeout=$(expr $timeout - 1)

    echo "waiting for the output file: $FILE_NAME"
  done
}

btp_icon_validators_hash() {
  URI=$ICON_ENDPOINT \
    HEIGHT=$(decimal2Hex $(cat $CONFIG_DIR/offset.icon)) \
    $BTPSIMPLE_BIN_DIR/iconvalidators | jq -r .hash
}

generate_relay_config() {
  validatorsHash=$(btp_icon_validators_hash)
  jq -n '
    .base_dir = $base_dir |
    .log_level = "debug" |
    .console_level = "trace" |
    .log_writer.filename = $log_writer_filename |
    .relays = [ $b2i_relay, $i2b_relay ]' \
    --arg base_dir "$BASE_DIR" \
    --arg log_writer_filename "$LOG_FILENAME" \
    --argjson b2i_relay "$(
      jq -n '
            .name = "b2i" |
            .src.address = $src_address |
            .src.endpoint = [ $src_endpoint ] |
            .src.options = $src_options |
            .src.offset = $src_offset |
            .dst.address = $dst_address |
            .dst.endpoint = [ $dst_endpoint ] |
            .dst.options = $dst_options |
            .dst.key_store = $dst_key_store |
            .dst.key_store.coinType = $dst_key_store_cointype |
            .dst.key_password = $dst_key_password ' \
        --arg src_address "$(cat $CONFIG_DIR/btp.bsc)" \
        --arg src_endpoint "$BSC_ENDPOINT" \
        --argjson src_offset "$(cat $CONFIG_DIR/offset.bsc)" \
        --argjson src_options "$(
          jq -n {}
        )" \
        --arg dst_address "$(cat $CONFIG_DIR/btp.icon)" \
        --arg dst_endpoint "$ICON_ENDPOINT" \
        --argfile dst_key_store "$GOLOOP_RPC_KEY_STORE" \
        --arg dst_key_store_cointype "icx" \
        --arg dst_key_password "$(cat $GOLOOP_RPC_KEY_SECRET)" \
        --argjson dst_options '{"step_limit":13610920010, "tx_data_size_limit":8192}'
    )" \
    --argjson i2b_relay "$(
      jq -n '
            .name = "i2b" |
            .src.address = $src_address |
            .src.endpoint = [ $src_endpoint ] |
            .src.offset = $src_offset |
            .src.options.verifier.blockHeight = $src_options_verifier_blockHeight |
            .src.options.verifier.validatorsHash = $src_options_verifier_validatorsHash |
            .dst.address = $dst_address |
            .dst.endpoint = [ $dst_endpoint ] |
            .dst.options = $dst_options |
            .dst.tx_data_size_limit = $dst_tx_data_size_limit |
            .dst.key_store = $dst_key_store |
            .dst.key_store.coinType = $dst_key_store_cointype |
            .dst.key_password = $dst_key_password ' \
        --arg src_address "$(cat $CONFIG_DIR/btp.icon)" \
        --arg src_endpoint "$ICON_ENDPOINT" \
        --argjson src_offset "$(cat $CONFIG_DIR/offset.icon)" \
        --argjson src_options_verifier_blockHeight "$(cat $CONFIG_DIR/offset.icon)" \
        --arg src_options_verifier_validatorsHash "$validatorsHash" \
        --arg dst_address "$(cat $CONFIG_DIR/btp.bsc)" \
        --arg dst_endpoint "$BSC_ENDPOINT" \
        --argfile dst_key_store "$BSC_KEY_STORE" \
        --arg dst_key_store_cointype "evm" \
        --arg dst_key_password "$(cat $BSC_SECRET)" \
        --argjson dst_tx_data_size_limit 8192 \
        --argjson dst_options '{"gas_limit":8000000}'
    )"
}

wait-for-it.sh $GOLOOP_RPC_ADMIN_URI
# run provisioning
provision
