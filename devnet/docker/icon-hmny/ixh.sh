#!/bin/bash

# source $ixh_env

localhosts=( localdckr )
remotehosts=( localnets )
allowedhosts=( ${localhosts[@]} ${remotehosts[@]} )

docker_user="ubuntu"
# docker_host="localdckr"
docker_host="localnets"
docker_port="5000"
docker_registry="$docker_host$([[ "$docker_port" == "" ]] && echo || echo ":$docker_port")"

btp_hmny_uri="http://$docker_host:9500"
btp_icon_uri="http://$docker_host:9080/api/v3"

ixh_dir=$PWD
ixh_tmp_dir=$ixh_dir/_ixh
ixh_build_dir=$ixh_tmp_dir/build
ixh_tests_dir=$ixh_tmp_dir/tests
ixh_env=$ixh_tmp_dir/ixh.env
ixh_src_dir=$ixh_dir/src

root_dir="$ixh_dir/../../.."

btp_icon_branch="v1.2.3"
btp_hmny_branch="v4.3.7"

btp_icon_config=$ixh_src_dir/icon.config.json
btp_icon_wallet=$ixh_src_dir/icon.wallet.json
btp_icon_wallet_password=gochain
btp_icon_step_limit=13610920001
btp_icon_nativecoin_symbol=ICX
btp_icon_nativecoin_bsh_svc_name=nativecoin

btp_hmny_wallet=$ixh_src_dir/hmny.wallet.json
btp_hmny_wallet_password=
btp_hmny_wallet_private_key=1f84c95ac16e6a50f08d44c7bde7aff8742212fda6e4321fde48bf83bef266dc
btp_hmny_wallet_address=0xA5241513DA9F4463F1d4874b548dFBAC29D91f34
# btp_hmny_wallet_address_one=one155jp2y76nazx8uw5sa94fr0m4s5aj8e5xm6fu3
btp_hmny_gas_limit=3000000000
btp_hmny_gas_price=30000000000
btp_hmny_nativecoin_symbol=ONE_DEV

# fd for verbose logs
mkdir -p $ixh_tmp_dir
exec 3<> $ixh_tmp_dir/ixh.log
echo $(date) >&3 # print current time

# override commands to write stderr to log file
goloop=$(which goloop)
ethkey=$(which ethkey)
truffle=$(which truffle)
function goloop() {
    echo "goloop:" >&3 && $goloop $@ 2>&3
}
function ethkey() {
    echo "ethkey:" >&3 && $ethkey $@ 2>&3
}
function truffle() {
    echo "truffle:" >&3 && $truffle $@ 2>&3
}

# log message to stderr
function log() {
    echo -e "$@" >&2 
}

function log_status() {
    [[ "$1" == 0 ]] && log -n " ✔" || log -n " ✘"
}

function dec2hex() {
    hex=$(echo "obase=16; ibase=10; ${@}"| bc)
    echo "0x${hex,,}"
}

