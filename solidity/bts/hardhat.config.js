require("@nomiclabs/hardhat-waffle");
require("hardhat-gas-reporter");
require('dotenv').config({path:__dirname+'/.env'})


// This is a sample Hardhat task. To learn how to create your own go to
// https://hardhat.org/guides/create-task.html
task("accounts", "Prints the list of accounts", async (taskArgs, hre) => {
  const accounts = await hre.ethers.getSigners();

  for (const account of accounts) {
    console.log(account.address);
  }
});

task("balances", "Prints the list of AVAX account balances", async () => {
  const accounts = await ethers.getSigners();

  for (const account of accounts) {
    balance = await ethers.provider.getBalance(account.address);
    console.log(account.address, "has balance", balance.toString());
  }
});

task("deploy-bts", "Deploy bts contract")
    .addParam("bmcaddress", "The address of the bmc contract")
    .setAction(async (taskArgs) => {
      const [deployer] = await hre.ethers.getSigners();
      console.log("Deploying contracts with the account:", deployer.address);

      //BTC Core
      const BTSCore = await hre.ethers.getContractFactory("BTSCoreV2");
      const bts = await BTSCore.deploy();
      await bts.deployed();
      console.log('BTSCore address:', bts.address);

      //BTSPeriphery
      const BTSPeriphery = await hre.ethers.getContractFactory("BTSPeriphery");
      const btsPeriphery = await BTSPeriphery.deploy();
      await btsPeriphery.deployed();
      await btsPeriphery.initialize(taskArgs.bmcaddress, bts.address);

      console.log('BTSPeriphery address:', btsPeriphery.address);
})

task("deploy-bts-core", "Deploy bts core contract", async (taskArgs, hre) => {
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying contracts with the account:", deployer.address);

  const BTSCoreV2 = await hre.ethers.getContractFactory("BTSCoreV2");
  const bts = await BTSCoreV2.deploy();
  await bts.deployed();

  console.log('BTSCoreV2 address:', bts.address);
});


task("deploy-bts-periphery", "Deploy bts periphery contract")
    .addParam("bmcaddress", "The address of the bmc contract")
    .addParam("btscoreaddr", "The bts cor contract address")
    .setAction(async (taskArgs) => {
      const [deployer] = await hre.ethers.getSigners();
      console.log("BMCPeriphery deploy contract account:", deployer.address);

      const BTSPeriphery = await hre.ethers.getContractFactory("BTSPeriphery");
      const bts = await BTSPeriphery.deploy();

      await bts.deployed();

      await bts.initialize(taskArgs.bmcaddress, taskArgs.btscoreaddr);

      console.log('BTSPeriphery address:', bts.address);
    });




const { RPC_URL, PRIVATE_KEY } = process.env;

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: {
    version: "0.8.2",
    settings: {
      optimizer: {
        enabled: true,
        runs: 10
      },
      evmVersion: "petersburg"
    }
  },
  defaultNetwork: "hardhat",
  networks: {
    hardhat: {
      //wei
      gasPrice: 100000000000,
      chainId: 553,
      network_id: process.env.BSC_NID,
      // allowUnlimitedContractSize: true,
      mining: {
        auto: true,
        interval: [3000, 6000]
      }
    },

    arctic: {
      url: RPC_URL,
      //wei
      gasPrice: 100000000000,
      blockGasLimit: 60000000000,
      chainId: 553,
      // allowUnlimitedContractSize: true,
      accounts: PRIVATE_KEY,
      timeout:250000,

      //Chain specific
      chainNetworkId:'0x61.bsc',
    },

    moonAlpha: {
      url: RPC_URL,
      //wei
      gasPrice: 1000000000,
      blockGasLimit: 30000000000,
      chainId: 1287,
      accounts: PRIVATE_KEY,
      timeout:250000
    }
  },

  gasReporter: {
    showMethodSig: true,
    enabled: true,
    //gwei
    gasPrice: 100,
    token: 'ICZ',
    currency: 'USD'
    // coinmarketcap:'7be4d35b-9f51-4dbe-830a-d840a2c7a165'
  },

  mocha: {
    timeout: 120000
  },


};

