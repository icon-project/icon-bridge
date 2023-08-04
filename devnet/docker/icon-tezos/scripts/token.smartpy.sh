#!/bin/bash
## smarpy service methods - start ###

# source utils.sh
# source ~/GoProjects/icon-bridge/devnet/docker/icon-bsc/scripts/rpc.sh
# source keystore.sh

export ICON_NET_ID=0x7
export ICON_NET_URI=https://berlin.net.solidwallet.io/api/v3

export BASE_DIR=$(echo $(pwd))/../../../..
export CONFIG_DIR=$BASE_DIR/devnet/docker/icon-tezos
export TEZOS_SETTER=$BASE_DIR/tezos-addresses
export JAVASCORE_DIR=$BASE_DIR/javascore
export SMARTPY_DIR=$BASE_DIR/smartpy
export CONTRACTS_DIR=$BASE_DIR
export TEZOS_BMC_NID=NetXnHfVqm9iesp.tezos
export ICON_BMC_NID=0x7.icon
export TZ_COIN_SYMBOL=XTZ
export TZ_FIXED_FEE=0
export TZ_NUMERATOR=0
export TZ_DECIMALS=6
export ICON_NATIVE_COIN_NAME=btp-0x7.icon-ICX
export ICON_SYMBOL=ICX
export ICON_FIXED_FEE=4300000000000000000
export ICON_NUMERATOR=100
export ICON_DECIMALS=18
export FEE_GATHERING_INTERVAL=43200
export RELAYER_ADDRESS=tz1ZPVxKiybvbV1GvELRJJpyE1xj1UpNpXMv


tz_lastBlock() {
    octez-client rpc get /chains/main/blocks/head/header
}

extract_chainHeight() {
    # cd $CONFIG_DIR
    local tz_block_height=$(tz_lastBlock | jq -r .level)
    echo $tz_block_height > tz.chain.height
}

ensure_config_dir() {
  echo ensuring config dir
  cd $CONFIG_DIR
  if [ ! -d _ixh ]; then 
    echo _ixh not found so creating one 
    mkdir _ixh 
  fi
  if [ ! -d $CONFIG_DIR/_ixh/tx ]; then 
    echo tx not found so creating one
    cd _ixh
    mkdir tx
  fi
}

