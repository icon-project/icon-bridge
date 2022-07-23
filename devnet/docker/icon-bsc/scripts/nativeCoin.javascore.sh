#!/bin/sh

# Parts of this code is adapted from https://github.com/icon-project/btp/blob/goloop2moonbeam/testnet/goloop2moonbeam/scripts
source env.variables.sh
source rpc.sh
source /iconbridge/bin/keystore.sh
source utils.sh
ensure_key_store alice.ks.json alice.secret

deploy_javascore_nativeCoin_BSH() {
  echo "deploying javascore Native coin BSH"
  cd $CONFIG_DIR
  IRC2_SERIALIZED_SCORE=$(xxd -p $CONTRACTS_DIR/javascore/irc2Tradeable-optimized.jar | tr -d '\n')
  goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/nativecoin-optimized.jar \
    --content_type application/java \
    --param _bmc=$(cat bmc.icon) \
    --param _serializedIrc2=$IRC2_SERIALIZED_SCORE \
    --param _name=ICX | jq -r . >tx.nativebsh.icon
  extract_scoreAddress tx.nativebsh.icon nativebsh.icon
  extract_scoreAddress tx.nativebsh.icon token_bsh.icon
}

bmc_javascore_addNativeService() {
  echo "Adding NativeCoin service into BMC"
  cd $CONFIG_DIR
  goloop rpc sendtx call --to $(cat bmc.icon) \
    --method addService \
    --param _svc=nativecoin \
    --param _addr=$(cat nativebsh.icon) | jq -r . >tx/addService.native.icon
  ensure_txresult tx/addService.native.icon
}

nativeBSH_javascore_register() {
  echo "Register Coin Name with NativeBSH"
  cd $CONFIG_DIR
  FEE_NUMERATOR=0x64
  FIXED_FEE=0x1388
  goloop rpc sendtx call --to $(cat nativebsh.icon) \
    --method register \
    --param _name=BNB \
    --param _symbol=BNB \
    --param _decimals=18 \
    --param _feeNumerator=${FEE_NUMERATOR} \
    --param _fixedFee=${FIXED_FEE} | jq -r . >tx/register.nativeCoin.icon
  ensure_txresult tx/register.nativeCoin.icon

  goloop rpc call --to $(cat nativebsh.icon) \
    --method coinAddress --param _coinName=BNB | sed -e 's/^"//' -e 's/"$//' >irc2TradeableToken.icon
}


nativeBSH_javascore_register_token() {
  echo "Register Coin Name with NativeBSH"
  cd $CONFIG_DIR
  FEE_NUMERATOR=0x64
  FIXED_FEE=0x1388
  goloop rpc sendtx call --to $(cat nativebsh.icon) \
    --method register \
    --param _addr=$(cat irc2_token.icon) \
    --param _name=${TOKEN_NAME} \
    --param _symbol=${TOKEN_SYM} \
    --param _decimals=${TOKEN_DECIMALS}  \
    --param _feeNumerator=${FEE_NUMERATOR} \
    --param _fixedFee=${FIXED_FEE} | jq -r . >tx/register.token.icon
  ensure_txresult tx/register.token.icon  
}


nativeBSH_javascore_setFeeRatio() {
  echo "Setting Fee ratio for NativeCoin"
  cd $CONFIG_DIR
  FEE_NUMERATOR=0x64
  FIXED_FEE=0x1388
  goloop rpc sendtx call --to $(cat nativebsh.icon) \
    --method setFeeRatio \
    --param _name=ICX \
    --param _feeNumerator=${FEE_NUMERATOR} \
    --param _fixedFee=${FIXED_FEE} | jq -r . >tx/setFeeRatio.nativebsh.icon
  ensure_txresult tx/setFeeRatio.nativebsh.icon
}

configure_javascore_NativeBSH_restrictor() {
  echo "configuring javascore Restrictor for TokenBSH"
  cd $CONFIG_DIR
  goloop rpc sendtx call --to $(cat nativebsh.icon) \
    --method addRestrictor \
    --param _address=$(cat restrictor.icon) | jq -r . >tx.configure.addRestrictor.nativebsh.icon
  ensure_txresult tx.configure.addRestrictor.nativebsh.icon

  weiAmount=$(coin2wei 10000)
  goloop rpc sendtx call --to $(cat restrictor.icon) \
    --method registerTokenLimit \
    --param _name=BNB \
    --param _symbol=BNB \
    --param _address=$(cat irc2TradeableToken.icon) \
    --param _limit=$weiAmount | jq -r . >tx.configure.registerTokenLimit.nativebsh.icon
  ensure_txresult tx.configure.registerTokenLimit.nativebsh.icon

  weiAmount=$(coin2wei 10000)
  goloop rpc sendtx call --to $(cat restrictor.icon) \
    --method registerTokenLimit \
    --param _name=ICX \
    --param _symbol=ICX \
    --param _address=$(cat nativebsh.icon) \
    --param _limit=$weiAmount | jq -r . >tx.configure.registerTokenLimit2.nativebsh.icon
  ensure_txresult tx.configure.registerTokenLimit2.nativebsh.icon
}