function hex2dec() {
    hex=${@#0x}
    echo "obase=10; ibase=16; ${hex^^}"| bc
}

function rel_path() {
    realpath --relative-to="$ixh_dir" "$1"
}

function hmny_jsonrpc() {
    curl "$btp_hmny_uri"  -s -X POST -H 'Content-Type: application/json' \
        -d "$(printf '{"jsonrpc": "2.0","id": 1,"method": "%s","params": [%s]}' "$1" "$2")"
}

function icon_jsonrpc() {
    curl "$btp_icon_uri"  -s -X POST -H 'Content-Type: application/json' \
        -d "$(printf '{"jsonrpc": "2.0","id": 1,"method": "%s","params": %s}' "$1" "$2")"
}

function icon_wait_tx() {
    local tx_hash=$1
    [[ ! $tx_hash ]] && return
    local uri=$btp_icon_uri
    local ret=1
    local stime=$(date +%s)
    [[ $verbose ]] && log -n "[tx=${tx_hash}]" || log -n "[tx=${tx_hash:0:10}]"
    for _ in $(seq 1 60); do
        sleep 1
        log -n "."
        tx=$(goloop rpc --uri "$uri" txresult "$tx_hash" 2>&1)
        status=$(2>/dev/null jq -r .status <<< $tx)
        if [[ "$status" == '0x0' ]]; then
            log_status 1
            [[ $verbose ]] && echo "$tx"
            break
        elif [[ "$status" == '0x1' ]]; then
            log_status 0
            echo "$tx"
            ret=0
            break
        fi
    done
    log -n " $(( $(date +%s) - stime ))s "
    return $ret
}

function icon_create_wallet() {
    local keystore=$1
    local password=$2
    local balance=$3
    _=$(goloop ks gen -o "$keystore" -p "$password")
    addr=$(jq -r .address "$keystore")
    tx_hash=$( # add balance
        goloop rpc \
            --uri "$btp_icon_uri" \
            sendtx transfer \
            --to "$addr" \
            --value "$balance" \
            --key_store "$btp_icon_wallet" \
            --key_password "$btp_icon_wallet_password" \
            --nid "$btp_icon_nid" \
            --step_limit "$btp_icon_step_limit" \
        | jq -r .)
    _=$(icon_wait_tx "$tx_hash")
    echo "$addr"
}

function icon_deploy_sc() {
    local sc_filepath=$1
    local params
    for i in "${@:2}"; do params="$params --param $i"; done
    tx_hash=$(
        goloop rpc \
            --uri "$btp_icon_uri" \
            sendtx deploy $sc_filepath \
            --key_store "$btp_icon_wallet" \
            --key_password "$btp_icon_wallet_password" \
            --nid "$btp_icon_nid" \
            --content_type application/java \
            --step_limit "$btp_icon_step_limit" \
            $params \
        | jq -r .)
    icon_wait_tx "$tx_hash"
}

function hmny_create_wallet() {
    local keystore=$1
    local password=$2
    local balance=$3
    { read addr; read body; read tx; } <<< $(
        docker run --rm $docker_registry/hmny:latest bash -c "
            to_addr=\$(hmy keys add bmr --passphrase-file <(echo '$password') | tail -n1 | cut -d: -f2 | xargs)
            hmy keys export-ks --passphrase-file <(echo '$password') bmr / > /dev/null 2>&1
            
            echo \$to_addr
            cat /\$to_addr.key && printf '\n'

            # transfer funds
            from_addr=\$(hmy keys import-private-key '$btp_hmny_wallet_private_key' root | tail -n1 | cut -d: -f2 | xargs)
            hmy -n '$btp_hmny_uri' \
                --no-pretty \
                transfer \
                --from=\$from_addr \
                --to=\$to_addr \
                --from-shard=0 \
                --to-shard=0 \
                --timeout=120 \
                --amount='$balance'
        "
    )
    cat > "$keystore" <<< $body
    status=$(2>/dev/null  jq -r '[.[]."blockchain-receipt".status][0]' <<< $tx)
    [[ "$status" == '0x1' ]] && log_status 0 || log_status 1
    echo "$addr"
}

function hmny_deploy_sc() {
    cd $1
    cp $ixh_src_dir/hmny.truffle-config.js truffle-config.js # replace original truffle-config.js
    yarn > /dev/null 2>&1 # download node_modules
    _truffle compile --all >&3 # compile all
    
    local stime=$(date +%s)
    for _ in $(seq 1 20); do # repeat until successfully deployed
        _truffle migrate --network hmny --skip-dry-run >&3
        if [ $? == 0 ]; then
            local imports=''
            local waiters=''
            for i in ${@:2}; do
                imports="$imports const $i = artifacts.require('$i');"
                waiters="$waiters await $i.deployed(); console.log($i.address);"
            done
            _truffle exec --network hmny <(echo "
                $imports
                module.exports = async function(cb) { $waiters; cb(); }
            ") | tail -n $(( $# - 1 ))
            log_status $?
            break
        fi
        log -n "."
    done
    log -n " $(( $(date +%s) - stime ))s "
}

function hmny_save_root_wallet() {
    local keystore=$1
    local password=$2
    local private_key=$3
    { read addr; read body; } <<< $(
        docker run --rm $docker_registry/hmny:latest bash -c "
            addr=\$(hmy keys import-private-key '$private_key' root | tail -n1 | cut -d: -f2 | xargs)
            hmy keys export-ks root / > /dev/null 2>&1
            echo \$addr
            cat /\$addr.key
    ")
    log_status $?
    cat > "$keystore" <<< $body
    echo $addr
}

function hmny_get_hmny_chain_status() {
    local lh=$(hmny_jsonrpc hmyv2_latestHeader)
    local bh=$(jq -r '.result.blockHash' <<< "$lh")
    local bn=$(jq -r '.result.blockNumber' <<< "$lh")
    local ep=$(hmny_jsonrpc hmyv2_getFullHeader "$(( $bn + 1 ))" | jq -r '.result.epoch')
    local lb=$(hmny_jsonrpc hmyv2_epochLastBlock "$(( $ep - 1 ))" | jq -r '.result')
    local ss=$(hmny_jsonrpc hmyv2_getFullHeader "$lb" | jq -r '.result.shardState')
    echo -e "$bn\n$bh\n$ep\n$ss\n"
}

# Ensure following tools are installed
# gradle, jdk@11.x, sdkman, goloop, docker, truffle@5.3.0, node@15.12.0, ethkey
function deploysc() {

    local init_start_time=$(date +%s)

    # build dir
    mkdir -p $ixh_build_dir
    [[ "$1" == "reset" ]] && rm -rf $ixh_build_dir/* # clean build when reset is enabled

    # create root wallets
    log "Wallet:"
    
    # icon
    log -n "    icon:" && log_status 0
    btp_icon_wallet_address=$(jq -r .address "$btp_icon_wallet")
    log

    # hmny
    log -n "    hmny: [$(rel_path "$btp_hmny_wallet")] "
    btp_hmny_wallet_address_one=$(hmny_save_root_wallet \
        "$btp_hmny_wallet" "$btp_hmny_wallet_password" "$btp_hmny_wallet_private_key")
    btp_hmny_wallet_address="0x$(jq -r .address "$btp_hmny_wallet")"
    log

    # prepare javascore build dir
    local ixh_jsc_dir=$ixh_build_dir/javascore
    cp -r $root_dir/javascore $ixh_jsc_dir

    # build javascores
    log "\nBuild: "
    log -n "    javascores:"
    cd $ixh_jsc_dir/bmc && \
        gradle optimizedJar > /dev/null 2>&1; log_status $?
    # cd $ixh_jsc_dir/bsh && \
    #     gradle optimizedJar > /dev/null 2>&1 && \
    #     gradle optimizedJarIRC2 > /dev/null 2>&1; log_status $?
    log

    log "\nDeploy: "

    btp_icon_nid=$(dec2hex $(cat "$btp_icon_config" | jq -r .nid))
    btp_icon_net="$btp_icon_nid.icon"

    btp_hmny_nid="0x2"
    btp_hmny_net="$btp_hmny_nid.hmny"

    # deploy
    log "icon"    

    # bmc
    log -n "    bmc: "
    r=$(icon_deploy_sc \
        $ixh_jsc_dir/bmc/build/libs/bmc-0.1.0-optimized.jar \
        _net="$btp_icon_net")
    btp_icon_bmc=$(jq -r .scoreAddress <<< $r)
    btp_icon_block_hash=$(jq -r .blockHash <<< $r)
    btp_icon_block_height=$(hex2dec $(jq -r .blockHeight <<< $r))
    log "$btp_icon_bmc"

    # irc31
    log -n "    irc31: "
    r=$(icon_deploy_sc \
        $ixh_jsc_dir/irc31-0.1.0-optimized.jar)
    btp_icon_irc31=$(jq -r .scoreAddress <<< $r)
    log "$btp_icon_irc31"

    # nativecoin bsh
    log -n "    nativecoin_bsh: "
    r=$(icon_deploy_sc \
        $ixh_jsc_dir/nativecoin-0.1.0-optimized.jar \
        _name="$btp_icon_nativecoin_symbol" \
        _bmc="$btp_icon_bmc" \
        _irc31="$btp_icon_irc31")
    btp_icon_nativecoin_bsh=$(jq -r .scoreAddress <<< $r)
    log "$btp_icon_nativecoin_bsh"

    # # token bsh
    # log -n "    bsh: "
    # r=$(icon_deploy_sc \
    #     $ixh_jsc_dir/bsh/build/libs/bsh-optimized.jar \
    #     _bmc="$btp_icon_bmc")
    # btp_icon_token_bsh=$(jq -r .scoreAddress <<< $r)
    # log "$btp_icon_token_bsh"

    # # irc2
    # log -n "    irc2: "
    # r=$(icon_deploy_sc \
    #     $ixh_jsc_dir/bsh/build/libs/irc2-optimized.jar \
    #     _name="$btp_hmny_nativecoin_symbol" \
    #     _symbol="$btp_hmny_nativecoin_symbol" \
    #     _decimals=2 \
    #     _initialSupply=10000)
    # btp_icon_irc2=$(jq -r .scoreAddress <<< $r)
    # log "$btp_icon_irc2"

    # icon btp address
    btp_icon_btp_address="btp://$btp_icon_net/$btp_icon_bmc"
    log "    btp: $btp_icon_btp_address"

    # hmny
    ixh_sol_dir=$ixh_build_dir/solidity
    cp -r $root_dir/solidity $ixh_sol_dir

    # deploy
    log "hmny"

    # before bmc
    { read btp_hmny_block_height; \
        read btp_hmny_block_hash; \
        read btp_hmny_block_epoch; \
        read btp_hmny_shard_state; } <<< "$(hmny_get_hmny_chain_status)"

    # bmc
    log -n "    bmc: "
    { read btp_hmny_bmc_management; read btp_hmny_bmc_periphery; } <<< $(
        BMC_BTP_NET="$btp_hmny_net" \
            hmny_deploy_sc $ixh_sol_dir/bmc BMCManagement BMCPeriphery)
    log "m=$btp_hmny_bmc_management, p=$btp_hmny_bmc_periphery"
    # TODO get hmny bmc block height and epoch in hex (0x...)

    # bsh
    log -n "    bsh: "
    { read btp_hmny_bsh_core; read btp_hmny_bsh_periphery; } <<< $(
        BSH_COIN_URL="https://github.com/icon/btp" \
        BSH_COIN_NAME="$btp_hmny_nativecoin_symbol" \
        BSH_COIN_FEE=10 \
        BSH_FIXED_FEE=500000 \
        BMC_PERIPHERY_ADDRESS="$btp_hmny_bmc_periphery" \
        BSH_SERVICE="$btp_icon_nativecoin_bsh_svc_name" \
            hmny_deploy_sc $ixh_sol_dir/bsh BSHCore BSHPeriphery)
    log "c=$btp_hmny_bsh_core, p=$btp_hmny_bsh_periphery"

    # tokenbsh
    # log -n "    tokenbsh: "
    # { read btp_hmny_bsh_core; read btp_hmny_bsh_periphery; } <<< $(
    #     BSH_TOKEN_FEE=10 \
    #     BMC_PERIPHERY_ADDRESS="$btp_hmny_bmc_periphery" \
    #     BSH_SERVICE="$btp_icon_nativecoin_bsh_svc_name" \
    #         hmny_deploy_sc $ixh_sol_dir/bsh BSHProxy BSHImpl BEP20TKN)
    # log "c=$btp_hmny_bsh_core, p=$btp_hmny_bsh_periphery"

    # hmny btp address
    btp_hmny_btp_address="btp://$btp_hmny_net/$btp_hmny_bmc_periphery"
    log "    btp: $btp_hmny_btp_address"

    # configuration
    log "\nConfiguring: "

    btp_icon_bmc_owner_wallet="$ixh_tmp_dir/bmc.owner.json"
    btp_icon_bmc_owner_wallet_password="1234"
    btp_icon_nativecoin_bsh_owner_wallet="$ixh_tmp_dir/nativecoin.icon.owner.json"
    btp_icon_nativecoin_bsh_owner_wallet_password="1234"
    btp_icon_bmr_owner_wallet="$ixh_tmp_dir/bmr.icon.json"
    btp_icon_bmr_owner_wallet_password="1234"
    btp_hmny_bmr_owner_wallet="$ixh_tmp_dir/bmr.hmny.json"
    btp_hmny_bmr_owner_wallet_password="1234"
    btp_h2i_relay_config="$ixh_tmp_dir/h2i.config.json"
    btp_i2h_relay_config="$ixh_tmp_dir/i2h.config.json"

    # icon: begin
    log "icon"

    # create and add bmc owner
    log -n "    create_wallet: [$(rel_path "$btp_icon_bmc_owner_wallet")] "
    btp_icon_bmc_owner=$(
        icon_create_wallet "$btp_icon_bmc_owner_wallet" \
        "$btp_icon_bmc_owner_wallet_password" 1000000000000000000000000000)

    log -n "\n    bmc_add_owner: [${btp_icon_bmc_owner:0:10}] "
    _=$(run_jstxcall "$btp_icon_bmc" addOwner 0 "_addr=$btp_icon_bmc_owner")

    # link hmny bmc to icon bmc
    log -n "\n    bmc_link_hmny_bmc: "
    log -n "\n        addLink: "
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
            run_jstxcall "$btp_icon_bmc" addLink 0 "_link=$btp_hmny_btp_address")
    log -n "\n        setLinkRxHeight: "
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
            run_jstxcall "$btp_icon_bmc" setLinkRxHeight 0 "_link=$btp_hmny_btp_address" "_height=$btp_hmny_block_height")
    log -n "\n        getLinkStatus: "
    btp_icon_rx_height=$(hex2dec $(run_jscall "$btp_icon_bmc" getStatus "_link=$btp_hmny_btp_address" | jq -r .rx_height) )
    log -n "rxHeight=$btp_icon_rx_height"

    # add bsh to bmc
    log -n "\n    bmc_add_nativecoin_bsh: "
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
            run_jstxcall "$btp_icon_bmc" addService 0 "_addr=$btp_icon_nativecoin_bsh" "_svc=$btp_icon_nativecoin_bsh_svc_name")

    # create and add nativecoin bsh owner
    log -n "\n    create_wallet: [$(rel_path "$btp_icon_nativecoin_bsh_owner_wallet")] "
    btp_icon_nativecoin_bsh_owner=$(
        icon_create_wallet "$btp_icon_nativecoin_bsh_owner_wallet" \
        "$btp_icon_nativecoin_bsh_owner_wallet_password" 1000000000000000000000000000)

    log -n "\n    nativecoin_bsh_add_owner: [${btp_icon_nativecoin_bsh_owner:0:10}] "
    _=$(run_jstxcall "$btp_icon_nativecoin_bsh" addOwner 0 "_addr=$btp_icon_nativecoin_bsh_owner")

    # register one_dev token
    log -n "\n    nativecoin_bsh_register_irc31: "
    _=$(WALLET=$btp_icon_nativecoin_bsh_owner_wallet \
        PASSWORD=$btp_icon_nativecoin_bsh_owner_wallet_password \
            run_jstxcall "$btp_icon_nativecoin_bsh" register 0 "_name=$btp_hmny_nativecoin_symbol")

    # register relay to bmc
    log -n "\n    create_wallet: [$(rel_path "$btp_icon_bmr_owner_wallet")] "
    btp_icon_bmr_owner=$(
        icon_create_wallet "$btp_icon_bmr_owner_wallet" \
        "$btp_icon_bmr_owner_wallet_password" 1000000000000000000000000000)

    log -n "\n    bmc_add_relay: "
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
            run_jstxcall "$btp_icon_bmc" addRelay 0 "_link=$btp_hmny_btp_address"  "_addr=$btp_icon_bmr_owner")

    # set nativecoinbsh as owner of irc31 token
    log -n "\n    irc31_add_owner: [${btp_icon_nativecoin_bsh:0:10}] "
    _=$(run_jstxcall "$btp_icon_irc31" addOwner 0 "_addr=$btp_icon_nativecoin_bsh")
    log

    # icon: end

    # hmny: begin
    log "hmny"

    function _run_sol() {
        run_sol $@ > /dev/null 2>&1
        log_status $?
    }

    cp $ixh_src_dir/hmny.truffle-config.js $ixh_sol_dir/bmc # replace original truffle-config.js
    cp $ixh_src_dir/hmny.truffle-config.js $ixh_sol_dir/bsh # replace original truffle-config.js
   
    # bmc
    # add bsh
    log -n "    bmc_add_bsh: "
    _run_sol bmc.BMCManagement.addService "'$btp_icon_nativecoin_bsh_svc_name','$btp_hmny_bsh_periphery'"

    # link icon to hmny
    log -n "\n    bmc_link_to_icon_bmc: "
    _run_sol bmc.BMCManagement.addLink "'$btp_icon_btp_address'"
    _run_sol bmc.BMCManagement.setLinkRxHeight "'$btp_icon_btp_address',$btp_icon_block_height"
    # TODO check: response should have one raw logs ?

    # add relay
    log -n "\n    create_wallet: [$(rel_path "$btp_hmny_bmr_owner_wallet")] "
    btp_hmny_bmr_owner=$(hmny_create_wallet "$btp_hmny_bmr_owner_wallet" \
        "$btp_hmny_bmr_owner_wallet_password" 100000000)

    btp_hmny_bmr_owner="0x$(jq -r .address < $btp_hmny_bmr_owner_wallet)"

    log -n "\n    bmc_add_relay: "
    _run_sol bmc.BMCManagement.addRelay "'$btp_icon_btp_address',['$btp_hmny_bmr_owner']"

    # bsh    
    # register icon nativecoin to hmny
    log -n "\n    bsh_register_coin: "
    _run_sol bsh.BSHCore.register "'$btp_icon_nativecoin_symbol'"

    # hmny: end

    # dump relevant variables to be used later
    echo > $ixh_env
    compgen -v | grep btp_ | while read l; do
        echo "$l=${!l}" >> $ixh_env
    done

    log "\n"
    log "deploysc completed in $(( $(date +%s) - $init_start_time ))s."
    log "important variables have been written to $ixh_env"

    # generate btp configs
    generate_relay_configs
}

function _truffle() {
    URI="${URI:-$btp_hmny_uri}" \
    BMC_BTP_NET="${BMC_BTP_NET:-$btp_hmny_net}" \
    PRIVATE_KEY="${PRIVATE_KEY:-$btp_hmny_wallet_private_key}" \
    GASLIMIT="${GASLIMIT:-$btp_hmny_gas_limit}" \
    GASPRICE="${GASPRICE:-$btp_hmny_gas_price}" \
        truffle $@
}

function generate_relay_configs() {
    btp_icon_link_status_rx_height=$btp_hmny_block_height
    btp_hmny_link_status_rx_height=$btp_icon_block_height

    # harmony to icon
    generate_relay_config \
        h2i \
        "icx" \
        "$btp_icon_bmr_owner_wallet" \
        "$btp_icon_bmr_owner_wallet_password" \
        "$btp_hmny_btp_address" \
        "$btp_hmny_uri" \
        "$btp_icon_btp_address" \
        "$btp_icon_uri/default" \
        "$btp_icon_link_status_rx_height" \
            > "$btp_h2i_relay_config"
    # icon to harmony
    generate_relay_config \
        i2h \
        "evm" \
        "$btp_hmny_bmr_owner_wallet" \
        "$btp_hmny_bmr_owner_wallet_password" \
        "$btp_icon_btp_address" \
        "$btp_icon_uri/default" \
        "$btp_hmny_btp_address" \
        "$btp_hmny_uri" \
        "$btp_hmny_link_status_rx_height" \
            > "$btp_i2h_relay_config"
}

function generate_relay_config() {
    prefix="$1"
    cointype="$2"
    keystore_filename="$3"
    key_password="$4"
    src_address="$5"
    src_endpoint="$6"
    dst_address="$7"
    dst_endpoint="$8"
    offset="$9"

    echo "{}" | jq '
        .base_dir = $base_dir |
        .log_level = "debug" |
        .console_level = "trace" |
        .log_forwarder.level = "info" |
        .log_writer.filename = $log_writer_filename |
        .key_store = $key_store |
        .key_store.coinType = $cointype |
        .key_password = $key_password |
        .offset = $offset |
        .src.address = $src_address |
        .src.endpoint = [ $src_endpoint ] |
        .dst.address = $dst_address |
        .dst.endpoint = [ $dst_endpoint ]' \
            --arg base_dir "run/$prefix" \
            --arg log_writer_filename "run/$prefix.log" \
            --argfile key_store "$keystore_filename" \
            --arg cointype "$cointype" \
            --arg key_password "$key_password" \
            --arg src_address "$src_address" \
            --arg src_endpoint "$src_endpoint" \
            --arg dst_address "$dst_address" \
            --arg dst_endpoint "$dst_endpoint" \
            --argjson offset "$offset"
}

function run_sol() {
    if [ $WALLET ]; then
        export PRIVATE_KEY=$(
            ethkey inspect --json --private --passwordfile <(echo "$PASSWORD") "$WALLET" \
        | jq -r .PrivateKey)
    fi

    dsm="$1"
    args="$2"
    IFS='.' read -ra dsm <<< "$dsm"
    cd "$ixh_build_dir/solidity/${dsm[0]}"
    smc="${dsm[1]}"
    mth="${dsm[2]}"
    _truffle exec --network hmny <(echo "
        const smc = artifacts.require('$smc');
        module.exports = async function (callback) {
            try {
                let res = await (await smc.deployed()).$mth($args);
                try {
                    console.log(JSON.stringify(res, null, 2));
                } catch(err) {
                    console.log(res);
                }
            } catch(err) {
                console.error(err);
            } finally { callback(); }
        }") | sed '1d' | sed '1d' # trim first 2 lines
    cd - > /dev/null
}

function run_jscall() {
    scaddr="$1"
    method="$2"
    params=""
    for i in "${@:3}"; do params="$params --param $i"; done
    goloop rpc \
        --uri "$btp_icon_uri" \
        call \
        --to "$scaddr" \
        --method "$method" $params
}

function run_jstxcall() {
    key_store=${WALLET:-$btp_icon_wallet}
    key_password=${PASSWORD:-$btp_icon_wallet_password}

    scaddr=$1
    method=$2
    value=$3
    params=""
    for i in "${@:4}"; do params="$params --param $i"; done
    tx_hash=$(goloop rpc \
        --uri "$btp_icon_uri" \
        sendtx call \
        --to "$scaddr" \
        --key_store "$key_store" \
        --key_password "$key_password" \
        --nid "$btp_icon_nid" \
        --step_limit "$btp_icon_step_limit" \
        --value "$value" \
        --method "$method" $params | jq -r .)
    icon_wait_tx $tx_hash
}

function docker_compose() {
    if [[ ! " ${localhosts[*]} " =~ " ${docker_host} " ]]; then
        export DOCKER_HOST="ssh://$docker_user@$docker_host"
    fi
    env_file=$(mktemp /tmp/ixh.env.XXXXX)
    echo "docker_registry=$docker_registry" > $env_file
    docker-compose -f $ixh_src_dir/docker-compose.yml --env-file $env_file $@
    rm $env_file
}

function stop() {
    docker_compose down
}

function start() {
    if [[ ! " ${localhosts[*]} " =~ " ${docker_host} " ]]; then
        docker_compose pull
    fi
    docker_compose up $@
}

function build_images() {
    repos_dir=$ixh_tmp_dir/repos
    mkdir -p $repos_dir
    
    log "building hmny"
    cd $repos_dir
    # git clone --single-branch --branch ${btp_icon_branch:-main} https://github.com/harmony-one/harmony
    cd $ixh_dir && docker build --build-arg SHARDING_HOST="$docker_host" \
        -f $ixh_src_dir/hmny.Dockerfile -t $docker_registry/hmny:latest .

    log "building icon"
    cd $repos_dir
    git clone --single-branch --branch ${btp_icon_branch:-master} https://github.com/icon-project/goloop
    cd goloop && make gochain-icon-image
    cd $ixh_dir && docker build --build-arg CONFIG_JSON="$(cat $btp_icon_config)" \
        -f $ixh_src_dir/icon.Dockerfile -t $docker_registry/icon:latest .
    
    cd $ixh_dir
}

function publish_images() {
    log "publishing hmny to $docker_registry"
    docker push $docker_registry/hmny:latest
    
    log "publishing icon to $docker_registry"
    docker push $docker_registry/icon:latest
}

function run_test() {
    export verbose=true
    func=$1
    args=( ${@:2} )
    case "$func" in
        iconGetBalance)
            wallet_address=${args[0]}
            params=$(echo '{}' | jq -c '.address = $address' --arg address $wallet_address)
            balance=$(icon_jsonrpc icx_getBalance "$params" | jq -r .result)
            hex2dec $balance
            ;;
        iconGetWrappedCoins)
            run_jscall "$btp_icon_nativecoin_bsh" coinNames
            ;;
        iconRegisterWrappedCoin)
            coinName=${args[0]}
            run_jstxcall "$btp_icon_nativecoin_bsh" register 0 "_name=$coinName"
            ;;
        iconGetWrappedCoinBalance)
            wallet_address=${args[0]}
            coinName=${args[1]}
            coinId=$(run_jscall "$btp_icon_nativecoin_bsh" coinId "_coinName=$coinName" | jq -r .)
            balance=$(run_jscall "$btp_icon_irc31" balanceOf "_owner=$wallet_address" "_id=$coinId" | jq -r .)
            hex2dec $balance
            ;;
        iconTransfer)
            to=${args[1]}
            echo "Not Implemented!" && exit 1
            ;;
        iconTransferNativeCoin)
            value=${args[0]}
            to=${args[1]}
            run_jstxcall "$btp_icon_nativecoin_bsh" transferNativeCoin $value _to=$to
            ;;
        iconTransferWrappedCoin)
            coinName=${args[0]}
            value=${args[1]}
            to=${args[2]}
            run_jstxcall "$btp_icon_nativecoin_bsh" transfer 0 _coinName=$coinName _value=$value _to=$to
            ;;
        iconGetBMCStatus)
            run_jscall "$btp_icon_bmc" getStatus "_link=$btp_hmny_btp_address"
            ;;
        iconBSHIsApprovedForAll)
            wallet_address=${args[0]}
            run_jscall "$btp_icon_irc31" isApprovedForAll "_owner=$wallet_address" "_operator=$btp_icon_nativecoin_bsh"
            ;;
        iconBSHSetApprovalForAll)
            approved=${args[0]:-1}
            run_jstxcall "$btp_icon_irc31" setApprovalForAll 0 "_operator=$btp_icon_nativecoin_bsh" "_approved=$approved"
            ;;
        hmnyGetBalance)
            wallet_address=${args[0]}
            hmny_jsonrpc hmyv2_getBalance "\"$wallet_address\"" | jq -r .result
            ;;
        hmnyGetWrappedCoins)
            run_sol bsh.BSHCore.coinNames
            ;;
        hmnyRegisterWrappedCoin)
            coinName=${args[0]}
            run_sol bsh.BSHCore.register "'$coinName'"
            ;;
        hmnyGetWrappedCoinBalance)
            wallet_address=${args[0]}
            coinName=${args[1]}
            run_sol bsh.BSHCore.getBalanceOf "'$wallet_address','$coinName'"
            ;;
        hmnyTransferNativeCoin)
            value=$(dec2hex ${args[0]})
            to=${args[1]}
            run_sol bsh.BSHCore.transferNativeCoin "'$to',{value:'$value'}"
            ;;
        hmnyTransferWrappedCoin)
            coinName=${args[0]}
            value=$(dec2hex ${args[1]})
            to=${args[2]}
            run_sol bsh.BSHCore.transfer "'$coinName','$value','$to'"
            ;;
        hmnyGetBMCStatus)
            run_sol bmc.BMCPeriphery.getStatus "'$btp_icon_btp_address'"
            ;;
        hmnyBSHIsApprovedForAll)
            wallet_address=${args[0]}
            run_sol bsh.BSHCore.isApprovedForAll "'$wallet_address','$btp_hmny_bsh_core'"
            ;;
        hmnyBSHSetApprovalForAll)
            approved=${args[0]:-1}
            approved=$([[ $approved == 0 ]] && echo false || echo true)
            run_sol bsh.BSHCore.setApprovalForAll "'$btp_hmny_bsh_core',$approved"
            ;;
        hmnyChainStatus)
            hmny_get_hmny_chain_status
            ;;
        *)
            log "Invalid test command: $cmd"
            echo "Usage: $func []"
            exit 1 
            ;;
    esac
}

function run_demo() {
    function tx_relay_wait() {
        sleep 45
    }

    btp_icon_test_wallet=$btp_icon_wallet
    btp_icon_test_wallet_address=$btp_icon_wallet_address
    btp_icon_test_wallet_password=$btp_icon_wallet_password
    btp_hmny_test_wallet=$btp_hmny_wallet
    btp_hmny_test_wallet_address=$btp_hmny_wallet_address
    btp_hmny_test_wallet_password=$btp_hmny_wallet_password


    function get_icon_balance() {
        run_test iconGetBalance $btp_icon_test_wallet_address
    }

    function get_hmny_balance() {
        run_test hmnyGetBalance $btp_hmny_test_wallet_address
    }
    
    function get_icon_wrapped_ONE_DEV() {
        run_test iconGetWrappedCoinBalance $btp_icon_test_wallet_address ONE_DEV
    }

    function get_hmny_wrapped_ICX() {
        hex=$(run_test hmnyGetWrappedCoinBalance $btp_hmny_test_wallet_address ICX | jq -r ._usableBalance)
        hex2dec "0x$hex"
    }

    function show_balances() {
        log
        log "Balance:"
        log "    Icon: $btp_icon_test_wallet_address"
        export icon_balance=$(get_icon_balance)
        log "        Native: $icon_balance"
        export icon_wrapped_ONE_DEV=$(get_icon_wrapped_ONE_DEV)
        log "        Wrapped (ONE_DEV): $icon_wrapped_ONE_DEV"
        log "    Hmny: $btp_hmny_test_wallet_address"
        export hmny_balance=$(get_hmny_balance)
        log "        Native: $hmny_balance"
        export hmny_wrapped_ICX=$(get_hmny_wrapped_ICX)
        log "        Wrapped (ICX): $hmny_wrapped_ICX"
        log
    }

    log "Icon Wrapped Coins:"
    log "    $(run_test iconGetWrappedCoins | jq -c .)"

    log "Hmny Wrapped Coins:"
    log "    $(run_test hmnyGetWrappedCoins | jq -c .)"

    show_balances

    ixh_nativecoin_transfer_amount=$(python3 -c "print($icon_balance//3)")
    log "TransferNativeCoin (Icon -> Hmny):"
    log "    amount=$ixh_nativecoin_transfer_amount"
    log -n "    "
    run_test iconTransferNativeCoin $ixh_nativecoin_transfer_amount "btp://$btp_hmny_net/$btp_hmny_test_wallet_address"
    log

    tx_relay_wait

    show_balances

    h2i_nativecoin_transfer_amount=$(python3 -c "print($hmny_balance//3)")
    log "TransferNativeCoin (Hmny -> Icon):"
    log "    amount=$h2i_nativecoin_transfer_amount"
    log -n "    "
    run_test hmnyTransferNativeCoin $h2i_nativecoin_transfer_amount "btp://$btp_icon_net/$btp_icon_test_wallet_address"
    log
    
    tx_relay_wait

    show_balances

    log "Approve Icon NativeCoinBSH"
    log -n "    "
    WALLET=$btp_icon_test_wallet PASSWORD=$btp_icon_test_wallet_password run_test iconBSHSetApprovalForAll 1
    log
    log "    Status: $(run_test iconBSHIsApprovedForAll $btp_icon_test_wallet_address)"

    log "Approve Hmny BSHCore"
    WALLET=$btp_hmny_test_wallet PASSWORD=$btp_hmny_test_wallet_password run_test hmnyBSHSetApprovalForAll 1
    log "    Status: $(run_test hmnyBSHIsApprovedForAll $btp_hmny_test_wallet_address)"
    log

    h2i_wrapped_ICX_transfer_amount=$(python3 -c "print($hmny_wrapped_ICX//2)")
    log "TransferWrappedCoin ICX (Hmny -> Icon):"
    log "    amount=$h2i_wrapped_ICX_transfer_amount"
    log -n "    "
    WALLET=$btp_hmny_test_wallet PASSWORD=$btp_hmny_test_wallet_password \
        run_test hmnyTransferWrappedCoin ICX $h2i_wrapped_ICX_transfer_amount "btp://$btp_icon_net/$btp_icon_test_wallet_address"

    tx_relay_wait

    show_balances

    ixh_wrapped_ONE_DEV_transfer_amount=$(python3 -c "print($icon_wrapped_ONE_DEV//2)")
    log "TransferWrapped Coin ONE_DEV (Icon -> Hmny):"
    log "    amount=$ixh_wrapped_ONE_DEV_transfer_amount"
    log -n "    "
    WALLET=$btp_icon_test_wallet PASSWORD=$btp_icon_test_wallet_password \
        run_test iconTransferWrappedCoin ONE_DEV $ixh_wrapped_ONE_DEV_transfer_amount "btp://$btp_hmny_net/$btp_hmny_test_wallet_address"
    log

    tx_relay_wait

    show_balances
}


# trap cleanup SIGINT SIGTERM

function usage() {
  echo "Usage: $(basename $0) [build|publish|deploysc|start|stop]"
  exit 1
}

if [ $# -gt 0 ]; then cmd=$1; else usage; fi

args=( ${@:2} )

case "$cmd" in
    start) 
        start ${args[@]}
        ;;

    stop)
        stop ${args[@]}
        ;;

    docker_compose)
        docker_compose ${args[@]}
        ;;

    build)
        build_images ${args[@]}
        ;;

    publish) 
        publish_images ${args[@]}
        ;;

    deploysc)
        if [[ ! " ${allowedhosts[*]} " =~ " ${docker_host} " ]]; then
            echo "docker_host: $docker_host not in allowedhosts!"
            exit 1
        fi
        deploysc ${args[@]}
        ;;

    generateRelayConfigs)
        . $ixh_env
        generate_relay_configs
        ;;

    clean)
        rm -rf $ixh_tmp_dir
        ;;

    fishenv)
        cat $ixh_env | sed 's/=/ /1' | sed -E 's/(.*)/set \1/'
        ;;

    bashenv)
        cat $ixh_env
        ;;

    sol)
        . $ixh_env
        run_sol ${args[@]}
        ;;

    jscall)
        . $ixh_env
        run_jscall ${args[@]}
        ;;

    jstxcall)
        . $ixh_env
        run_jstxcall ${args[@]}
        ;;

    test)
        . $ixh_env
        run_test ${args[@]}
        ;;
        
    demo)
        . $ixh_env
        run_demo ${args[@]}
        ;;
    
    *)
        log "Invalid command: $cmd"
        usage 
        ;;
esac

# close log file descriptor (fd 3)
exec 3>&-