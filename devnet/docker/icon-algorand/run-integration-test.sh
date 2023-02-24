#!/bin/bash

MSG_BEF_TEST=$(goloop rpc call --to $(echo $(cat cache/icon_dbsh_addr) | cut -d '"' -f 2) \
                --uri http://localhost:9080/api/v3 \
                --method getLastReceivedMessage | xxd -r -p)

ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) PRIVATE_KEY=$(cat cache/algo_private_key) dbsh-call-send-service-message ../../../pyteal/teal $(cat cache/bmc_app_id) $(cat cache/dbsh_app_id)
sleep 10

MSG_AFT_TEST=$(goloop rpc call --to $(echo $(cat cache/icon_dbsh_addr) | cut -d '"' -f 2) \
    --uri http://localhost:9080/api/v3 \
    --method getLastReceivedMessage | xxd -r -p)

if [ $MSG_BEF_TEST -eq $MSG_AFT_TEST ]
then
      echo "Dummy BSH didn't receive the message from Algorand"
      exit 1
fi