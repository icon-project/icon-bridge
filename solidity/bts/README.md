# BTP Token Service (BTS)

BTP Token Service is a service handler smart contract which handles token transfer operations. It handles transfer of nativecoin and ERC20 tokens cross chain.
The BTS Contract in solidity is split into 2 sub contracts, BTSCore and BTSPeriphery. This was done because of the size limitation for solidity contracts.

* [BTSCore](BTSCore.md)
    - This contract is used to handle coin transferring service
* [BTSPeriphery](BTSPeriphery.md)
    - This contract is used to handle communications among BTP Message Center (BMC) and BTSCore contract.
    - It also maintains information about token limit and blacklist.


## Set up
Node >= 10.x
```
$ node --version
v15.12.0
```
Install tools
```
$ npm install --global yarn truffle@5.3.0
npm install @truffle/hdwallet-provider
npm install dotenv
npm install @openzeppelin/contracts
npm install @openzeppelin/contracts-upgradeable
```
Install dependencies
```
$ yarn
```

## Test
1. Run in a background process or seperate terminal window
```
$ docker run --rm -d -p 9933:9933 -p 9944:9944 purestake/moonbeam:v0.9.2 --dev --ws-external --rpc-external
```
2. Compile contracts
```
$ yarn contract:compile
```
3. Run unit and integration test
```
$ yarn test
```
-  Run specific test
```
$ yarn test:unit
$ yarn test:integration
```
