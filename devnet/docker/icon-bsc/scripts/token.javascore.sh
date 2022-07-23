#!/bin/sh
######################################## javascore service methods - start ######################################
source utils.sh
source rpc.sh
source keystore.sh
# Parts of this code is adapted from https://github.com/icon-project/btp/blob/goloop2moonbeam/testnet/goloop2moonbeam/scripts
goloop_lastblock() {
  goloop rpc lastblock
}

deploy_javascore_bmc() {
  echo "deploying javascore BMC"
  cd $CONFIG_DIR
  goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/bmc.jar \
    --content_type application/java \
    --param _net=$(cat net.btp.icon) | jq -r . >tx.icon.bmc
  extract_scoreAddress tx.icon.bmc btp.icon.bmc
  echo "btp://$(cat net.btp.icon)/$(cat btp.icon.bmc)" >btp.icon.btp.address
  btp_icon_block_height=$(goloop_lastblock | jq -r .height)
  echo $btp_icon_block_height > btp.icon.block.height
  echo $(URI=$ICON_ENDPOINT HEIGHT=$(decimal2Hex $(($btp_icon_block_height - 1))) $ICONBRIDGE_BIN_DIR/iconvalidators | jq -r .hash) > btp.icon.validators.hash
}

deploy_javascore_bsr() {
  echo "deploying javascore bsr"
  cd $CONFIG_DIR
  goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/bsr.jar \
    --content_type application/java | jq -r . >tx.icon.bsr
  extract_scoreAddress tx.icon.bsr btp.icon.bsr
} 

deploy_javascore_bts() {
  echo "deploying javascore bts"
  cd $CONFIG_DIR
  local irc2Tradeable_score=$(xxd -p $CONTRACTS_DIR/javascore/irc2Tradeable.jar | tr -d '\n')
  goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/bts.jar \
    --content_type application/java \
    --param _name="ICX" \
    --param _bmc=$(cat btp.icon.bmc) \
    --param _serializedIrc2="$irc2Tradeable_score" | jq -r . > tx.icon.bts
  extract_scoreAddress tx.icon.bts btp.icon.bts
}

deploy_javascore_irc2() {
  echo "deploying javascore IRC2Token TICX"
  cd $CONFIG_DIR
  goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/irc2.jar \
    --content_type application/java \
    --param _name="TICX" \
    --param _symbol="TICX" \
    --param _initialSupply="0x186a0" \
    --param _decimals="0x12" | jq -r . >tx.icon.ticx
  extract_scoreAddress tx.icon.ticx btp.icon.ticx
}

configure_javascore_add_bmc_owner() {
  echo "bmc Add Owner"
  local btp_icon_bmc_owner=$(cat $GOLOOP_RPC_KEY_STORE | jq -r .address)
  cd $CONFIG_DIR
  local is_owner=$(goloop rpc call \
    --to $(cat btp.icon.bmc) \
    --method isOwner \
    --param _addr=$btp_icon_bmc_owner | jq -r .)
  if [ "$is_owner" == "0x0" ]; then
    goloop rpc sendtx call --to $(cat btp.icon.bmc) \
    --method addOwner \
    --param _addr=$btp_icon_bmc_owner > tx/addbmcowner.icon
     ensure_txresult tx/addbmcowner.icon
  fi
}

configure_javascore_bmc_setFeeAggregator() {
  echo "bmc setFeeAggregator"
  cd $CONFIG_DIR
  ensure_key_store fa.ks.json fa.secret
  FA=$(cat fa.ks.json | jq -r .address)
  goloop rpc sendtx call --to $(cat btp.icon.bmc) \
    --method setFeeAggregator \
    --param _addr=${FA} | jq -r . >tx/setFeeAggregator.icon
  ensure_txresult tx/setFeeAggregator.icon

  goloop rpc sendtx call --to $(cat btp.icon.bmc) \
    --method setFeeGatheringTerm \
    --param _value=1000 | jq -r . >tx/setFeeGatheringTerm.icon
  ensure_txresult tx/setFeeGatheringTerm.icon
}

configure_javascore_add_bts() {
  echo "bmc add bts"
  cd $CONFIG_DIR
  goloop rpc sendtx call --to $(cat btp.icon.bmc) \
    --method addService \
    --value 0 \
    --param _addr=$(cat btp.icon.bts) \
    --param _svc=$SVC_NAME | jq -r . >tx/addService.icon
  ensure_txresult tx/addService.icon
}

