#!/usr/bin/env -S bash -eET -o pipefail
#-O inherit_errexit

function dec2hex() {
    hex=$(echo "obase=16; ibase=10; ${@}" | bc)
    echo "0x${hex,,}"
}

function hex2dec() {
    hex=${@#0x}
    echo "obase=10; ibase=16; ${hex^^}" | bc
}

function repeat() {
    for i in $(seq 1 $2); do echo -n "$1"; done
}

function rel_path() {
    realpath --relative-to="$ixh_dir" "$1"
}

# echo message to stderr
function log() {
    local prefix="$(date '+%Y-%m-%d %H:%M:%S') $(repeat '    ' $((${#FUNCNAME[@]} - 2)))"
    echo -e "$prefix$@" >&2
}

function log_status() {
    [[ "$1" == 0 ]] && log " ✔" || log " ✘"
}

function log_stack() {
    local cmd=${FUNCNAME[1]}
    if [[ $# > 0 ]]; then cmd="$@"; fi
    local prefix="$(date '+%Y-%m-%d %H:%M:%S') $(repeat '    ' $((${#FUNCNAME[@]} - 3)))"
    echo -e "$prefix$cmd():${BASH_LINENO[1]}" >&2
}

function require() {
    log_stack
    [[ -z "$1" ]] && log "$3" && exit 1
    [[ -z "$2" ]] || log "$2"
}

function require_integer() {
    log_stack
    local integer=
    [[ "$1" =~ ^[0-9]+$ ]] && integer=$1
    require "$integer" "$2" "$3 (invalid integer:'$1')"
}

function require_address() {
    log_stack
    local address=
    [[ "$1" =~ ^[0hc][xX][0-9a-fA-F]{40}$ ]] && address=$1
    require "$address" "$2" "$3 (invalid address:'$1')"
}

function require_existsdir() {
    log_stack
    local dir=
    [[ -d "$1" ]] && dir=$1
    require "$dir" "$2" "$3 (dir does not exist:'$1')"
}

# override commands to write stderr to log file
shopt -s expand_aliases
function run() {
    cmd="$1"
    args=("${@:2}")
    log_stack "$cmd"
    # local indent="$(repeat '    ' $((${#FUNCNAME[@]} - 1)))"
    # $cmd "${args[@]}" 2> >(sed "s/^/$indent/" >&2) | tee >(sed "s/^/$indent/" >&2)
    local prefix="$(date '+%Y-%m-%d %H:%M:%S') $(repeat '    ' $((${#FUNCNAME[@]} - 1)))"
    { { $cmd "${args[@]}" 2>&3 | tee >(sed "s/^/$prefix/" >&2); } 3>&1 >&4 | sed "s/^/$prefix/" >&2; } 4>&1
}
alias jq="run jq"
alias yarn="run yarn"
alias curl="run curl"
alias gradle="run gradle"
alias goloop="run goloop"
alias ethkey="run ethkey"
alias truffle="run truffle"
alias docker="run docker"

# ethkey_get_private_key() <wallet.json> <password>
function ethkey_get_private_key() {
    log_stack
    if [ -z "$1" ]; then
        log "invalid <wallet.json>"
        exit 1
    fi
    ethkey inspect --json --private \
        --passwordfile <(echo "$2") "$1" | jq -r .PrivateKey
}

# hmny_jsonrpc [method] [arguments_in_json]
function hmny_jsonrpc() {
    log_stack
    curl "$btp_hmny_uri" -s -X POST \
        -H 'Content-Type:application/json' \
        -d "$(jq <<<{} -c \
            '.id=1|.jsonrpc="2.0"|.method=$method|.params=$params' \
            --arg method "$1" --argjson params "$2")"
}

# icon_jsonrpc [method] [arguments_in_json]
function icon_jsonrpc() {
    log_stack
    curl "$btp_icon_uri" -s -X POST \
        -H 'Content-Type:application/json' \
        -d "$(jq <<<{} -c \
            '.id=1|.jsonrpc="2.0"|.method=$method|.params=$params' \
            --arg method "$1" --argjson params "$2")"
}

# WALLET=<wallet.json> \
# PASSWORD=<password> \
#     xxxx_transfer [address] [balance]
function validate_transfer() {
    log_stack
    WALLET=${WALLET:-}
    require "$WALLET" "WALLET='$WALLET'" "invalid WALLET='$WALLET'"
    local address=$1
    require_address "$address" "address='$address'" "failed"
    local balance=$2
    require_integer "$balance" "" "must be an integer: $balance"
}

function icon_wait_tx() {
    log_stack
    local ret=1
    local tx_hash=$1
    [[ -z $tx_hash ]] && return
    log "[txh=${tx_hash}]"
    while :; do
        goloop rpc \
            --uri "$btp_icon_uri" \
            txresult "$tx_hash" &>/dev/null && break || sleep 1
    done
    local txr=$(goloop rpc --uri "$btp_icon_uri" txresult "$tx_hash" 2>/dev/null)
    local status=$(jq <<<"$txr" -r .status)
    [[ "$status" == 0x0 ]] && echo $txr
    [[ "$status" == 0x1 ]] && echo $txr && ret=0
    return $ret
}

function icon_callsc() {
    log_stack

    local address=$1
    require_address "$address" "address: $address" "icon_callsc"

    local method=$2
    require "$method" "method: $method" "invalid method: '$method'"

    local params=()
    for i in "${@:3}"; do params+=("--param $i"); done

    goloop rpc \
        --uri "$btp_icon_uri" \
        call \
        --to "$address" \
        --method "$method" ${params[@]}
}

# WALLET=<wallet.json> PASSWORD=<password> \
#     icon_sendtx_call [address] [method] [value] [params]..."
function icon_sendtx_call() {
    log_stack

    WALLET=${WALLET:-}
    PASSWORD=${PASSWORD:-}

    require "$WALLET" "WALLET='$WALLET'" "invalid WALLET='$WALLET'"

    local address="$1"
    require_address "$address" "address: $address" "icon_sendtx_call"

    local method="$2"
    require "$method" "method: $method" "invalid method: '$method'"

    local value="$3"
    [[ -z "$value" ]] ||
        require_integer "$value" "value: $value" "invalid value: '$value'"

    local params=()
    for i in "${@:4}"; do params+=("--param $i"); done

    local tx_hash=$(
        goloop rpc \
            --uri "$btp_icon_uri" \
            sendtx call \
            --to "$address" \
            --key_store "$WALLET" \
            --key_password "$PASSWORD" \
            --nid "$btp_icon_nid" \
            --step_limit "$btp_icon_step_limit" \
            --value "$value" \
            --method "$method" \
            ${params[@]} | jq -r .
    )
    icon_wait_tx "$tx_hash"
}

# WALLET=<wallet.json> PASSWORD=<password> \
#      icon_sendtx_deploy <sc.jar> [params]...
function icon_sendtx_deploy() {
    log_stack

    local scfile=$1
    require "$scfile" "$scfile" "invalid sc.jar: '$scfile'"

    local params=()
    for i in "${@:2}"; do params+=("--param $i"); done

    tx_hash=$(
        goloop rpc \
            --uri "$btp_icon_uri" \
            sendtx deploy $scfile \
            --key_store "$WALLET" \
            --key_password "$PASSWORD" \
            --nid "$btp_icon_nid" \
            --content_type application/java \
            --step_limit "$btp_icon_step_limit" \
            ${params[@]} | jq -r .
    )
    icon_wait_tx "$tx_hash"
}

# WALLET=<wallet.json> PASSWORD=<password> \
#      icon_sendtx_transfer [address] [value]
function icon_sendtx_transfer() {
    log_stack

    WALLET=${WALLET:-}
    PASSWORD=${PASSWORD:-}

    require "$WALLET" "WALLET='$WALLET'" "invalid WALLET='$WALLET'"

    local address="$1"
    require_address "$address" "address: $address" "icon_sendtx_transfer"

    local value="$2"
    require "$value" "value: $value" "icon_sendtx_transfer: invalid value: '$value'"

    local tx_hash=$(
        goloop rpc \
            --uri "$btp_icon_uri" \
            sendtx transfer \
            --to "$address" \
            --value "$value" \
            --key_store "$WALLET" \
            --key_password "$PASSWORD" \
            --nid "$btp_icon_nid" \
            --step_limit "$btp_icon_step_limit" | jq -r .
    )
    icon_wait_tx "$tx_hash"
}

function icon_transfer() {
    log_stack
    type=icon validate_transfer "$@"
    local address=$1
    local balance=$2
    if [ $balance == 0 ]; then return 0; fi
    local txr=$(icon_sendtx_transfer "$address" $balance)
    local status=$(jq <<<"$txr" -r .status)
    [[ "$status" == 0x1 ]] || status=""
    require "$status" "" "icon_transfer: failed to transfer balance to $address!"
}

function icon_deploysc() {
    log_stack
    icon_sendtx_deploy "$@"
}

# icon_create_wallet [keystore] [password] [balance]
function icon_create_wallet() {
    log_stack

    local keystore=$1
    local password=$2
    local balance=$3

    require "$keystore" "keystore: $keystore" "icon_create_wallet: invalid keystore: $keystore"

    goloop &>/dev/null ks gen -o "$keystore" -p "$password"

    local address=$(jq -r .address "$keystore")
    require_address "$address" "address: $address" "icon_create_wallet: failed"

    icon_transfer "$address" "$balance"

    echo "$address"
}

function hmny_save_god_wallet() {
    log_stack
    local keystore=$1
    local password=$2
    local private_key=$3
    {
        read address
        read keystore_content
    } <<<$(
        echo "
            hmy keys import-private-key '$private_key' root
            keystore=\$(hmy keys export-ks root / 2>/dev/null | cut -d/ -f2)
            echo \$keystore | cut -d. -f1
            cat /\$keystore
        " |
            docker run -i --rm --network=host \
                $docker_registry/hmny:latest /bin/bash
    )
    cat >"$keystore" <<<$keystore_content
    echo $address
}

function _truffle() {
    log_stack
    PRIVATE_KEY="${PRIVATE_KEY:-$btp_hmny_dummy_private_key}" \
        NETWORK_ID="${NETWORK_ID:-$btp_hmny_nid}" \
        URI="${URI:-$btp_hmny_uri}" \
        BMC_BTP_NET="${BMC_BTP_NET:-$btp_hmny_net}" \
        GASLIMIT="${GASLIMIT:-$btp_hmny_gas_limit}" \
        GASPRICE="${GASPRICE:-$btp_hmny_gas_price}" \
        truffle "$@"
}

# _truffle_exec [contract.method] [arguments]
function _truffle_exec() {
    log_stack

    WALLET=${WALLET:-}
    PASSWORD=${PASSWORD:-}

    if [ $WALLET ]; then
        export PRIVATE_KEY=$(ethkey_get_private_key "$WALLET" "$PASSWORD")
    fi

    scdir="$1"
    require_existsdir "$scdir" "scdir: $scdir" "invalid scdir: '$scdir'"

    IFS='.' read -ra dsm <<<"$2"
    if [ "${#dsm[@]}" != 2 ]; then
        log "_truffle_exec: invalid contract.method: $2"
        return 1
    fi
    smc="${dsm[0]}"
    mth="${dsm[1]}"

    args="${3:-}"

    cd $scdir
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
    cd $OLDPWD
}

function is_address() {
    [[ "$1" =~ ^[0hc][xX][0-9a-fA-F]{40}$ ]]
}

function hmny_bech32_address() {
    log_stack
    local address=$1
    if [[ "$address" =~ ^[0hc][xX][0-9a-fA-F]{40}$ ]]; then
        address=$(
            echo "hmy utility addr-to-bech32 $address" |
                docker run -i --rm --network=host \
                    $docker_registry/hmny:latest /bin/bash
        )
    fi
    echo $address
}

function hmny_transfer() {
    log_stack
    type=hmny validate_transfer "$@"

    WALLET=${WALLET:-}
    PASSWORD=${PASSWORD:-}

    local address=$(hmny_bech32_address $1)
    local balance=$(bc <<<"scale=18;$2/10^18")

    local private_key=$(ethkey_get_private_key "$WALLET" "$PASSWORD")
    local txr=$(
        echo "
            src_address=\$(hmy keys import-private-key \
                '$private_key' root | tail -n1 | cut -d: -f2 | xargs)
            hmy -n '$btp_hmny_uri' --no-pretty  \
                transfer \
                --from=\$src_address --to=$address \
                --chain-id=$btp_hmny_chain_id \
                --from-shard=0 --to-shard=0 --timeout=120 \
                --amount=$balance
        " |
            docker run -i --rm --network=host \
                $docker_registry/hmny:latest /bin/bash
    )
    local status=$(jq <<<"$txr" -r '[.[]."blockchain-receipt".status][0]')
    [[ "$status" == 0x1 ]] || status=""
    require "$status" "" "hmny_transfer: failed to transfer balance to $address!"
}

function hmny_deploysc() {
    log_stack

    if [ $WALLET ]; then
        export PRIVATE_KEY=$(ethkey_get_private_key "$WALLET" "$PASSWORD")
    fi
    require "$PRIVATE_KEY" "PRIVATE_KEY: '$PRIVATE_KEY'" "invalid PRIVATE_KEY='$PRIVATE_KEY'"

    scdir=$1
    require "$scdir" "scdir: $scdir" "hmny_deploysc: invalid scdir: '$scdir'"
    cd $scdir

    # replace original truffle-config.js
    cp $ixh_src_dir/hmny.truffle-config.js truffle-config.js
    yarn install --silent &>/dev/null # download node_modules
    _truffle compile --all --quiet &>/dev/null

    for _ in $(# repeat until successfully deployed
        seq 1 20
    ); do
        _truffle migrate --network hmny --compile-none --skip-dry-run >/dev/null
        if [ $? == 0 ]; then
            local imports=""
            local waiters=""
            for i in "${@:2}"; do
                imports="$imports const $i = artifacts.require('$i');"
                waiters="$waiters await $i.deployed(); console.log($i.address);"
            done
            _truffle exec --network hmny <(echo "
                $imports
                module.exports = async function(cb) { $waiters; cb(); }
            ") | tail -n $(($# - 1))
            break
        fi
    done
    cd $OLDPWD
}

# hmny_create_wallet [keystore] [password] [balance]
function hmny_create_wallet() {
    log_stack

    local keystore=$1
    local password=$2
    local balance=$3

    require "$keystore" "keystore: $keystore" "invalid keystore: $keystore"

    local address
    local keystore_content

    {
        read address
        read keystore_content
    } <<<$(
        echo "
            hmy keys add mykey --passphrase-file <(cat<<<'1234') 2>&1 > /dev/null
            keystore=\$(hmy keys export-ks --passphrase-file \
                <(cat<<<'1234') mykey / 2>/dev/null | cut -d/ -f2)
            echo \$keystore | cut -d. -f1
            cat /\$keystore && echo
        " |
            docker run -i --rm --network=host \
                $docker_registry/hmny:latest /bin/bash
    )

    address="0x$(jq <<<$keystore_content -r .address)"
    require_address "$address" "address=$address" "failed"

    # write keystore_content to given file
    cat <<<$keystore_content >$keystore

    hmny_transfer $address $balance

    echo $address
}

# run_sol [dir.contract.method] [arguments]
function run_sol() {
    log_stack
    IFS='.' read -ra dsm <<<"$1"
    if [ "${#dsm[@]}" != 3 ]; then
        log "run_sol: invalid dir.contract.method: $1"
        return 1
    fi
    _truffle_exec "$ixh_build_dir/solidity/${dsm[0]}" "${dsm[1]}.${dsm[2]}" "${@:2}"
}

function hmny_get_hmny_chain_status() {
    log_stack
    local bn=$(hmny_jsonrpc hmyv2_blockNumber "[]" | jq -r .result)
    ((bn--))
    local lh=$(hmny_jsonrpc hmyv2_getBlockByNumber "[$bn,{}]" | jq -r .result.hash)
    local fh=$(hmny_jsonrpc hmyv2_getFullHeader "[$(($bn + 1))]")
    local ep=$(jq <<<"$fh" -r '.result.epoch')
    local bitmap=$(jq <<<"$fh" -r '.result.lastCommitBitmap')
    local signature=$(jq <<<"$fh" -r '.result.lastCommitSignature')
    local lb=$(hmny_jsonrpc hmyv2_epochLastBlock "[$(($ep - 1))]" | jq -r '.result')
    local ss=$(hmny_jsonrpc hmyv2_getFullHeader "[$lb]" | jq -r '.result.shardState')
    echo -e "$bn\n$lh\n$ep\n$ss\n$bitmap\n$signature"
}

function ensure_wallet_minimum_balance() {
    log_stack
    [[ -z "$btp_icon_wallet" ]] && icon_create_wallet \
        "$btp_icon_wallet" "$btp_icon_wallet_password" $btp_icon_wallet_minimum_balance

    [[ -z "$btp_hmny_wallet" ]] && hmny_create_wallet \
        "$btp_hmny_wallet" "$btp_hmny_wallet_password" $btp_hmny_wallet_minimum_balance

    ibal=$(run_exec iconGetBalance $btp_icon_wallet)
    if [[ "$ibal" < "$btp_icon_wallet_minimum_balance" ]]; then
        # transfer balance from god wallet
        WALLET=$btp_icon_god_wallet PASSWORD=$btp_icon_god_wallet_password \
            icon_transfer "$btp_icon_wallet_address" $btp_icon_wallet_minimum_balance
        log
    fi
    hbal=$(run_exec hmnyGetBalance $btp_hmny_wallet)
    if [[ "$hbal" < "$btp_hmny_wallet_minimum_balance" ]]; then
        # transfer balance from god wallet
        WALLET=$btp_hmny_god_wallet PASSWORD=$btp_hmny_god_wallet_password \
            hmny_transfer "$btp_hmny_wallet_address" $btp_hmny_wallet_minimum_balance
    fi
}

# Ensure following tools are installed
# gradle, jdk@11.x, sdkman, goloop, docker, truffle@5.3.0, node@15.12.0, yarn, ethkey
function deploysc() {
    log_stack

    export init_start_time=$(date +%s)
    function save_important_variables() {
        echo >>$ixh_env
        compgen -v | grep btp_ | while read l; do
            echo "$l='${!l}'" >>$ixh_env
        done
        log "\n"
        log "deploysc completed in $(($(date +%s) - $init_start_time))s."
        log "important variables have been written to $ixh_env"
    }
    trap "save_important_variables" EXIT

    # build dir
    mkdir -p $ixh_build_dir

    # create root wallets
    log "Wallet:"

    # icon
    log "icon: [$(rel_path "$btp_icon_wallet")]"
    btp_icon_wallet_address=$(jq -r .address "$btp_icon_wallet")

    # hmny
    log "hmny: [$(rel_path "$btp_hmny_wallet")]"
    btp_hmny_wallet_address="0x$(jq -r .address "$btp_hmny_wallet")"

    # prepare javascore build dir
    cp -r $root_dir/javascore $ixh_jsc_dir

    # build javascores
    log "Build: "
    log "javascores:"

    cd $ixh_jsc_dir
    gradle clean
    gradle bmc:optimizedJar
    gradle bsr:optimizedJar
    gradle bts:optimizedJar
    # gradle fee_aggregation:optimizedJar

    mkdir -p dist
    cp bmc/build/libs/bmc-0.1.0-optimized.jar dist/bmc.jar
    cp bsr/build/libs/restrictions-0.1.0-optimized.jar dist/bsr.jar
    cp bts/build/libs/bts-0.1.0-optimized.jar dist/bts.jar
    cp lib/irc2Tradeable-0.1.0-optimized.jar dist/irc2Tradeable.jar
    # cp fee_aggregation/build/libs/fee_aggregation-1.0-optimized.jar dist/fee_aggregator.jar

    cd $ixh_build_dir
    git clone https://github.com/icon-project/java-score-examples.git
    cd java-score-examples
    gradle irc2-token:clean
    gradle irc2-token:optimizedJar
    cp irc2-token/build/libs/irc2-token-0.9.1-optimized.jar $ixh_jsc_dir/dist/irc2.jar

    log

    log "Deploy: "

    # deploy
    log "icon"

    # contracts: being
    # bmc
    if [ -z "$btp_icon_bmc" ]; then
        log "bmc:"
        r=$(WALLET=$btp_icon_wallet \
            PASSWORD=$btp_icon_wallet_password \
            icon_deploysc \
            $ixh_jsc_dir/dist/bmc.jar \
            _net="$btp_icon_net")
        btp_icon_bmc=$(jq -r .scoreAddress <<<$r)
        btp_icon_block_hash=$(jq -r .blockHash <<<$r)
        btp_icon_block_height=$(hex2dec $(jq -r .blockHeight <<<$r))
    fi

    btp_icon_validators_hash=$(
        URI=$btp_icon_uri \
            HEIGHT=$(dec2hex $(($btp_icon_block_height - 1))) \
            $ixh_dir/src/iconvalidators | jq -r .hash
    )

    # bsr
    if [ -z "$btp_icon_bsr" ]; then
        log "bsr:"
        r=$(WALLET=$btp_icon_wallet \
            PASSWORD=$btp_icon_wallet_password \
            icon_deploysc \
            $ixh_jsc_dir/dist/bsr.jar)
        btp_icon_bsr=$(jq -r .scoreAddress <<<$r)
    fi

    # bts
    if [ -z "$btp_icon_bts" ]; then
        log "bts: "
        irc2Tradeable_score=$(xxd -p $ixh_jsc_dir/dist/irc2Tradeable.jar | tr -d '\n')
        r=$(WALLET=$btp_icon_wallet \
            PASSWORD=$btp_icon_wallet_password \
            icon_deploysc \
            $ixh_jsc_dir/dist/bts.jar \
            _name="ICX" \
            _bmc="$btp_icon_bmc" \
            _serializedIrc2="$irc2Tradeable_score")
        btp_icon_bts=$(jq -r .scoreAddress <<<$r)
    fi

    ## bts: irc2 (TICX)
    if [ -z "$btp_icon_ticx" ]; then
        log "irc2: "
        r=$(WALLET=$btp_icon_wallet \
            PASSWORD=$btp_icon_wallet_password \
            icon_deploysc \
            $ixh_jsc_dir/dist/irc2.jar \
            _name="TICX" \
            _symbol="TICX" \
            _decimals="0x12" \
            _initialSupply="0x186a0") # 100000 TICX
        btp_icon_ticx=$(jq -r .scoreAddress <<<$r)
        log "$btp_icon_ticx"
    fi

    # # fee aggregation
    # log "fee aggregation: "
    # r=$(WALLET=$btp_icon_wallet \
    #     PASSWORD=$btp_icon_wallet_password \
    #     icon_deploysc \
    #     $ixh_jsc_dir/dist/fee_aggregator.jar \
    #     _cps_address="$btp_icon_cps_address" \
    #     _band_protocol_address="$btp_icon_band_protocol_address")
    # btp_icon_fee_aggregator=$(jq -r .scoreAddress <<<$r)
    # log "$btp_icon_fee_aggregator"

    # contracts: end

    # icon btp address
    btp_icon_btp_address="btp://$btp_icon_net/$btp_icon_bmc"
    log "btp: $btp_icon_btp_address"

    # configuration: begin
    if [ -z "$btp_icon_bmc_owner_wallet" ]; then
        btp_icon_bmc_owner_wallet="$ixh_tmp_dir/bmc.owner.json"
        btp_icon_bmc_owner_wallet_password="1234"

        # create and add bmc owner
        log "create_wallet: [$(rel_path "$btp_icon_bmc_owner_wallet")] "
        btp_icon_bmc_owner=$(
            WALLET=$btp_icon_wallet PASSWORD=$btp_icon_wallet_password \
                icon_create_wallet "$btp_icon_bmc_owner_wallet" \
                "$btp_icon_bmc_owner_wallet_password" $btp_icon_bmc_owner_balance
        )
    fi
    btp_icon_bmc_owner=$(jq -r .address <$btp_icon_bmc_owner_wallet)

    local is_owner=$(icon_callsc "$btp_icon_bmc" \
        isOwner "_addr=$btp_icon_bmc_owner" | jq -r .)
    if [ "$is_owner" == "0x0" ]; then
        log "bmc_add_owner: [${btp_icon_bmc_owner:0:10}] "
        _=$(WALLET=$btp_icon_wallet PASSWORD=$btp_icon_wallet_password \
            icon_sendtx_call "$btp_icon_bmc" addOwner 0 "_addr=$btp_icon_bmc_owner")
    fi

    if [ ! -z "$btp_icon_fee_aggregator" ]; then
        log "bmc_set_fee_aggregator:"
        _=$(WALLET=$btp_icon_bmc_owner_wallet \
            PASSWORD=$btp_icon_bmc_owner_wallet_password \
            icon_sendtx_call \
            "$btp_icon_bmc" setFeeAggregator 0 "_addr=$btp_icon_fee_aggregator")

        log "bmc_set_fee_gathering_term:"
        _=$(WALLET=$btp_icon_bmc_owner_wallet \
            PASSWORD=$btp_icon_bmc_owner_wallet_password \
            icon_sendtx_call \
            "$btp_icon_bmc" setFeeGatheringTerm 0 "_value=1000") # every 1000 blocks
    fi

    # add bts to bmc
    log "bmc_add_bts: "
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
        icon_sendtx_call "$btp_icon_bmc" addService 0 "_addr=$btp_icon_bts" "_svc=$btp_bts_svc_name")

    if [ -z "$btp_icon_bts_owner_wallet" ]; then
        btp_icon_bts_owner_wallet="$ixh_tmp_dir/bts.owner.json"
        btp_icon_bts_owner_wallet_password="1234"

        # create and add bts owner
        log "create_wallet: [$(rel_path "$btp_icon_bts_owner_wallet")] "
        btp_icon_bts_owner=$(
            WALLET=$btp_icon_wallet PASSWORD=$btp_icon_wallet_password \
                icon_create_wallet "$btp_icon_bts_owner_wallet" \
                "$btp_icon_bts_owner_wallet_password" $btp_icon_bts_owner_balance
        )
    fi
    btp_icon_bts_owner=$(jq -r .address <$btp_icon_bts_owner_wallet)

    local is_owner=$(icon_callsc "$btp_icon_bts" \
        isOwner "_addr=$btp_icon_bts_owner" | jq -r .)
    if [ "$is_owner" == "0x0" ]; then
        log "bts_add_owner: [${btp_icon_bts_owner:0:10}] "
        _=$(WALLET=$btp_icon_wallet PASSWORD=$btp_icon_wallet_password \
            icon_sendtx_call "$btp_icon_bts" addOwner 0 "_addr=$btp_icon_bts_owner")
    fi

    # set bsr in bts
    log "bts_set_bsr: "
    _=$(WALLET=$btp_icon_bts_owner_wallet \
        PASSWORD=$btp_icon_bts_owner_wallet_password \
        icon_sendtx_call "$btp_icon_bts" addRestrictor 0 "_address=$btp_icon_bsr")

    # bts set fee ratio
    log "bts set ICX fee: "
    _=$(WALLET=$btp_icon_bts_owner_wallet \
        PASSWORD=$btp_icon_bts_owner_wallet_password \
        icon_sendtx_call "$btp_icon_bts" \
        setFeeRatio 0 _name="ICX" _feeNumerator=$(dec2hex $btp_bts_fee_numerator) _fixedFee=$(dec2hex $btp_bts_fixed_fee))

    # configuration: end

    # hmny
    cp -r $root_dir/solidity $ixh_sol_dir

    # deploy
    log "hmny"

    cp $ixh_src_dir/hmny.truffle-config.js $ixh_sol_dir/bmc # replace original truffle-config.js
    cp $ixh_src_dir/hmny.truffle-config.js $ixh_sol_dir/bts # replace original truffle-config.js

    # before bmc
    {
        read btp_hmny_block_height
        read btp_hmny_block_hash
        read btp_hmny_block_epoch
        read btp_hmny_shard_state
        read btp_hmny_verifier_commit_bitmap
        read btp_hmny_verifier_commit_signature
    } <<<"$(hmny_get_hmny_chain_status)"

    # bmc
    if [ -z "$btp_hmny_bmc_periphery" ] || [ -z "$btp_hmny_bmc_management"]; then
        log "bmc: "
        {
            read btp_hmny_bmc_management
            read btp_hmny_bmc_periphery
        } <<<$(
            WALLET=$btp_hmny_wallet \
                PASSWORD=$btp_hmny_wallet_password \
                BMC_BTP_NET="$btp_hmny_net" \
                hmny_deploysc $ixh_sol_dir/bmc BMCManagement BMCPeriphery
        )
        log "m=$btp_hmny_bmc_management, p=$btp_hmny_bmc_periphery"
    # TODO get hmny bmc block height and epoch in hex (0x...)
    fi

    # bts
    if [ -z "$btp_hmny_bts_core" ] || [ -z "$btp_hmny_bts_periphery"]; then
        log "bts: "
        {
            read btp_hmny_bts_core
            read btp_hmny_bts_periphery
            read btp_hmny_tone
        } <<<$(
            WALLET=$btp_hmny_wallet \
                PASSWORD=$btp_hmny_wallet_password \
                BSH_COIN_URL="https://github.com/icon/btp" \
                BSH_COIN_NAME="ONE" \
                BSH_COIN_FEE=$btp_bts_fee_numerator \
                BSH_FIXED_FEE=$btp_bts_fixed_fee \
                BMC_PERIPHERY_ADDRESS="$btp_hmny_bmc_periphery" \
                BSH_SERVICE="$btp_bts_svc_name" \
                hmny_deploysc $ixh_sol_dir/bts BTSCore BTSPeriphery HRC20
        )
        log "core=$btp_hmny_bts_core, periphery=$btp_hmny_bts_periphery, erc20=$btp_hmny_tone"
    fi

    # hmny btp address
    btp_hmny_btp_address="btp://$btp_hmny_net/$btp_hmny_bmc_periphery"
    log "btp: $btp_hmny_btp_address"

    # configuration: begin

    log "bmc_add_bts: "
    WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
        run_sol >/dev/null \
        bmc.BMCManagement.addService \
        "'$btp_bts_svc_name','$btp_hmny_bts_periphery'"

    log "bts setFeeRatio:"
    WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
        run_sol >/dev/null \
        bts.BTSCore.setFeeRatio "'ONE','$btp_bts_fee_numerator','$btp_bts_fixed_fee'"

    # configuration: end

    log "Configure Links: "

    # icon: begin
    log "icon"

    # link hmny bmc to icon bmc
    log "BMC: Add Link to HMNY BMC: "
    log "addLink: "
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
        icon_sendtx_call "$btp_icon_bmc" addLink 0 "_link=$btp_hmny_btp_address")
    log "setLinkRxHeight: "
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
        icon_sendtx_call "$btp_icon_bmc" setLinkRxHeight 0 "_link=$btp_hmny_btp_address" "_height=$btp_hmny_block_height")
    log "getLinkStatus: "
    btp_icon_rx_height=$(hex2dec $(icon_callsc "$btp_icon_bmc" getStatus "_link=$btp_hmny_btp_address" | jq -r .rx_height))
    log "btp_icon_rx_height=$btp_icon_rx_height"

    # register relay in bmc
    if [ -z "$btp_icon_bmr_owner" ]; then
        btp_icon_bmr_owner_wallet="$ixh_tmp_dir/bmr.icon.json"
        btp_icon_bmr_owner_wallet_password="1234"

        log "create_wallet: [$(rel_path "$btp_icon_bmr_owner_wallet")] "
        btp_icon_bmr_owner=$(
            WALLET=$btp_icon_wallet PASSWORD=$btp_icon_wallet_password \
                icon_create_wallet "$btp_icon_bmr_owner_wallet" \
                "$btp_icon_bmr_owner_wallet_password" $btp_icon_bmr_owner_balance
        )
        btp_icon_bmr_owner=$(jq -r .address <$btp_icon_bmr_owner_wallet)
    fi

    log "bmc_add_relay: "
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
        icon_sendtx_call "$btp_icon_bmc" addRelay 0 "_link=$btp_hmny_btp_address" "_addr=$btp_icon_bmr_owner")

    log "bts: Register ONE: " # nativecoin
    _=$(WALLET=$btp_icon_bts_owner_wallet \
        PASSWORD=$btp_icon_bts_owner_wallet_password \
        icon_sendtx_call "$btp_icon_bts" register 0 \
        "_name=ONE" "_symbol=ONE" "_decimals=0x12" \
        "_feeNumerator=$(dec2hex $btp_bts_fee_numerator)" \
        "_fixedFee=$(dec2hex $btp_bts_fixed_fee)")

    btp_icon_one=$(icon_callsc \
        "$btp_icon_bts" coinAddress \
        "_coinName=ONE" | jq -r)
    log "btp_icon_one: $btp_icon_one"

    log "bts: Register IRC2 (TICX):" # pre-existing token
    _=$(WALLET=$btp_icon_wallet PASSWORD=$btp_icon_wallet_password \
        icon_sendtx_call "$btp_icon_bts" register 0 \
        "_name=TICX" "_symbol=TICX" "_decimals=0x12" \
        "_feeNumerator=$(dec2hex $btp_bts_fee_numerator)" \
        "_fixedFee=$(dec2hex $btp_bts_fixed_fee)" \
        "_addr=$btp_icon_ticx")
    log "bts: registered: $(icon_callsc "$btp_icon_bts" tokenNames | jq -r .)"

    log "bts: Register IRC2 (TONE):" # hmny's erc20 token
    _=$(WALLET=$btp_icon_wallet PASSWORD=$btp_icon_wallet_password \
        icon_sendtx_call "$btp_icon_bts" register 0 \
        "_name=TONE" "_symbol=TONE" "_decimals=0x12" \
        "_feeNumerator=$(dec2hex $btp_bts_fee_numerator)" \
        "_fixedFee=$(dec2hex $btp_bts_fixed_fee)")
    log "bts: registered: $(icon_callsc "$btp_icon_bts" tokenNames | jq -r .)"

    btp_icon_tone=$(icon_callsc \
        "$btp_icon_bts" coinAddress \
        "_coinName=TONE" | jq -r)
    log "btp_icon_tone: $btp_icon_tone"
    # icon: end

    # hmny: begin
    log "hmny"

    log "BMC: Add Link to ICON BMC: "
    WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
        run_sol >/dev/null \
        bmc.BMCManagement.addLink \
        "'$btp_icon_btp_address'"
    WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
        run_sol >/dev/null \
        bmc.BMCManagement.setLinkRxHeight \
        "'$btp_icon_btp_address',$btp_icon_block_height"
    # TODO check: response should have one raw logs ?

    # add relay
    if [ -z "$btp_hmny_bmr_owner" ]; then
        btp_hmny_bmr_owner_wallet="$ixh_tmp_dir/bmr.hmny.json"
        btp_hmny_bmr_owner_wallet_password="1234"

        log "create_wallet: [$(rel_path "$btp_hmny_bmr_owner_wallet")] "
        btp_hmny_bmr_owner=$(
            WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
                hmny_create_wallet "$btp_hmny_bmr_owner_wallet" \
                "$btp_hmny_bmr_owner_wallet_password" $btp_hmny_bmr_owner_balance
        )
        btp_hmny_bmr_owner="0x$(jq -r .address <$btp_hmny_bmr_owner_wallet)"
    fi

    log "bmc_add_relay: "
    WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
        run_sol >/dev/null \
        bmc.BMCManagement.addRelay \
        "'$btp_icon_btp_address',['$btp_hmny_bmr_owner']"

    log "bts: Register ICX: " # nativecoin
    WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
        run_sol >/dev/null \
        bts.BTSCore.register \
        "'ICX','ICX',18,'$btp_bts_fee_numerator','$btp_bts_fixed_fee','0x0000000000000000000000000000000000000000'"

    btp_hmny_icx=$(
        WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
            run_sol 2>/dev/null \
            bts.BTSCore.coinId "'ICX'" | jq -r .
    )

    log "bts: Register ERC20 (TONE):" # pre-existing token
    WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
        run_sol >/dev/null \
        bts.BTSCore.register \
        "'TONE','TONE',18,'$btp_bts_fee_numerator','$btp_bts_fixed_fee','$btp_hmny_tone'"

    log "bts: Register ERC20 (TICX):" # icon's IRC2 token
    WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
        run_sol >/dev/null \
        bts.BTSCore.register \
        "'TICX','TICX',18,'$btp_bts_fee_numerator','$btp_bts_fixed_fee','0x0000000000000000000000000000000000000000'"

    btp_hmny_ticx=$(
        WALLET=$btp_hmny_wallet PASSWORD=$btp_hmny_wallet_password \
            run_sol 2>/dev/null \
            bts.BTSCore.coinId "'TICX'" | jq -r .
    )

    # hmny: end

    # generate btp config
    generate_relay_config >$ixh_tmp_dir/bmr.config.json
}

