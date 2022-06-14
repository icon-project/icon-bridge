#!/bin/sh
source env.variables.sh
source nativeCoin.solidity.sh
source nativeCoin.javascore.sh

TOKENS_TRANSFER_AMOUNT=${1:-8}
NATIVE_COIN_NAME="ICX"

# ensure alice user keystore creation
printf "\n\nStep 1: creating/ensuring Alice keystore\n"
source /btpsimple/bin/keystore.sh
ensure_key_store alice.ks.json alice.secret

#Check Alice's balance before deposit
printf "\n\nStep 2: Alice's ICX balance before BTP Transfer\n"
get_alice_balance

printf "\n\nStep 3 Bob's $NATIVE_COIN_NAME balance before BTP Native transfer \n"
get_bob_ICX_balance
echo "$BOB_BALANCE"

# TODO:#bob approve transfer
#initiate Transfer from BSC to ICON from BSH
printf "\n\nStep 4: BOB Initiates BTP Native coin transfer of $TOKENS_TRANSFER_AMOUNT ($NATIVE_COIN_NAME) to Alice\n"
bsc_init_wrapped_native_btp_transfer "$NATIVE_COIN_NAME" "$TOKENS_TRANSFER_AMOUNT" >>$CONFIG_DIR/tx.native.bsc_icon.transfer.$NATIVE_COIN_NAME


#Check alice balance after 20s
printf "\n\nStep 6: Alice ICX Balance after BTP token transfer\n"
check_alice_native_balance_with_wait


printf "\n\n Bob's $NATIVE_COIN_NAME balance after BTP Native transfer \n"
get_bob_ICX_balance
echo "$BOB_BALANCE"
