#!/bin/bash
## smarpy service methods - start ###

# source utils.sh
# source prc.sh
# source keystore.sh

tz_lastBlock() {
    octez-client rpc get /chains/main/blocks/head/header
}

extract_chainHeight() {
    # cd $CONFIG_DIR
    local tz_block_height=$(tz_lastBlock | jq -r .level)
    echo $tz_block_height > tz.chain.height
}

deploy_smartpy_bmc_management(){
    if [ ! -f tz.addr.bmcmanagementbtp ]; then
        echo "deploying bmc_management"
        extract_chainHeight
        cd ~/GoProjects/icon-bridge/smartpy/bmc
        npm run compile bmc_management
        local deploy=$(npm run deploy bmc_management @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > tz.addr.bmc_management
        echo "btp://0x63.tezos/$(cat tz.addr.bmc)" > tz.addr.bmcmanagementbtp
    fi
}

deploy_smartpy_bmc_periphery(){
    if [ ! -f tz.addr.bmcperipherybtp ]; then
        echo "deploying bmc_periphery"
        cd ~/GoProjects/icon-bridge/smartpy/bmc
        npm run compile bmc_periphery
        local deploy=$(npm run deploy bmc_periphery @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > tz.addr.bmc_periphery
        echo "btp://0x63.tezos/$(cat tz.addr.bmc_periphery)" > tz.addr.bmcperipherybtp
    fi
}

deploy_smartpy_bts_periphery(){
    if [ ! -f tz.addr.bts_periphery ]; then
        echo "deploying bts_periphery"
        cd ~/GoProjects/icon-bridge/smartpy/bts
        npm run compile bts_periphery
        local deploy=$(npm run deploy bts_periphery @GHOSTNET)
        sleep 5
        deploy=${deploy#*::}
        echo $deploy > tz.addr.bts_periphery
    fi
}

deploy_smartpy_bts_core(){
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




# bts core
# bts owner manager




deploy_smartpy_bts_owner_manager