function generate_relay_config() {
    log_stack
    local btp_icon_link_status_rx_height=$btp_hmny_block_height
    local btp_hmny_link_status_rx_height=$btp_icon_block_height

    jq <<<{} '
    .base_dir = $base_dir |
    .log_level = "debug" |
    .console_level = "trace" |
    .log_writer.filename = $log_writer_filename |
    .relays = [ $h2i_relay, $i2h_relay ]' \
        --arg base_dir "bmr" \
        --arg log_writer_filename "bmr/bmr.log" \
        --argjson h2i_relay "$(
            jq <<<{} '
            .name = "h2i" |
            .src.address = $src_address |
            .src.endpoint = [ $src_endpoint ] |
            .src.options = $src_options |
            .src.offset = $src_offset |
            .dst.address = $dst_address |
            .dst.endpoint = [ $dst_endpoint ] |
            .dst.options = $dst_options |
            .dst.key_store = $dst_key_store |
            .dst.key_store.coinType = $dst_key_store_cointype |
            .dst.key_password = $dst_key_password ' \
                --arg src_address "$btp_hmny_btp_address" \
                --arg src_endpoint "$btp_hmny_uri" \
                --argjson src_offset "$btp_icon_link_status_rx_height" \
                --argjson src_options "$(
                    jq <<<{} '
                    .syncConcurrency = 100 |
                    .verifier.blockHeight = $verifier_block_height |
                    .verifier.commitBitmap = $verifier_commit_bitmap |
                    .verifier.commitSignature = $verifier_commit_signature ' \
                        --argjson verifier_block_height "$btp_hmny_block_height" \
                        --arg verifier_commit_bitmap "$btp_hmny_verifier_commit_bitmap" \
                        --arg verifier_commit_signature "$btp_hmny_verifier_commit_signature"
                )" \
                --arg dst_address "$btp_icon_btp_address" \
                --arg dst_endpoint "$btp_icon_uri" \
                --argfile dst_key_store "$btp_icon_bmr_owner_wallet" \
                --arg dst_key_store_cointype "icx" \
                --arg dst_key_password "$btp_icon_bmr_owner_wallet_password" \
                --argjson dst_options '{"step_limit":2500000000,"tx_data_size_limit":8192}'
        )" \
        --argjson i2h_relay "$(
            jq <<<{} '
            .name = "i2h" |
            .src.address = $src_address |
            .src.endpoint = [ $src_endpoint ] |
            .src.options = $src_options |
            .src.offset = $src_offset |
            .dst.address = $dst_address |
            .dst.endpoint = [ $dst_endpoint ] |
            .dst.options = $dst_options |
            .dst.key_store = $dst_key_store |
            .dst.key_store.coinType = $dst_key_store_cointype |
            .dst.key_password = $dst_key_password ' \
                --arg src_address "$btp_icon_btp_address" \
                --arg src_endpoint "$btp_icon_uri" \
                --argjson src_offset "$btp_hmny_link_status_rx_height" \
                --argjson src_options "$(
                    jq <<<{} '
                    .verifier.blockHeight = $verifier_block_height |
                    .verifier.validatorsHash = $verifier_validators_hash ' \
                        --argjson verifier_block_height "$btp_icon_block_height" \
                        --arg verifier_validators_hash "$btp_icon_validators_hash"
                )" \
                --arg dst_address "$btp_hmny_btp_address" \
                --arg dst_endpoint "$btp_hmny_uri" \
                --argfile dst_key_store "$btp_hmny_bmr_owner_wallet" \
                --arg dst_key_store_cointype "evm" \
                --arg dst_key_password "$btp_hmny_bmr_owner_wallet_password" \
                --argjson dst_options '{"gas_limit":80000000,"boost_gas_price":1.5,"tx_data_size_limit":8192}'
        )"
}

