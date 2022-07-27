const HRC20 = artifacts.require("HRC20");

module.exports = async function (deployer, network) {
    if (network !== "development") {
        console.log('Start deploy HRC20 ');
        await deployer.deploy(HRC20);
        console.log('Finish deploy HRC20 ');
    }
};
