#!/bin/bash

MSG_BEF_TEST=$(goloop rpc call --to $(echo $(cat cache/icon_dbsh_addr) | cut -d '"' -f 2) \
                --uri http://localhost:9080/api/v3/icon \
                --method getLastReceivedMessage | xxd -r -p)

ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) PRIVATE_KEY=$(cat cache/algo_private_key) dbsh-call-send-service-message ../../../pyteal/teal $(cat cache/bmc_app_id) $(cat cache/dbsh_app_id)
sleep 10

MSG_AFT_TEST=$(goloop rpc call --to $(echo $(cat cache/icon_dbsh_addr) | cut -d '"' -f 2) \
    --uri http://localhost:9080/api/v3/icon \
    --method getLastReceivedMessage | xxd -r -p)

if [ $MSG_BEF_TEST -eq $MSG_AFT_TEST ]
then
      echo "Dummy BSH didn't receive the message from Algorand"
      exit 1
fi

# TODO: check received message on Algorand after get I2A relayer running.

TXN_ID=$(
    goloop rpc sendtx call --to $(cat cache/icon_dbsh_addr) \
        --method sendServiceMessage \
        --value 0 \
        --step_limit=3000000000 \
        --uri http://localhost:9080/api/v3/icon \
        --key_store icon.keystore.json --key_password gochain \
        --nid=$(cat cache/nid)
)

./../../algorand/scripts/wait_for_transaction.sh $TXN_ID