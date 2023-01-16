const BTSPeriphery = artifacts.require("BTSPeriphery");
const BTSCore = artifacts.require("BTSCore");
const BTSOwnerManager = artifacts.require("BTSOwnerManager");
const {deployProxy} = require("@openzeppelin/truffle-upgrades");

module.exports = async function (deployer, network) {
    if (network !== "development") {
        console.log('Start deploy Proxy BTSOwnerManager ');

        await deployProxy(
            BTSOwnerManager,
            { deployer }
        );

        const btsOwnerManager = await BTSOwnerManager.deployed();

        console.log('Start deploy Proxy BTSCore ');

        await deployProxy(
            BTSCore,
            [
                process.env.BSH_COIN_NAME,
                process.env.BSH_COIN_FEE,
                process.env.BSH_FIXED_FEE,
                btsOwnerManager.address
            ],
            {deployer}
        );

        console.log('Start deploy Proxy BTSPeriphery ');
        await deployProxy(
            BTSPeriphery,
            [process.env.BMC_PERIPHERY_ADDRESS, BTSCore.address],
            {deployer}
        );

        console.log('Start deploy BTSCore ');
        const btsCore = await BTSCore.deployed();

        console.log('Updating BTS periphery');
        await btsCore.updateBTSPeriphery(BTSPeriphery.address);

        console.log("BTS Owner Manager deployed, address: " + btsOwnerManager.address);
        console.log('BTSCore: ' + btsCore.address);
        console.log('BTSPeriphery: ' + BTSPeriphery.address);

    }
};
