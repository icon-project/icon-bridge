#!/bin/bash
## smarpy service methods - start ###

# source utils.sh
# source prc.sh
# source keystore.sh

export CONFIG_DIR=~/GoProjects/icon-bridge/smartpy  
export TEZOS_BMC_NID="NetXnHfVqm9iesp.tezos"

tz_lastBlock() {
    octez-client rpc get /chains/main/blocks/head/header
}

extract_chainHeight() {
    # cd $CONFIG_DIR
    local tz_block_height=$(tz_lastBlock | jq -r .level)
    echo $tz_block_height > tz.chain.height
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
    if [ ! -f tz.addr.btsperipherybtp ]; then
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

configure_smartpy_bmc_management_set








# bts core
# bts owner manager



# deploy_smartpy_bmc_management
# deploy_smartpy_bmc_periphery
# deploy_smartpy_bts_periphery
# deploy_smartpy_bts_core
# deploy_smartpy_bts_owner_manager
configure_smartpy_bmc_management_set_bmc_periphery