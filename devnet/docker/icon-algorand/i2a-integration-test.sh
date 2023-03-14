#!/bin/bash

echo "Start i2a integration test"

TRANSFER_AMOUNT=1000000

# Get Algorand receiver address
ALGORAND_RECEIVER_ADDRESS=$(PRIVATE_KEY=$(cat cache/algo_receiver_private_key) get-public-key-hex)

# Transfer Test Token from minter to sender
TXN_ID=$(goloop rpc sendtx call --method transfer --to $(cat cache/icon_test_token_addr) \
  --value 0 \
  --param _to=$(cat sender.keystore.json | jq -r '.address') \
  --param _value=5000000 \
  --step_limit=3000000000 \
  --uri http://localhost:9080/api/v3/icon \
  --key_store test_token_minter.keystore.json --key_password gochain \
  --nid=$(cat cache/nid)
)
./../../algorand/scripts/wait_for_transaction.sh $TXN_ID

# Transfer Test Token to bridge
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

ICON_ESCROW_BALANCE=$(goloop rpc call \
--from hx4b1a15d6781912a0285f1bfc47752f27cf54615b \
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

# Wait for transfer BTP message to Algorand
sleep 60

ALGO_RECEIVER_BALANCE=$(PRIVATE_KEY=$(cat cache/algo_receiver_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-asset-holding-amount $(cat cache/algo_wrapped_token_id))

if [ "$ALGO_RECEIVER_BALANCE" != "$TRANSFER_AMOUNT" ]
then
      echo "Algorand receiver Wrapped Token Balance should be equal to transfer amount"
      exit 1
fi

echo "i2a integration test finish successfully"