configure_javascore_add_bts_owner() {
  echo "Add bts Owner"
  local btp_icon_bts_owner=$(cat $GOLOOP_RPC_KEY_STORE | jq -r .address)
  cd $CONFIG_DIR
  local is_owner=$(goloop rpc call \
    --to $(cat btp.icon.bts) \
    --method isOwner \
    --param _addr="$btp_icon_bts_owner" | jq -r .)
  if [ "$is_owner" == "0x0" ]; then
    goloop rpc sendtx call --to $(cat btp.icon.bts) \
    --method addOwner \
    --param _addr=$btp_icon_bts_owner  | jq -r . >tx/addBtsOwner.icon
    ensure_txresult tx/addBtsOwner.icon
  fi
}

configure_javascore_set_bsr() {
  echo "bts set bsr"
  cd $CONFIG_DIR
  goloop rpc sendtx call --to $(cat btp.icon.bts) \
    --method addRestrictor \
    --param _address=$(cat btp.icon.bsr) | jq -r . >tx/addRestrictor.icon
  ensure_txresult tx/addRestrictor.icon
}


configure_javascore_bts_setICXFee() {
  echo "bts set ICX fee"
  local btp_bts_fee_numerator=100
  local btp_bts_fixed_fee=5000
  cd $CONFIG_DIR
  goloop rpc sendtx call --to $(cat btp.icon.bts) \
    --method setFeeRatio \
    --param _name="ICX" \
    --param _feeNumerator=$(decimal2Hex $btp_bts_fee_numerator) \
    --param _fixedFee=$(decimal2Hex $btp_bts_fixed_fee) | jq -r . >tx/setICXFee.icon
  ensure_txresult tx/setICXFee.icon
}

configure_javascore_addLink() {
  echo "BMC: Add Link to BSC BMC:"
  cd $CONFIG_DIR
  goloop rpc sendtx call --to $(cat btp.icon.bmc) \
    --method addLink \
    --param _link=$(cat btp.bsc.btp.address) | jq -r . >tx/addLink.icon
  ensure_txresult tx/addLink.icon

  goloop rpc sendtx call --to $(cat btp.icon.bmc) \
    --method setLinkRxHeight \
    --param _link=$(cat btp.bsc.btp.address) \
    --param _height=$(cat btp.bsc.block.height)| jq -r . >tx/setLinkRxHeight.icon
  ensure_txresult tx/setLinkRxHeight.icon
}

configure_bmc_javascore_addRelay() {
  echo "Adding bsc Relay"
  local btp_icon_bmr_owner=$(cat $GOLOOP_RPC_KEY_STORE | jq -r .address)
  cd $CONFIG_DIR
  goloop rpc sendtx call --to $(cat btp.icon.bmc) \
    --method addRelay \
    --param _link=$(cat btp.bsc.btp.address) \
    --param _addr=${btp_icon_bmr_owner} | jq -r . >tx/addRelay.icon
  ensure_txresult tx/addRelay.icon
}

configure_javascore_register_bnb() {
  echo "Register BNB"
  cd $CONFIG_DIR
  local btp_bts_fee_numerator=100
  local btp_bts_fixed_fee=5000
  goloop rpc sendtx call --to $(cat btp.icon.bts) \
    --method register \
    --param _name="BNB" \
    --param _symbol="BNB" \
    --param _decimals="0x12" \
    --param _feeNumerator=$(decimal2Hex $btp_bts_fee_numerator) \
    --param _fixedFee=$(decimal2Hex $btp_bts_fixed_fee) | jq -r . >tx/register.coin.bnb
  ensure_txresult tx/register.coin.bnb
}

get_btp_icon_bnb() {
  echo "Get BTP Icon BNB Addr"
  cd $CONFIG_DIR
  goloop rpc sendtx call --to $(cat btp.icon.bts) \
    --method "coinAddress" \
    --param _coinName="BNB" | jq -r . >tx/bnb.coin.icon
  ensure_txresult tx/bnb.coin.icon
}

configure_javascore_register_ticx() {
  echo "Register TICX"
  cd $CONFIG_DIR
  local btp_bts_fee_numerator=100
  local btp_bts_fixed_fee=5000
  goloop rpc sendtx call --to $(cat btp.icon.bts) \
    --method register \
    --param _name=TICX \
    --param _symbol=TICX \
    --param _decimals=0x12 \
    --param _addr=$(cat btp.icon.ticx) \
    --param _feeNumerator=$(decimal2Hex $btp_bts_fee_numerator) \
    --param _fixedFee=$(decimal2Hex $btp_bts_fixed_fee) | jq -r . >tx/register.coin.ticx
  ensure_txresult tx/register.coin.ticx
}

configure_javascore_register_tbnb() {
  echo "Register TBNB"
  cd $CONFIG_DIR
  local btp_bts_fee_numerator=100
  local btp_bts_fixed_fee=5000
  goloop rpc sendtx call --to $(cat btp.icon.bts) \
    --method register \
    --param _name=TBNB \
    --param _symbol=TBNB \
    --param _decimals=0x12 \
    --param _feeNumerator=$(decimal2Hex $btp_bts_fee_numerator) \
    --param _fixedFee=$(decimal2Hex $btp_bts_fixed_fee) | jq -r . >tx/register.coin.tbnb
  ensure_txresult tx/register.coin.tbnb
}

