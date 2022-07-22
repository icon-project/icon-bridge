#!/bin/sh
source utils.sh
# Parts of this code is adapted from https://github.com/icon-project/btp/blob/goloop2moonbeam/testnet/goloop2moonbeam/scripts

eth_blocknumber() {
  curl -s -X POST $BSC_RPC_URI --header 'Content-Type: application/json' \
    --data-raw '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[], "id": 1}' | jq -r .result | xargs printf "%d\n"
}

deploy_solidity_bmc() {
  echo "deploying solidity bmc"
  cd $CONTRACTS_DIR/solidity/bmc
  cp $ICONBRIDGE_BASE_DIR/bin/env ./.env
  rm -rf contracts/test build .openzeppelin
  BMC_BTP_NET=$BSC_BMC_NET \
    truffle migrate --network bsc --compile-all
  eth_blocknumber > $CONFIG_DIR/btp.bsc.block.height
  generate_metadata "BMC"
} 



deploy_solidity_bts() {
  echo "deploying solidity bts"
  cd $CONTRACTS_DIR/solidity/bts
  cp $ICONBRIDGE_BASE_DIR/bin/env ./.env
  rm -rf contracts/test build .openzeppelin

  BSH_COIN_NAME="BNB" \
  BSH_COIN_FEE=100 \
  BSH_FIXED_FEE=5000 \
  BMC_PERIPHERY_ADDRESS="$(cat $CONFIG_DIR/btp.bsc.bmc.periphery)" \
  BSH_SERVICE=$SVC_NAME \
  truffle migrate --compile-all --network bsc 
  generate_metadata "BTS"
}

configure_solidity_add_bts_service() {
  echo "adding ${SVC_NAME} service into BMC"
  cd $CONTRACTS_DIR/solidity/bmc
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method addService --name $SVC_NAME --addr $(cat $CONFIG_DIR/btp.bsc.bts.periphery))
  echo "$tx" >$CONFIG_DIR/tx/addService.bsc
}
 
 
configure_solidity_set_fee_ratio() {
  echo "SetFee Ratio"
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method setFeeRatio --name BNB --feeNumerator 100 --fixedFee 5000)
  echo "$tx" >$CONFIG_DIR/tx/setFee.bsc
}
 
add_icon_link() {
  echo "adding icon link $(cat $CONFIG_DIR/btp.icon.btp.address)"
  cd $CONTRACTS_DIR/solidity/bmc
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method addLink --link $(cat $CONFIG_DIR/btp.icon.btp.address) --blockInterval 3000 --maxAggregation 2 --delayLimit 3)
  echo "$tx" >$CONFIG_DIR/tx/addLink.bsc
}

set_link_height() {
  echo "set link height"
  cd $CONTRACTS_DIR/solidity/bmc
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method setLinkRxHeight --link $(cat $CONFIG_DIR/btp.icon.btp.address) --height $(cat $CONFIG_DIR/btp.icon.block.height))
  echo "$tx" >$CONFIG_DIR/tx/setLinkRxHeight.bsc
}

add_icon_relay() {
  echo "adding icon link"
  BSC_RELAY_USER=$(cat $BSC_KEY_STORE | jq -r .address)
  cd $CONTRACTS_DIR/solidity/bmc
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method addRelay --link $(cat $CONFIG_DIR/btp.icon.btp.address) --addr "0x${BSC_RELAY_USER}")
  echo "$tx" >$CONFIG_DIR/tx/addRelay.bsc
} 


bsc_register_icx() {
  echo "bts: Register ICX: "
  local btp_bts_fee_numerator=100
  local btp_bts_fixed_fee=5000
  echo "Registering ICX"
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method register --name "ICX" --symbol "ICX" --decimals 18 --addr "0x0000000000000000000000000000000000000000" --feeNumerator ${btp_bts_fee_numerator} --fixedFee ${btp_bts_fixed_fee})
  echo "$tx" >$CONFIG_DIR/tx/register.icx.bsc
}

