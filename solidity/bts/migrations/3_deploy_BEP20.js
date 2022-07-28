
const BEP20TKN = artifacts.require("ERC20TKN");

module.exports = async function (deployer, network) {
  if (network !== "development") {
    console.log('Start deploy BEP20TKN ');
    await deployer.deploy(BEP20TKN);
    console.log('Start deploy BEP20TKN ');
  }
};