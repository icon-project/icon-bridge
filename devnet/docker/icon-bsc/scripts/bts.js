const BTSCore = artifacts.require("BTSCore");
const BTSOwnerManager = artifacts.require("BTSOwnerManager");
module.exports = async function (callback) {
    try {
      var argv = require('minimist')(process.argv.slice(2), { string: ['addr'] });
      const btsCore = await BTSCore.deployed();
      const btsOwnerManager = await BTSOwnerManager.deployed();
      switch (argv["method"]) {
        case "setFeeRatio":
            tx = await btsCore.setFeeRatio(argv.name, argv.feeNumerator, `${argv.fixedFee}`);
            console.log(tx);
            break;
        case "register":
            tx = await btsCore.register(argv.name, argv.symbol, argv.decimals, argv.feeNumerator, `${argv.fixedFee}`, argv.addr);
            console.log(tx);
            break;
        case "coinId":
            tx = await btsCore.coinId(argv.coinName);
            console.log("coinId:",tx);
            break;      
        case "addBTSOwnerManager":
            console.log("Add BTSOwnerManager")
            tx = await btsCore.setBTSOwnerManager(argv.addr)
            console.log(tx)
            break;
        case "updateCoinDb":
            tx = await btsCore.updateCoinDb()
            console.log(tx)
            break;
        case "addOwner":
            console.log("Add bts owner ", argv.addr)
            tx = await btsOwnerManager.addOwner(argv.addr)
            console.log(tx)
            break;
        case "removeOwner":
            console.log("Remove bts owner ", argv.addr)
            tx = await btsOwnerManager.removeOwner(argv.addr)
            console.log(tx)
            break;
        case "getOwners":
            console.log("Fetch current bts owners ")
            tx = await btsOwnerManager.getOwners()
            console.log("BTS Contract owners are: ",tx)
            break;
        case "isOwner":
            console.log("Addr ", argv.addr)
            tx = await btsOwnerManager.isOwner(argv.addr)
            console.log("IsOwner:",tx)
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

