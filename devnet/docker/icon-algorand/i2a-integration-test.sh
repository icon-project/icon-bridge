#!/bin/bash

echo "Start i2a integration test"

SENDER_INITIAL_AMOUNT=5000000
TRANSFER_AMOUNT=1000000
ICON_SENDER_ADDRESS=$(cat sender.keystore.json | jq -r '.address')

echo "Create ICON receiver account"
goloop ks gen --out receiver.keystore.json
ICON_RECEIVER_ADDRESS=$(cat receiver.keystore.json | jq -r '.address')

echo "Get Algorand receiver address"
ALGORAND_RECEIVER_ADDRESS=$(PRIVATE_KEY=$(cat cache/algo_receiver_private_key) get-public-key-hex)

echo "Transfer Test Token from minter to sender"
TXN_ID=$(goloop rpc sendtx call --method transfer --to $(cat cache/icon_test_token_addr) \
  --value 0 \
  --param _to=$ICON_SENDER_ADDRESS \
  --param _value=$SENDER_INITIAL_AMOUNT \
  --step_limit=3000000000 \
  --uri http://localhost:9080/api/v3/icon \
  --key_store test_token_minter.keystore.json --key_password gochain \
  --nid=$(cat cache/nid)
)
./../../algorand/scripts/wait_for_transaction.sh $TXN_ID

echo "Transfer Test Token to bridge"
TXN_ID=$(goloop rpc sendtx call --method transfer --to $(cat cache/icon_test_token_addr) \
  --value 0 \
  --param _to=$(cat cache/icon_escrow_addr) \
  --param _value=$TRANSFER_AMOUNT \
  --param _data=$ALGORAND_RECEIVER_ADDRESS \
  --step_limit=3000000000 \
  --uri http://localhost:9080/api/v3/icon \
  --key_store sender.keystore.json --key_password gochain \
  --nid=$(cat cache/nid)
)
./../../algorand/scripts/wait_for_transaction.sh $TXN_ID

echo "Check ICON Escrow balance"
ICON_ESCROW_BALANCE=$(goloop rpc call \
--from $(cat icon.keystore.json | jq -r '.address') \
--to $(cat cache/icon_test_token_addr) \
--method balanceOf \
--param _owner=$(cat cache/icon_escrow_addr) \
--uri http://localhost:9080/api/v3/icon
)

if [ $(printf "%d\n" $(echo $ICON_ESCROW_BALANCE | cut -d '"' -f 2)) != $TRANSFER_AMOUNT ]
then
      echo "Escrow balance should be equal to transfer amount"
      exit 1
fi

echo "Wait 60 seconds for transfer BTP message to Algorand"
sleep 60

echo "Check Algorand receiver account balance"
ALGO_RECEIVER_BALANCE=$(PRIVATE_KEY=$(cat cache/algo_receiver_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-asset-holding-amount $(cat cache/algo_wrapped_token_id))

if [ "$ALGO_RECEIVER_BALANCE" != "$TRANSFER_AMOUNT" ]
then
      echo "Algorand receiver Wrapped Token Balance should be equal to transfer amount"
      exit 1
fi

echo "Transfer some amount of Wrapped Token to another wallet in Algorand"
# get test account address
ALGORAND_TEST_ACCOUNT=$(PRIVATE_KEY=$(cat cache/algo_private_key) get-algorand-address)
# Opt in to asset
ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) PRIVATE_KEY=$(cat cache/algo_private_key) algorand-send-asset $(cat cache/algo_wrapped_token_id) $ALGORAND_TEST_ACCOUNT 0
# Transfer Wrapped Token from Algorand receiver account to algorand test account
ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) PRIVATE_KEY=$(cat cache/algo_receiver_private_key) algorand-send-asset $(cat cache/algo_wrapped_token_id) $ALGORAND_TEST_ACCOUNT 5000

echo "Burn Algorand Wrapped Token"
PRIVATE_KEY=$(cat cache/algo_receiver_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) algorand-burn-token ../../../pyteal/teal/reserve $(cat cache/bmc_app_id) $(cat cache/reserve_app_id) $ICON_RECEIVER_ADDRESS $(cat cache/algo_wrapped_token_id) 5000 

echo "Wait 30 seconds for transfer BTP message to ICON"
sleep 30

echo "Check ICON receiver balance"
ICON_RECEIVER_BALANCE=$(goloop rpc call \
--from $(cat icon.keystore.json | jq -r '.address') \
--to $(cat cache/icon_test_token_addr) \
--method balanceOf \
--param _owner=$ICON_RECEIVER_ADDRESS \
--uri http://localhost:9080/api/v3/icon
)

if [ $(printf "%d\n" $(echo $ICON_RECEIVER_BALANCE | cut -d '"' -f 2)) != 5000 ]
then
      echo "ICON receiver balance is not correct"
      exit 1
fi

echo "i2a integration test finish successfully"