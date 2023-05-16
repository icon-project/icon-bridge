#!/bin/bash
export PATH=$PATH:~/go/bin

echo "Starting i2a integration test"

TRANSFER_AMOUNT=1000000
ALGORAND_RECEIVER_ADDRESS=$(PRIVATE_KEY=$(cat cache/algo_receiver_private_key) get-public-key-hex)

ALGO_RECEIVER_BALANCE=$(PRIVATE_KEY=$(cat cache/algo_receiver_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-asset-holding-amount $(cat cache/algo_wrapped_token_id))

echo "Balance of Test Token on Algorand receiver address $ALGORAND_RECEIVER_ADDRESS is $ALGO_RECEIVER_BALANCE"
echo "Press enter to start transfer..."
read
echo "Executing token transfer from Icon to Algorand..."

echo "Transfer Test Token from minter to sender"
TXN_ID=$(goloop rpc sendtx call --method transfer --to $(cat cache/icon_test_token_addr) \
  --value 0 \
  --param _to=$(cat sender.keystore.json | jq -r '.address') \
  --param _value=5000000 \
  --step_limit=3000000000 \
  --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex \
  --key_store test_token_minter.keystore.json --key_password gochain \
  --nid=0x2
)
./../../algorand/scripts/wait_for_testnet_txn.sh $TXN_ID

echo "Transfer Test Token to bridge"
TXN_ID=$(goloop rpc sendtx call --method transfer --to $(cat cache/icon_test_token_addr) \
  --value 0 \
  --param _to=$(cat cache/icon_escrow_addr) \
  --param _value=$TRANSFER_AMOUNT \
  --param _data=$ALGORAND_RECEIVER_ADDRESS \
  --step_limit=3000000000 \
  --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex \
  --key_store sender.keystore.json --key_password gochain \
  --nid=0x2
)
./../../algorand/scripts/wait_for_testnet_txn.sh $TXN_ID

echo "Transaction $TXN_ID sent a message to mint $TRANSFER_AMOUNT Test Token on $ALGORAND_RECEIVER_ADDRESS"

ICON_ESCROW_BALANCE=$(goloop rpc call \
--from hx4b1a15d6781912a0285f1bfc47752f27cf54615b \
--to $(cat cache/icon_test_token_addr) \
--method balanceOf \
--param _owner=$(cat cache/icon_escrow_addr) \
--uri https://lisbon.net.solidwallet.io/api/v3/icon_dex
)

sleep 10

ALGO_RECEIVER_BALANCE=$(PRIVATE_KEY=$(cat cache/algo_receiver_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-asset-holding-amount $(cat cache/algo_wrapped_token_id))

echo "Transfer Complete"
echo "$ALGORAND_RECEIVER_ADDRESS Test Token balance: $ALGO_RECEIVER_BALANCE"
