#!/bin/bash

export PATH=$PATH:~/go/bin

TXN_ID=$1
FIELD=$2

END_TIME=$(($(date +%s) + 30))

while [ $(date +%s) -lt $END_TIME ]; do
  TXN_RESULT=$(goloop rpc txresult $(echo $TXN_ID | cut -d '"' -f 2) \
  --uri https://lisbon.net.solidwallet.io/api/v3 | jq .$FIELD| cut -d '"' -f 2)

  if [ -n "$TXN_RESULT" ]; then
    break
  fi
  sleep 1
done

if [ -z "$TXN_RESULT" ]; then
  echo "The transaction $TXN_ID was not confirmed after 30 seconds."
  exit 1
fi

echo $TXN_RESULT