#!/bin/sh
######################################## javascore service methods - start ######################################
source utils.sh
source rpc.sh
source keystore.sh
goloop_lastblock() {
  goloop rpc lastblock
}

deploy_javascore_bmc() {
  cd $CONFIG_DIR

  if [ ! -f icon.addr.bmcbtp ]; then
    echo "deploying javascore BMC"
    icon_block_height=$(goloop_lastblock | jq -r .height)
    echo $icon_block_height > icon.chain.height
    echo $(URI=$ICON_ENDPOINT HEIGHT=$(decimal2Hex $(($icon_block_height - 1))) $ICONBRIDGE_BIN_DIR/iconvalidators | jq -r .hash) > icon.chain.validators
    goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/bmc.jar \
      --content_type application/java \
      --param _net=$(cat net.btp.icon) | jq -r . >tx/tx.icon.bmc
    sleep 5
    extract_scoreAddress tx/tx.icon.bmc icon.addr.bmc
    echo "btp://$(cat net.btp.icon)/$(cat icon.addr.bmc)" >icon.addr.bmcbtp
  fi
}

deploy_javascore_bts() {
  echo "deploying javascore bts"
  cd $CONFIG_DIR
  if [ ! -f icon.addr.bts ]; then
    local bts_fee_numerator=100
    local bts_fixed_fee=5000
    goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/bts.jar \
      --content_type application/java \
      --param _name="${ICON_NATIVE_COIN_NAME}" \
      --param _bmc=$(cat icon.addr.bmc) \
      --param _decimals="0x12" \
      --param _feeNumerator=$(decimal2Hex $bts_fee_numerator) \
      --param _fixedFee=$(decimal2Hex $bts_fixed_fee) \
      --param _serializedIrc2=$(xxd -p $CONTRACTS_DIR/javascore/irc2Tradeable.jar | tr -d '\n') | jq -r . > tx/tx.icon.bts
    sleep 5
    extract_scoreAddress tx/tx.icon.bts icon.addr.bts
  fi
}

