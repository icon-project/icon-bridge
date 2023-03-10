echo "Sending msg from algo to icon"

A2I_BEF_TEST=$(goloop rpc call --to $(echo $(cat cache/icon_dbsh_addr) | cut -d '"' -f 2) \
                --uri http://localhost:9080/api/v3/icon \
                --method getLastReceivedMessage | xxd -r -p)

echo $A2I_BEF_TEST
ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) PRIVATE_KEY=$(cat cache/algo_private_key) dbsh-call-send-service-message ../../../pyteal/teal $(cat cache/bmc_app_id) $(cat cache/dbsh_app_id)
sleep 25

A2I_AFT_TEST=$(goloop rpc call --to $(echo $(cat cache/icon_dbsh_addr) | cut -d '"' -f 2) \
    --uri http://localhost:9080/api/v3/icon \
    --method getLastReceivedMessage | xxd -r -p)

echo $A2I_AFT_TEST

if [ "$A2I_BEF_TEST" = "$A2I_AFT_TEST" ]
then
    echo "Dummy BSH didn't receive the message from Algorand"
    exit 1
fi

echo "Sending msg from icon to algo"

I2A_BEF_TEST=$(ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-global-state-by-key $(cat cache/dbsh_app_id) last_received_message)
echo $I2A_BEF_TEST

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
sleep 60

I2A_AFT_TEST=$(ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-global-state-by-key $(cat cache/dbsh_app_id) last_received_message)
echo $I2A_AFT_TEST

if [ "$I2A_BEF_TEST" = "$I2A_AFT_TEST" ]
then
    echo "Dummy BSH didn't receive the message from Icon"
    exit 1
fi



