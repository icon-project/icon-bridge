const BTSPeriphery = artifacts.require("BTSPeriphery");
const BTSCore = artifacts.require("BTSCore");
const { deployProxy } = require("@openzeppelin/truffle-upgrades");

module.exports = async function (deployer, network) {
    if (network !== "development") {
        await deployProxy(
            BTSCore,
            [
                process.env.BSH_COIN_NAME,
                parseInt(process.env.BSH_COIN_FEE),
                parseInt(process.env.BSH_FIXED_FEE),
            ],
            { deployer }
        );
        await deployProxy(
            BTSPeriphery,
            [process.env.BMC_PERIPHERY_ADDRESS, BTSCore.address],
            { deployer }
        );
        const btsCore = await BTSCore.deployed();
        await btsCore.updateBTSPeriphery(BTSPeriphery.address);
    }
};
