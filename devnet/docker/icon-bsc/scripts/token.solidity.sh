#!/bin/sh
source utils.sh
# Parts of this code is adapted from https://github.com/icon-project/btp/blob/goloop2moonbeam/testnet/goloop2moonbeam/scripts
deploy_solidity_bmc() {
  echo "deploying solidity bmc"
  cd $CONTRACTS_DIR/solidity/bmc
  cp $BTPSIMPLE_BASE_DIR/bin/env ./.env
  rm -rf contracts/test build .openzeppelin  
  BMC_BTP_NET=$BSC_BMC_NET \
    truffle migrate --network bsc --compile-all

  generate_metadata "BMC"
}

deploy_solidity_tokenBSH_BEP20() {
  echo "deploying solidity Token BSH"
  cd $CONTRACTS_DIR/solidity/TokenBSH
  cp $BTPSIMPLE_BASE_DIR/bin/env ./.env
  rm -rf contracts/test build .openzeppelin
  #npm install --legacy-peer-deps  
  SVC_NAME=TokenBSH
  
  BSH_TOKEN_FEE=1 \
    BMC_PERIPHERY_ADDRESS=$BMC_ADDRESS \
    BSH_SERVICE=$SVC_NAME \
    truffle migrate --compile-all --network bsc

  generate_metadata "TOKEN_BSH"
}

add_icon_link() {
  echo "adding icon link $(cat $CONFIG_DIR/btp.icon)"
  cd $CONTRACTS_DIR/solidity/bmc
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method addLink --link $(cat $CONFIG_DIR/btp.icon) --blockInterval 3000 --maxAggregation 2 --delayLimit 3)
  echo "$tx" >$CONFIG_DIR/tx/addLink.bsc
}

add_icon_relay() {
  echo "adding icon link $(cat $CONFIG_DIR/btp.icon)"
  BSC_RELAY_USER=$(cat $CONFIG_DIR/bsc.ks.json | jq -r .address)
  cd $CONTRACTS_DIR/solidity/bmc
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method addRelay --link $(cat $CONFIG_DIR/btp.icon) --addr "0x${BSC_RELAY_USER}")
  echo "$tx" >$CONFIG_DIR/tx/addRelay.bsc
}

bsc_addService() {
  echo "adding ${SVC_NAME} service into BMC"
  cd $CONTRACTS_DIR/solidity/bmc
  BSH_IMPL_ADDRESS=$(cat $CONFIG_DIR/token_bsh.impl.bsc)
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method addService --name $SVC_NAME --addr "$BSH_IMPL_ADDRESS")
  echo "$tx" >$CONFIG_DIR/tx/addService.bsc
}

bsc_registerToken() {
  echo "Registering ${TOKEN_NAME} into tokenBSH"
  cd $CONTRACTS_DIR/solidity/bsh
  BEP20_TKN_ADDRESS=$(cat $CONFIG_DIR/bep20_token.bsc)
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bsh.token.js \
    --method registerToken --name $TOKEN_NAME --symbol $TOKEN_SYM --addr "$BEP20_TKN_ADDRESS" --feeNumerator 100 --fixedFee 50000)
  echo "$tx" >$CONFIG_DIR/tx/register.token.bsc
}

bsc_updateRxSeq() {
  cd $CONTRACTS_DIR/solidity/bmc
  truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method updateRxSeq --link $(cat $CONFIG_DIR/btp.icon) --value 1
}

token_bsc_fundBSH() {
  echo "Funding solidity BSH"
  cd $CONTRACTS_DIR/solidity/bsh
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bsh.token.js \
    --method fundBSH --addr $(cat $CONFIG_DIR/token_bsh.proxy.bsc) --amount 1000)
  echo "$tx" >$CONFIG_DIR/tx/fundBSH.bsc
}

deposit_token_for_bob() {
  echo "Funding BOB"
  cd $CONTRACTS_DIR/solidity/bsh
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bsh.token.js \
    --method fundBOB --addr $(get_bob_address) --amount $1)
  echo "$tx" >$CONFIG_DIR/fundBOB.bsc
}

