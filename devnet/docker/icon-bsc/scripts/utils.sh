#!/bin/sh

create_ensure_bob_account() {
  cd ${CONFIG_DIR}
  if [ ! -f bob.btp.address ]; then
    eth address:random >bob.ks.json
    echo "btp://$BSC_NID/$(get_bob_address)" >bob.btp.address
  fi
}

get_alice_address() {
  cat $CONFIG_DIR/alice.ks.json | jq -r .address
}

get_bob_address() {
  #cat $CONFIG_DIR/bsc.ks.json | jq -r .address
  echo 0x$(cat $CONFIG_DIR/bob.ks.json | jq -r .address)
}

hex2int() {
  input=$1
  input=$(echo $input | sed 's/^0x//g')
  input=$(uppercase $input)
  echo "ibase=16; $input" | bc
}

decimal2Hex() {
  printf '0x%x\n' $1
}

PRECISION=18
COIN_UNIT=$((10 ** $PRECISION))

coin2wei() {
  amount=$1
  printf 'scale=0; %s * %s / 1\n' $COIN_UNIT $amount | bc -l
}

wei2coin() {
  amount=$1
  printf 'scale=%s; %s / %s\n' $PRECISION $amount $COIN_UNIT | bc
}

uppercase() {
  input=$1
  printf '%s\n' "$input" | awk '{ print toupper($0) }'
}

create_contracts_address_json() {
  TYPE="${1-solidity}"
  NAME="$2"
  VALUE="$3"
  if test -f "$CONFIG_DIR/addresses.json"; then
    echo "appending address.json"
    objJSON="{\"$NAME\":\"$VALUE\"}"
    cat $CONFIG_DIR/addresses.json | jq --arg type "$TYPE" --argjson jsonString "$objJSON" '.[$type] += $jsonString' >$CONFIG_DIR/addresses.json
  else
    echo "creating address.json"
    objJSON="{\"$TYPE\":{\"$NAME\":\"$VALUE\"}}"
    jq -n --argjson jsonString "$objJSON" '$jsonString' >$CONFIG_DIR/addresses.json
    wait_for_file $CONFIG_DIR/addresses.json
  fi
}

extractAddresses() {
  if [ $# -lt 2 ]; then
    echo "Usage: extractAddresses [TYPE="javascore/solidity"] [NAME="bmc/TokenBSH/IRC2/NativeBSH"]"
    exit 1
  fi
  TYPE=$1
  NAME=$2
  echo $(cat $CONFIG_DIR/addresses.json | jq -r .$TYPE.$NAME)
}

generate_addresses_json() {
  jq -n '
    .javascore.bmc = $bmc |
    .javascore.TokenBSH = $token_bsh |
    .javascore.IRC2 = $irc2_token |
    .javascore.NativeBSH = $native_bsh |
    .javascore.BNB = $bnb_token |
    .solidity.BMCPeriphery = $bmc_periphery |
    .solidity.BMCManagement = $bmc_management |
    .solidity.BSHImpl = $bsh_impl |
    .solidity.BSHProxy = $bsh_proxy |
    .solidity.BEP20TKN = $bep20_token |
    .solidity.BSHPeriphery = $bsh_periphery |
    .solidity.BSHCore = $bsh_core' \
    --arg bmc "$(cat $CONFIG_DIR/bmc.icon)" \
    --arg token_bsh "$(cat $CONFIG_DIR/token_bsh.icon)" \
    --arg irc2_token "$(cat $CONFIG_DIR/irc2_token.icon)" \
    --arg native_bsh "$(cat $CONFIG_DIR/nativebsh.icon)" \
    --arg bnb_token "$(cat $CONFIG_DIR/irc2TradeableToken.icon)" \
    --arg bmc_periphery "$(cat $CONFIG_DIR/bmc.periphery.bsc)" \
    --arg bmc_management "$(cat $CONFIG_DIR/bmc.bsc)" \
    --arg bsh_impl "$(cat $CONFIG_DIR/token_bsh.impl.bsc)" \
    --arg bsh_proxy "$(cat $CONFIG_DIR/token_bsh.proxy.bsc)" \
    --arg bep20_token "$(cat $CONFIG_DIR/bep20_token.bsc)" \
    --arg bsh_periphery "$(cat $CONFIG_DIR/bsh.periphery.bsc)" \
    --arg bsh_core "$(cat $CONFIG_DIR/bsh.core.bsc)"
}

create_abi() {
  NAME=$1
  echo "Generating abi file from ./build/contracts/$NAME.json"
  [ ! -d $CONFIG_DIR/abi ] && mkdir -p $CONFIG_DIR/abi
  cat "./build/contracts/$NAME.json" | jq -r .abi >$CONFIG_DIR/abi/$NAME.json
  wait_for_file $CONFIG_DIR/abi/$NAME.json
}