deploy_javascore_token() {
  echo "deploying javascore IRC2Token " $2
  cd $CONFIG_DIR
  if [ ! -f icon.addr.$2 ]; then
    goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/irc2.jar \
      --content_type application/java \
      --param _name="$1" \
      --param _symbol=$2 \
      --param _initialSupply="0x186a0" \
      --param _decimals="0x12" | jq -r . >tx/tx.icon.$2
    sleep 5
    extract_scoreAddress tx/tx.icon.$2 icon.addr.$2
  fi
}

 
configure_javascore_add_bmc_owner() {
  echo "bmc Add Owner"
  echo $CONFIG_DIR/keystore/icon.bmc.wallet.json
  local icon_bmc_owner=$(cat $CONFIG_DIR/keystore/icon.bmc.wallet.json | jq -r .address)
  cd $CONFIG_DIR
  local is_owner=$(goloop rpc call \
    --to $(cat icon.addr.bmc) \
    --method isOwner \
    --param _addr=$icon_bmc_owner | jq -r .)
  if [ "$is_owner" == "0x0" ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
    --method addOwner \
    --param _addr=$icon_bmc_owner | jq -r . > tx/addbmcowner.icon
    sleep 3
    ensure_txresult tx/addbmcowner.icon
  fi
}

configure_javascore_bmc_setFeeAggregator() {
  echo "bmc setFeeAggregator"
  cd $CONFIG_DIR
  FA=$(cat $CONFIG_DIR/keystore/icon.fa.wallet.json | jq -r .address)
  goloop rpc sendtx call --to $(cat icon.addr.bmc) \
    --method setFeeAggregator \
    --param _addr=${FA} | jq -r . >tx/setFeeAggregator.icon
  sleep 3
  ensure_txresult tx/setFeeAggregator.icon

  goloop rpc sendtx call --to $(cat icon.addr.bmc) \
    --method setFeeGatheringTerm \
    --param _value=$FEE_GATHERING_INTERVAL | jq -r . >tx/setFeeGatheringTerm.icon
  sleep 3
  ensure_txresult tx/setFeeGatheringTerm.icon
}

configure_javascore_add_bts() {
  echo "bmc add bts"
  cd $CONFIG_DIR
  local hasBTS=$(goloop rpc call \
    --to $(cat icon.addr.bmc) \
    --method getServices | jq -r .bts)
  if [ "$hasBTS" == "null" ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
      --method addService \
      --value 0 \
      --param _addr=$(cat icon.addr.bts) \
      --param _svc="bts" | jq -r . >tx/addService.icon
    sleep 3
    ensure_txresult tx/addService.icon
  fi
  sleep 5
}

configure_javascore_add_bts_owner() {
  echo "Add bts Owner"
  local icon_bts_owner=$(cat $CONFIG_DIR/keystore/icon.bts.wallet.json | jq -r .address)
  cd $CONFIG_DIR
  local is_owner=$(goloop rpc call \
    --to $(cat icon.addr.bts) \
    --method isOwner \
    --param _addr="$icon_bts_owner" | jq -r .)
  if [ "$is_owner" == "0x0" ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bts) \
    --method addOwner \
    --param _addr=$icon_bts_owner  | jq -r . >tx/addBtsOwner.icon
    sleep 3
    ensure_txresult tx/addBtsOwner.icon
  fi
}

configure_javascore_bts_setICXFee() {
  echo "bts set fee" $ICON_NATIVE_COIN_SYM
  local bts_fee_numerator=100
  local bts_fixed_fee=5000
  cd $CONFIG_DIR
  goloop rpc sendtx call --to $(cat icon.addr.bts) \
    --method setFeeRatio \
    --param _name="${ICON_NATIVE_COIN_NAME}" \
    --param _feeNumerator=$(decimal2Hex $bts_fee_numerator) \
    --param _fixedFee=$(decimal2Hex $bts_fixed_fee) | jq -r . >tx/setICXFee.icon
  sleep 3
  ensure_txresult tx/setICXFee.icon
}

configure_javascore_addLink() {
  echo "BMC: Add Link to BSC BMC:"
  cd $CONFIG_DIR
  if [ ! -f icon.configure.addLink ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
      --method addLink \
      --param _link=$(cat bsc.addr.bmcbtp) | jq -r . >tx/addLink.icon
    sleep 3
    ensure_txresult tx/addLink.icon
    echo "addedLink" > icon.configure.addLink
  fi
}

configure_javascore_setLinkHeight() {
  echo "BMC: SetLinkHeight"
  cd $CONFIG_DIR
  if [ ! -f icon.configure.setLink ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
      --method setLinkRxHeight \
      --param _link=$(cat bsc.addr.bmcbtp) \
      --param _height=$(cat bsc.chain.height)| jq -r . >tx/setLinkRxHeight.icon
    sleep 3
    ensure_txresult tx/setLinkRxHeight.icon
    echo "setLink" > icon.configure.setLink
  fi
}

configure_bmc_javascore_addRelay() {
  echo "Adding bsc Relay"
  local icon_bmr_owner=$(cat $CONFIG_DIR/keystore/icon.bmr.wallet.json | jq -r .address)
  cd $CONFIG_DIR
  if [ ! -f icon.configure.addRelay ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
      --method addRelay \
      --param _link=$(cat bsc.addr.bmcbtp) \
      --param _addr=${icon_bmr_owner} | jq -r . >tx/addRelay.icon
    sleep 3
    ensure_txresult tx/addRelay.icon
    echo "addRelay" > icon.configure.addRelay
  fi
}


configure_javascore_register_native_token() {
  echo "Register Native Token " $2
  cd $CONFIG_DIR
  local bts_fee_numerator=100
  local bts_fixed_fee=5000
  if [ ! -f icon.register.coin$2 ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bts) \
      --method register \
      --param _name="$1" \
      --param _symbol=$2 \
      --param _decimals=0x12 \
      --param _addr=$(cat icon.addr.$2) \
      --param _feeNumerator=$(decimal2Hex $bts_fee_numerator) \
      --param _fixedFee=$(decimal2Hex $bts_fixed_fee) | jq -r . >tx/register.coin.$2
    sleep 5
    ensure_txresult tx/register.coin.$2
    echo "registered "$2 > icon.register.coin$2
  fi
}


configure_javascore_register_wrapped_coin() {
  echo "Register Wrapped Coin " $2
  cd $CONFIG_DIR
  local bts_fee_numerator=100
  local bts_fixed_fee=5000
  if [ ! -f icon.register.coin$2 ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bts) \
      --method register \
      --param _name="$1" \
      --param _symbol=$2 \
      --param _decimals=0x12 \
      --param _feeNumerator=$(decimal2Hex $bts_fee_numerator) \
      --param _fixedFee=$(decimal2Hex $bts_fixed_fee) | jq -r . >tx/register.coin.$2
    sleep 5
    ensure_txresult tx/register.coin.$2
    echo $2 > icon.register.coin$2
  fi
}

get_btp_icon_coinId() {
  echo "Get BTP Icon Addr " $2
  cd $CONFIG_DIR
  goloop rpc call --to $(cat icon.addr.bts) \
    --method coinId \
    --param _coinName="$1" | jq -r . >tx/icon.coinId.$2
  if [ "$(cat $CONFIG_DIR/tx/icon.coinId.$2)" == "null" ];
  then
    echo "Error Gettting  CoinAddress icon."$2
    return 1
  else 
    cat $CONFIG_DIR/tx/icon.coinId.$2 >$CONFIG_DIR/icon.addr.coin$2
  fi
  sleep 5
}