function parallel_cmds() {
    for cmd in "$@"; do {
        $cmd &
        pid=$!
        PID_LIST+=" $pid"
    }; done
    trap "kill $PID_LIST" SIGINT
    wait $PID_LIST
}

# exposed commands
function docker_compose() {
    log_stack
    if [ "$docker_host" != "localhost" ]; then
        export DOCKER_HOST="ssh://$docker_user@$docker_host"
    fi

    local env_file=$(mktemp /tmp/ixh.env.XXXXX)
    echo "docker_registry=$docker_registry" >$env_file
    if [ -f $ixh_dir/.env ]; then
        cat $ixh_dir/.env >>$env_file
        echo >>$env_file
    fi

    local func=$1
    local args=("${@:2}")
    case "$func" in
    bmr)
        if [ -f $ixh_tmp_dir/bmr.config.json ]; then
            echo "bmr_config_json='$(
                cat $ixh_tmp_dir/bmr.config.json
            )'" >>$env_file
        fi
        docker-compose \
            -f $ixh_src_dir/docker-compose.bmr.yml \
            --env-file $env_file "${args[@]}"
        ;;
    nodes)
        docker-compose \
            -f $ixh_src_dir/docker-compose.nodes.yml \
            --env-file $env_file "${args[@]}"
        ;;
    *)
        echo "docker_compose [bmr|nodes]"
        exit 1
        ;;
    esac
    rm $env_file
}

