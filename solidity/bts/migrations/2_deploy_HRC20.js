const HRC20 = artifacts.require("HRC20");

module.exports = async function (deployer, network) {
    if (network !== "development") {
        await deployer.deploy(HRC20);
    }
};
