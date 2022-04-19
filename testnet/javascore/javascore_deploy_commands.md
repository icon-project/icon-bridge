Step 1: Env setup
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
BTP Berlin:
```
URL="https://berlin.net.solidwallet.io/api/v3"
keystore="./berlinKeystore"
password="btpweb3labs"
nid="0x7" 
```


local btp docker:
```
URL="http://0.0.0.0:9080/api/v3/icon"
keystore="./data/goloop.keystore.json"
password="76a86e131c5572e9"
nid="0x954aa3"
```

From the testnet/javascore folder <br>
Check if the keystore wallet has enough balance to perform future transactions

```
./goloop rpc balance $(jq -r .address $keystore) --uri $URL
```



Step 2: Deploying the contracts & Configure
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Check if `var` folder exists or create one

1. Deploy BMC:

Change the _net param with proper source network id

```
./goloop rpc sendtx deploy artifacts/bmc-optimized.jar \
    --uri $URL --nid $nid\
    --key_store $keystore --key_password $password \
    --step_limit 10000000000 \
    --content_type application/java \
    --param _net="0x7.icon"  | jq -r . > var/tx.bmc


./goloop rpc --uri $URL txresult $(cat var/tx.bmc) | jq -r .scoreAddress > var/bmc
```

---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

2. DEPLOY BSH:

```
./goloop rpc --uri $URL sendtx deploy artifacts/bsh-optimized.jar \
    --key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --content_type application/java \
    --param _bmc=$(cat var/bmc)  | jq -r . > var/tx.bsh

./goloop rpc --uri $URL txresult $(cat var/tx.bsh) | jq -r .scoreAddress > var/bsh
```


Token deploy & fund scenario
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

3. DEPLOY IRC2 :

```
./goloop rpc --uri $URL sendtx deploy artifacts/irc2-token-optimized.jar \
    --key_store $keystore --key_password $password \
    --nid $nid --step_limit 10000000000 \
    --content_type application/java \
    --param _name="ETH" \
    --param _symbol="ETH" \
    --param _decimals="0x12" \
    --param _initialSupply="0x989680" | jq -r . > var/tx.irc2

./goloop rpc --uri $URL txresult  $(cat var/tx.irc2) | jq -r .scoreAddress > var/irc2

./goloop rpc --uri $URL txresult $(cat var/tx.irc2)
 ```


---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

4. Register IRC 2 Token with BSH:

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/bsh) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method register \
    --param name=ETH \
    --param symbol=ETH \
    --param decimals=0x12 \
    --param feeNumerator=0x64 \
    --param address=$(cat var/irc2) | jq -r . > var/tx.irc2.register

./goloop rpc --uri $URL txresult $(cat var/tx.irc2.register)

./goloop rpc --uri $URL  call --to $(cat var/bsh) \
    --method tokenNames  
```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

5. addService: Register BSH Service with BSC: 

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/bmc) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method addService \
    --param _svc=TokenBSH \
    --param _addr=$(cat var/bsh) 

./goloop rpc --uri $URL  call --to $(cat var/bmc) \
    --method getServices 
```

Mint IRC2 tokens to Alice
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

1. create new keystore - Alice

```
./goloop ks gen -o ./alice.json
```

2. transfer ICX balance from owner to Alice for transaction balance

```
./goloop rpc sendtx transfer \
--to $(jq -r .address alice.json) \
--value 190136500125000000000 \
--key_store $keystore \
--key_password $password \
--nid $nid \
--step_limit 13610920010 \
--uri $URL | jq -r . > var/tx.alice.transfer
    
./goloop rpc --uri $URL txresult $(cat var/tx.alice.transfer)
```

3. check alice balance

```
./goloop rpc balance $(jq -r .address alice.json) --uri $URL
```

4. Mint tokens to Alice

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/irc2) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method transfer  \
    --param _to=$(jq -r .address alice.json) \
    --param _value=100  | jq -r . > var/tx.alice.transfer.irc2

./goloop rpc --uri $URL  call --to $(cat var/irc2) \
    --method balanceOf \
    --param _owner=$(jq -r .address alice.json)
```

5. Mint tokens to TokenBSH

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/irc2) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method transfer  \
    --param _to=$(cat var/bsh) \
    --param _value=10000  | jq -r . > var/tx.bsh.transfer.irc2

./goloop rpc --uri $URL  call --to $(cat var/irc2) \
    --method balanceOf \
    --param _owner=$(cat var/bsh)
```

NativeCoin deploy & fund scenario:
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

1. deploy_javascore_nativeCoin_BSH


```
IRC2_SERIALIZED_SCORE=$(xxd -p artifacts/irc2Tradeable-0.1.0-optimized.jar | tr -d '\n')

```

```
./goloop rpc --uri $URL sendtx deploy artifacts/nativecoin-optimized.jar \
    --key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --content_type application/java \
    --param _bmc=$(cat var/bmc) \
    --param _serializedIrc2=0x$IRC2_SERIALIZED_SCORE \
    --param _name=ICX | jq -r . > var/tx.bsh.native

./goloop rpc --uri $URL txresult $(cat var/tx.bsh.native) | jq -r .scoreAddress > var/bsh.native
```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

3. bmc_javascore_addNativeService

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/bmc) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method addService \
    --param _svc=nativecoin \
    --param _addr=$(cat var/bsh.native)  | jq -r . > var/tx.addservice.native

./goloop rpc --uri $URL txresult $(cat var/tx.addservice.native) 


./goloop rpc --uri $URL  call --to $(cat var/bmc) \
    --method getServices 
```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

4. nativeBSH_javascore_register

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/bsh.native) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method register \
    --param _name=BNB \
    --param _symbol=BNB \
    --param _decimals=18 | jq -r . > var/tx.nativecoin.register

./goloop rpc --uri $URL txresult $(cat var/tx.nativecoin.register) 
```

