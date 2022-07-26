
const BEP20TKN = artifacts.require("ERC20TKN");
const { deployProxy } = require("@openzeppelin/truffle-upgrades");


module.exports = async function (deployer, network) {
  if (network != "development") {
    tx = await deployProxy(
      BEP20TKN, 
      [
          process.env.BSH_COIN_NAME,
          process.env.BSH_COIN_SYMBOL,
          process.env.BSH_DECIMALS,
          process.env.BSH_INITIAL_SUPPLY
      ],
      { deployer }
    );
  }
};

