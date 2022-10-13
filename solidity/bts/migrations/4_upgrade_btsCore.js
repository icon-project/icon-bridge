const {upgradeProxy} = require('@openzeppelin/truffle-upgrades');
module.exports = async function (deployer, network) {
    if ((network !== "development") && (network !=="soliditycoverage")) {
        const argv = require('minimist')(process.argv.slice(2), {string: ['proxyAddr']});
        const BTSCoreUpgrade = artifacts.require(argv.btsCore);
        const existing = argv.proxyAddr

        console.log('Existing BTSCore Contract', existing);
        const instance = await upgradeProxy(existing, BTSCoreUpgrade, {deployer});
        console.log("Upgraded BTSCoreUpgraded Proxy", instance.address);
    }
};