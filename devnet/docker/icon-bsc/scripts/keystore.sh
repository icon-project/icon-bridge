#!/bin/sh

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
    goloop ks gen --out ${KEY_STORE}tmp -p $(cat ${KEY_SECRET}) > /dev/null 2>&1
    cat ${KEY_STORE}tmp | jq -r . > ${KEY_STORE}
    rm ${KEY_STORE}tmp

  fi
  echo ${KEY_STORE}
}

ensure_bsc_key_store() {
  if [ $# -lt 2 ] ; then
    echo "Usage: ensure_key_store KEYSTORE_PATH SECRET_PATH"
    return 1
  fi

  local KEY_STORE_PATH=$1
  local KEY_SECRET_PATH=$(ensure_key_secret $2)
  if [ ! -f "${KEY_STORE_PATH}" ]; then
    mkdir -p $ICONBRIDGE_CONFIG_DIR/keystore
    ethkey generate --passwordfile $KEY_SECRET_PATH --json tmp
    cat tmp | jq -r . > $KEY_STORE_PATH
    ethkey inspect --json --private --passwordfile $KEY_SECRET_PATH $KEY_STORE_PATH | jq -r .PrivateKey > ${KEY_STORE_PATH}.priv
    rm tmp
    # tr -dc A-Fa-f0-9 </dev/urandom | head -c 64 > $ICONBRIDGE_CONFIG_DIR/keystore/$(basename ${KEY_STORE_PATH}).priv
    # tmpPath=$(geth account import --datadir $ICONBRIDGE_CONFIG_DIR --password $KEY_SECRET_PATH $ICONBRIDGE_CONFIG_DIR/keystore/$(basename ${KEY_STORE_PATH}).priv | sed -e "s/^Address: {//" -e "s/}//")
    # fileMatch=$(find $ICONBRIDGE_CONFIG_DIR/keystore -type f -name '*'$tmpPath)
    # cat $fileMatch | jq -r . > $KEY_STORE_PATH
    # rm $fileMatch
  fi
  echo ${KEY_STORE_PATH}
}

ensure_empty_key_secret() {
  if [ $# -lt 1 ] ; then
    echo "Usage: ensure_key_secret SECRET_PATH"
    return 1
  fi
  local KEY_SECRET=$1
  if [ ! -f "${KEY_SECRET}" ]; then
    mkdir -p $(dirname ${KEY_SECRET})
    echo -n '' > ${KEY_SECRET}
  fi
  echo ${KEY_SECRET}
}