function stop() {
    log_stack
    docker_compose "${1:-}" down "${@:2}"
}

function start() {
    log_stack
    if [ "$docker_host" != "localhost" ]; then
        docker_compose "${1:-}" pull
    fi
    docker_compose "${1:-}" up "${@:2}"
}

function build_images() {
    log_stack

    image="${1:-}"

    repos_dir=$ixh_tmp_dir/repos
    mkdir -p $repos_dir

    function build_bmr() {
        log "building bmr"
        cd $root_dir
        docker \
            build \
            -f $ixh_src_dir/bmr.Dockerfile \
            -t $docker_registry/bmr:latest .
        cd $ixh_dir
    }
    function build_icon() {
        log "building icon"
        cd $repos_dir
        if [ -d goloop ]; then
            cd goloop
            git fetch
            git checkout ${btp_icon_branch:-master}
        else
            git clone --single-branch \
                --branch ${btp_icon_branch:-master} \
                https://github.com/icon-project/goloop
            cd goloop
        fi
        make gochain-icon-image
        docker <$ixh_src_dir/icon.Dockerfile \
            build \
            --build-arg CONFIG_JSON="$(cat $btp_icon_config)" \
            -t $docker_registry/icon:latest -
        cd $ixh_dir
    }
    function build_hmny() {
        log "building hmny"
        cd $repos_dir
        # git clone --single-branch \
        #     --branch ${btp_icon_branch:-main} \
        #     https://github.com/harmony-one/harmony
        # cd harmony
        docker <$ixh_src_dir/hmny.Dockerfile \
            build \
            --build-arg SHARDING_HOST="$docker_host" \
            -t $docker_registry/hmny:latest -
        cd $ixh_dir
    }

    case "$image" in
    bmr) build_bmr ;;
    icon) build_icon ;;
    hmny) build_hmny ;;
    *)
        build_hmny
        build_icon
        build_bmr
        ;;
    esac
}

