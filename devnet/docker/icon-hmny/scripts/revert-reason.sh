#!/bin/bash

# This is for geth

# Fetching revert reason -- https://ethereum.stackexchange.com/questions/48383/how-to-receive-revert-reason-for-past-transactions

TX=$1
URI=${2:-http://localnets:9500}
SCRIPT=" tx = eth.getTransaction( \"$TX\" ); tx.data = tx.input; eth.call(tx, tx.blockNumber)"

geth --exec "$SCRIPT" attach $URI | cut -d '"' -f 2 | cut -c139- | xxd -r -p
echo
