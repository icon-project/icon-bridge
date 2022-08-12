#!/bin/sh
BUILD_DIR=$(echo "$(cd "$(dirname "../../../../../")"; pwd)"/build)
BASE_DIR=$(echo "$(cd "$(dirname "../../")"; pwd)")

export ICONBRIDGE_CONFIG_DIR=$BASE_DIR/_ixh
export ICONBRIDGE_CONTRACTS_DIR=$BUILD_DIR/contracts
export ICONBRIDGE_SCRIPTS_DIR=$BASE_DIR/scripts
export ICONBRIDGE_BIN_DIR=$BASE_DIR

export CONFIG_DIR=${CONFIG_DIR:-${ICONBRIDGE_CONFIG_DIR}}
export CONTRACTS_DIR=${CONTRACTS_DIR:-${ICONBRIDGE_CONTRACTS_DIR}}
export SCRIPTS_DIR=${SCRIPTS_DIR:-${ICONBRIDGE_SCRIPTS_DIR}}

###################################################################################

export ICON_ENDPOINT='http://localhost:9080/api/v3/default'
#'https://lisbon.net.solidwallet.io/api/v3/icon_dex'
export ICON_KEY_STORE=$ICONBRIDGE_CONFIG_DIR/keystore/icon.god.wallet.json
export ICON_SECRET=$ICONBRIDGE_CONFIG_DIR/keystore/icon.god.wallet.secret
export ICON_NATIVE_COIN=('ICX')
export ICON_NATIVE_TOKEN_SYM=('sICX' 'bnUSD')
export ICON_NATIVE_TOKEN_NAME=('Staked ICX' 'Balanced Dollar')
export ICON_WRAPPED_COIN_SYM=('BNB' 'BUSD' 'USDT' 'USDC' 'BTCB' 'ETH')
export ICON_WRAPPED_COIN_NAME=('Binance Coin' 'BUSD Token' 'Tether USD' 'USD Coin' 'BTCB Token' 'Ethereum Token')


export GOLOOP_RPC_STEP_LIMIT=5000000000
export GOLOOP_RPC_NID='0x5b9a77'
GOLOOPCHAIN=${GOLOOPCHAIN:-'goloop'}
export GOLOOP_RPC_URI=$ICON_ENDPOINT
export GOLOOP_RPC_KEY_STORE=$ICON_KEY_STORE
export GOLOOP_RPC_KEY_SECRET=$ICON_SECRET

###################################################################################

export BSC_ENDPOINT='http://localhost:8545'
#'https://data-seed-prebsc-1-s1.binance.org:8545'
export BSC_RPC_URI=$BSC_ENDPOINT
export BSC_KEY_STORE=$ICONBRIDGE_CONFIG_DIR/keystore/bsc.god.wallet.json
export BSC_SECRET=$ICONBRIDGE_CONFIG_DIR/keystore/bsc.god.wallet.secret
export BSC_NID=${BSC_NID:-'97'}
export BSC_BMC_NET=${BSC_BMC_NET:-'0x61.bsc'}

export BSC_NATIVE_COIN=('BNB')
export BSC_NATIVE_TOKEN_SYM=('BUSD' 'USDT' 'USDC' 'BTCB' 'ETH')
export BSC_NATIVE_TOKEN_NAME=('BUSD Token' 'Tether USD' 'USD Coin' 'BTCB Token' 'Ethereum Token')
export BSC_WRAPPED_COIN_SYM=('ICX' 'sICX' 'bnUSD')
export BSC_WRAPPED_COIN_NAME=('ICON Coin' 'Staked ICX' 'Balanced Dollars')

###################################################################################
export INIT_ADDRESS_PATH=$ICONBRIDGE_CONFIG_DIR/init_address.json