get_btp_icon_tbnb() {
  echo "Get BTP Icon TBNB Addr"
  cd $CONFIG_DIR
  goloop rpc sendtx call --to $(cat btp.icon.bts) \
    --method coinAddress \
    --param _coinName=TBNB | jq -r . >tx/tbnb.coin.icon
  ensure_txresult tx/tbnb.coin.icon
}

bsh_javascore_balance() {
  cd $CONFIG_DIR
  if [ $# -lt 1 ]; then
    echo "Usage: bsh_balance [EOA=$(rpceoa)]"
    return 1
  fi

  local EOA=$(rpceoa $1)
  echo "Balance of user $EOA"
  goloop rpc call --to "$(extractAddresses "javascore" "TokenBSH")" \
    --method balanceOf \
    --param _owner=$EOA \
    --param _coinName=$TOKEN_NAME
}

bsh_javascore_transfer() {
  cd $CONFIG_DIR
  if [ $# -lt 2 ]; then
    echo "Usage: bsh_transfer [VAL=0x10] [EOA=$(rpceoa)]"
    return 1
  fi
  local VAL=${1:-0x10}
  local EOA=$2
  local FROM=$(rpceoa $GOLOOP_RPC_KEY_STORE)
  echo "Transfering $VAL wei to: $EOA from: $FROM "
  TX=$(
    goloop rpc sendtx call --to "$(extractAddresses "javascore" "TokenBSH")" \
      --method transfer \
      --param _coinName=${TOKEN_NAME} \
      --param _value=$VAL \
      --param _to=btp://$BSC_BMC_NET/$EOA | jq -r .
  )
  ensure_txresult $TX
}

irc2_javascore_balance() {
  cd $CONFIG_DIR
  if [ $# -lt 1 ]; then
    echo "Usage: irc2_balance [EOA=$(rpceoa)]"
    return 1
  fi
  local EOA=$(rpceoa $1)
  balance=$(goloop rpc call --to "$(extractAddresses "javascore" "IRC2")" \
    --method balanceOf \
    --param _owner=$EOA | jq -r .)
  balance=$(hex2int $balance)
  balance=$(wei2coin $balance)
  echo "Balance: $balance"
}

check_alice_token_balance_with_wait() {
  echo "Checking Alice's balance..."
  cd $CONFIG_DIR
  ALICE_INITIAL_BAL=$(irc2_javascore_balance alice.ks.json)
  COUNTER=60
  while true; do
    printf "."
    if [ $COUNTER -le 0 ]; then
      printf "\nError: timed out while getting Alice's Balance: Balance unchanged\n"
      echo $ALICE_CURR_BAL
      exit 1
    fi
    sleep 3
    COUNTER=$(expr $COUNTER - 3)
    ALICE_CURR_BAL=$(irc2_javascore_balance alice.ks.json)
    if [ "$ALICE_CURR_BAL" != "$ALICE_INITIAL_BAL" ]; then
      printf "\nBTP Transfer Successfull! \n"
      break
    fi
  done
  echo "Alice's Balance after BTP transfer: $ALICE_CURR_BAL ETH"
}

irc2_javascore_transfer() {
  cd $CONFIG_DIR
  if [ $# -lt 1 ]; then
    echo "Usage: irc2_transfer [VAL=0x10] [EOA=Address of Token-BSH]"
    return 1
  fi
  local VAL=${1:-0x10}
  local EOA=$(rpceoa ${2:-"$(extractAddresses "javascore" "TokenBSH")"})
  local FROM=$(rpceoa $GOLOOP_RPC_KEY_STORE)
  echo "Transfering $VAL wei to: $EOA from: $FROM "
  TX=$(
    goloop rpc sendtx call --to "$(extractAddresses "javascore" "IRC2")" \
      --method transfer \
      --param _to=$EOA \
      --param _value=$VAL | jq -r .
  )
  ensure_txresult $TX
}

token_icon_fundBSH() {
  echo "funding BSH with 1000ETH tokens"
  cd $CONFIG_DIR
  weiAmount=$(coin2wei 1000)
  echo "Wei Amount: $weiAmount"
  irc2_javascore_transfer "$weiAmount"
  #echo "$tx" >tx/fundBSH.icon
  #ensure_txresult tx/fundBSH.icon
}

rpceoa() {
  local EOA=${1:-${GOLOOP_RPC_KEY_STORE}}
  if [ "$EOA" != "" ] && [ -f "$EOA" ]; then
    echo $(cat $EOA | jq -r .address)
  else
    echo $EOA
  fi
}

########################################################### javascore service methods - END #####################################################################