function publish_images() {
    log_stack

    image="${1-}"

    function publish_bmr() {
        log "publishing bmr to $docker_registry"
        docker push $docker_registry/bmr:latest
    }
    function publish_icon() {
        log "publishing icon to $docker_registry"
        docker push $docker_registry/icon:latest
    }
    function publish_hmny() {
        log "publishing hmny to $docker_registry"
        docker push $docker_registry/hmny:latest
    }

    case "$image" in
    bmr) publish_bmr ;;
    icon) publish_icon ;;
    hmny) publish_hmny ;;
    *)
        publish_bmr
        publish_icon
        publish_hmny
        ;;
    esac

}

function run_exec() {
    log_stack
    export verbose=true
    func=$1
    args=("${@:2}")
    case "$func" in
    iconGetBalance)
        wallet_address=${args[0]}
        params=$(jq <<<{} -c '.address=$address' --arg address $wallet_address)
        balance=$(icon_jsonrpc icx_getBalance "$params" | jq -r .result)
        hex2dec $balance
        ;;
    iconGetWrappedCoins)
        icon_callsc "$btp_icon_bts" coinNames
        ;;
    iconRegisterWrappedCoin)
        coinName=${args[0]}
        icon_sendtx_call "$btp_icon_bts" register 0 "_name=$coinName"
        ;;
    iconGetWrappedCoinBalance)
        wallet_address=${args[0]}
        coinName=${args[1]}
        # icon_callsc "$btp_icon_bts" balanceOf "_owner=$wallet_address" "_coinName=$coinName" | jq -r .
        coinAddress=$(icon_callsc "$btp_icon_bts" coinAddress "_coinName=$coinName" | jq -r .)
        icon_callsc "$coinAddress" balanceOf "_owner=$wallet_address" | jq -r .
        # hex2dec $balance
        ;;
    iconTransfer)
        address=${args[0]}
        amount=${args[1]}
        icon_transfer $address $amount
        ;;
    iconTransferNativeCoin)
        value=${args[0]}
        to=${args[1]}
        icon_sendtx_call "$btp_icon_bts" transferNativeCoin $value "_to=$to"
        ;;
    iconTransferWrappedCoin)
        coinName=${args[0]}
        value=${args[1]}
        to=${args[2]}
        icon_sendtx_call "$btp_icon_bts" transfer 0 "_coinName=$coinName" "_value=$value" "_to=$to"
        ;;
    iconGetBMCStatus)
        icon_callsc "$btp_icon_bmc" getStatus "_link=$btp_hmny_btp_address"
        ;;
    iconBSHApprove)
        coinName=${args[0]}
        spender=${args[1]}
        amount=${args[2]}
        coinAddress=$(icon_callsc "$btp_icon_bts" coinAddress "_coinName=$coinName" | jq -r .)
        icon_sendtx_call "$coinAddress" approve 0 "spender=$spender" "amount=$amount"
        ;;
    iconBSHAllowance)
        coinName=${args[0]}
        owner=${args[1]}
        spender=${args[2]}
        coinAddress=$(icon_callsc "$btp_icon_bts" coinAddress "_coinName=$coinName" | jq -r .)
        icon_callsc "$coinAddress" allowance "owner=$owner" "spender=$spender"
        ;;
    hmnyGetBalance)
        wallet_address=${args[0]}
        hmny_jsonrpc hmyv2_getBalance "[\"$wallet_address\"]" | python -c 'import json;print(json.loads(input())["result"])'
        ;;
    hmnyGetWrappedCoins)
        run_sol bts.BTSCore.coinNames
        ;;
    hmnyRegisterWrappedCoin)
        coinName=${args[0]}
        run_sol bts.BTSCore.register "'$coinName'"
        ;;
    hmnyGetWrappedCoinBalance)
        wallet_address=${args[0]}
        coinName=${args[1]}
        run_sol bts.BTSCore.getBalanceOf "'$wallet_address','$coinName'"
        ;;
    hmnyTransferNativeCoin)
        value=$(dec2hex ${args[0]})
        to=${args[1]}
        run_sol bts.BTSCore.transferNativeCoin "'$to',{value:'$value'}"
        ;;
    hmnyTransferWrappedCoin)
        coinName=${args[0]}
        value=$(dec2hex ${args[1]})
        to=${args[2]}
        run_sol bts.BTSCore.transfer "'$coinName','$value','$to'"
        ;;
    hmnyGetBMCStatus)
        run_sol bmc.BMCPeriphery.getStatus "'$btp_icon_btp_address'"
        ;;
    hmnyBSHIsApprovedForAll)
        wallet_address=${args[0]}
        run_sol bts.BTSCore.isApprovedForAll "'$wallet_address','$btp_hmny_bts_core'"
        ;;
    hmnyBSHSetApprovalForAll)
        approved=${args[0]:-1}
        approved=$([[ $approved == 0 ]] && echo false || echo true)
        run_sol bts.BTSCore.setApprovalForAll "'$btp_hmny_bts_core',$approved"
        ;;
    hmnyBSHApprove)
        coinName=${args[0]}
        spender=${args[1]}
        amount=${args[2]}
        coinAddress=$(run_sol bts.BTSCore.coinId "'$coinName'" | jq -r .)
        WALLET=${WALLET:-}
        PASSWORD=${PASSWORD:-}
        if [ $WALLET ]; then
            export PRIVATE_KEY=$(ethkey_get_private_key "$WALLET" "$PASSWORD")
        fi
        cd $ixh_sol_dir/bts
        _truffle exec --network hmny <(echo "
        const erc20t = artifacts.require('ERC20Tradable');
        module.exports = async function (callback) {
            try {
                const t = await erc20t.at('$coinAddress');
                let res = await t.approve('$spender','$amount');
                try {
                    console.log(JSON.stringify(res, null, 2));
                } catch(err) {
                    console.log(res);
                }
            } catch(err) {
                console.error(err);
            } finally { callback(); }
        }") | sed '1d' | sed '1d' # trim first 2 lines
        ;;
    hmnyBSHAllowance)
        coinName=${args[0]}
        owner=${args[1]}
        spender=${args[2]}
        coinAddress=$(run_sol bts.BTSCore.coinId "'$coinName'" | jq -r .)
        cd $ixh_sol_dir/bts
        _truffle exec --network hmny <(echo "
        const erc20t = artifacts.require('ERC20Tradable');
        module.exports = async function (callback) {
            try {
                const t = await erc20t.at('$coinAddress');
                let res = await t.allowance('$owner','$spender');
                try {
                    console.log(JSON.stringify(res, null, 2));
                } catch(err) {
                    console.log(res);
                }
            } catch(err) {
                console.error(err);
            } finally { callback(); }
        }") | sed '1d' | sed '1d' | jq -r . # trim first 2 lines
        ;;

    hmnyChainStatus)
        hmny_get_hmny_chain_status
        ;;
    # iconDeployWTS)
    #     # local oracle=${args[0]:-cx900e2d17c38903a340a0181523fa2f720af9a798} # sejong
    #     local oracle=${args[0]:-cx36a55b74aca43a9db9d5a8fc876c76d04daa85a2} # berlin
    #     scdir="$ixh_jsc_dir/wonderland"
    #     cd $scdir && gradle optimizedJar && cd $OLDPWD
    #     scfile="$scdir/build/libs/wts-0.0.1-optimized.jar"
    #     address=$(
    #         icon_deploysc $scfile \
    #             "_bmc=$btp_icon_bmc" "_net=$btp_hmny_net" \
    #             "_oracle=$oracle" | jq -r .scoreAddress
    #     )
    #     echo "icon wps: $address"
    #     icon_sendtx_call "$btp_icon_bmc" removeService 0 "_svc=WonderlandTokenSaleService"
    #     icon_sendtx_call "$btp_icon_bmc" addService 0 "_addr=$address" "_svc=WonderlandTokenSaleService"
    #     ;;
    *)
        log "invalid run command: $func"
        exit 1
        ;;
    esac
}

