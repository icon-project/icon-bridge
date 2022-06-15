#!/bin/sh
source env.variables.sh
source rpc.sh
source provision.sh
source nativeCoin.javascore.sh
source nativeCoin.solidity.sh

# ensure alice user keystore creation
printf "\n\nStep 1: creating/ensuring Alice keystore \n"
source /btpsimple/bin/keystore.sh
ensure_key_store alice.ks.json alice.secret

printf "\n\nStep 2 Bob's BNB balance before BTP Native transfer \n"
get_bob_BNB_balance
echo "$BOB_BNB_BALANCE"

printf "\n\nStep 3 Alice's BNB balance before BTP Native transfer \n"
get_alice_wrapped_native_balance "BNB"

#transfer 1 BNB from Alice to BSC BOB
rpcks alice.ks.json alice.secret
BNB_TRANSER_AMOUNT=$(coin2wei ${1:-1})
printf "\n\nStep 4: Alice Initiates BTP Native transfer of $(wei2coin $BNB_TRANSER_AMOUNT) BNB to BOB \n"
transfer_BNB_from_Alice_to_Bob $BNB_TRANSER_AMOUNT >>$CONFIG_DIR/tx.bnb.native.icon_bsc.transfer
wait_for_file $CONFIG_DIR/tx.bnb.native.icon_bsc.transfer

#get Bob's balance after BTP transfer with wait
printf "\n\nStep 5: Bob's BNB balance after Deposit \n"
get_Bob_BNB_Balance_with_wait
