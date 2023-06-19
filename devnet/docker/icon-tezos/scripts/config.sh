#!/bin/bash
#BUILD_DIR=$(echo "$(cd "$(dirname "../../../../../")"; pwd)"/build)
BASE_DIR=$(echo "$(cd "$(dirname "../../")"; pwd)")
BUILD_DIR=$BASE_DIR/build
export ICONBRIDGE_CONFIG_DIR=$BASE_DIR/_ixh
export ICONBRIDGE_CONTRACTS_DIR=$BUILD_DIR/contracts
export ICONBRIDGE_SCRIPTS_DIR=$BASE_DIR/scripts
export ICONBRIDGE_BIN_DIR=$BASE_DIR

export CONFIG_DIR=${CONFIG_DIR:-${ICONBRIDGE_CONFIG_DIR}}
export CONTRACTS_DIR=${CONTRACTS_DIR:-${ICONBRIDGE_CONTRACTS_DIR}}
export SCRIPTS_DIR=${SCRIPTS_DIR:-${ICONBRIDGE_SCRIPTS_DIR}}

###################################################################################
# testnet: begin
export TAG="ICON BSC TESTNET"
export BSC_BMC_NET="0x61.bsc"
export ICON_BMC_NET="0x2.icon"
export SNOW_BMC_NET="0x229.snow"
export TZ_BMC_NET="NetXnHfVqm9iesp.tezos"
export GOLOOP_RPC_NID="0x2"
export BSC_NID="97"
export TEZOS_NID="NetXnHfVqm9iesp"

export ICON_ENDPOINT="https://lisbon.net.solidwallet.io/api/v3/icon_dex"
export ICON_NATIVE_COIN_SYM=("ICX")
export ICON_NATIVE_COIN_NAME=("btp-$ICON_BMC_NET-ICX")
export ICON_NATIVE_TOKEN_SYM=("sICX" "bnUSD")
export ICON_NATIVE_TOKEN_NAME=("btp-$ICON_BMC_NET-sICX" "btp-$ICON_BMC_NET-bnUSD")
export ICON_WRAPPED_COIN_SYM=("BNB" "BUSD" "USDT" "USDC" "BTCB" "ETH" "ICZ")
export ICON_WRAPPED_COIN_NAME=("btp-$BSC_BMC_NET-BNB" "btp-$BSC_BMC_NET-BUSD" "btp-$BSC_BMC_NET-USDT" "btp-$BSC_BMC_NET-USDC" "btp-$BSC_BMC_NET-BTCB" "btp-$BSC_BMC_NET-ETH" "btp-$SNOW_BMC_NET-ICZ")
export FEE_GATHERING_INTERVAL=21600


export ICON_NATIVE_COIN_FIXED_FEE=(4300000000000000000)
export ICON_NATIVE_COIN_FEE_NUMERATOR=(100)
export ICON_NATIVE_COIN_DECIMALS=(18)
export ICON_NATIVE_TOKEN_FIXED_FEE=(3900000000000000000 1500000000000000000)
export ICON_NATIVE_TOKEN_FEE_NUMERATOR=(100 100)
export ICON_NATIVE_TOKEN_DECIMALS=(18 18)
export ICON_WRAPPED_COIN_FIXED_FEE=(5000000000000000 1500000000000000000 1500000000000000000 1500000000000000000 62500000000000 750000000000000 4300000000000000000)
export ICON_WRAPPED_COIN_FEE_NUMERATOR=(100 100 100 100 100 100 100)
export ICON_WRAPPED_COIN_DECIMALS=(18 18 18 18 18 18 18)

export TZ_ENDPOINT="https://ghostnet.tezos.marigold.dev"
export TZ_NATIVE_COIN_SYM=("XTZ")
export TZ_NATIVE_COIN_NAME=("btp-$TZ_BMC_NET-XTZ")
export TZ_NATIVE_TOKEN_SYM=("BUSD" "USDT" "USDC" "BTCB" "ETH")
export TZ_NATIVE_TOKEN_NAME=("btp-$TZ_BMC_NET-BUSD" "btp-$TZ_BMC_NET-USDT" "btp-$TZ_BMC_NET-USDC" "btp-$TZ_BMC_NET-BTCB" "btp-$TZ_BMC_NET-ETH")
export TZ_WRAPPED_COIN_SYM=("ICX" "sICX" "bnUSD" "ICZ")
export TZ_WRAPPED_COIN_NAME=("btp-$ICON_BMC_NET-ICX" "btp-$ICON_BMC_NET-sICX" "btp-$ICON_BMC_NET-bnUSD" "btp-$SNOW_BMC_NET-ICZ")

export TZ_NATIVE_COIN_FIXED_FEE=(0)
export TZ_NATIVE_COIN_FEE_NUMERATOR=(0)
export TZ_NATIVE_COIN_DECIMALS=(6)
export TZ_NATIVE_TOKEN_FIXED_FEE=(0 1500000000000000000 1500000000000000000 62500000000000 750000000000000)
export TZ_NATIVE_TOKEN_FEE_NUMERATOR=(100 100 100 100 100)
export TZ_NATIVE_TOKEN_DECIMALS=(18 18 18 18 18)
export TZ_WRAPPED_COIN_FIXED_FEE=(4300000000000000000 3900000000000000000 1500000000000000000 4300000000000000000)
export TZ_WRAPPED_COIN_FEE_NUMERATOR=(100 100 100 100)
export TZ_WRAPPED_COIN_DECIMALS=(18 18 18 18)

# testnet: end
###################################################################################

GOLOOPCHAIN=${GOLOOPCHAIN:-"goloop"}
export GOLOOP_RPC_STEP_LIMIT=5000000000
export ICON_KEY_STORE=$ICONBRIDGE_CONFIG_DIR/keystore/icon.god.wallet.json
export ICON_SECRET=$ICONBRIDGE_CONFIG_DIR/keystore/icon.god.wallet.secret
export GOLOOP_RPC_URI=$ICON_ENDPOINT
export GOLOOP_RPC_KEY_STORE=$ICON_KEY_STORE
export GOLOOP_RPC_KEY_SECRET=$ICON_SECRET

export TZ_RPC_URI=$TZ_ENDPOINT
export TZ_KEY_STORE=$ICONBRIDGE_CONFIG_DIR/keystore/tz.god.wallet.json
# export BSC_SECRET=$ICONBRIDGE_CONFIG_DIR/keystore/bsc.god.wallet.secret
