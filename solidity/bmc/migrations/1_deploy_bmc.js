const BMCManagement = artifacts.require('BMCManagement');
const BMCPeriphery = artifacts.require('BMCPeriphery');
const { deployProxy } = require('@openzeppelin/truffle-upgrades');

module.exports = async function (deployer, network) {
  if (network !== "development") {
    console.log('Start deploy Proxy BMCManagement ');
    await deployProxy(BMCManagement, { deployer });

    console.log('Start deploy Proxy BMCPeriphery ');
    await deployProxy(BMCPeriphery, [process.env.BMC_BTP_NET, BMCManagement.address], { deployer });

    console.log('Start deploy BMCManagement ');
    const bmcManagement = await BMCManagement.deployed();
    console.log('BMCManagement deployed, address: ' + bmcManagement.address);

    console.log('Set BMC Periphery ');
    console.log('BMC Periphery, address: ' + BMCPeriphery.address);
    await bmcManagement.setBMCPeriphery(BMCPeriphery.address);
  }
};