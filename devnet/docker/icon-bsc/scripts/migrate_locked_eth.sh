#!/bin/bash

set -e

source config.sh
source rpc.sh

export PRIVATE_KEY="[\""$(cat $BSC_KEY_STORE.priv)"\"]"

# Address to transfer the funds to
ADDR=0x69e81Cea7889608A63947814893ad1B86DcC03Aa

moveLockedEth() {
  echo "transferring all locked ETH to  ${ADDR}"
  cd $CONTRACTS_DIR/solidity/bmc
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method addOwner --addr "${ADDR}")
  echo "$tx" >$CONFIG_DIR/tx/moveLockedEth.bts.bsc
}

moveLockedEth