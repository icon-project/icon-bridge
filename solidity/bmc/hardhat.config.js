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

task("deploy-bmc", "Deploy bmc management contract", async (taskArgs, hre) => {
      const [deployer] = await hre.ethers.getSigners();
      console.log("Deploying contracts with the account:", deployer.address);

      //BMCManagement
      const BMCManagement = await hre.ethers.getContractFactory("BMCManagement");
      const bmc = await BMCManagement.deploy();
      await bmc.deployed();
      console.log('BMCManagement address:', bmc.address);

      //BMCPeriphery
      const BMCPeriphery = await hre.ethers.getContractFactory("BMCPeriphery");
      const bmcPeriphery = await BMCPeriphery.deploy();
      await bmcPeriphery.deployed();
      await bmcPeriphery.initialize(hre.network.config.chainNetworkId, bmc.address);

      console.log('BMCPeriphery address:', bmcPeriphery.address);
})

task("deploy-bmc-management", "Deploy bmc management contract", async (taskArgs, hre) => {
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying contracts with the account:", deployer.address);

  const BMCManagement = await hre.ethers.getContractFactory("BMCManagement");
  const bmc = await BMCManagement.deploy();

  await bmc.deployed();

  console.log('BMCManagement address:', bmc.address);
});


task("deploy-bmc-periphery", "Deploy bmc periphery contract")
    .addParam("chainnetworkid", "The network Id")
    .addParam("bmcmanagementaddr", "The contract management address")
    .setAction(async (taskArgs) => {
      const [deployer] = await hre.ethers.getSigners();
      console.log("BMCPeriphery deploy contract account:", deployer.address);

      const BMCPeriphery = await hre.ethers.getContractFactory("BMCPeriphery");
      const bmc = await BMCPeriphery.deploy();

      await bmc.deployed();

      await bmc.initialize(taskArgs.chainnetworkid, taskArgs.bmcmanagementaddr);

      console.log('BMCPeriphery address:', bmc.address);
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