get_coinID_icx() {
  echo "getCoinID ICX"
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method coinId --coinName "ICX")
  echo "$tx" >$CONFIG_DIR/tx/coinID.icx
}

bsc_register_tbnb() {
  echo "bts: Register TBNB: "
  local btp_bts_fee_numerator=100
  local btp_bts_fixed_fee=5000
  local addr=$(cat $CONFIG_DIR/btp.bsc.tbnb) 
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method register --name "TBNB" --symbol "TBNB" --decimals 18 --addr $addr --feeNumerator $btp_bts_fee_numerator --fixedFee ${btp_bts_fixed_fee})
  echo "$tx" >$CONFIG_DIR/tx/register.tbnb.bsc
}

get_coinID_tbnb() {
  echo "getCoinID TBNB"
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method coinId --coinName "TBNB")
  echo "$tx" >$CONFIG_DIR/tx/coinID.tbnb
}

bsc_register_eth() {
  echo "bts: Register ETH: "
  local btp_bts_fee_numerator=100
  local btp_bts_fixed_fee=5000
  local addr=$(cat $CONFIG_DIR/btp.bsc.eth) 
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method register --name "ETH" --symbol "ETH" --decimals 18 --addr $addr --feeNumerator $btp_bts_fee_numerator --fixedFee ${btp_bts_fixed_fee})
  echo "$tx" >$CONFIG_DIR/tx/register.eth.bsc
}

get_coinID_eth() {
  echo "getCoinID ETH"
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method coinId --coinName "ETH")
  echo "$tx" >$CONFIG_DIR/tx/coinID.eth
}

bsc_register_ticx() {
  local btp_bts_fee_numerator=100
  local btp_bts_fixed_fee=5000
  echo "Registering TICX"
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method register --name "TICX" --symbol "TICX" --decimals 18 --addr "0x0000000000000000000000000000000000000000" --feeNumerator ${btp_bts_fee_numerator} --fixedFee ${btp_bts_fixed_fee})
  echo "$tx" >$CONFIG_DIR/tx/register.ticx.bsc
}

get_coinID_ticx() {
  echo "getCoinID TICX"
  cd $CONTRACTS_DIR/solidity/bts
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bts.js \
    --method coinId --coinName "TICX")
  echo "$tx" >$CONFIG_DIR/tx/coinID.ticx
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

    local BMC_ADDRESS=$(jq -r '.networks[] | .address' build/contracts/BMCPeriphery.json)
    echo btp://$BSC_BMC_NET/"${BMC_ADDRESS}" >$CONFIG_DIR/btp.bsc.btp.address
    echo "${BMC_ADDRESS}" >$CONFIG_DIR/btp.bsc.bmc.periphery
    wait_for_file $CONFIG_DIR/btp.bsc.bmc.periphery

    jq -r '.networks[] | .address' build/contracts/BMCManagement.json >$CONFIG_DIR/btp.bsc.bmc.management
    wait_for_file $CONFIG_DIR/btp.bsc.bmc.management
    echo "DONE."
    ;;

  \
    BTS)
    echo "################### Generating BTS  Solidity metadata ###################"

    jq -r '.networks[] | .address' build/contracts/BTSCore.json >$CONFIG_DIR/btp.bsc.bts.core
    wait_for_file $CONFIG_DIR/btp.bsc.bts.core
    jq -r '.networks[] | .address' build/contracts/BTSPeriphery.json >$CONFIG_DIR/btp.bsc.bts.periphery
    wait_for_file $CONFIG_DIR/btp.bsc.bts.periphery
    jq -r '.networks[] | .address' build/contracts/HRC20.json >$CONFIG_DIR/btp.bsc.tbnb
    wait_for_file $CONFIG_DIR/btp.bsc.tbnb
    jq -r '.networks[] | .address' build/contracts/ERC20TKN.json >$CONFIG_DIR/btp.bsc.eth
    wait_for_file $CONFIG_DIR/btp.bsc.eth

    echo "DONE."
    ;;
  *)
    echo "Invalid option for generating meta data"
    ;;
  esac
}