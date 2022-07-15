
const BEP20TKN = artifacts.require("ERC20TKN");

module.exports = async function (deployer, network) {
  if (network != "development") {
    await deployer.deploy(BEP20TKN);
  }
};