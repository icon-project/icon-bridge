# BTP Token Service (BTS)

BTP Token Service is a service handler smart contract which handles token transfer operations. It handles transfer of nativecoin and ERC20 tokens cross chain.
The BTS Contract in solidity is split into 2 sub contracts, BTSCore and BTSPeriphery. This was done because of the size limitation for solidity contracts.

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
$ npm install
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

### Deploy contracts

1. Compile contracts
```
$ npx hardhat compile
```

2. Deploy all BTS contracts

Network by default is Local. List of networks in the hardhat.config.js
The bmc address you can put as a command parameter 'bmcaddress' or in the env file 'PERIPHERY_ADDRESS'

```
$ npx hardhat deploy-bts
```
or
```
$ npx hardhat deploy-bts --bmcaddress <BMC Address 0x........>
```
or
```
$ npx hardhat deploy-bts --network arctic --bmcaddress <BMC Address 0x........>
```

3. Deploy only one contract (Don't do it if you deploy every thing on the step 2)

Deploy only BTS core contract

```
$ npx hardhat deploy-bts --bmcaddress <BMC Address 0x........>
```
or
```
$ npx hardhat deploy-bts --network arctic --bmcaddress <BMC Address 0x........>
```

Deploy only BTS periphery contract

```
$ npx hardhat deploy-bts-periphery --bmcaddress <BMC management address 0x........> --btscoreaddr <BTS core address 0x........>
```
or
```
$ npx hardhat deploy-bts-periphery --network arctic --bmcaddress <BMC management address 0x........> --btscoreaddr <BTS core address 0x........>
```