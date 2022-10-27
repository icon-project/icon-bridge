#!/bin/bash
set -e

source config.sh
source keystore.sh
source rpc.sh

setup_linked_account() {
  echo "check god keys..."
  ## ICON 
  local godKeysFilename=$(basename $ICON_KEY_STORE)
  local godKeystorePath=${linkfrom}/keystore/${godKeysFilename}
  if [ ! -f "${godKeystorePath}" ]; then 
    echo "Keystore does not exist inside liked path: Expected file: "$godKeystorePath
  else 
    echo "Copying icon keystores from path "${linkfrom}/keystore to $ICONBRIDGE_CONFIG_DIR/keystore/
    mkdir -p $ICONBRIDGE_CONFIG_DIR/keystore
    cp -r ${linkfrom}/keystore/icon.* $ICONBRIDGE_CONFIG_DIR/keystore/
  fi 

  ## Arctic
    if [ ! -f "${BSC_KEY_STORE}" ]; then
      ensure_bsc_key_store $BSC_KEY_STORE $BSC_SECRET
      echo "Do not Panic..."
      echo "Missing BSC God Wallet on the required path. One has been created "$BSC_KEY_STORE
      echo "Fund this newly created wallet and rerun again " 
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
}

setup_account() {
    echo "check god keys..."
    if [ ! -f "${ICON_KEY_STORE}" ]; then
        ensure_key_store $ICON_KEY_STORE $ICON_SECRET
        echo "Do not Panic..."
        echo "Missing ICON God Wallet on the required path. One has been created "$ICON_KEY_STORE
        echo "Fund this newly created wallet and rerun the same command again" 
        exit 0
    fi
    if [ ! -f "${BSC_KEY_STORE}" ]; then
        ensure_bsc_key_store $BSC_KEY_STORE $BSC_SECRET
        echo "Do not Panic..."
        echo "Missing BSC God Wallet on the required path. One has been created "$BSC_KEY_STORE
        echo "Fund this newly created wallet and rerun again " 
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
}

linked_deploysc() {
    source token.javascore.sh
    source token.solidity.sh

    if [ ! -d $BUILD_DIR ]; then 
      echo "Do not Panic..."
      echo "Build Artifacts have not been created. Expected on path "$BUILD_DIR 
      echo "Run make buildsc to do so. Check README.md for more"
      exit 0
    fi

    echo "Start. Wait time 5 seconds"
    sleep 5

    if [ ! -f "${linkfrom}/icon.deploy.all" ]; then 
      echo "Missing file on expected path "${linkfrom}/icon.deploy.all
      exit 0
    else 
      if [ ! -f $ICONBRIDGE_CONFIG_DIR/icon.deploy.all ]; then 
        echo "Copying icon metadata from path "${linkfrom}
        cp -r ${linkfrom}/icon.* $ICONBRIDGE_CONFIG_DIR/
        rm $ICONBRIDGE_CONFIG_DIR/icon.configure.*
        extract_chain_height_and_validator
      fi
    fi

    echo "$GOLOOP_RPC_NID.icon" >$CONFIG_DIR/net.btp.icon #0x240fa7.icon
    mkdir -p $CONFIG_DIR/tx


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
    generate_e2e_config >$CONFIG_DIR/e2e.config.json
    wait_for_file $CONFIG_DIR/bmr.config.json
}