ensure_tezos_keystore(){
    echo "ensuring key store"
    cd $(echo $SMARTPY_DIR/bmc)
    if [ -f .env ]; then
        echo ".env found"
        octez-client forget address bmcOwner --force
        octez-client gen keys bmcOwner
        local keystore=$(echo $(octez-client show address bmcOwner -S))
        local keystoreClone=$keystore
        keystore_secret=${keystore#*Secret Key: unencrypted:}
        keystore_hash=${keystoreClone#*Hash: }
        keystore_hash=${keystore_hash%% *}
        echo $keystore_hash > tz.bmc.wallet
        echo $keystore_secret > .env
    fi

    cd $SMARTPY_DIR/bts
    if [ -f .env ]; then
        echo ".env found"
        octez-client forget address btsOwner --force
        octez-client gen keys btsOwner
        local keystore=$(echo $(octez-client show address btsOwner -S))
        local keystoreClone=$keystore
        keystore_secret=${keystore#*Secret Key: unencrypted:}
        keystore_hash=${keystoreClone#*Hash: }
        keystore_hash=${keystore_hash%% *}
        echo $keystore_hash > tz.bts.wallet
        echo $keystore_secret > .env
    fi

}

ensure_key_secret() {
  if [ $# -lt 1 ] ; then
    echo "Usage: ensure_key_secret SECRET_PATH"
    return 1
  fi
  local KEY_SECRET=$1
  if [ ! -f "${KEY_SECRET}" ]; then
    mkdir -p $(dirname ${KEY_SECRET})
    echo -n $(openssl rand -hex 20) > ${KEY_SECRET}
  fi
  echo ${KEY_SECRET}
}

ensure_key_store() {
  if [ $# -lt 2 ] ; then
    echo "Usage: ensure_key_store KEYSTORE_PATH SECRET_PATH"
    return 1
  fi
  local KEY_STORE=$1
  local KEY_SECRET=$(ensure_key_secret $2)
  if [ ! -f "${KEY_STORE}" ]; then
    echo should not reach here
    goloop ks gen --out ${KEY_STORE}tmp -p $(cat ${KEY_SECRET}) > /dev/null 2>&1
    cat ${KEY_STORE}tmp | jq -r . > ${KEY_STORE}
    rm ${KEY_STORE}tmp

  fi
  echo ${KEY_STORE}
}

fund_it_flag() {
  cd $CONFIG_DIR
  if [ ! -f fundit.flag ]; then 
    echo Fund the recently created wallet and run the script once again
    echo icon bmc wallet:      $(cat $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json | jq -r .address)
    echo icon bts wallet:      $(cat $CONFIG_DIR/_ixh/keystore/icon.bts.wallet.json | jq -r .address)
    echo icon bmr wallet:      $(cat $CONFIG_DIR/_ixh/keystore/icon.bmr.wallet.json | jq -r .address)
    echo icon fa wallet :      $(cat $CONFIG_DIR/_ixh/keystore/icon.fa.wallet.json | jq -r .address)

    # echo tz bmc wallet  :      $(cat $SMARTPY_DIR/bmc/tz.bmc.wallet)
    # echo tz bts wallet  :      $(cat $SMARTPY_DIR/bts/tz.bts.wallet)

    echo fund it flag > fundit.flag
    exit 0
  fi
}

deploy_smartpy_bmc_management(){
    cd $(echo $SMARTPY_DIR/bmc)
    if [ ! -f $CONDIG_DIR/_ixh/tz.addr.bmcmanagementbtp ]; then
        echo "deploying bmc_management"
        extract_chainHeight
        npm run compile bmc_management
        local deploy=$(npm run deploy bmc_management @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > $CONFIG_DIR/_ixh/tz.addr.bmc_management
        cd $CONFIG_DIR/_ixh
        echo "btp://$(echo $TEZOS_BMC_NID)/$(cat tz.addr.bmc_management)" > $CONFIG_DIR/_ixh/tz.addr.bmcmanagementbtp
    fi
}

deploy_smartpy_bmc_periphery(){
    cd $(echo $SMARTPY_DIR/bmc)
    if [ ! -f $CONDIG_DIR/_ixh/tz.addr.bmcperipherybtp ]; then
        echo "deploying bmc_periphery"
        npm run compile bmc_periphery
        local deploy=$(npm run deploy bmc_periphery @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > $CONFIG_DIR/_ixh/tz.addr.bmc_periphery
        cd $CONFIG_DIR/_ixh
        echo "btp://$(echo $TEZOS_BMC_NID)/$(cat tz.addr.bmc_periphery)" > $CONFIG_DIR/_ixh/tz.addr.bmcperipherybtp
    fi
}

deploy_smartpy_bts_periphery(){
    cd $(echo $SMARTPY_DIR/bts)
    if [ ! -f $CONDIG_DIR/_ixh/tz.addr.bts_periphery ]; then
        echo "deploying bts_periphery"
        npm run compile bts_periphery
        local deploy=$(npm run deploy bts_periphery @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > $CONFIG_DIR/_ixh/tz.addr.bts_periphery
    fi
}

deploy_smartpy_bts_core(){
    cd $(echo $SMARTPY_DIR/bts)
    if [ ! -f $CONDIG_DIR/_ixh/tz.addr.bts_core ]; then
        echo "deploying bts_core"
        cd ~/GoProjects/icon-bridge/smartpy/bts
        npm run compile bts_core
        local deploy=$(npm run deploy bts_core @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > $CONFIG_DIR/_ixh/tz.addr.bts_core
    fi
}

deploy_smartpy_bts_owner_manager(){
    cd $(echo $SMARTPY_DIR/bts)
    if [ ! -f $CONDIG_DIR/_ixh/tz.addr.bts_owner_manager ]; then
        echo "deploying bts_owner_manager"
        cd ~/GoProjects/icon-bridge/smartpy/bts
        npm run compile bts_owner_manager
        local deploy=$(npm run deploy bts_owner_manager @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > $CONFIG_DIR/_ixh/tz.addr.bts_owner_manager
    fi 
}

ensure_txresult() {
  OLD_SET_OPTS=$(set +o)
  set +e
  local TX=$1

  if [ -f "${TX}" ]; then
    TX=$(cat ${TX})
  fi

  sleep 2
  RESULT=$(goloop rpc txresult ${TX} --uri $ICON_NET_URI)
  RET=$?
  echo $RESULT
  while [ "$RET" != "0" ] && [ "$(echo $RESULT | grep -E 'Executing|Pending')" == "$RESULT" ]; do
    sleep 1
    RESULT=$(goloop rpc txresult ${TX} --rui $ICON_NET_URI)
    RET=$?
    echo $RESULT
  done
  eval "${OLD_SET_OPTS}"
  set -e
  if [ "$RET" != "0" ]; then
    echo $RESULT
    return $RET
  else
    STATUS=$(echo $RESULT | jq -r .status)
    if [ "$STATUS" == "0x1" ]; then
      return 0
    else
      echo $RESULT
      return 1
    fi
  fi
}

extract_scoreAddress() {
  local TX=$1
  local ADDR=$2

  RESULT=$(ensure_txresult $TX)
  RET=$?

  if [ "$RET" != "0" ]; then
    echo $RESULT
    return $RET
  else
    SCORE=$(echo $RESULT | jq -r .scoreAddress)
    echo $SCORE | tee ${ADDR}
  fi
}


configure_smartpy_bmc_management_set_bmc_periphery() {
    echo "Adding BMC periphery in bmc management"
    cd $(echo $CONFIG_DIR/_ixh)

    local bmc_periphery=$(echo $(cat tz.addr.bmc_periphery))
    echo $bmc_periphery

    local bmc_management=$(echo $(cat tz.addr.bmc_management))
    echo $bmc_management

    local ocBR=\'\"
    local coBR=\"\'
    local arg=$(echo $(echo $ocBR$(echo $bmc_periphery$(echo $coBR))))

    echo $arg

    # octez-client transfer 0 from bmcOwner to KT1BE6ohnjunYd1C96yPaThwNvFZu4TN8iBy --entrypoint set_bmc_periphery --arg '"KT1JX3Z3AQnf6oDae87Z9mw1g4jhB38tAGQY"' --burn-cap 1
    echo octez-client transfer 0 from bmcOwner to $(echo $bmc_management) --entrypoint set_bmc_periphery --arg $(echo $arg) --burn-cap 1
}

configure_dotenv() {
    echo "Configuring .env file for running the setter script"
    cd $(echo $CONFIG_DIR/_ixh)
    local bmc_periphery=$(echo $(cat tz.addr.bmc_periphery))
    local bmc_management=$(echo $(cat tz.addr.bmc_management))
    local bmc_height=$(echo $(cat tz.chain.height))
    local icon_bmc_height=$(echo $(cat icon.chain.height))
    local icon_bmc=$(echo $(cat icon.addr.bmc))
    echo $bmc_periphery

    local bts_core=$(echo $(cat tz.addr.bts_core))
    local bts_owner_manager=$(echo $(cat tz.addr.bts_owner_manager))
    local bts_periphery=$(echo $(cat tz.addr.bts_periphery))

    cd $SMARTPY_DIR/bmc
    local env=$(cat .env)
    env=${env#*=}
    local secret_deployer=$(echo "secret_deployer=$(echo $env)")
    
    cd $(echo $TEZOS_SETTER)
    go mod tidy
    if [ -f .env ]; then
        echo ".env exists so removing"
        rm .env
    fi
    touch .env
    local output=.env


    local TZ_NETWORK=$(echo "TZ_NETWORK=$(echo $TEZOS_BMC_NID)")
    local ICON_NETWORK=$(echo "ICON_NETWORK=$(echo $ICON_BMC_NID)")
    local TEZOS_NATIVE_COIN_NAME=$(echo "TZ_NATIVE_COIN_NAME=btp-$(echo $TEZOS_BMC_NID)-XTZ")
    local TEZOS_SYMBOL=$(echo "TZ_SYMBOL=$(echo $TZ_COIN_SYMBOL)")
    local TEZ_FIXED_FEE=$(echo "TZ_FIXED_FEE=$(echo $TZ_FIXED_FEE)")

    local TEZ_NUMERATOR=$(echo "TZ_NUMERATOR=$(echo $TZ_NUMERATOR)")
    local TEZ_DECIMALS=$(echo "TZ_DECIMALS=$(echo $TZ_DECIMALS)")
    local IC_NATIVE_COIN_NAME=$(echo "ICON_NATIVE_COIN_NAME=$(echo $ICON_NATIVE_COIN_NAME)")

    local IC_SYMBOL=$(echo "ICON_SYMBOL=$(echo $ICON_SYMBOL)")

    local IC_FIXED_FEE=$(echo "ICON_FIXED_FEE=$(echo $ICON_FIXED_FEE)")
    
    local IC_NUMERATOR=$(echo "ICON_NUMERATOR=$(echo $ICON_NUMERATOR)")
    local IC_DECIMALS=$(echo "ICON_DECIMALS=$(echo $ICON_DECIMALS)")

    local BMC_PERIPHERY=$(echo "BMC_PERIPHERY=$(echo $bmc_periphery)") 
    local BMC_MANAGEMENT=$(echo "BMC_MANAGEMENT=$(echo $bmc_management)")
    local BMC_HEIGHT=$(echo "bmc_periphery_height=$(echo $bmc_height)")
    
    local BTS_PERIPHERY=$(echo "BTS_PERIPHERY=$(echo $bts_periphery)")
    local BTS_CORE=$(echo "BTS_CORE=$(echo $bts_core)")
    local BTS_OWNER_MANAGER=$(echo "BTS_OWNER_MANAGER=$(echo $bts_owner_manager)")
    local RELAY_ADDRESS=$(echo "RELAYER_ADDRESS=$(echo $RELAYER_ADDRESS)")
    local ICON_BMC=$(echo "ICON_BMC=$(echo $icon_bmc)")
    local ICON_RX_HEIGHT=$(echo "ICON_RX_HEIGHT=$(echo $icon_bmc_height)")


    echo $secret_deployer>>$output

    echo $TZ_NETWORK>>$output
    echo $ICON_NETWORK>>$output
    echo $TEZOS_NATIVE_COIN_NAME>>$output
    echo $TEZOS_SYMBOL>>$output
    echo $TEZ_FIXED_FEE>>$output
    echo $TEZ_NUMERATOR>>$output
    echo $TEZ_DECIMALS>>$output
    echo $IC_NATIVE_COIN_NAME>>$output
    echo $IC_SYMBOL>>$output
    echo $IC_FIXED_FEE>>$output
    echo $IC_NUMERATOR>>$output
    echo $IC_DECIMALS>>$output
    echo $BMC_PERIPHERY>>$output 
    echo $BMC_MANAGEMENT>>$output 
    echo $BMC_HEIGHT>>$output 
    echo $BTS_PERIPHERY>>$output
    echo $BTS_CORE>>$output
    echo $BTS_OWNER_MANAGER>>$output
    echo $RELAY_ADDRESS>>$output
    echo $ICON_BMC>>$output
    echo $ICON_RX_HEIGHT>>$output
}

run_tezos_setters(){
    cd $(echo $TEZOS_SETTER)
    go run main.go
}


# build java scores
build_javascores(){
  echo in java score
  cd $JAVASCORE_DIR
  ./gradlew bmc:optimizedJar
  ./gradlew bts:optimizedJar

  # irc2-token
  if [ ! -f irc2.jar ]; then 
    git clone https://github.com/icon-project/java-score-examples.git
    cd java-score-examples
    ./gradlew irc2-token:clean
    ./gradlew irc2-token:optimizedJar
    cp irc2-token/build/libs/irc2-token-0.9.1-optimized.jar $JAVASCORE_DIR/irc2.jar
    rm -rf $JAVASCORE_DIR/java-score-examples
  fi

  cp bmc/build/libs/bmc-optimized.jar $JAVASCORE_DIR/bmc.jar
  cp bts/build/libs/bts-optimized.jar $JAVASCORE_DIR/bts.jar
  cp lib/irc2Tradeable-0.1.0-optimized.jar $JAVASCORE_DIR/irc2Tradeable.jar
}



# deploy java scores 

goloop_lastblock() {
  goloop rpc lastblock --uri $ICON_NET_URI
}

extract_chain_height_and_validator() {
    cd $CONFIG_DIR/_ixh
    local icon_block_height=$(goloop_lastblock | jq -r .height)
    echo $icon_block_height > icon.chain.height
}

deploy_javascore_bmc() {
  cd $CONFIG_DIR/_ixh
  if [ ! -f icon.addr.bmcbtp ]; then
    echo "deploying javascore BMC"
    extract_chain_height_and_validator
    cd _ixh
    goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/bmc.jar \
      --content_type application/java \
      --param _net=$ICON_BMC_NID \
      --key_store $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json \
      --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.secret)) \
      --nid $ICON_NET_ID \
      --step_limit 4000000000 \
      --uri $ICON_NET_URI | jq -r . >tx/tx.icon.bmc
    sleep 5
    echo $(pwd)
    extract_scoreAddress tx/tx.icon.bmc icon.addr.bmc
    echo "btp://$(echo $ICON_BMC_NID)/$(cat icon.addr.bmc)" >icon.addr.bmcbtp
  fi
}

deploy_javascore_bts() {
  echo "deploying javascore bts"
  cd $CONFIG_DIR/_ixh
  if [ ! -f icon.addr.bts ]; then
    #local bts_fee_numerator=100
    #local bts_fixed_fee=5000
    goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/bts.jar \
      --content_type application/java \
      --param _name=$ICON_NATIVE_COIN_NAME \
      --param _bmc=$(cat icon.addr.bmc) \
      --param _decimals=$(decimal2Hex $3) \
      --param _feeNumerator=$(decimal2Hex $2) \
      --param _fixedFee=$(decimal2Hex $1) \
      --param _serializedIrc2=$(xxd -p $CONTRACTS_DIR/javascore/irc2Tradeable.jar | tr -d '\n') \
      --key_store $CONFIG_DIR/_ixh/keystore/icon.bts.wallet.json \
      --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bts.wallet.secret)) \
      --step_limit 4000000000 \
      --nid $ICON_NET_ID \
      --uri $ICON_NET_URI  | jq -r . > tx/tx.icon.bts
    sleep 5
    extract_scoreAddress tx/tx.icon.bts icon.addr.bts
  fi
}

deploy_javascore_token() {
  echo "deploying javascore IRC2Token " $2
  cd $CONFIG_DIR/_ixh
  if [ ! -f icon.addr.$2 ]; then
    goloop rpc sendtx deploy $CONTRACTS_DIR/javascore/irc2.jar \
      --content_type application/java \
      --param _name="$1" \
      --param _symbol=$2 \
      --param _initialSupply="0x5f5e100" \
      --param _decimals="0x12" \
      --key_store $CONFIG_DIR/_ixh/keystore/icon.fa.wallet.json \
      --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.fa.wallet.secret)) \
      --nid $ICON_NET_ID \
      --step_limit 4000000000 \
      --uri $ICON_NET_URI | jq -r . >tx/tx.icon.$2
    sleep 5
    extract_scoreAddress tx/tx.icon.$2 icon.addr.$2
  fi
}


configure_javascore_add_bmc_owner() {
  echo "bmc Add Owner"
  echo $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json
  local icon_bmc_owner=$(cat $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json | jq -r .address)
  cd $CONFIG_DIR/_ixh
  local is_owner=$(goloop rpc call \
    --to $(cat icon.addr.bmc) \
    --method isOwner \
    --param _addr=$icon_bmc_owner \
    --uri $ICON_NET_URI | jq -r .)
  if [ "$is_owner" == "0x0" ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
    --method addOwner \
    --param _addr=$icon_bmc_owner \
    --key_store $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json \
    --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.secret)) \
    --step_limit 1000000000 \
    --nid $ICON_NET_ID \
    --uri $ICON_NET_URI | jq -r . > tx/addbmcowner.icon
    sleep 3
    ensure_txresult tx/addbmcowner.icon
  fi
}

configure_javascore_bmc_setFeeAggregator() {
  echo "bmc setFeeAggregator"
  cd $CONFIG_DIR/_ixh
  local FA=$(cat $CONFIG_DIR/_ixh/keystore/icon.fa.wallet.json | jq -r .address)
  goloop rpc sendtx call --to $(cat icon.addr.bmc) \
    --method setFeeAggregator \
    --param _addr=${FA} \
    --key_store $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json \
    --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.secret)) \
    --step_limit 1000000000 \
    --nid $ICON_NET_ID \
    --uri $ICON_NET_URI | jq -r . >tx/setFeeAggregator.icon
  sleep 3
  ensure_txresult tx/setFeeAggregator.icon

  goloop rpc sendtx call --to $(cat icon.addr.bmc) \
    --method setFeeGatheringTerm \
    --param _value=$FEE_GATHERING_INTERVAL \
    --key_store $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json \
    --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.secret)) \
    --step_limit 1000000000 \
    --nid $ICON_NET_ID \
    --uri $ICON_NET_URI | jq -r . >tx/setFeeGatheringTerm.icon
  sleep 3
  ensure_txresult tx/setFeeGatheringTerm.icon
}

configure_javascore_add_bts() {
  echo "bmc add bts"
  cd $CONFIG_DIR/_ixh
  local hasBTS=$(goloop rpc call \
    --to $(cat icon.addr.bmc) \
    --method getServices \
    --uri $ICON_NET_URI | jq -r .bts)
  if [ "$hasBTS" == "null" ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
      --method addService \
      --value 0 \
      --param _addr=$(cat icon.addr.bts) \
      --param _svc="bts" \
      --key_store $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json \
      --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.secret)) \
      --step_limit 1000000000 \
      --nid $ICON_NET_ID \
      --uri $ICON_NET_URI | jq -r . >tx/addService.icon
    sleep 3
    ensure_txresult tx/addService.icon
  fi
  sleep 5
}

