const HDWalletProvider = require("@truffle/hdwallet-provider");

module.exports = {
    db: { enabled: false },
    compilers: {
        solc: {
            version: "0.8.0",
            settings: {
                evmVersion: "petersburg",
                optimizer: { enabled: true, runs: 10 },
            },
        },
    },
    networks: {
        hmny: {
            network_id: parseInt(process.env.NETWORK_ID),
            provider: () => {
                if (!process.env.PRIVATE_KEY) {
                    throw new Error(`PRIVATE_KEY not provided!`);
                }
                return new HDWalletProvider({
                    privateKeys: [process.env.PRIVATE_KEY],
                    providerOrUrl: process.env.URI,
                    derivationPath: `m/44'/1023'/0'/0/`,
                });
            },
            skipDryRun: true,
        },
    },
};