deploysc() {
    if [ ! -d $BUILD_DIR ]; then 
      echo "Do not Panic..."
      echo "Build Artifacts have not been created. Expected on path "$BUILD_DIR 
      echo "Run make buildsc to do so. Check README.md for more"
      exit 0
    fi
    echo "Start "
    sleep 15
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
    generate_e2e_config >$CONFIG_DIR/e2e.config.json
    wait_for_file $CONFIG_DIR/bmr.config.json
    
    echo "Smart contracts have been deployed "
    echo "You can now run the relay with make runrelaysrc OR make runrelayimg"
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
            .name = "s2i" |
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
        --arg src_options_verifier_validatorData "$(cat $CONFIG_DIR/bsc.chain.validatorData)" \
        --arg dst_address "$(cat $CONFIG_DIR/icon.addr.bmcbtp)" \
        --arg dst_endpoint "$ICON_ENDPOINT" \
        --argfile dst_key_store "$CONFIG_DIR/keystore/icon.bmr.wallet.json" \
        --arg dst_key_store_cointype "icx" \
        --arg dst_key_password "$(cat $CONFIG_DIR/keystore/icon.bmr.wallet.secret)" \
        --argjson dst_options '{"step_limit":2500000000, "tx_data_size_limit":8192,"balance_threshold":"10000000000000000000"}'
    )" \
    --argjson i2b_relay "$(
      jq -n '
            .name = "i2s" |
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

generate_e2e_config() {
    jq -n '
    .log_level = "debug" |
    .console_level = "info" |
    .log_writer.filename = "_ixh_e2e.log" |
    .fee_aggregator = $fee_aggregator |
    .enable_experimental_features = false |
    .chains = [ $icon_config, $bsc_config ]' \
    --arg fee_aggregator "$(cat $CONFIG_DIR/keystore/icon.fa.wallet.json | jq -r .address)" \
    --argjson icon_config "$(
        jq -n '
            .name = "ICON" |
            .url = $url |
            .contract_addresses = $contract_addresses |
            .native_coin = $native_coin |
            .native_tokens = $native_tokens | 
            .wrapped_coins = $wrapped_coins |
            .god_wallet_keystore_path = $god_wallet_keystore_path |
            .god_wallet_secret_path = $god_wallet_secret_path |
            .bts_owner_keystore_path = $bts_owner_keystore_path |
            .bts_owner_secret_path = $bts_owner_secret_path | 
            .network_id = $network_id | 
            .gas_limit = $gas_limit' \
            --arg url "$ICON_ENDPOINT" \
            --argjson contract_addresses "$(
                jq -n '
                .BTS = $bts_address | 
                .BMC = $bmc_address' \
                 --arg bts_address "$(cat $CONFIG_DIR/icon.addr.bts)" \
                 --arg bmc_address "$(cat $CONFIG_DIR/icon.addr.bmc)" \
            )" \
            --arg native_coin "${ICON_NATIVE_COIN_NAME[0]}" \
            --argjson native_tokens  "$(
                    jq --compact-output --null-input '$ARGS.positional' --args -- "${ICON_NATIVE_TOKEN_NAME[@]}"
                )"\
            --argjson wrapped_coins  "$(
                    jq --compact-output --null-input '$ARGS.positional' --args -- "${ICON_WRAPPED_COIN_NAME[@]}"
                )"\
            --arg god_wallet_keystore_path $(echo $CONFIG_DIR/keystore/icon.god.wallet.json) \
            --arg god_wallet_secret_path $(echo $CONFIG_DIR/keystore/icon.god.wallet.secret) \
            --arg bts_owner_keystore_path $(echo $CONFIG_DIR/keystore/icon.god.wallet.json) \
            --arg bts_owner_secret_path $(echo $CONFIG_DIR/keystore/icon.god.wallet.secret) \
            --arg network_id $ICON_BMC_NET \
            --argjson gas_limit '{
                "TransferNativeCoinIntraChainGasLimit":150000,
                "TransferTokenIntraChainGasLimit":300000,
                "ApproveTokenInterChainGasLimit":800000,
                "TransferCoinInterChainGasLimit":2500000,
                "TransferBatchCoinInterChainGasLimit":4000000,
                "DefaultGasLimit":5000000
            }'
        )" \
    --argjson bsc_config "$(
        jq -n '
            .name = "SNOW" |
            .url = $url |
            .contract_addresses = $contract_addresses |
            .native_coin = $native_coin |
            .native_tokens = $native_tokens | 
            .wrapped_coins = $wrapped_coins |
            .god_wallet_keystore_path = $god_wallet_keystore_path |
            .god_wallet_secret_path = $god_wallet_secret_path |
            .bts_owner_keystore_path = $bts_owner_keystore_path |
            .bts_owner_secret_path = $bts_owner_secret_path | 
            .network_id = $network_id | 
            .gas_limit = $gas_limit' \
            --arg url "$BSC_ENDPOINT" \
            --argjson contract_addresses "$(
                jq -n '
                .BTS = $bts_address | 
                .BTSPeriphery = $bts_periphery_address | 
                .BMCPeriphery = $bmc_periphery_address' \
                 --arg bts_address "$(cat $CONFIG_DIR/bsc.addr.btscore)" \
                 --arg bts_periphery_address "$(cat $CONFIG_DIR/bsc.addr.btsperiphery)" \
                 --arg bmc_periphery_address "$(cat $CONFIG_DIR/bsc.addr.bmcperiphery)" \
            )" \
            --arg native_coin "${BSC_NATIVE_COIN_NAME[0]}" \
            --argjson native_tokens  "$(
                    jq --compact-output --null-input '$ARGS.positional' --args -- "${BSC_NATIVE_TOKEN_NAME[@]}"
                )"\
            --argjson wrapped_coins  "$(
                    jq --compact-output --null-input '$ARGS.positional' --args -- "${BSC_WRAPPED_COIN_NAME[@]}"
                )"\
            --arg god_wallet_keystore_path $(echo $CONFIG_DIR/keystore/bsc.god.wallet.json) \
            --arg god_wallet_secret_path $(echo $CONFIG_DIR/keystore/bsc.god.wallet.secret) \
            --arg bts_owner_keystore_path $(echo $CONFIG_DIR/keystore/bsc.god.wallet.json) \
            --arg bts_owner_secret_path $(echo $CONFIG_DIR/keystore/bsc.god.wallet.secret) \
            --arg network_id $BSC_BMC_NET \
            --argjson gas_limit '{
                "TransferNativeCoinIntraChainGasLimit":5000000,
                "TransferTokenIntraChainGasLimit":5000000,
                "ApproveTokenInterChainGasLimit":5000000,
                "TransferCoinInterChainGasLimit":5000000,
                "TransferBatchCoinInterChainGasLimit":5000000,
                "DefaultGasLimit":5000000
            }'
        )" 
}

#wait-for-it.sh $GOLOOP_RPC_ADMIN_URI
# run provisioning
echo "start..."
if [ ${#linkfrom} != 0 ]; then 
  if [ ! -d ${linkfrom} ]; then
    echo "path ${linkfrom} does not exist. Use absolute path"
    exit 0
  fi
  setup_linked_account
  linked_deploysc
else 
  setup_account
  deploysc
fi