configure_javascore_add_bts_owner() {
  echo "Add bts Owner"
  local icon_bts_owner=$(cat $CONFIG_DIR/_ixh/keystore/icon.bts.wallet.json | jq -r .address)
  cd $CONFIG_DIR/_ixh
  local is_owner=$(goloop rpc call \
    --to $(cat icon.addr.bts) \
    --method isOwner \
    --param _addr="$icon_bts_owner" \
    --uri $ICON_NET_URI | jq -r .)
  if [ "$is_owner" == "0x0" ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bts) \
    --method addOwner \
    --param _addr=$icon_bts_owner \
    --key_store $CONFIG_DIR/_ixh/keystore/icon.bts.wallet.json \
    --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bts.wallet.secret)) \
    --step_limit 1000000000 \
    --nid $ICON_NET_ID \
    --uri $ICON_NET_URI | jq -r . >tx/addBtsOwner.icon
    sleep 3
    ensure_txresult tx/addBtsOwner.icon
  fi
}

configure_javascore_bts_setICXFee() {
  echo "bts set fee" $ICON_SYMBOL
  #local bts_fee_numerator=100
  #local bts_fixed_fee=5000
  cd $CONFIG_DIR/_ixh
  goloop rpc sendtx call --to $(cat icon.addr.bts) \
    --method setFeeRatio \
    --param _name=$ICON_NATIVE_COIN_NAME \
    --param _feeNumerator=$(decimal2Hex $2) \
    --param _fixedFee=$(decimal2Hex $1) \
    --key_store $CONFIG_DIR/_ixh/keystore/icon.bts.wallet.json \
    --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bts.wallet.secret)) \
    --nid $ICON_NET_ID \
    --step_limit 1000000000 \
    --uri $ICON_NET_URI | jq -r . >tx/setICXFee.icon
  sleep 3
  ensure_txresult tx/setICXFee.icon
}


