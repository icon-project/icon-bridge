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

    // Set default mocha options here, use special reporters etc.
    mocha: {
        // timeout: 100000
    },

    // Configure your compilers
    compilers: {
        solc: {
            version: "0.8.0", // Fetch exact version from solc-bin (default: truffle's version)
            // docker: true,        // Use "0.5.1" you've installed locally with docker (default: false)
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
    api_keys: {
        bscscan: process.env.BSC_SCAN_API_KEY,
    },
};