deposit_ICX_for_Alice() {
  get_alice_balance
  echo "Depositing $(wei2coin $ICX_DEPOSIT_AMOUNT) ICX to Alice"
  cd ${CONFIG_DIR}
  goloop rpc sendtx transfer \
    --to $(get_alice_address) \
    --value $ICX_DEPOSIT_AMOUNT | jq -r . >tx/deposit.alice
  ensure_txresult tx/deposit.alice
}

transfer_ICX_from_Alice_to_Bob() {
  ICX_TRANSER_AMOUNT=$1
  echo "Transfer $(wei2coin $ICX_TRANSER_AMOUNT) ICX from Alice to Bob"
  cd ${CONFIG_DIR}
  LAST_BOCK=$(goloop_lastblock)
  LAST_HEIGHT=$(echo $LAST_BOCK | jq -r .height)
  LAST_HASH=0x$(echo $LAST_BOCK | jq -r .block_hash)
  echo "goloop height:$LAST_HEIGHT hash:$LAST_HASH"
  echo "$(get_bob_address)"
  echo "$BSC_BMC_NET,$ICX_TRANSER_AMOUNT "
  goloop rpc sendtx call \
    --to "$(extractAddresses "javascore" "NativeBSH")" --method transferNativeCoin \
    --param _to=btp://$BSC_BMC_NET/$(get_bob_address) --value $ICX_TRANSER_AMOUNT \
    --key_store alice.ks.json --key_secret alice.secret |
    jq -r . >tx/Alice2Bob.transfer
  ensure_txresult tx/Alice2Bob.transfer
}

transfer_BNB_from_Alice_to_Bob() {
  BNB_TRANSER_AMOUNT=$1
  echo "Transfer $(wei2coin $BNB_TRANSER_AMOUNT) BNB from Alice to Bob"
  cd ${CONFIG_DIR}

  goloop rpc sendtx call \
    --to "$(extractAddresses "javascore" "BNB")" --method approve \
    --param spender="$(extractAddresses "javascore" "NativeBSH")" \
    --param amount=$BNB_TRANSER_AMOUNT \
    --key_store alice.ks.json --key_secret alice.secret |
    jq -r . >tx/Alice2Bob.approve.BNB

  goloop rpc sendtx call \
    --to "$(extractAddresses "javascore" "NativeBSH")" --method transfer \
    --param _coinName="BNB" \
    --param _to=btp://$BSC_BMC_NET/$(get_bob_address) \
    --param _value=$BNB_TRANSER_AMOUNT \
    --key_store alice.ks.json --key_secret alice.secret |
    jq -r . >tx/Alice2Bob.transfer.BNB
  ensure_txresult tx/Alice2Bob.transfer.BNB
}

get_alice_balance() {
  balance=$(goloop rpc balance $(get_alice_address) | jq -r)
  balance=$(hex2int $balance)
  balance=$(wei2coin $balance)
  echo "Alice's balance: $balance (ICX)"
}

get_alice_wrapped_native_balance() {
  cd $CONFIG_DIR

  local EOA=$(rpceoa alice.ks.json)

  balance=$(goloop rpc call --to $(extractAddresses "javascore" "BNB") \
    --method balanceOf \
    --param _owner=$EOA | jq -r)

  balance=$(hex2int $balance)
  balance=$(wei2coin $balance)
  echo "$balance ($1)"
}

check_alice_wrapped_native_balance_with_wait() {
  echo "Checking Alice's balance..."

  cd $CONFIG_DIR
  ALICE_INITIAL_BAL=$(get_alice_wrapped_native_balance $1)
  COUNTER=30
  while true; do
    printf "."
    if [ $COUNTER -le 0 ]; then
      printf "\nError: timed out while getting Alice's Balance: Balance unchanged\n"
      echo $ALICE_CURR_BAL
      exit 1
    fi
    sleep 3
    COUNTER=$(expr $COUNTER - 3)
    ALICE_CURR_BAL=$(get_alice_wrapped_native_balance $1)
    if [ "$ALICE_CURR_BAL" != "$ALICE_INITIAL_BAL" ]; then
      printf "\nBTP Transfer Successfull! \n"
      break
    fi
  done
  echo "Alice's Balance after BTP transfer: $ALICE_CURR_BAL"
}

check_alice_native_balance_with_wait() {
  echo "Checking Alice's balance..."

  cd $CONFIG_DIR
  ALICE_INITIAL_BAL=$(get_alice_balance)
  COUNTER=30
  while true; do
    printf "."
    if [ $COUNTER -le 0 ]; then
      printf "\nError: timed out while getting Alice's Balance: Balance unchanged\n"
      echo $ALICE_CURR_BAL
      exit 1
    fi
    sleep 3
    COUNTER=$(expr $COUNTER - 3)
    ALICE_CURR_BAL=$(get_alice_balance)
    if [ "$ALICE_CURR_BAL" != "$ALICE_INITIAL_BAL" ]; then
      printf "\nBTP Transfer Successfull! \n"
      break
    fi
  done
  echo "Alice's Balance after BTP transfer: $ALICE_CURR_BAL"
}

goloop_lastblock() {
  goloop rpc lastblock
}