configure_javascore_addLink() {
  echo "BMC: Add Link to BSC BMC:"
  cd $CONFIG_DIR/_ixh
  if [ ! -f icon.configure.addLink ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
      --method addLink \
      --param _link=$(cat tz.addr.bmcperipherybtp) \
      --key_store $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json \
      --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.secret)) \
      --nid $ICON_NET_ID \
      --step_limit 1000000000 \
      --uri $ICON_NET_URI | jq -r . > addLink.icon
      
    sleep 3
    echo "addedLink" > icon.configure.addLink
  fi
}

configure_javascore_setLinkHeight() {
  echo "BMC: SetLinkHeight"
  cd $CONFIG_DIR/_ixh
  if [ ! -f icon.configure.setLink ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
      --method setLinkRxHeight \
      --param _link=$(cat tz.addr.bmcperipherybtp) \
      --param _height=$(cat tz.chain.height) \
      --key_store $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json \
      --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.secret)) \
      --nid $ICON_NET_ID \
      --step_limit 1000000000 \
      --uri $ICON_NET_URI | jq -r . > setLinkRxHeight.icon
      
    sleep 3
    echo "setLink" > icon.configure.setLink
  fi
}

configure_bmc_javascore_addRelay() {
  echo "Adding bsc Relay"
  local icon_bmr_owner=$(cat $CONFIG_DIR/_ixh/keystore/icon.bmr.wallet.json | jq -r .address)
  echo $icon_bmr_owner
  sleep 5
  echo "Starting"
  cd $CONFIG_DIR/_ixh
  if [ ! -f icon.configure.addRelay ]; then
    goloop rpc sendtx call --to $(cat icon.addr.bmc) \
      --method addRelay \
      --param _link=$(cat tz.addr.bmcperipherybtp) \
      --param _addr=${icon_bmr_owner} \
      --key_store $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json \
      --key_password $(echo $(cat $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.secret)) \
      --nid $ICON_NET_ID \
      --step_limit 1000000000 \
      --uri $ICON_NET_URI | jq -r . > addRelay.icon

    sleep 3
    echo "addRelay" > icon.configure.addRelay
  fi
}

