#!/bin/bash
set -e

# Parts of this code is adapted from https://github.com/icon-project/btp/blob/goloop2moonbeam/testnet/goloop2moonbeam/scripts
source env.variables.sh

source rpc.sh


provision() {
  cp -r $ICONBRIDGE_BASE_DIR/keys/* $ICONBRIDGE_CONFIG_DIR
  cp $ICONBRIDGE_CONFIG_DIR/env $ICONBRIDGE_BIN_DIR/env

  if [ ! -f $ICONBRIDGE_CONFIG_DIR/provision ]; then
    echo "start provisioning..."
    sleep 10
    echo "$GOLOOP_RPC_NID.icon" >net.btp.icon #0x240fa7.icon
    mkdir -p $ICONBRIDGE_CONFIG_DIR/tx

    source token.javascore.sh
    source token.solidity.sh

    #deploy icon
    deploy_javascore_bmc
    deploy_javascore_bsr
    deploy_javascore_bts
    deploy_javascore_irc2
    deploy_solidity_bmc
    deploy_solidity_bts
    #deploy bsc


    generate_addresses_json >$ICONBRIDGE_CONFIG_DIR/addresses.json
    cp $ICONBRIDGE_CONFIG_DIR/addresses.json $SCRIPTS_DIR/
    cp $ICONBRIDGE_CONFIG_DIR/addresses.json $ICONBRIDGE_CONTRACTS_DIR/solidity/bts/
    cp $ICONBRIDGE_CONFIG_DIR/addresses.json $ICONBRIDGE_CONTRACTS_DIR/solidity/bmc/

    #configure icon
    echo "CONFIGURE ICON"
    configure_javascore_add_bmc_owner
    configure_javascore_bmc_setFeeAggregator
    configure_javascore_add_bts
    configure_javascore_add_bts_owner
    configure_javascore_set_bsr
    configure_javascore_bts_setICXFee
    #configure bsc    
    echo "CONFIGURE BSC"
    configure_solidity_add_bts_service
    configure_solidity_set_fee_ratio

    #Link icon
    echo "LINK ICON"
    configure_javascore_addLink
    configure_bmc_javascore_addRelay
    configure_javascore_register_bnb
    get_btp_icon_bnb
    configure_javascore_register_ticx
    configure_javascore_register_tbnb
    get_btp_icon_tbnb

    #Link bsc
    echo "LINK BSC"
    add_icon_link
    set_link_height
    add_icon_relay
    bsc_register_icx
    get_coinID_icx
    bsc_register_tbnb
    bsc_register_ticx
    get_coinID_ticx

    # token_bsc_fundBSH
    # token_icon_fundBSH

    generate_relay_config >$ICONBRIDGE_CONFIG_DIR/bmr.config.json
    wait_for_file $ICONBRIDGE_CONFIG_DIR/bmr.config.json



    touch $ICONBRIDGE_CONFIG_DIR/provision
    echo "provision is now complete"
  else
    prepare_solidity_env
  fi
}

prepare_solidity_env() {

  cp $ICONBRIDGE_CONFIG_DIR/env $ICONBRIDGE_CONTRACTS_DIR/solidity/bmc/.env
  cp $ICONBRIDGE_CONFIG_DIR/env $ICONBRIDGE_CONTRACTS_DIR/solidity/bts/.env

  cp $ICONBRIDGE_CONFIG_DIR/addresses.json $SCRIPTS_DIR/

  if [ ! -f $ICONBRIDGE_CONTRACTS_DIR/solidity/bts/build/contracts/BTSCore.json ]; then
    cd $ICONBRIDGE_CONTRACTS_DIR/solidity/bts/
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



generate_relay_config() {
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
        --arg src_address "$(cat $CONFIG_DIR/btp.bsc.btp.address)" \
        --arg src_endpoint "$BSC_ENDPOINT" \
        --argjson src_offset "$(cat $CONFIG_DIR/btp.bsc.block.height)" \
        --argjson src_options "$(
          jq -n {}
        )" \
        --arg dst_address "$(cat $CONFIG_DIR/btp.icon.btp.address)" \
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
        --arg src_address "$(cat $CONFIG_DIR/btp.icon.btp.address)" \
        --arg src_endpoint "$ICON_ENDPOINT" \
        --argjson src_offset "$(cat $CONFIG_DIR/btp.icon.block.height)" \
        --argjson src_options_verifier_blockHeight "$(cat $CONFIG_DIR/btp.icon.block.height)" \
        --arg src_options_verifier_validatorsHash "$(cat $CONFIG_DIR/btp.icon.validators.hash)" \
        --arg dst_address "$(cat $CONFIG_DIR/btp.bsc.btp.address)" \
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
