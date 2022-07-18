const BTSCore = artifacts.require("BTSCore");
const address = require('./addresses.json');
module.exports = async function (callback) {
    try {
      var argv = require('minimist')(process.argv.slice(2), { string: ['addr'] });
      const btsCore = await BTSCore.deployed();
      switch (argv["method"]) {
        case "setFeeRatio":
            tx = await btsCore.setFeeRatio(argv.name, argv.feeNumerator, argv.fixedFee);
            console.log(tx);
            break;
        case "register":
            tx = await btsCore.register(argv.name, argv.symbol, argv.decimals, argv.feeNumerator, argv.fixedFee, argv.addr);
            console.log(tx);
            break;
        case "coinId":
            tx = await btsCore.coinId(argv.coinName);
            console.log(tx);
            break;
        default:
            console.error("Bad input for method, ", argv)
      }
    }
    catch (error) {
        console.log(error)
    }
    callback()
}

