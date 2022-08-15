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

function hex2int() {
    hex=${@#0x}
    echo "obase=10; ibase=16; ${hex^^}" | bc
}

function decimal2Hex() {
    hex=$(echo "obase=16; ibase=10; ${@}" | bc)
    echo "0x$(tr [A-Z] [a-z] <<< "$hex")"
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
echo "{"
echo "    \"javascore\": {"
for v in "${ICON_NATIVE_TOKEN_SYM[@]}"
do
    echo "        " \"$v\" : \"$(cat $CONFIG_DIR/icon.addr.coin$v)\",
done
for v in "${ICON_WRAPPED_COIN_SYM[@]}"
do
    echo "        " \"$v\" : \"$(cat $CONFIG_DIR/icon.addr.coin$v)\",
done
echo "        " \"bmc\": \"$(cat $CONFIG_DIR/icon.addr.bmc)\",
echo "        " \"bts\": \"$(cat $CONFIG_DIR/icon.addr.bts)\"
echo "    },"
echo "    \"solidity\": {"
for v in "${BSC_NATIVE_TOKEN_SYM[@]}"
do
    echo "        " \"$v\" : \"$(cat $CONFIG_DIR/bsc.addr.coin$v)\",
done
for v in "${BSC_WRAPPED_COIN_SYM[@]}"
do
    echo "        " \"$v\" : \"$(cat $CONFIG_DIR/bsc.addr.coin$v)\",
done
echo "        " \"BMCManagement\": \"$(cat $CONFIG_DIR/bsc.addr.bmcmanagement)\",
echo "        " \"BMCPeriphery\": \"$(cat $CONFIG_DIR/bsc.addr.bmcperiphery)\",
echo "        " \"BTSCore\": \"$(cat $CONFIG_DIR/bsc.addr.btscore)\",
echo "        " \"BTSPeriphery\": \"$(cat $CONFIG_DIR/bsc.addr.btsperiphery)\"
echo "    }"
echo "}"
  #jq -n $str
  # jq -n '
  #   $str
  #   .javascore.bmc = $bmc |
  #   .javascore.bts = $bts |
  #   .solidity.BMCPeriphery = $bmc_periphery |
  #   .solidity.BMCManagement = $bmc_management |
  #   .solidity.BTSCore = $bts_core | 
  #   .solidity.BTSPeriphery = $bts_periphery' \
  #   --arg bmc "$(cat $CONFIG_DIR/icon.addr.bmc)" \
  #   --arg bts "$(cat $CONFIG_DIR/icon.addr.bts)" \
  #   --arg bmc_periphery "$(cat $CONFIG_DIR/bsc.addr.bmcperiphery)" \
  #   --arg bmc_management "$(cat $CONFIG_DIR/bsc.addr.bmcmanagement)" \
  #   --arg bts_periphery "$(cat $CONFIG_DIR/bsc.addr.btsperiphery)" \
  #   --arg bts_core "$(cat $CONFIG_DIR/bsc.addr.btscore)"
}

create_abi() {
  NAME=$1
  echo "Generating abi file from ./build/contracts/$NAME.json"
  [ ! -d $CONFIG_DIR/abi ] && mkdir -p $CONFIG_DIR/abi
  cat "./build/contracts/$NAME.json" | jq -r .abi >$CONFIG_DIR/abi/$NAME.json
  wait_for_file $CONFIG_DIR/abi/$NAME.json
}
