## Set up
Node >= 10.x
```
$ node --version
v15.12.0
```
#### Install tools
```
$ npm install --global yarn truffle@5.3.0
npm install @truffle/hdwallet-provider
npm install dotenv
npm install @openzeppelin/contracts
npm install @openzeppelin/contracts-upgradeable
```
#### Install dependencies
```
$ npm install
```


## Deploy contracts

1. Compile contracts
```
$ npx hardhat compile
```

2. Deploy all BMC contracts

Network by default is Local. List of networks in the hardhat.config.js  

```
$ npx hardhat deploy-bmc
```
or
```
$ npx hardhat deploy-bmc --network arctic
```

3. Deploy only one contract (Don't do it if you deploy every thing on the step 2)

Deploy only BMC management contract

```
$ npx hardhat deploy-bmc-management
```
or
```
$ npx hardhat deploy-bmc-management --network arctic
```

Deploy only BMC periphery contract

```
$ npx hardhat deploy-bmc-periphery --chainnetworkid <0x61.bsc> --bmcmanagementaddr <BMC management address 0x........>
```
or
```
$ npx hardhat deploy-bmc-periphery --network arctic --chainnetworkid <0x61.bsc> --bmcmanagementaddr <BMC management address 0x........>
```