#!/bin/bash
export PATH=$PATH:~/go/bin

echo "Start a2i integration test"

AMOUNT=100
SENDER_BALANCE_BEFORE_TEST=$(PRIVATE_KEY=$(cat cache/algo_minter_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-asset-holding-amount $(cat cache/algo_test_asset_id))

# Create receiver account on ICON
goloop ks gen --out receiver.keystore.json
ICON_RECEIVER_ADDRESS=$(cat receiver.keystore.json | jq -r '.address')

# Transfer Asset to ICON
PRIVATE_KEY=$(cat cache/algo_minter_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) algorand-deposit-token ../../../pyteal/teal/escrow $(cat cache/bmc_app_id) $(cat cache/escrow_app_id) $ICON_RECEIVER_ADDRESS $(cat cache/algo_test_asset_id) $AMOUNT 

sleep 10

# Get Wrap Test Token (WTT) balance of receiver account
BALANCE=$(goloop rpc call \
  --from $(cat icon.keystore.json | jq -r '.address') \
  --to $(cat cache/icon_wtt_addr) \
  --method balanceOf \
  --param _owner=$ICON_RECEIVER_ADDRESS \
  --uri http://localhost:9080/api/v3/icon
)

# Check if receiver WTT balance is equal to amount sent from Algorand sender account
if [ $(printf "%d\n" $(echo $BALANCE | cut -d '"' -f 2)) != $AMOUNT ]
then
      echo "Balance is not equal to amount that we sent from Algorand"
      exit 1
fi

# Get Algorand sender public key for release assets
ALGORAND_RECEIVER_ADDRESS=$(PRIVATE_KEY=$(cat cache/algo_minter_private_key) get-public-key-hex)

# Burn WTT from ICON receiver account and send message to Algorand
TXN_ID=$(goloop rpc sendtx call --method burn --to $(cat cache/icon_wtt_addr) \
  --value 0 \
  --param _amount=$AMOUNT \
  --param algoPubKey=$ALGORAND_RECEIVER_ADDRESS \
  --step_limit=3000000000 \
  --uri http://localhost:9080/api/v3/icon \
  --key_store receiver.keystore.json --key_password gochain \
  --nid=$(cat cache/nid)
)

# Wait for executing transaction
./../../algorand/scripts/wait_for_transaction.sh $TXN_ID

sleep 60

# get Sender asset balance by asset id and sender private key
SENDER_BALANCE_AFTER_TEST=$(PRIVATE_KEY=$(cat cache/algo_minter_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-asset-holding-amount $(cat cache/algo_test_asset_id))

if [ "$SENDER_BALANCE_BEFORE_TEST" != "$SENDER_BALANCE_AFTER_TEST" ]
then
      echo "Sender Asset balance after test should be equal sender asset balance before test"
      exit 1
fi

echo "a2i integration test finish successfully"
