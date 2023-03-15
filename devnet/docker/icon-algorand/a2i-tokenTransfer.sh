#!/bin/bash

echo "Starting a2i integration test"

AMOUNT=2500

goloop ks gen --out receiver.keystore.json
ICON_RECEIVER_ADDRESS=$(cat receiver.keystore.json | jq -r '.address')
SENDER_BALANCE_BEFORE_TEST=$(PRIVATE_KEY=$(cat cache/algo_minter_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-asset-holding-amount $(cat cache/algo_test_asset_id))

BALANCE=$(goloop rpc call \
  --from $(cat icon.keystore.json | jq -r '.address') \
  --to $(cat cache/icon_wtt_addr) \
  --method balanceOf \
  --param _owner=$ICON_RECEIVER_ADDRESS \
  --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex
)
echo "Balance of Test Token on Icon receiver address $ICON_RECEIVER_ADDRESS is $BALANCE"
echo "Press enter to start transfer..."
read
echo "Executing token transfer from Algorand to Icon..."

PRIVATE_KEY=$(cat cache/algo_minter_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) algorand-deposit-token ../../../pyteal/teal/escrow $(cat cache/bmc_app_id) $(cat cache/escrow_app_id) $ICON_RECEIVER_ADDRESS $(cat cache/algo_test_asset_id) $AMOUNT 

sleep 60

echo "Transfer Complete"
BALANCE=$(goloop rpc call \
  --from $(cat icon.keystore.json | jq -r '.address') \
  --to $(cat cache/icon_wtt_addr) \
  --method balanceOf \
  --param _owner=$ICON_RECEIVER_ADDRESS \
  --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex
)


echo "Balance of Test Token on Icon receiver address $ICON_RECEIVER_ADDRESS is $BALANCE"