function run_test() {
    log_stack

    local test_dir="/tmp/btp_test"
    mkdir -p $test_dir
    # rm -rf $test_dir/* # cleanup

    local func=$1
    local args=("${@:2}")

    # . $ixh_env

    case "$func" in
    icon_transfer)
        address=${args[0]}
        ibal=$(run_exec iconGetBalance $address)
        WALLET=$btp_icon_god_wallet \
            PASSWORD=$btp_icon_god_wallet_password \
            icon_transfer $address 1
        nbal=$(run_exec iconGetBalance $address)
        [[ "$nbal" == "$(bc <<<"$ibal+1")" ]] && echo "success" || echo "failed"
        ;;

    hmny_transfer)
        address=${args[0]}
        ibal=$(run_exec hmnyGetBalance $address)
        WALLET=$btp_hmny_god_wallet \
            PASSWORD=$btp_hmny_god_wallet_password \
            hmny_transfer $address 1
        nbal=$(run_exec hmnyGetBalance $address)
        [[ "$nbal" == "$(bc <<<"$ibal+1")" ]] && echo "success" || echo "failed"
        ;;

    icon_create_wallet)
        wallet="$test_dir/icon.wallet.json"
        address=$(
            WALLET=$btp_icon_god_wallet \
                PASSWORD=$btp_icon_god_wallet_password \
                icon_create_wallet $wallet "1234" 0
        )
        require_address "$address" "" "success: $address ($wallet)" "failed!"
        ;;

    hmny_create_wallet)
        wallet="$test_dir/hmny.wallet.json"
        address=$(
            WALLET=$btp_hmny_god_wallet \
                PASSWORD=$btp_hmny_god_wallet_password \
                hmny_create_wallet $wallet "1234" 0
        )
        require_address "$address" "success: $address ($wallet)" "failed!"
        ;;

    icon_deploysc)
        scfile="$ixh_src_dir/testsc/HelloWorld.jar"
        address=$(
            WALLET=$btp_icon_god_wallet \
                PASSWORD=$btp_icon_god_wallet_password \
                icon_deploysc $scfile "name=icon" | jq -r .scoreAddress
        )
        log
        require_address "$address" "" "failed"
        name=$(icon_callsc "$address" name | jq -r .)
        [[ $name == icon ]] || name=""
        require "$name" "success" "failed"
        ;;

    hmny_deploysc)
        scdir="$ixh_tmp_dir/HelloWorld"
        rm -rf $scdir &&
            cp -r "$ixh_src_dir/testsc/HelloWorld" $scdir
        address=$(WALLET=$btp_hmny_god_wallet \
            PASSWORD=$btp_hmny_god_wallet_password \
            NAME="hmny" \
            hmny_deploysc $scdir HelloWorld)
        require_address "$address" "sc: $address" "failed to deploy"
        name=$(_truffle_exec $scdir HelloWorld.name | jq -r .)
        [[ "$name" == hmny ]] || name=""
        require "$name" "success" "failed"
        ;;

    hmny_drain_wallets)
        address=${args[0]}
        hmny_wallets_dir=$ixh_dir/testnet/hmny_wallets
        for i in $(ls $hmny_wallets_dir); do
            wallet="$hmny_wallets_dir/$i"
            from=$(echo $i | cut -d. -f1)
            if [[ "$from" == "one1a57qzygzqjpu2lpdwa9r3qa72jauytqkzsh95u" ]]; then
                echo $from
                run_exec hmnyGetBalance $from 2>/dev/null
                WALLET=$wallet PASSWORD=1234 \
                    hmny_transfer $address 9999370000000000000
            fi
        done
        ;;

    hmny_create_wallets)
        count=${args[0]}
        hmny_wallets_dir=$ixh_dir/testnet/hmny_wallets
        mkdir -p $hmny_wallets_dir
        wallets=$(
            echo "
            for (( i = 0; i < $count; i++ )); do
                hmy keys add "key\$i" --passphrase-file <(cat<<<'1234') 2>&1 > /dev/null
                keystore=\$(hmy keys export-ks --passphrase-file \
                    <(cat<<<'1234') "key\$i" / 2>/dev/null | cut -d/ -f2)
                echo \$keystore | cut -d. -f1
                cat /\$keystore && echo
            done
        " | docker run -i --rm --network=host \
                $docker_registry/hmny:latest /bin/bash
        )
        local address=
        for i in $wallets; do
            [[ -z $address ]] && address=$i || {
                echo "$i" >"$hmny_wallets_dir/$address.json"
                address=
            }
        done
        ;;

    *)
        log "invalid test: $func"
        ;;
    esac

}

