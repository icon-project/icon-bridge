#!/bin/bash
## smarpy service methods - start ###

# source utils.sh
# source prc.sh
# source keystore.sh

export CONFIG_DIR=~/GoProjects/icon-bridge/smartpy  
export TEZOS_SETTER=~/tezos-addresses
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
export RELAYER_ADDRESS=tz1ZPVxKiybvbV1GvELRJJpyE1xj1UpNpXMv


tz_lastBlock() {
    octez-client rpc get /chains/main/blocks/head/header
}

extract_chainHeight() {
    # cd $CONFIG_DIR
    local tz_block_height=$(tz_lastBlock | jq -r .level)
    echo $tz_block_height > tz.chain.height
}

ensure_tezos_keystore(){
    echo "ensuring key store"
    cd $(echo $CONFIG_DIR/bmc)
    if [ -f .env ]; then
        echo ".env found"
        octez-client gen keys bmcbtsOwner
        echo $(octez-client show address bmcbtsOwner -S)
    fi
}

deploy_smartpy_bmc_management(){
    cd $(echo $CONFIG_DIR/bmc)
    if [ ! -f tz.addr.bmcmanagementbtp ]; then
        echo "deploying bmc_management"
        extract_chainHeight
        npm run compile bmc_management
        local deploy=$(npm run deploy bmc_management @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > tz.addr.bmc_management
        echo "btp://$(echo $TEZOS_BMC_NID)/$(cat tz.addr.bmc_management)" > tz.addr.bmcmanagementbtp
    fi
}

deploy_smartpy_bmc_periphery(){
    cd $(echo $CONFIG_DIR/bmc)
    if [ ! -f tz.addr.bmcperipherybtp ]; then
        echo "deploying bmc_periphery"
        npm run compile bmc_periphery
        local deploy=$(npm run deploy bmc_periphery @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > tz.addr.bmc_periphery
        echo "btp://$(echo $TEZOS_BMC_NID)/$(cat tz.addr.bmc_periphery)" > tz.addr.bmcperipherybtp
    fi
}

deploy_smartpy_bts_periphery(){
    cd $(echo $CONFIG_DIR/bts)
    if [ ! -f tz.addr.bts_periphery ]; then
        echo "deploying bts_periphery"
        npm run compile bts_periphery
        local deploy=$(npm run deploy bts_periphery @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > tz.addr.bts_periphery
    fi
}

deploy_smartpy_bts_core(){
    cd $(echo $CONFIG_DIR/bts)
    if [ ! -f tz.addr.bts_core ]; then
        echo "deploying bts_core"
        cd ~/GoProjects/icon-bridge/smartpy/bts
        npm run compile bts_core
        local deploy=$(npm run deploy bts_core @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > tz.addr.bts_core
    fi
}

deploy_smartpy_bts_owner_manager(){
    cd $(echo $CONFIG_DIR/bts)
    if [ ! -f tz.addr.bts_owner_manager ]; then
        echo "deploying bts_owner_manager"
        cd ~/GoProjects/icon-bridge/smartpy/bts
        npm run compile bts_owner_manager
        local deploy=$(npm run deploy bts_owner_manager @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > tz.addr.bts_owner_manager
    fi 
}


configure_smartpy_bmc_management_set_bmc_periphery() {
    echo "Adding BMC periphery in bmc management"
    cd $(echo $CONFIG_DIR/bmc)

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
    cd $(echo $CONFIG_DIR/bmc)
    local bmc_periphery=$(echo $(cat tz.addr.bmc_periphery))
    local bmc_management=$(echo $(cat tz.addr.bmc_management))
    local bmc_height=$(echo $(cat tz.chain.height))
    local icon_bmc_height=$(echo $(cat iconbmcheight))
    local icon_bmc=$(echo $(cat iconbmc))
    echo $bmc_periphery

    cd $(echo $CONFIG_DIR/bts)
    local bts_core=$(echo $(cat tz.addr.bts_core))
    local bts_owner_manager=$(echo $(cat tz.addr.bts_owner_manager))
    local bts_periphery=$(echo $(cat tz.addr.bts_periphery))
    local env=$(cat .env)
    env=${env#*=}
    local secret_deployer=$(echo "secret_deployer=$(echo $env)")
    
    cd $(echo $TEZOS_SETTER)
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

# bts core
# bts owner manager


# ensure_tezos_keystore
deploy_smartpy_bmc_management
deploy_smartpy_bmc_periphery
deploy_smartpy_bts_periphery
deploy_smartpy_bts_core
deploy_smartpy_bts_owner_manager
configure_dotenv
run_tezos_setters
