Step 1: Env setup
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

```
network="bscTestnet"
```

Update `.env`  with the testnet wallet private key with funds for all the `bmc`, `bsh`, & `TokenBSH` projects.


STEP 2: Deploy Contracts
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

1. BMC Deploy

From the root folder execute the following commands:<br/>
Only for the first time to clean up artifacts folder:

```
cd solidity/bmc/
rm -rf .openzeppelin
rm -rf build
yarn install
```

Deploy Command

Note: change the proper network ID in `BMC_PRA_NET`

```
BMC_PRA_NET=0x61.bsc truffle migrate --network $network --compile-all
```

Copy the deployed address to a file

```
jq -r '.networks[] | .address' build/contracts/BMCPeriphery.json >../../testnet/solidity/var/bmc.periphery.bsc

jq -r '.networks[] | .address' build/contracts/BMCManagement.json >../../testnet/solidity/var/bmc.management.bsc
```

Note:

Substitute BMCPeriphery address from above outfiles into the btp address in `addlink` back in javascore deployment scripts

Verify:

Note: Replace the Addresses from the logs, not the output files

```
 truffle run verify BMCManagement@0xc2655eCDDd7320665ffda10c4Bbdc02842D7EbA6 --network $network

 truffle run verify BMCPeriphery@0xDa1CE8349ea6aFc2aDBc8895503369A620686ef5 --network $network
```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------


2. Token BSH deploy

From the root folder execute the following commands:<br/>
Only for the first time to clean up artifacts folder:

```
network="bscTestnet"
cd solidity/TokenBSH
rm -rf .openzeppelin
rm -rf build
yarn install
```

Deploy Command:

Note: change the BMC address in `BMC_PERIPHERY_ADDRESS`

```
BSH_TOKEN_FEE=1 BMC_PERIPHERY_ADDRESS=0x121A1AAd623AF68162B1bD84c44234Bc3a3562a9  BSH_SERVICE=TokenBSH truffle migrate -f 2 --to 2 --compile-all --network $network
```

```
jq -r '.networks[] | .address' build/contracts/BSHImpl.json >../../testnet/solidity/var/bsh.impl.bsc

jq -r '.networks[] | .address' build/contracts/BSHProxy.json >../../testnet/solidity/var/bsh.proxy.bsc
```


Verify:

Note: Replace the Addresses from the logs, not the output files

```
 truffle run verify BSHProxy@0xDA5eE5f3cc7a98a9615F655fEd3B1a97197a2521 --network $network

 truffle run verify BSHImpl@0xEd2B59A9F160408D5a1A384e01ec1B16F7C6892E --network $network
```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

3. BEP20 Deploy

```
truffle migrate -f 3 --network $network --compile-none 
```

```
jq -r '.networks[] | .address' build/contracts/BEP20TKN.json >../../testnet/solidity/var/bep20.bsc
```


Verify:

Note: Replace the Addresses from the logs, not the output files

```
 truffle run verify BEP20TKN@0x0df55BeBdF518C3c2937498E4459414a4Ac124DD --network $network
```
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

4. native BSH Deploy

```
network="bscTestnet"
cd solidity/bsh
rm -rf .openzeppelin
rm -rf build

yarn install
```


Deploy Command: <br>
Note: change the BMC address in `BMC_PERIPHERY_ADDRESS`


```
BSH_COIN_URL=https://www.binance.com/en/ \
    BSH_COIN_NAME=BNB \
    BSH_COIN_FEE=100 \
    BSH_FIXED_FEE=50000 \
    BMC_PERIPHERY_ADDRESS=0x121A1AAd623AF68162B1bD84c44234Bc3a3562a9 \
    BSH_SERVICE=nativecoin \
    truffle migrate --compile-all --network $network
```

```
    jq -r '.networks[] | .address' build/contracts/BSHCore.json > ../../testnet/solidity/var/bsh.core.bsc

    jq -r '.networks[] | .address' build/contracts/BSHPeriphery.json > ../../testnet/solidity/var/bsh.periphery.bsc
```

```
 truffle run verify BSHCore@0x2fBa3e9211e5327D2e90dAEd70Ef7307c2B31C99 --network $network

 truffle run verify BSHPeriphery@0xd7637F4AA7BC3648E7741Ac1e555d43d8BBd8921 --network $network
 ```

Provisioning scripts:
---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

The scripts to provision the BMC & BSH contracts are in the `scripts` folder.
1. `yarn install`
2. edit addresses.json and copy the deployed addresses from the `var/` folder and paste them into the addresses.json file for respective contract addresses
3. edit `.env` file to add RPC WS URL & HTTP URL along with private key
3. `node provision.js` (edit if the network is not 0x7.icon)


Note:

RPC details:
https://docs.binance.org/smart-chain/developer/rpc.html

