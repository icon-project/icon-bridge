const HDWalletProvider = require('@truffle/hdwallet-provider');
require('dotenv').config()

console.log('BSC_RPC_URI: ' + process.env.BSC_RPC_URI);

module.exports = {
  networks: {
    development: {
      provider: () => new HDWalletProvider({
        privateKeys: JSON.parse(process.env.PRIVATE_KEY),
        providerOrUrl: "http://localhost:9545",
      }),
      network_id: '*'
    },
    bsc: {
      provider: () => new HDWalletProvider({
        privateKeys: JSON.parse(process.env.PRIVATE_KEY),
        providerOrUrl: process.env.BSC_RPC_URI,
      }),
      network_id: process.env.BSC_NID,
      skipDryRun: true,
      networkCheckTimeout: 1000000,
      timeoutBlocks: 200,
      gasPrice: 20000000000,
    },
    hmny: {
      provider: () =>
          new HDWalletProvider({
              privateKeys: [process.env.PRIVATE_KEY],
              providerOrUrl: process.env.RPC_URL,
          }),
      network_id: process.env.NETWORK_ID,
      skipDryRun: true,
      networkCheckTimeout: 1000000,
      timeoutBlocks: 200,
      gasPrice: 20000000000,
    },
  },
  compilers: {
    solc: {
      version: "0.8.0",
      settings: {
        optimizer: {
          enabled: true,
          runs: 10
        },
        evmVersion: "petersburg"
      }
    }
  },
  plugins: ["truffle-plugin-verify", "@chainsafe/truffle-plugin-abigen"],
  api_keys: {
    bscscan: process.env.BSC_SCAN_API_KEY
  }
};