token_approveTransfer() {
  cd $CONTRACTS_DIR/solidity/bsh
  truffle exec --network bsc "$SCRIPTS_DIR"/bsh.token.js \
    --method approve --addr $(cat $CONFIG_DIR/token_bsh.proxy.bsc) --amount $1 --from "$(get_bob_address)"
}

bsc_init_btp_transfer() {
  ICON_NET=$(cat $CONFIG_DIR/net.btp.icon)
  ALICE_ADDRESS=$(get_alice_address)
  BTP_TO="btp://$ICON_NET/$ALICE_ADDRESS"
  cd $CONTRACTS_DIR/solidity/bsh
  truffle exec --network bsc "$SCRIPTS_DIR"/bsh.token.js \
    --method transfer --to $BTP_TO --amount $1 --from "$(get_bob_address)"
}

calculateTransferFee() {
  cd $CONTRACTS_DIR/solidity/TokenBSH
  truffle exec --network bsc "$SCRIPTS_DIR"/bsh.token.js \
    --method calculateTransferFee --amount $1
}

get_Bob_Token_Balance() {
  cd $CONTRACTS_DIR/solidity/bsh
  BSC_USER=$(get_bob_address)
  BOB_BALANCE=$(truffle exec --network bsc "$SCRIPTS_DIR"/bsh.token.js \
    --method getBalance --addr $BSC_USER)
}

get_Bob_Token_Balance_with_wait() {
  echo "Checking Bob's Balance after BTP transfer:"
  get_Bob_Token_Balance
  BOB_INITIAL_BAL=$BOB_BALANCE
  COUNTER=30
  while true; do
    printf "."
    if [ $COUNTER -le 0 ]; then
      printf "\nError: timed out while getting Bob's Balance: Balance unchanged\n"
      echo $BOB_CURRENT_BAL
      exit 1
    fi
    sleep 3
    COUNTER=$(expr $COUNTER - 3)
    get_Bob_Token_Balance
    BOB_CURRENT_BAL=$BOB_BALANCE
    if [ "$BOB_CURRENT_BAL" != "$BOB_INITIAL_BAL" ]; then
      printf "\nBTP Transfer Successfull! \n"
      break
    fi
  done
  echo "Bob's Balance after BTP transfer: $BOB_CURRENT_BAL ETH"
}

generate_metadata() {
  option=$1
  case "$option" in

  BMC)
    echo "###################  Generating BMC Solidity metadata ###################"

    BMC_ADDRESS=$(jq -r '.networks[] | .address' build/contracts/BMCPeriphery.json)
    echo btp://$BSC_BMC_NET/"${BMC_ADDRESS}" >$CONFIG_DIR/btp.bsc
    echo "${BMC_ADDRESS}" >$CONFIG_DIR/bmc.periphery.bsc
    wait_for_file $CONFIG_DIR/bmc.periphery.bsc

    jq -r '.networks[] | .address' build/contracts/BMCManagement.json >$CONFIG_DIR/bmc.bsc
    wait_for_file $CONFIG_DIR/bmc.bsc

    create_abi "BMCPeriphery"
    create_abi "BMCManagement"
    echo "DONE."
    ;;


  TOKEN_BSH)
    echo "################### Generating Token BSH & BEP20  Solidity metadata ###################"

    # BSH_IMPL_ADDRESS=$(jq -r '.networks[] | .address' build/contracts/BSHImpl.json)
    # jq -r '.networks[] | .address' build/contracts/BSHImpl.json >$CONFIG_DIR/token_bsh.impl.bsc
    # jq -r '.networks[] | .address' build/contracts/BSHProxy.json >$CONFIG_DIR/token_bsh.proxy.bsc

    # wait_for_file $CONFIG_DIR/token_bsh.impl.bsc
    # wait_for_file $CONFIG_DIR/token_bsh.proxy.bsc

    jq -r '.networks[] | .address' build/contracts/ERC20TKN.json >$CONFIG_DIR/bep20_token.bsc
    wait_for_file $CONFIG_DIR/bep20_token.bsc

    #create_abi "BSHProxy"
    #create_abi "BSHImpl"
    create_abi "ERC20TKN"
    echo "DONE."
    ;;

  *)
    echo "Invalid option for generating meta data"
    ;;
  esac
}