function run_demo() {
    log_stack
    function tx_relay_wait() {
        sleep 45
    }

    btp_icon_step_limit=10000000

    # create and fund demo wallets
    btp_icon_demo_wallet="$ixh_src_dir/icon.demo.wallet.json"
    btp_icon_demo_wallet_address="$(jq -r .address $btp_icon_demo_wallet)"
    btp_icon_demo_wallet_password="1234"
    btp_hmny_demo_wallet="$ixh_src_dir/hmny.demo.wallet.json"
    btp_hmny_demo_wallet_address="0x$(jq -r .address $btp_hmny_demo_wallet)"
    btp_hmny_demo_wallet_password="1234"

    function get_icon_ICX_balance() {
        run_exec iconGetBalance $btp_icon_demo_wallet_address
    }

    function get_icon_ONE_balance() {
        balance=$(run_exec iconGetWrappedCoinBalance $btp_icon_demo_wallet_address ONE)
        hex2dec $balance
    }

    function get_icon_TICX_balance() {
        balance=$(run_exec iconGetWrappedCoinBalance $btp_icon_demo_wallet_address TICX)
        hex2dec $balance
    }

    function get_icon_TONE_balance() {
        balance=$(run_exec iconGetWrappedCoinBalance $btp_icon_demo_wallet_address TONE)
        hex2dec $balance
    }

    function get_hmny_ONE_balance() {
        run_exec hmnyGetBalance $btp_hmny_demo_wallet_address
    }

    function get_hmny_ICX_balance() {
        balance=$(run_exec hmnyGetWrappedCoinBalance $btp_hmny_demo_wallet_address ICX | jq -r ._usableBalance)
        hex2dec "0x$balance"
    }

    function get_hmny_TONE_balance() {
        balance=$(run_exec hmnyGetWrappedCoinBalance $btp_hmny_demo_wallet_address TONE | jq -r ._usableBalance)
        hex2dec "0x$balance"
    }

    function get_hmny_TICX_balance() {
        balance=$(run_exec hmnyGetWrappedCoinBalance $btp_hmny_demo_wallet_address TICX | jq -r ._usableBalance)
        hex2dec "0x$balance"
    }

    function fund_demo_wallets() {
        echo
        echo "Funding demo wallets..."

        local icx_target=10000000000000000000
        local irc2_target=10000000000000000000
        local one_target=10000000000000000000
        local erc20_target=10000000000000000000

        local bal=0

        echo -n "    ICON ($btp_icon_demo_wallet_address): "

        bal=$(get_icon_ICX_balance)
        bal=$(echo "scale=18;$icx_target-$bal" | bc)
        if (($(echo "$bal > 0" | bc -l))); then
            WALLET=$btp_icon_wallet \
                PASSWORD=$btp_icon_wallet_password \
                icon_transfer $btp_icon_demo_wallet_address "$bal" # make 250 ICX
        else
            bal=0
        fi
        echo -n "+$(echo "scale=2;$bal/10^18" | bc) ICX"

        bal=$(get_icon_TICX_balance)
        bal=$(echo "scale=18;$irc2_target-$bal" | bc)
        if (($(echo "$bal > 0" | bc -l))); then
            WALLET=$btp_icon_wallet \
                PASSWORD=$btp_icon_wallet_password \
                icon_sendtx_call >/dev/null \
                "$btp_icon_ticx" transfer 0 \
                "_to=$btp_icon_demo_wallet_address" \
                "_value=$bal"
        else
            bal=0
        fi
        echo ", +$(echo "scale=2;$bal/10^18" | bc) TICX"

        echo -n "    HMNY ($btp_hmny_demo_wallet_address): "

        bal=$(get_hmny_ONE_balance)
        bal=$(echo "scale=18;$one_target-$bal" | bc)
        if (($(echo "$bal > 0" | bc -l))); then
            WALLET=$btp_hmny_wallet \
                PASSWORD=$btp_hmny_wallet_password \
                hmny_transfer $btp_hmny_demo_wallet_address $bal
        else
            bal=0
        fi
        echo -n "+$(echo "scale=2;$bal/10^18" | bc) ONE"

        bal=$(get_hmny_TONE_balance)
        bal=$(echo "scale=18;$erc20_target-$bal" | bc)
        if (($(echo "$bal > 0" | bc -l))); then

            WALLET=$btp_hmny_wallet \
                PASSWORD=$btp_hmny_wallet_password \
                run_sol >/dev/null \
                bts.HRC20.transfer \
                "'$btp_hmny_demo_wallet_address','$bal'"
        else
            bal=0
        fi
        echo ", +$(echo "scale=2;$bal/10^18" | bc) TONE"

        echo
    }

    function format_token() {
        echo "scale=2;$1/10^18" | bc
    }

    function show_balances() {
        echo
        echo "Balance:"
        echo "    ICON: $btp_icon_demo_wallet_address"
        local icon_balance=$(get_icon_ICX_balance)
        echo "        ICX: $(format_token $icon_balance)"
        local icon_TICX=$(get_icon_TICX_balance)
        echo "        TICX: $(format_token $icon_TICX)"
        local icon_ONE=$(get_icon_ONE_balance)
        echo "        ONE: $(format_token $icon_ONE)"
        local icon_TONE=$(get_icon_TONE_balance)
        echo "        TONE: $(format_token $icon_TONE)"
        echo "    HMNY: $btp_hmny_demo_wallet_address"
        local hmny_balance=$(get_hmny_ONE_balance)
        echo "        ONE: $(format_token $hmny_balance)"
        local hmny_TONE=$(get_hmny_TONE_balance)
        echo "        TONE: $(format_token $hmny_TONE)"
        local hmny_ICX=$(get_hmny_ICX_balance)
        echo "        ICX: $(format_token $hmny_ICX)"
        local hmny_TICX=$(get_hmny_TICX_balance)
        echo "        TICX: $(format_token $hmny_TICX)"
        echo
    }

    function show_token_names() {
        echo "ICON:"
        echo "    $(run_exec iconGetWrappedCoins | jq -c .)"
        echo "HMNY:"
        echo "    $(run_exec hmnyGetWrappedCoins | jq -c .)"
    }

    show_token_names
    fund_demo_wallets
    show_balances

    i2h_ICX_transfer_amount=3000000000000000000 # 3 ICX
    echo "Transfer Native ICX (ICON -> HMNY):"
    echo "    amount=$(format_token $i2h_ICX_transfer_amount)"
    echo -n "    "
    WALLET=$btp_icon_demo_wallet \
        PASSWORD=$btp_icon_demo_wallet_password \
        run_exec iconTransferNativeCoin \
        $i2h_ICX_transfer_amount \
        "btp://$btp_hmny_net/$btp_hmny_demo_wallet_address" >/dev/null
    echo

    tx_relay_wait
    show_balances

    h2i_ONE_transfer_amount=3000000000000000000 # 3 ONE
    echo "Transfer Native ONE (HMNY -> ICON):"
    echo "    amount=$(format_token $h2i_ONE_transfer_amount)"
    echo -n "    "
    WALLET=$btp_hmny_demo_wallet \
        PASSWORD=$btp_hmny_demo_wallet_password \
        run_exec hmnyTransferNativeCoin \
        $h2i_ONE_transfer_amount \
        "btp://$btp_icon_net/$btp_icon_demo_wallet_address" >/dev/null
    echo

    tx_relay_wait
    show_balances

    h2i_ICX_transfer_amount=1000000000000000000 # 1 ICX
    echo "Transfer ICX (HMNY -> ICON):"
    echo "    amount=$(format_token $h2i_ICX_transfer_amount)"
    echo -n "    "
    WALLET=$btp_hmny_demo_wallet \
        PASSWORD=$btp_hmny_demo_wallet_password \
        run_exec hmnyBSHApprove "ICX" \
        "$btp_hmny_bts_core" "$h2i_ICX_transfer_amount" >/dev/null
    WALLET=$btp_hmny_demo_wallet \
        PASSWORD=$btp_hmny_demo_wallet_password \
        run_exec hmnyTransferWrappedCoin \
        ICX \
        $h2i_ICX_transfer_amount \
        "btp://$btp_icon_net/$btp_icon_demo_wallet_address" >/dev/null

    tx_relay_wait
    show_balances

    i2h_ONE_transfer_amount=1000000000000000000 # 1 ONE
    echo "Transfer ONE (ICON -> HMNY):"
    echo "    amount=$(format_token $i2h_ONE_transfer_amount)"
    echo -n "    "
    WALLET=$btp_icon_demo_wallet \
        PASSWORD=$btp_icon_demo_wallet_password \
        run_exec iconBSHApprove "ONE" \
        "$btp_icon_bts" "$i2h_ONE_transfer_amount" >/dev/null
    WALLET=$btp_icon_demo_wallet PASSWORD=$btp_icon_demo_wallet_password \
        run_exec iconTransferWrappedCoin \
        ONE \
        $i2h_ONE_transfer_amount \
        "btp://$btp_hmny_net/$btp_hmny_demo_wallet_address" >/dev/null
    echo

    tx_relay_wait
    show_balances

    i2h_TICX_transfer_amount=3000000000000000000 # 3 TICX
    echo "Transfer TICX (ICON -> HMNY):"
    echo "    amount=$(format_token $i2h_TICX_transfer_amount)"
    echo -n "    "
    WALLET=$btp_icon_demo_wallet \
        PASSWORD=$btp_icon_demo_wallet_password \
        icon_sendtx_call >/dev/null \
        "$btp_icon_ticx" transfer 0 \
        "_to=$btp_icon_bts" \
        "_value=$i2h_TICX_transfer_amount"
    WALLET=$btp_icon_demo_wallet \
        PASSWORD=$btp_icon_demo_wallet_password \
        icon_sendtx_call >/dev/null \
        "$btp_icon_bts" transfer 0 \
        "_coinName=TICX" \
        "_value=$i2h_TICX_transfer_amount" \
        "_to=btp://$btp_hmny_net/$btp_hmny_demo_wallet_address"
    echo

    tx_relay_wait
    show_balances

    h2i_TONE_transfer_amount=3000000000000000000 # 3 TONE
    echo "Transfer TONE (HMNY -> ICON):"
    echo "    amount=$(format_token $h2i_TONE_transfer_amount)"
    echo -n "    "
    WALLET=$btp_hmny_demo_wallet \
        PASSWORD=$btp_hmny_demo_wallet_password \
        run_exec hmnyBSHApprove "TONE" \
        "$btp_hmny_bts_core" "$h2i_TONE_transfer_amount" >/dev/null
    WALLET=$btp_hmny_demo_wallet \
        PASSWORD=$btp_hmny_demo_wallet_password \
        run_sol >/dev/null \
        bts.BTSCore.transfer \
        "'TONE','$h2i_TONE_transfer_amount','btp://$btp_icon_net/$btp_icon_demo_wallet_address'"

    tx_relay_wait
    show_balances

    h2i_TICX_transfer_amount=1000000000000000000 # 1 TICX
    echo "Transfer TICX (HMNY -> ICON):"
    echo "    amount=$(format_token $h2i_TICX_transfer_amount)"
    echo -n "    "
    WALLET=$btp_hmny_demo_wallet \
        PASSWORD=$btp_hmny_demo_wallet_password \
        run_exec hmnyBSHApprove "TICX" \
        "$btp_hmny_bts_core" "$h2i_TICX_transfer_amount" >/dev/null
    WALLET=$btp_hmny_demo_wallet \
        PASSWORD=$btp_hmny_demo_wallet_password \
        run_sol >/dev/null \
        bts.BTSCore.transfer \
        "'TICX','$h2i_TICX_transfer_amount','btp://$btp_icon_net/$btp_icon_demo_wallet_address'"

    tx_relay_wait
    show_balances

    i2h_TONE_transfer_amount=1000000000000000000 # 1 TONE
    echo "Transfer TONE (ICON -> HMNY):"
    echo "    amount=$(format_token $i2h_TONE_transfer_amount)"
    echo -n "    "
    WALLET=$btp_icon_demo_wallet \
        PASSWORD=$btp_icon_demo_wallet_password \
        run_exec iconBSHApprove "TONE" \
        "$btp_icon_bts" "$i2h_TONE_transfer_amount" >/dev/null
    WALLET=$btp_icon_demo_wallet \
        PASSWORD=$btp_icon_demo_wallet_password \
        icon_sendtx_call >/dev/null \
        "$btp_icon_bts" transfer 0 \
        "_coinName=TICX" \
        "_value=$i2h_TONE_transfer_amount" \
        "_to=btp://$btp_hmny_net/$btp_hmny_demo_wallet_address"
    echo

    tx_relay_wait
    show_balances
}

