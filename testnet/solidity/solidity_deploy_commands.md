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
 truffle run verify BMCManagement@0xa135CE1aD1B5240ff4bb3044Bbb96CAbC438CD60 --network $network

 truffle run verify BMCPeriphery@0x99601b8f614b69f33bA5DE07eEE774f18D5DA051 --network $network
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
BSH_TOKEN_FEE=1 BMC_PERIPHERY_ADDRESS=0x45a0D0cda9e9Fb8e745B91104ca6444DC151D5A7  BSH_SERVICE=TokenBSH truffle migrate -f 2 --to 2 --compile-all --network $network
```

```
jq -r '.networks[] | .address' build/contracts/BSHImpl.json >../../testnet/solidity/var/bsh.impl.bsc

jq -r '.networks[] | .address' build/contracts/BSHProxy.json >../../testnet/solidity/var/bsh.proxy.bsc
```


Verify:

Note: Replace the Addresses from the logs, not the output files

```
 truffle run verify BSHProxy@0xd943b89A2694422B35cE1fAebAF40577E79C8BE0 --network $network

 truffle run verify BSHImpl@0xDe54f169f4a9721c4b7cd185B54D329c102202E6 --network $network
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
 truffle run verify BEP20TKN@0x77D94B32A15660E6A6bb587cc4c86f49456d3993 --network $network
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
    BMC_PERIPHERY_ADDRESS=0x45a0D0cda9e9Fb8e745B91104ca6444DC151D5A7 \
    BSH_SERVICE=nativecoin \
    truffle migrate --compile-all --network $network
```

```
    jq -r '.networks[] | .address' build/contracts/BSHCore.json > ../../testnet/solidity/var/bsh.core.bsc

    jq -r '.networks[] | .address' build/contracts/BSHPeriphery.json > ../../testnet/solidity/var/bsh.periphery.bsc
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

