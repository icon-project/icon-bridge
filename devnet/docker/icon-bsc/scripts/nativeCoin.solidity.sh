#!/bin/sh

source rpc.sh
source utils.sh
# Parts of this code is adapted from https://github.com/icon-project/btp/blob/goloop2moonbeam/testnet/goloop2moonbeam/scripts

deploy_solidity_nativeCoin_BSH() {
  echo "deploying solidity Native BSH"
  cd $CONTRACTS_DIR/solidity/bsh
  cp $ICONBRIDGE_BASE_DIR/bin/env ./.env
  rm -rf contracts/test build .openzeppelin
  NODE_ENV=docker BSH_COIN_URL=https://ethereum.org/en/ \
    BSH_COIN_NAME=BNB \
    BSH_COIN_FEE=100 \
    BSH_FIXED_FEE=50000 \
    BMC_PERIPHERY_ADDRESS=$(cat $CONFIG_DIR/bmc.periphery.bsc) \
    BSH_SERVICE=nativecoin \
    truffle migrate --compile-all --network bsc

  generate_native_metadata "BSH"
}

bmc_solidity_addNativeService() {
  echo "adding ${SVC_NAME} service into BMC"
  cd $CONTRACTS_DIR/solidity/bmc
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bmc.js \
    --method addService --name nativecoin --addr "$BSH_PERIPHERY_ADDRESS")
  echo "$tx" >$CONFIG_DIR/tx/addService.native.bsc
}

nativeBSH_solidity_register() {
  echo "Register Coin Name with NativeBSH"
  cd $CONTRACTS_DIR/solidity/bsh
  tx=$(truffle exec --network bsc "$SCRIPTS_DIR"/bsh.nativeCoin.js \
    --method register --name "ICX" --symbol "ICX" --decimals 18 --feeNumerator 100 --fixedFee 50000)
  echo "$tx" >$CONFIG_DIR/tx/register.nativeCoin.bsc
}

bsc_init_native_btp_transfer() {
  ICON_NET=$(cat $CONFIG_DIR/net.btp.icon)
  ALICE_ADDRESS=$(get_alice_address)
  BTP_TO="btp://$ICON_NET/$ALICE_ADDRESS"
  cd $CONTRACTS_DIR/solidity/bsh
  truffle exec --network bsc "$SCRIPTS_DIR"/bsh.nativeCoin.js \
    --method transferNativeCoin --to $BTP_TO --amount $1 --from $(get_bob_address)
}

bsc_init_wrapped_native_btp_transfer() {
  ICON_NET=$(cat $CONFIG_DIR/net.btp.icon)
  ALICE_ADDRESS=$(get_alice_address)
  BTP_TO="btp://$ICON_NET/$ALICE_ADDRESS"
  cd $CONTRACTS_DIR/solidity/bsh
  truffle exec --network bsc "$SCRIPTS_DIR"/bsh.nativeCoin.js \
    --method transferWrappedNativeCoin --to $BTP_TO --coinName $1 --amount $2 --from $(get_bob_address)
}

get_bob_ICX_balance() {
  cd $CONTRACTS_DIR/solidity/bsh
  BOB_BALANCE=$(truffle exec --network bsc "$SCRIPTS_DIR"/bsh.nativeCoin.js \
    --method getBalanceOf --addr $(get_bob_address) --name "ICX")
}

get_Bob_ICX_Balance_with_wait() {
  printf "Checking Bob's Balance after BTP transfer \n"
  get_bob_ICX_balance
  BOB_INITIAL_BAL=$BOB_BALANCE
  COUNTER=30
  while true; do
    printf "."
    if [ $COUNTER -le 0 ]; then
      printf "\n Error: timed out while getting Bob's Balance: Balance unchanged \n"
      echo "$BOB_CURRENT_BAL"
      exit 1
    fi
    sleep 3
    COUNTER=$(expr $COUNTER - 3)
    get_bob_ICX_balance
    BOB_CURRENT_BAL=$BOB_BALANCE
    if [ "$BOB_CURRENT_BAL" != "$BOB_INITIAL_BAL" ]; then
      printf "\n BTP Native Transfer Successfull! \n"
      break
    fi
  done
  echo "Bob's Balance after BTP Native transfer: $BOB_CURRENT_BAL"
}

get_bob_BNB_balance() {
  cd $CONTRACTS_DIR/solidity/bsh
  BOB_BNB_BALANCE=$(truffle exec --network bsc "$SCRIPTS_DIR"/bsh.nativeCoin.js \
    --method getBalanceOf --addr $(get_bob_address) --name "BNB")
}

get_Bob_BNB_Balance_with_wait() {
  printf "Checking Bob's Balance after BTP transfer \n"
  get_bob_BNB_balance
  BOB_INITIAL_BAL=$BOB_BNB_BALANCE
  COUNTER=60
  while true; do
    printf "."
    if [ $COUNTER -le 0 ]; then
      printf "\n Error: timed out while getting Bob's Balance: Balance unchanged \n"
      echo "$BOB_CURRENT_BAL"
      exit 1
    fi
    sleep 3
    COUNTER=$(expr $COUNTER - 3)
    get_bob_BNB_balance
    BOB_CURRENT_BAL=$BOB_BNB_BALANCE
    if [ "$BOB_CURRENT_BAL" != "$BOB_INITIAL_BAL" ]; then
      printf "\n BTP Native Transfer Successfull! \n"
      break
    fi
  done
  echo "Bob's Balance after BTP Native transfer: $BOB_CURRENT_BAL"
}

generate_native_metadata() {
  option=$1
  case "$option" in
  BSH)
    echo "################### Generating Native BSH Solidity metadata ###################"

    BSH_PERIPHERY_ADDRESS=$(jq -r '.networks[] | .address' build/contracts/BSHPeriphery.json)
    jq -r '.networks[] | .address' build/contracts/BSHCore.json >$CONFIG_DIR/bsh.core.bsc
    jq -r '.networks[] | .address' build/contracts/BSHPeriphery.json >$CONFIG_DIR/bsh.periphery.bsc

    wait_for_file $CONFIG_DIR/bsh.core.bsc
    wait_for_file $CONFIG_DIR/bsh.periphery.bsc

    create_abi "BSHPeriphery"
    create_abi "BSHCore"
    echo "DONE."
    ;;

  *)
    echo "Invalid option for generating meta data"
    ;;
  esac
}
