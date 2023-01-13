const BTSOwnerManager = artifacts.require("BTSOwnerManager");
const { deployProxy } = require("@openzeppelin/truffle-upgrades");

module.exports = async function (deployer, network) {
    if (network !== "development") {
        console.log('Start deploy Proxy BTSOwnerManager ');

        await deployProxy(
            BTSOwnerManager,
            { deployer }
        );

        const btsOwnerManager = await BTSOwnerManager.deployed();
        console.log("BTS Owner Manager deployed, address: " + btsOwnerManager.address);
    }
};
