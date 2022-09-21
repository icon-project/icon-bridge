#!/bin/bash

source config.sh
source utils.sh

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
            --arg bts_owner_keystore_path $(echo $CONFIG_DIR/keystore/icon.bts.wallet.json) \
            --arg bts_owner_secret_path $(echo $CONFIG_DIR/keystore/icon.bts.wallet.secret) \
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
            .name = "BSC" |
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
            --arg bts_owner_keystore_path $(echo $CONFIG_DIR/keystore/bsc.bts.wallet.json) \
            --arg bts_owner_secret_path $(echo $CONFIG_DIR/keystore/bsc.bts.wallet.secret) \
            --arg network_id $BSC_BMC_NET \
            --argjson gas_limit '{
                "TransferNativeCoinIntraChainGasLimit":25000,
                "TransferTokenIntraChainGasLimit":60000,
                "ApproveTokenInterChainGasLimit":50000,
                "TransferCoinInterChainGasLimit":700000,
                "TransferBatchCoinInterChainGasLimit":900000,
                "DefaultGasLimit":5000000
            }'
        )" 
}
