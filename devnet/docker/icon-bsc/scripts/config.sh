#!/bin/sh
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
export BSC_BMC_NET=${BSC_BMC_NET:-"0x61.bsc"}
export ICON_BMC_NET=${ICON_BMC_NET:-"0x2.icon"}
export GOLOOP_RPC_NID="0x2"
export BSC_NID="97"

export ICON_ENDPOINT="https://lisbon.net.solidwallet.io/api/v3"
export ICON_NATIVE_COIN_SYM=("ICX")
export ICON_NATIVE_COIN_NAME=("btp-$ICON_BMC_NET-ICX")
export ICON_NATIVE_TOKEN_SYM=("sICX" "bnUSD")
export ICON_NATIVE_TOKEN_NAME=("btp-$ICON_BMC_NET-sICX" "btp-$ICON_BMC_NET-bnUSD")
export ICON_WRAPPED_COIN_SYM=("BNB" "BUSD" "USDT" "USDC" "BTCB" "ETH")
export ICON_WRAPPED_COIN_NAME=("btp-$BSC_BMC_NET-BNB" "btp-$BSC_BMC_NET-BUSD" "btp-$BSC_BMC_NET-USDT" "btp-$BSC_BMC_NET-USDC" "btp-$BSC_BMC_NET-BTCB" "btp-$BSC_BMC_NET-ETH")
export FEE_GATHERING_INTERVAL=1000

export BSC_ENDPOINT="https://data-seed-prebsc-1-s1.binance.org:8545"
export BSC_NATIVE_COIN_SYM=("BNB")
export BSC_NATIVE_COIN_NAME=("btp-$BSC_BMC_NET-BNB")
export BSC_NATIVE_TOKEN_SYM=("BUSD" "USDT" "USDC" "BTCB" "ETH")
export BSC_NATIVE_TOKEN_NAME=("btp-$BSC_BMC_NET-BUSD" "btp-$BSC_BMC_NET-USDT" "btp-$BSC_BMC_NET-USDC" "btp-$BSC_BMC_NET-BTCB" "btp-$BSC_BMC_NET-ETH")
export BSC_WRAPPED_COIN_SYM=("ICX" "sICX" "bnUSD")
export BSC_WRAPPED_COIN_NAME=("btp-$ICON_BMC_NET-ICX" "btp-$ICON_BMC_NET-sICX" "btp-$ICON_BMC_NET-bnUSD")
# testnet: end

# # mainnet: begin
# export BSC_BMC_NET=${BSC_BMC_NET:-"0x38.bsc"}
# export ICON_BMC_NET=${ICON_BMC_NET:-"0x1.icon"}
# export GOLOOP_RPC_NID="0x1"
# export BSC_NID="56"

# export ICON_ENDPOINT="https://ctz.solidwallet.io/api/v3"
# export ICON_NATIVE_COIN_SYM=("ICX")
# export ICON_NATIVE_COIN_NAME=("btp-$ICON_BMC_NET-ICX")
# export ICON_NATIVE_TOKEN_SYM=("sICX" "bnUSD")
# export ICON_NATIVE_TOKEN_NAME=("btp-$ICON_BMC_NET-sICX" "btp-$ICON_BMC_NET-bnUSD")
# export ICON_WRAPPED_COIN_SYM=("BNB" "BUSD" "USDT" "USDC" "BTCB" "ETH")
# export ICON_WRAPPED_COIN_NAME=("btp-$BSC_BMC_NET-BNB" "btp-$BSC_BMC_NET-BUSD" "btp-$BSC_BMC_NET-USDT" "btp-$BSC_BMC_NET-USDC" "btp-$BSC_BMC_NET-BTCB" "btp-$BSC_BMC_NET-ETH")
# export FEE_GATHERING_INTERVAL=1000

# export BSC_ENDPOINT="https://bsc-dataseed.binance.org"
# export BSC_NATIVE_COIN_SYM=("BNB")
# export BSC_NATIVE_COIN_NAME=("btp-$BSC_BMC_NET-BNB")
# export BSC_NATIVE_TOKEN_SYM=("BUSD" "USDT" "USDC" "BTCB" "ETH")
# export BSC_NATIVE_TOKEN_NAME=("btp-$BSC_BMC_NET-BUSD" "btp-$BSC_BMC_NET-USDT" "btp-$BSC_BMC_NET-USDC" "btp-$BSC_BMC_NET-BTCB" "btp-$BSC_BMC_NET-ETH")
# export BSC_WRAPPED_COIN_SYM=("ICX" "sICX" "bnUSD")
# export BSC_WRAPPED_COIN_NAME=("btp-$ICON_BMC_NET-ICX" "btp-$ICON_BMC_NET-sICX" "btp-$ICON_BMC_NET-bnUSD")

# export INIT_ADDRESS_PATH=$BASE_DIR/mainnet/addresses.json
# # mainnet: end


GOLOOPCHAIN=${GOLOOPCHAIN:-"goloop"}
export GOLOOP_RPC_STEP_LIMIT=5000000000
export ICON_KEY_STORE=$ICONBRIDGE_CONFIG_DIR/keystore/icon.god.wallet.json
export ICON_SECRET=$ICONBRIDGE_CONFIG_DIR/keystore/icon.god.wallet.secret
export GOLOOP_RPC_URI=$ICON_ENDPOINT
export GOLOOP_RPC_KEY_STORE=$ICON_KEY_STORE
export GOLOOP_RPC_KEY_SECRET=$ICON_SECRET

export BSC_RPC_URI=$BSC_ENDPOINT
export BSC_KEY_STORE=$ICONBRIDGE_CONFIG_DIR/keystore/bsc.god.wallet.json
export BSC_SECRET=$ICONBRIDGE_CONFIG_DIR/keystore/bsc.god.wallet.secret