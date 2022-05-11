const HDWalletProvider = require("@truffle/hdwallet-provider");

module.exports = {
    networks: {
        development: {
            provider: () =>
                new HDWalletProvider({
                    privateKeys: [],
                    providerOrUrl: "http://localhost:9500",
                }),
            network_id: "*",
        },
        hmyLocal: {
            provider: () =>
                new HDWalletProvider({
                    privateKeys: [],
                    providerOrUrl: "ws://localhost:9500",
                    chainId: 97,
                }),
            network_id: "02",
        },
        hmyTestnet: {
            provider: () =>
                new HDWalletProvider({
                    privateKeys: [process.env.PRIVATE_KEY],
                    providerOrUrl: process.env.RPC_URL,
                    chainId: 97,
                }),
            network_id: 97,
            skipDryRun: true,
            networkCheckTimeout: 1000000,
            timeoutBlocks: 200,
            gasPrice: 20000000000,
        },
    },
    mocha: {
        // timeout: 100000
    },
    compilers: {
        solc: {
            version: "0.8.0",
            settings: {
                // See the solidity docs for advice about optimization and evmVersion
                optimizer: {
                    enabled: true,
                    runs: 10,
                },
                evmVersion: "petersburg",
            },
        },
    },
    db: {
        enabled: false,
    },
    plugins: ["truffle-plugin-verify", "@chainsafe/truffle-plugin-abigen"],
};
