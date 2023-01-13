# #!/bin/bash

# set -e
# source utils.sh
# source config.sh

# ROOT_DIR=$(echo "$(
#   cd "$(dirname "../../../../../")"
#   pwd
# )")

# export PRIVATE_KEY="[\""$(cat $BSC_KEY_STORE.priv)"\"]"

# copy_migrations() {
#   echo "copying btscore migration start"
#   cd $ROOT_DIR/solidity
  
#   if [ ! -f $ROOT_DIR/solidity/bts/contracts/${1}.sol ]; then
#     echo "Contract BTSCore to upgrade to: ${1} does not exist"
#     exit 0
#   fi

#   if [ ! -f $ROOT_DIR/solidity/bts/migrations/4_upgrade_btsCore.js ]; then
#     echo "Migration script does not exist"
#     exit 0
#   fi

#   cp bts/contracts/${1}.sol $CONTRACTS_DIR/solidity/bts/contracts/${1}.sol
#   cp bts/migrations/4_upgrade_btsCore.js $CONTRACTS_DIR/solidity/bts/migrations/4_upgrade_btsCore.js
#   echo "copy btscore upgrade migration successfully"
# }


# upgrade_solidity_bts() {
#   echo "Upgrading solidity btsCore"
#   cd $CONTRACTS_DIR/solidity/bts
#   if [ ! -f $CONFIG_DIR/bsc.addr.btscore ]; then
#     echo "BTSCore address file bsc.addr.btscore does not exist"
#     exit
#   fi

#   if [ ! -f $CONFIG_DIR/bsc.addr.btscore.upgrade ]; then
#     truffle compile --all
#     set +e
#     local status="retry"

#     echo "Check if ${1} contract exists: "
#     if [ ! -f $CONTRACTS_DIR/solidity/bts/build/contracts/${1}.json ]; then
#       echo "Contract BTSCore to upgrade to: ${1} not compiled"
#       exit 0
#     fi
#     echo "${1} exists"

#     proxyBTSCore=$(jq -r '.networks[] | .address' $CONTRACTS_DIR/solidity/bts/build/contracts/BTSCore.json)
#     deployedBTSCore=$(cat $CONFIG_DIR/bsc.addr.btscore)

#     if [ "$proxyBTSCore" != "$deployedBTSCore" ]; then
#       echo "Address not verified"
#       exit 0
#     fi

#     for i in $(seq 1 20); do
#       truffle migrate --compile-none --network bsc --f 4 --to 4 --btsCore ${1} --proxyAddr ${proxyBTSCore}
#       if [ $? == 0 ]; then
#         status="ok"
#         break
#       fi
#       echo "Retry: "$i
#     done
#     set -e
#     if [ "$status" == "retry" ]; then
#       echo "BTSCore Upgrade Failed after retry"
#       exit 1
#     fi
#     echo 'BTSCore Proxy Address after upgrade'
#     jq -r '.networks[] | .address' build/contracts/BTSCore.json
#     echo -n "btscoreupgraded" >$CONFIG_DIR/bsc.addr.btscore.upgrade
#   fi
# }

# post_upgrade() {
#   if [ ! -f $CONFIG_DIR/bsc.addr.btscore.upgrade ]; then
#     echo "BTSCoreUpgrade address file bsc.addr.btscore.upgrade does not exist"
#     exit
#   fi
#   npx -v
#   if [ $? == 0 ]; then
#     export BSC_SCAN_API_KEY=${2}
#     npx truffle run verify ${1} --network bsc
#   else
#     echo "npx not installed."
#   fi

# }

# if [ $# -eq 0 ]; then
#   echo "No arguments supplied: Pass --help for details."
# elif [[ $1 == "--upgradeTo" && $3 == "--apiKey" ]]; then
#   echo "Start BTSCore Upgrade "
#   copy_migrations $2
#   upgrade_solidity_bts ${2}
#   post_upgrade ${2} ${4}
#   echo "Done"
# else
#   echo "Invalid argument: Pass --upgradeTo ContractName --apiKey <BSC API KEY> to upgrade"
#   echo "Example: ./upgrade_btsCore.sh --upgradeTo BTSCoreV2 --apiKey 2QXI8QMAYX38IT3U3336QM69J3KXVJE9Ias"
# fi
