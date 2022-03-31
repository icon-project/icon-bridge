const { TruffleProvider } = require("@harmony-js/core");
module.exports = {
  db: { enabled: false },
  compilers: {
    solc: {
      version: "0.7.6",
      settings: {
        evmVersion: "petersburg",
        optimizer: { enabled: true, runs: 10 },
      },
    },
  },
  networks: {
    hmny: {
      network_id: 2,
      provider: () => {
        const truffleProvider = new TruffleProvider(
          process.env.URI,
          { derivationPath: `m/44'/1023'/0'/0/` },
          { shardID: 0, chainId: 2 },
          { gasLimit: process.env.GASLIMIT, gasPrice: process.env.GASPRICE }
        );
        const newAcc = truffleProvider.addByPrivateKey(process.env.PRIVATE_KEY);
        truffleProvider.setSigner(newAcc);
        return truffleProvider;
      },
    },
  },
};
