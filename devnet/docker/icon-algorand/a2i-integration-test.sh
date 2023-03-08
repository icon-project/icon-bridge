#!/bin/bash

AMOUNT=100

# Create receiver account on ICON
goloop ks gen --out receiver.keystore.json
RECEIVER_ADDRESS=$(cat receiver.keystore.json | jq -r '.address')
          
# Transfer Asset to ICON
PRIVATE_KEY=$(cat cache/algo_minter_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) algorand-deposit-token ../../../pyteal/teal/escrow $(cat cache/bmc_app_id) $(cat cache/escrow_app_id) $RECEIVER_ADDRESS $(cat cache/algo_test_asset_id) $AMOUNT 

sleep 10

BALANCE=$(goloop rpc call \
  --from $(cat icon.keystore.json | jq -r '.address') \
  --to $(cat cache/icon_wtt_addr) \
  --method balanceOf \
  --param _owner=$RECEIVER_ADDRESS \
  --uri http://localhost:9080/api/v3/icon
)

if [ $(printf "%d\n" $(echo $BALANCE | cut -d '"' -f 2)) != $AMOUNT ]
then
      echo "Balance is not equal to amount that we sent from Algorand"
      exit 1
fi