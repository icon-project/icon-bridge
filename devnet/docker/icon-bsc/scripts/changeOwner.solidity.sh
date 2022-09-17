#!/bin/bash

set -e

source config.sh
source rpc.sh

export PRIVATE_KEY="[\""$(cat $BSC_KEY_STORE.priv)"\"]"

# Address of owner to add/remove
ADDR=0xdb23ace5d4cb14682af9fd85feb499f76edaea6b

addOwnerBMC() {
  isOwnerBMC ${1}
  echo "adding ${1} owner to BMC"
  cd $CONTRACTS_DIR/solidity/bmc
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method addOwner --addr "${1}")
  echo "$tx" >$CONFIG_DIR/tx/addOwner.bmc.bsc
  isOwnerBMC ${1}
}

removeOwnerBMC() {
  isOwnerBMC ${1}
  echo "removing ${1} owner from BMC"
  cd $CONTRACTS_DIR/solidity/bmc
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method removeOwner --addr "${1}")
  echo "$tx" >$CONFIG_DIR/tx/removeOwner.bmc.bsc
  isOwnerBMC ${1}
}

isOwnerBMC() {
  cd $CONTRACTS_DIR/solidity/bmc
  truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method isOwner --addr ${1}
}

addOwnerBTS() {
  getOwnerBTS
  isOwnerBTS ${1}
  echo "adding ${1} owner to BTS"
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method addOwner --addr "${1}")
  echo "$tx" >$CONFIG_DIR/tx/addOwner.bts.bsc
  getOwnerBTS
  isOwnerBTS ${1}
}

removeOwnerBTS() {
  getOwnerBTS
  isOwnerBTS ${1}
  echo "removing ${1} owner from BTS"
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method removeOwner --addr "${1}")
  echo "$tx" >$CONFIG_DIR/tx/removeOwner.bts.bsc
  getOwnerBTS
  isOwnerBTS ${1}
}

isOwnerBTS() {
  cd $CONTRACTS_DIR/solidity/bts
  truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method isOwner --addr ${1}
}

getOwnerBTS() {
  cd $CONTRACTS_DIR/solidity/bts
  truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method getOwners
}

if [ $# -eq 0 ]; then
  echo "No arguments supplied: Pass --help for details"
elif [ $1 == "--show-bts" ]; then
  echo "Current Owners of BTS are: "
  getOwnerBTS
elif [ $1 == "--show-bmc" ]; then
  isOwnerBMC ${ADDR}
elif [ $1 == "--add-bts" ]; then
  addOwnerBTS ${ADDR}
elif [ $1 == "--remove-bts" ]; then
  removeOwnerBTS ${ADDR}
elif [ $1 == "--add-bmc" ]; then
  addOwnerBMC ${ADDR}
elif [ $1 == "--remove-bmc" ]; then
  removeOwnerBMC ${ADDR}
else
  echo "Invalid argument: "
  echo "Valid arguments: "
  echo "--show-bmc: Show BMC Owners"
  echo "--show-bts: Show BTS Owners"
  echo "--add-bmc: Add BMC Owner"
  echo "--add-bts: Add BTS Owner"
  echo "--remove-bmc: Remove BMC Owner"
  echo "--remove-bts: Remove BTS Owner"
fi
