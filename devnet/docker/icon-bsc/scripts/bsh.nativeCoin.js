const BSHCore = artifacts.require("BSHCore");
const ERC20 = artifacts.require("ERC20Tradable");
module.exports = async function (callback) {
  try {
    var argv = require('minimist')(process.argv.slice(2), { string: ['addr', 'from'] });
    const bshCore = await BSHCore.deployed();
    let tx;
    switch (argv["method"]) {
      case "register":
        console.log("registerCoin", argv.name)
        tx = await bshCore.register(argv.name, argv.symbol, argv.decimals);
        //console.log(await bshCore.coinNames())
        console.log(tx)
        break;
      case "getBalanceOf":
        var balance = await bshCore.getBalanceOf(argv.addr, argv.name);
        //console.log("balance of user:" + argv.addr + " = " + balance._usableBalance);
        var bal = await web3.utils.fromWei(balance._usableBalance, "ether")
        console.log(bal)
        break;
      case "transferNativeCoin":
        console.log("Init BTP native transfer of " + web3.utils.toWei("" + argv.amount, 'ether') + " wei to " + argv.to)
        tx = await bshCore.transferNativeCoin(argv.to, { from: argv.from, value: web3.utils.toWei("" + argv.amount, 'ether') })
        console.log(tx)
        break;
      case "transferWrappedNativeCoin":
        var coinAddress = await bshCore.coinId(argv.coinName);
        console.log(argv.coinName+" deployed address:"+ coinAddress)

        console.log("Approving native BSH to use Bob's tokens")
        var wrappedToken = await ERC20.at(coinAddress);
        tx = await wrappedToken.approve(bshCore.address, web3.utils.toWei("" + argv.amount, 'ether'), { from: argv.from })

        console.log("Init BTP wrapped native transfer of " + web3.utils.toWei("" + argv.amount, 'ether') + " wei to " + argv.to)
        tx = await bshCore.transfer(argv.coinName, web3.utils.toWei("" + argv.amount, 'ether'), argv.to, { from: argv.from})
        console.log(tx)

        console.log("User ("+argv.name+") balance after transfer")
        var balance = await wrappedToken.balanceOf(argv.from);
        break;
      case "coinId":
        var coinAddress = await bshCore.coinId(argv.name);
        console.log(coinAddress)
        break;
      case "approve":
        console.log("Approving native BSH to use Bob's tokens")
        var wrappedToken = await ERC20.at(argv.coinAddress);
        var balance = await wrappedToken.approve(argv.addr, web3.utils.toWei("" + argv.amount, 'ether'), { from: argv.from })
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