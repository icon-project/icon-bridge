const {upgradeProxy} = require('@openzeppelin/truffle-upgrades');
module.exports = async function (deployer, network) {
    if (network !== "develop") {
        const argv = require('minimist')(process.argv.slice(2), {string: ['upgradeTo']});
        
        const ContractName = argv.contract;
        const ContractUpgrade = artifacts.require(argv.upgradeTo);

        const existing = artifacts.require(ContractName).address;

        console.log(`Existing ${ContractName} Contract Address :>`, existing);
        const instance = await upgradeProxy(existing, ContractUpgrade, {deployer});
        console.log(`Upgraded  ${argv.upgradeTo} Upgraded Proxy Address :> `, instance.address);
    }
};