function usage() {
    echo "Usage: $(basename $0) [build|publish|deploysc|start|stop]"
    exit 1
}

if [ $# -gt 0 ]; then
    cmd=$1
else
    usage
fi
args=("${@:2}")

################################## init: begin #######################################

docker_user="ubuntu"
docker_host="localnets"
docker_port="5000"
docker_registry="$docker_host"
[[ -z $docker_port ]] || docker_registry+=":$docker_port"

ixh_dir=$PWD
ixh_tmp_dir=$ixh_dir/_ixh
ixh_build_dir=$ixh_tmp_dir/build
ixh_tests_dir=$ixh_tmp_dir/tests
ixh_env=$ixh_tmp_dir/ixh.env
ixh_src_dir=$ixh_dir/src
ixh_sol_dir=$ixh_build_dir/solidity
ixh_jsc_dir=$ixh_build_dir/javascore

root_dir="$ixh_dir/../../.."

btp_icon_branch="v1.2.3"
btp_hmny_branch="v4.3.7"

# hmny dummy wallet: used for smart contract calls (zero balance)
btp_hmny_dummy_private_key=a49152cea2bd63cc8dddebc7f7699b9f0b2bc770af67554f1c54894b683b9f4a

# common configuration
btp_icon_bmc_owner_balance=50000000000000000000 # 50 ICX
btp_icon_bts_owner_balance=50000000000000000000 # 50 ICX
btp_icon_step_limit=3500000000
btp_bts_svc_name=bts
btp_bts_fee_numerator=100
btp_bts_fixed_fee=200000000000000000
btp_icon_bmr_owner_balance=50000000000000000000 # 50 ICX
btp_hmny_bmr_owner_balance=10000000000000000000 # 10 ONE
btp_hmny_gas_limit=80000000                     # equal to block gas limit
btp_hmny_gas_price=30000000000                  # 30 gwei

# localnet: begin
btp_icon_config=$ixh_src_dir/icon.config.json

btp_icon_god_wallet=$ixh_src_dir/icon.god.wallet.json # at least 100 ICX
btp_icon_god_wallet_address=$(jq <$btp_icon_god_wallet -r .address 2>/dev/null)
btp_icon_god_wallet_password=gochain

btp_hmny_god_wallet=$ixh_src_dir/hmny.god.wallet.json # at least 100 ONE
btp_hmny_god_wallet_password=
btp_hmny_god_wallet_private_key=1f84c95ac16e6a50f08d44c7bde7aff8742212fda6e4321fde48bf83bef266dc
btp_hmny_god_wallet_address=0xA5241513DA9F4463F1d4874b548dFBAC29D91f34
btp_hmny_god_wallet_address_bech32=one155jp2y76nazx8uw5sa94fr0m4s5aj8e5xm6fu3

btp_icon_nid=$(dec2hex $(cat "$btp_icon_config" | jq -r .nid 2>/dev/null))
btp_icon_uri="http://$docker_host:9080/api/v3/default"
btp_hmny_nid=0x6357d2e0
btp_hmny_uri="http://$docker_host:9500"
btp_hmny_chain_id=2
btp_icon_fee_aggregator='hx62f0e50312629bbb4201200bbd201f840780b025'
# localnet: end

# # testnet: begin
# btp_icon_god_wallet=$ixh_dir/testnet/icon.god.wallet.json # at least 100 ICX
# btp_icon_god_wallet_address=hxff0ea998b84ab9955157ab27915a9dc1805edd35
# btp_icon_god_wallet_password=gochain

# btp_hmny_god_wallet=$ixh_dir/testnet/hmny.god.wallet.json # at least 100 ONE
# btp_hmny_god_wallet_private_key=0xd104bd9d3acaff111d52dad5bedac0eaeba059af5c2c5fa6c4bc5e7e53cfe424
# btp_hmny_god_wallet_address=0xedce30ac360b30134d9cc0880d621d97d3a4c517
# btp_hmny_god_wallet_address_bech32=one1ah8rptpkpvcpxnvuczyq6csajlf6f3ghs8ekym
# btp_hmny_god_wallet_password=

# btp_icon_nid="0x7" # 0x53 (sejong)
# btp_icon_uri="https://berlin.net.solidwallet.io/api/v3/icon_dex"
# btp_hmny_nid=0x6357d2e0
# btp_hmny_uri="https://rpc.s0.b.hmny.io"
# btp_hmny_chain_id=2
# btp_icon_fee_aggregator='hx62f0e50312629bbb4201200bbd201f840780b025'
# # testnet: end

# wallets for deploysc/tests
btp_icon_wallet=${btp_icon_wallet:-$btp_icon_god_wallet}
btp_icon_wallet_address=${btp_icon_wallet_address:-$btp_icon_god_wallet_address}
btp_icon_wallet_password=${btp_icon_wallet_password:-$btp_icon_god_wallet_password}
btp_hmny_wallet=${btp_hmny_wallet:-$btp_hmny_god_wallet}
btp_hmny_wallet_address=${btp_hmny_wallet_address:-$btp_hmny_god_wallet_address}
btp_hmny_wallet_password=${btp_hmny_wallet_password:-$btp_hmny_god_wallet_password}
btp_icon_wallet_minimum_balance=100000000000000000000 #100 ICX
btp_hmny_wallet_minimum_balance=100000000000000000000 #100 ONE

# icon/hmny network ids
btp_icon_nid=${btp_icon_nid:-}
btp_icon_net="$btp_icon_nid.icon"

btp_hmny_nid=${btp_hmny_nid:-$(hmny_jsonrpc eth_chainId '[]' | jq -r .result)}
btp_hmny_net="$btp_hmny_nid.hmny"

# create tmp dir
mkdir -p $ixh_tmp_dir

# require "$btp_icon_nid" "icon_nid: $btp_icon_nid" "invalid icon_nid: $btp_icon_nid"
# require "$btp_hmny_nid" "hmny_nid: $btp_hmny_nid" "invalid hmny_nid: $btp_hmny_nid"

# overrides
# btp_icon_bmc_owner_wallet="$ixh_tmp_dir/bmc.owner.json"
# btp_icon_bmc_owner_wallet_password="1234"
# btp_icon_bts_owner_wallet="$ixh_tmp_dir/bts.icon.owner.json"
# btp_icon_bts_owner_wallet_password="1234"

btp_icon_bmc_owner_wallet="$btp_icon_wallet"
btp_icon_bmc_owner_wallet_password="$btp_icon_wallet_password"
btp_icon_bts_owner_wallet="$btp_icon_wallet"
btp_icon_bts_owner_wallet_password="$btp_icon_wallet_password"
# btp_icon_bmr_owner="hxfaff7dfd515d7f2b43270d5977b7587a65a48972"
# btp_hmny_bmr_owner="0x4617eae515f629867ca6b486662e3ee65888937c"
# btp_icon_bmr_owner_wallet=
# btp_hmny_bmr_owner_wallet=

################################## init: end #######################################

case "$cmd" in
start)
    start "${args[@]}"
    ;;

stop)
    stop "${args[@]}"
    ;;

docker_compose)
    docker_compose "${args[@]}"
    ;;

build)
    build_images "${args[@]}"
    ;;

publish)
    publish_images "${args[@]}"
    ;;

deploysc)
    if [ "${args[0]}" == "reset" ]; then
        # clean build
        rm -rf $ixh_build_dir/*
        echo >$ixh_env
    elif [ -f $ixh_env ]; then
        . $ixh_env
    fi
    deploysc "${args[@]:2}"
    ;;

generate_relay_config)
    . $ixh_env
    generate_relay_config >$ixh_tmp_dir/bmr.config.json
    ;;

aws_bmr_secrets)
    cat $ixh_tmp_dir/bmr.config.json |
        jq ".relays[] | { key_store: .dst.key_store, secret: .dst.key_password }"
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
    run_sol "${args[@]}"
    ;;

jscall)
    . $ixh_env
    icon_callsc "${args[@]}"
    ;;

jstxcall)
    . $ixh_env
    icon_sendtx_call "${args[@]}"
    ;;

exec)
    . $ixh_env
    run_exec "${args[@]}"
    ;;

test)
    run_test "${args[@]}"
    ;;

demo)
    . $ixh_env
    run_demo "${args[@]}"
    ;;

set_fee_aggregator)
    . $ixh_env

    echo "fee_aggregator: $(icon_callsc "$btp_icon_bmc" getFeeAggregator)"
    # configure
    log "bmc_set_fee_aggregator:"
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
        icon_sendtx_call \
        "$btp_icon_bmc" setFeeAggregator 0 "_addr=$btp_icon_fee_aggregator")
    echo "fee_aggregator: $(icon_callsc "$btp_icon_bmc" getFeeAggregator)"

    echo "fee_gathering_term: $(icon_callsc "$btp_icon_bmc" getFeeGatheringTerm)"
    log "bmc_set_fee_gathering_term:"
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
        icon_sendtx_call \
        "$btp_icon_bmc" setFeeGatheringTerm 0 "_value=1000") # every 1000 blocks
    echo "fee_gathering_term: $(icon_callsc "$btp_icon_bmc" getFeeGatheringTerm)"
    ;;

send_fee_gathering)
    . $ixh_env

    log "bmc_send_fee_gathering:"
    _=$(WALLET=$btp_icon_bmc_owner_wallet \
        PASSWORD=$btp_icon_bmc_owner_wallet_password \
        icon_sendtx_call \
        "$btp_icon_bmc" sendFeeGathering 0)

    ;;

*)
    log "Invalid command: $cmd"
    usage
    ;;
esac
