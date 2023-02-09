#!/bin/env bash

export PATH=$PATH:~/go/bin

DEPLOY_TXN_ID=$1
end_time=$(($(date +%s) + 30))

while [[ $(date +%s) -lt $end_time ]]; do
  ICON_BMC_ADDR=$(goloop rpc txresult $(echo $DEPLOY_TXN_ID | cut -d '"' -f 2) --uri http://localhost:9080/api/v3 | grep '"scoreAddress"' | cut -d '"' -f 4)
  if [ -n "$ICON_BMC_ADDR" ]; then
    break
  fi
  sleep 1
done

if [ -z "$ICON_BMC_ADDR" ]; then
  echo "The ICON_BMC_ADDR environment variable is empty, there was an error deploying the BMC."
  exit 1
fi

echo $ICON_BMC_ADDR