function decimal2Hex() {
    hex=$(echo "obase=16; ibase=10; ${@}" | bc)
    echo "0x$(tr [A-Z] [a-z] <<< "$hex")"
}

ensure_config_dir
ensure_key_store $CONFIG_DIR/_ixh/keystore/icon.bts.wallet.json $CONFIG_DIR/_ixh/keystore/icon.bts.wallet.secret
ensure_key_store $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.json $CONFIG_DIR/_ixh/keystore/icon.bmc.wallet.secret
ensure_key_store $CONFIG_DIR/_ixh/keystore/icon.bmr.wallet.json $CONFIG_DIR/_ixh/keystore/icon.bmr.wallet.secret
ensure_key_store $CONFIG_DIR/_ixh/keystore/icon.fa.wallet.json $CONFIG_DIR/_ixh/keystore/icon.fa.wallet.secret
# ensure_tezos_keystore
fund_it_flag

build_javascores
deploy_javascore_bmc
deploy_javascore_bts 0 0 18
deploy_javascore_token

configure_javascore_add_bmc_owner
configure_javascore_add_bts
configure_javascore_add_bts_owner
configure_javascore_bmc_setFeeAggregator
configure_javascore_bts_setICXFee 0 0



# # tezos configuration
deploy_smartpy_bmc_management
deploy_smartpy_bmc_periphery
deploy_smartpy_bts_periphery
deploy_smartpy_bts_core
deploy_smartpy_bts_owner_manager
configure_dotenv
run_tezos_setters

# # # icon configuration of tezos
configure_javascore_addLink
configure_javascore_setLinkHeight
configure_bmc_javascore_addRelay

