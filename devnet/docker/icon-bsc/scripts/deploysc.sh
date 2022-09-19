#!/bin/bash
set -e

# Parts of this code is adapted from https://github.com/icon-project/btp/blob/goloop2moonbeam/testnet/goloop2moonbeam/scripts
source config.sh
source keystore.sh
source rpc.sh


deploysc() {
    
    echo "start..."
    echo "check god keys..."
    if [ ! -f "${ICON_KEY_STORE}" ]; then
        ensure_key_store $ICON_KEY_STORE $ICON_SECRET
        echo "Fund newly created wallet " $ICON_KEY_STORE
        exit 0
    fi
    if [ ! -f "${BSC_KEY_STORE}" ]; then
        ensure_bsc_key_store $BSC_KEY_STORE $BSC_SECRET
        echo "Fund newly created wallet " $BSC_KEY_STORE
        exit 0
    fi
    export PRIVATE_KEY="[\""$(cat $BSC_KEY_STORE.priv)"\"]"
    # add owners
    echo "List/Create user accounts"
    ensure_key_store $CONFIG_DIR/keystore/icon.bts.wallet.json $CONFIG_DIR/keystore/icon.bts.wallet.secret
    ensure_key_store $CONFIG_DIR/keystore/icon.bmc.wallet.json $CONFIG_DIR/keystore/icon.bmc.wallet.secret
    ensure_key_store $CONFIG_DIR/keystore/icon.bmr.wallet.json $CONFIG_DIR/keystore/icon.bmr.wallet.secret
    ensure_key_store $CONFIG_DIR/keystore/icon.fa.wallet.json $CONFIG_DIR/keystore/icon.fa.wallet.secret

    ensure_bsc_key_store $CONFIG_DIR/keystore/bsc.bts.wallet.json $CONFIG_DIR/keystore/bsc.bts.wallet.secret
    ensure_bsc_key_store $CONFIG_DIR/keystore/bsc.bmc.wallet.json $CONFIG_DIR/keystore/bsc.bmc.wallet.secret
    ensure_bsc_key_store $CONFIG_DIR/keystore/bsc.bmr.wallet.json $CONFIG_DIR/keystore/bsc.bmr.wallet.secret

    echo "$GOLOOP_RPC_NID.icon" >$CONFIG_DIR/net.btp.icon #0x240fa7.icon
    mkdir -p $CONFIG_DIR/tx

    source token.javascore.sh
    source token.solidity.sh


    if [ ! -f $CONFIG_DIR/bsc.deploy.all ]; then
      echo "Deploy solidity"
      sleep 2
      deploy_solidity_bmc
      deploy_solidity_bts "${BSC_NATIVE_COIN_FIXED_FEE[0]}" "${BSC_NATIVE_COIN_FEE_NUMERATOR[0]}" "${BSC_NATIVE_COIN_DECIMALS[0]}"

      if [ -n "${INIT_ADDRESS_PATH}" ];
      then
        if [ ! -f $INIT_ADDRESS_PATH ]; then
          echo "No file found on "$INIT_ADDRESS_PATH
          return 1
        fi
        for i in "${!BSC_NATIVE_TOKEN_SYM[@]}"
        do
          addr=$(cat $INIT_ADDRESS_PATH | jq -r .solidity.${BSC_NATIVE_TOKEN_SYM[$i]})
          if [ "$addr" != "null" ]; 
          then
            echo -n $addr > $CONFIG_DIR/bsc.addr.${BSC_NATIVE_TOKEN_SYM[$i]}
          else 
            echo "BSC Token does not exist on address file" ${BSC_NATIVE_TOKEN_SYM[$i]}
            return 1
          fi
        done
      else 
        for i in "${!BSC_NATIVE_TOKEN_SYM[@]}"
        do
            deploy_solidity_token "${BSC_NATIVE_TOKEN_NAME[$i]}" "${BSC_NATIVE_TOKEN_SYM[$i]}"
        done              
      fi
      echo "CONFIGURE BSC"
      configure_solidity_add_bmc_owner
      configure_solidity_add_bts_service
      configure_solidity_set_fee_ratio "${BSC_NATIVE_COIN_FIXED_FEE[0]}" "${BSC_NATIVE_COIN_FEE_NUMERATOR[0]}"
      configure_solidity_add_bts_owner
      echo "Register BSC Tokens"
      for i in "${!BSC_NATIVE_TOKEN_SYM[@]}"
      do
          bsc_register_native_token "${BSC_NATIVE_TOKEN_NAME[$i]}" "${BSC_NATIVE_TOKEN_SYM[$i]}" "${BSC_NATIVE_TOKEN_FIXED_FEE[$i]}" "${BSC_NATIVE_TOKEN_FEE_NUMERATOR[$i]}" "${BSC_NATIVE_TOKEN_DECIMALS[$i]}"
          get_coinID "${BSC_NATIVE_TOKEN_NAME[$i]}" "${BSC_NATIVE_TOKEN_SYM[$i]}"
      done
      for i in "${!BSC_WRAPPED_COIN_SYM[@]}"
      do
          bsc_register_wrapped_coin "${BSC_WRAPPED_COIN_NAME[$i]}" "${BSC_WRAPPED_COIN_SYM[$i]}" "${BSC_WRAPPED_COIN_FIXED_FEE[$i]}" "${BSC_WRAPPED_COIN_FEE_NUMERATOR[$i]}" "${BSC_WRAPPED_COIN_DECIMALS[$i]}"
          get_coinID "${BSC_WRAPPED_COIN_NAME[$i]}" "${BSC_WRAPPED_COIN_SYM[$i]}"
      done
      echo "deployedSol" > $CONFIG_DIR/bsc.deploy.all 
    fi

    if [ ! -f $CONFIG_DIR/icon.deploy.all ]; then
      echo "Deploy Javascore"
      sleep 2
      deploy_javascore_bmc
      deploy_javascore_bts "${ICON_NATIVE_COIN_FIXED_FEE[0]}" "${ICON_NATIVE_COIN_FEE_NUMERATOR[0]}" "${ICON_NATIVE_COIN_DECIMALS[0]}"

      if [ -n "${INIT_ADDRESS_PATH}" ];
      then
        if [ ! -f $INIT_ADDRESS_PATH ]; then
          echo "No file found on "$INIT_ADDRESS_PATH
          return 1
        fi
        for i in "${!ICON_NATIVE_TOKEN_SYM[@]}"
        do
          addr=$(cat $INIT_ADDRESS_PATH | jq -r .javascore.${ICON_NATIVE_TOKEN_SYM[$i]})
          if [ "$addr" != "null" ]; 
          then
            echo -n $addr > $CONFIG_DIR/icon.addr.${ICON_NATIVE_TOKEN_SYM[$i]}
          else 
            echo "ICON Token ${ICON_NATIVE_TOKEN_SYM[$i]} does not exist on address file"
            return 1
          fi
        done
      else 
        for i in "${!ICON_NATIVE_TOKEN_SYM[@]}"
        do
            deploy_javascore_token "${ICON_NATIVE_TOKEN_NAME[$i]}" "${ICON_NATIVE_TOKEN_SYM[$i]}"
        done           
      fi
      echo "CONFIGURE ICON"
      configure_javascore_add_bmc_owner
      configure_javascore_bmc_setFeeAggregator
      configure_javascore_add_bts
      configure_javascore_add_bts_owner
      configure_javascore_bts_setICXFee "${ICON_NATIVE_COIN_FIXED_FEE[0]}" "${ICON_NATIVE_COIN_FEE_NUMERATOR[0]}"
      echo "Register ICON Tokens"
      for i in "${!ICON_NATIVE_TOKEN_SYM[@]}"
      do
          configure_javascore_register_native_token "${ICON_NATIVE_TOKEN_NAME[$i]}" "${ICON_NATIVE_TOKEN_SYM[$i]}" "${ICON_NATIVE_TOKEN_FIXED_FEE[$i]}" "${ICON_NATIVE_TOKEN_FEE_NUMERATOR[$i]}" "${ICON_NATIVE_TOKEN_DECIMALS[$i]}"
          get_btp_icon_coinId "${ICON_NATIVE_TOKEN_NAME[$i]}" "${ICON_NATIVE_TOKEN_SYM[$i]}"
      done
      for i in "${!ICON_WRAPPED_COIN_SYM[@]}"
      do
          configure_javascore_register_wrapped_coin "${ICON_WRAPPED_COIN_NAME[$i]}" "${ICON_WRAPPED_COIN_SYM[$i]}" "${ICON_WRAPPED_COIN_FIXED_FEE[$i]}" "${ICON_WRAPPED_COIN_FEE_NUMERATOR[$i]}" "${ICON_WRAPPED_COIN_DECIMALS[$i]}"
          get_btp_icon_coinId "${ICON_WRAPPED_COIN_NAME[$i]}" "${ICON_WRAPPED_COIN_SYM[$i]}"
      done
      echo "deployedJavascore" > $CONFIG_DIR/icon.deploy.all 
    fi

    if [ ! -f $CONFIG_DIR/link.all ]; then
      echo "LINK ICON"
      configure_javascore_addLink
      configure_javascore_setLinkHeight
      configure_bmc_javascore_addRelay
      echo "LINK BSC"
      add_icon_link
      set_link_height
      add_icon_relay
      echo "linked" > $CONFIG_DIR/link.all
    fi

    generate_addresses_json >$CONFIG_DIR/addresses.json  
    generate_relay_config >$CONFIG_DIR/bmr.config.json
    wait_for_file $CONFIG_DIR/bmr.config.json
    echo "Done deploying"
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
    --arg base_dir "bmr" \
    --arg log_writer_filename "bmr/bmr.log" \
    --argjson b2i_relay "$(
      jq -n '
            .name = "b2i" |
            .src.address = $src_address |
            .src.endpoint = [ $src_endpoint ] |
            .src.options.verifier.blockHeight = $src_options_verifier_blockHeight |
            .src.options.verifier.parentHash = $src_options_verifier_parentHash |
            .src.options.syncConcurrency = 100 |
            .src.offset = $src_offset |
            .dst.address = $dst_address |
            .dst.endpoint = [ $dst_endpoint ] |
            .dst.options = $dst_options |
            .dst.key_store = $dst_key_store |
            .dst.key_store.coinType = $dst_key_store_cointype |
            .dst.key_password = $dst_key_password ' \
        --arg src_address "$(cat $CONFIG_DIR/bsc.addr.bmcbtp)" \
        --arg src_endpoint "$BSC_ENDPOINT" \
        --argjson src_offset "$(cat $CONFIG_DIR/bsc.chain.height)" \
        --argjson src_options_verifier_blockHeight "$(cat $CONFIG_DIR/bsc.chain.height)" \
        --arg src_options_verifier_parentHash "$(cat $CONFIG_DIR/bsc.chain.parentHash)" \
        --arg dst_address "$(cat $CONFIG_DIR/icon.addr.bmcbtp)" \
        --arg dst_endpoint "$ICON_ENDPOINT" \
        --argfile dst_key_store "$CONFIG_DIR/keystore/icon.bmr.wallet.json" \
        --arg dst_key_store_cointype "icx" \
        --arg dst_key_password "$(cat $CONFIG_DIR/keystore/icon.bmr.wallet.secret)" \
        --argjson dst_options '{"step_limit":2500000000, "tx_data_size_limit":8192,"balance_threshold":"10000000000000000000"}'
    )" \
    --argjson i2b_relay "$(
      jq -n '
            .name = "i2b" |
            .src.address = $src_address |
            .src.endpoint = [ $src_endpoint ] |
            .src.offset = $src_offset |
            .src.options.verifier.blockHeight = $src_options_verifier_blockHeight |
            .src.options.verifier.validatorsHash = $src_options_verifier_validatorsHash |
            .src.options.syncConcurrency = 100 |
            .dst.address = $dst_address |
            .dst.endpoint = [ $dst_endpoint ] |
            .dst.options = $dst_options |
            .dst.tx_data_size_limit = $dst_tx_data_size_limit |
            .dst.key_store = $dst_key_store |
            .dst.key_store.coinType = $dst_key_store_cointype |
            .dst.key_password = $dst_key_password ' \
        --arg src_address "$(cat $CONFIG_DIR/icon.addr.bmcbtp)" \
        --arg src_endpoint "$ICON_ENDPOINT" \
        --argjson src_offset "$(cat $CONFIG_DIR/icon.chain.height)" \
        --argjson src_options_verifier_blockHeight "$(cat $CONFIG_DIR/icon.chain.height)" \
        --arg src_options_verifier_validatorsHash "$(cat $CONFIG_DIR/icon.chain.validators)" \
        --arg dst_address "$(cat $CONFIG_DIR/bsc.addr.bmcbtp)" \
        --arg dst_endpoint "$BSC_ENDPOINT" \
        --argfile dst_key_store "$CONFIG_DIR/keystore/bsc.bmr.wallet.json" \
        --arg dst_key_store_cointype "evm" \
        --arg dst_key_password "$(cat $CONFIG_DIR/keystore/bsc.bmr.wallet.secret)" \
        --argjson dst_tx_data_size_limit 8192 \
        --argjson dst_options '{"gas_limit":24000000,"tx_data_size_limit":8192,"balance_threshold":"100000000000000000000","boost_gas_price":1.0}'
    )"
}

#wait-for-it.sh $GOLOOP_RPC_ADMIN_URI
# run provisioning
deploysc