Get the address of the deployed IRC2 token factory by Coin Name
```
./goloop rpc --uri $URL  call --to $(cat var/bsh.native) \
        --method coinAddress --param _coinName=BNB | sed -e 's/^"//' -e 's/"$//' > var/irc2TradeableToken.icon
```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

5. nativeBSH_javascore_setFeeRatio

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/bsh.native) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method setFeeRatio \
    --param _feeNumerator=100 | jq -r . > var/tx.setFeeRatio.nativebsh 
    
./goloop rpc --uri $URL txresult $(cat var/tx.setFeeRatio.nativebsh) 
```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


Configure BMC
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Note: better to finish the deployment of the solidity contracts before configuring BMC to keep the sync close to the latest block numbers

1. Add LINKS:

Register BTP Address of BSC BMC on ICON chain

`Note:Deploy BMC on bsc chain and then get the BMCPeriphery contract address to substitute along the network id to form _link`

```
cat ../solidity/var/bmc.periphery.bsc
```
```
./goloop rpc --uri $URL  sendtx call --to $(cat var/bmc) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method addLink \
    --param _link=btp://0x61.bsc/0x121A1AAd623AF68162B1bD84c44234Bc3a3562a9 | jq -r . > var/addLinks.tx.bmc

./goloop rpc --uri $URL  call --to $(cat var/bmc) \
    --method getLinks 
```

Remove LINKS:

Important: 
only to use this command if you want to remove the existing link from BMC
```
./goloop rpc --uri $URL  sendtx call --to $(cat var/bmc) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method removeLink \
    --param _link=btp://0x61.bsc/0x121A1AAd623AF68162B1bD84c44234Bc3a3562a9 
```

---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

2. AddRelay

Change the proper _link & wallet _addr
`Note:BMCPeriphery contract address to substitute`

```
./goloop rpc  --uri $URL sendtx call \
--to $(cat var/bmc) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method addRelay \
    --param _link=btp://0x61.bsc/0x121A1AAd623AF68162B1bD84c44234Bc3a3562a9 \
    --param _addr=hx681a290ecf0e460998d6bebe7b3da7589ed6b3db  | jq -r . > var/addRelay.tx.bmc

./goloop rpc --uri $URL txresult $(cat var/addRelay.tx.bmc)
```

Check status after adding the relay

```
./goloop rpc --uri $URL  call --to $(cat var/bmc) \
    --method getRelays\
    --param _link=btp://0x61.bsc/0x121A1AAd623AF68162B1bD84c44234Bc3a3562a9

./goloop rpc --uri $URL  call --to $(cat var/bmc) \
    --method getStatus \
    --param _link=btp://0x61.bsc/0x121A1AAd623AF68162B1bD84c44234Bc3a3562a9
```

Remove Relay

Important: only to use this command if you want to remove the existing relay from BMC

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/bmc) \
--key_store $keystore --key_password $password \
    --nid $nid --step_limit 13610920010 \
    --method removeRelay \
    --param _link=btp://0x61.bsc/0x45a0D0cda9e9Fb8e745B91104ca6444DC151D5A7 
```




Deposit & Initiate BTP Token Transfer Scenario
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


Note: Make sure the solidity contracts are deployed, configured properly & that the relay for both sides are running. before running this


1. Deposit token from Alice to TOKENBSH

Important: use alice keystore from here on out

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/irc2) \
--key_store ./alice.json --key_password gochain \
    --nid $nid --step_limit 13610920010 \
    --method transfer  \
    --param _to=$(cat var/bsh) \
    --param _value=0x05  | jq -r . > var/tx.deposit.bsh 

./goloop rpc --uri $URL txresult $(cat var/tx.deposit.bsh ) 
```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

2. check the balance of alice in BSH

```
./goloop rpc --uri $URL  call --to $(cat var/bsh) \
    --method getBalance  \
    --param user=$(jq -r .address ./alice.json) \
    --param tokenName=ETH
 ```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

3. Initiate BTP transfer on BSH

 change the address of `to` to user address on BSC

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/bsh) \
--key_store ./alice.json --key_password gochain \
    --nid $nid --step_limit 13610920010 \
    --method transfer  \
    --param tokenName=ETH \
    --param value=0x05  \
    --param to=btp://0x61.bsc/0x0baEAd25fe0346B76C73e84c083bb503c14309F1  | jq -r . > var/transfer_tx.bsh

./goloop rpc --uri $URL txresult $(cat var/transfer_tx.bsh)
```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

4. check the balance of alice in BSH should have the amount in locked balance

```
./goloop rpc --uri $URL  call --to $(cat var/bsh) \
    --method getBalance  \
    --param user=$(jq -r .address alice.json) \
    --param tokenName=ETH 
 
```


Deposit & Initiate BTP Native Transfer Scenario
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


1. Initiate BTP transfer on Native BSH from Alice to BOB

 change the address of `to` to user address on BSC

```
./goloop rpc --uri $URL  sendtx call --to $(cat var/bsh.native) \
--key_store ./alice.json --key_password gochain \
    --nid $nid --step_limit 13610920010 \
    --method transferNativeCoin  \
    --value=0x05  \
    --param _to=btp://0x61.bsc/0x0baEAd25fe0346B76C73e84c083bb503c14309F1  | jq -r . > var/nativetransfer_tx.bsh

./goloop rpc --uri $URL txresult $(cat var/nativetransfer_tx.bsh)
```

Notes:
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
RPC details:
https://www.icondev.io/introduction/the-icon-network/testnet


Create a new wallet:
`goloop ks gen --out berlinKeystore -p "password" > /dev/null 2>&1`
