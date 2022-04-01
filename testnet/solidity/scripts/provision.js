const Web3 = require('web3');
const HDWalletProvider = require("@truffle/hdwallet-provider");
require('dotenv').config();
var Personal = require('web3-eth-personal')

const BMCManagement = require('../../../solidity/bmc/build/contracts/BMCManagement.json');
const BMCPeriphery = require('../../../solidity/bmc/build/contracts/BMCPeriphery.json');
const BSHImpl = require('../../../solidity/TokenBSH/build/contracts/BSHImpl.json');
const BEP20TKN = require('../../../solidity/TokenBSH/build/contracts/BEP20TKN.json');
const BSHProxy = require('../../../solidity/TokenBSH/build/contracts/BSHProxy.json');
const BSHCore = require('../../../solidity/bsh/build/contracts/BSHCore.json');
const BSHPeriphery = require('../../../solidity/bsh/build/contracts/BSHPeriphery.json');
const address = require('./addresses.json');

//To change for different network

const ws_url = process.env.RPC_WS_URL
let link = "btp://0x7.icon/" + address.javascore.bmc
console.log('link:', link)
const gasAmount = "4712388"

const provider = new HDWalletProvider(
  process.env.PRIVATE_KEY,
  ws_url,
);

const init = async () => {
  let web3 = new Web3(provider);
  const networkId = await web3.eth.net.getId();
  console.log(networkId)
  const myAccount = web3.eth.accounts.privateKeyToAccount(process.env.PRIVATE_KEY);
  console.log(myAccount.address)
  let owner = myAccount.address
  var balance = await web3.eth.getBalance(myAccount.address);
  console.log(balance)
  var eth = await web3.utils.fromWei(balance);
  console.log(eth)

  const bmcManagement = await new web3.eth.Contract(
    BMCManagement.abi,
    address.solidity.BMCManagement,
    { from: owner, gas: gasAmount }
  );

  const bmcPeriphery = await new web3.eth.Contract(
    BMCPeriphery.abi,
    address.solidity.BMCPeriphery,
    { from: owner, gas: gasAmount }
  );

  const bshImpl = await new web3.eth.Contract(
    BSHImpl.abi,
    address.solidity.BSHImpl,
    { from: owner, gas: gasAmount }
  );

  const bshProxy = await new web3.eth.Contract(
    BSHProxy.abi,
    address.solidity.BSHProxy,
    { from: owner, gas: gasAmount }
  );

  const bshCore = await new web3.eth.Contract(
    BSHCore.abi,
    address.solidity.BSHCore,
    { from: owner, gas: gasAmount }
  );

  const bshPeriphery = await new web3.eth.Contract(
    BSHPeriphery.abi,
    address.solidity.BSHPeriphery,
    { from: owner, gas: gasAmount }
  );

  const bep20tkn = await new web3.eth.Contract(
    BEP20TKN.abi,
    address.solidity.BEP20TKN,
    { from: owner, gas: gasAmount }
  );

  bmcManagement.address = address.solidity.BMCManagement
  bmcPeriphery.address = address.solidity.BMCPeriphery
  bshImpl.address = address.solidity.BSHImpl
  bshProxy.address = address.solidity.BSHProxy
  bep20tkn.address = address.solidity.BEP20TKN
  bshCore.address = address.solidity.BSHCore
  bshPeriphery.address = address.solidity.BSHPeriphery


  try {
    //await bmcManagement.methods.removeService("TokenBSH").send({from:owner}) ;   

    let tx = await bmcManagement.methods.addService("TokenBSH", bshImpl.address).send({ from: owner });
    console.log("Add Service TokenBSH: ", tx);

    tx = await bmcManagement.methods.addService("nativecoin", bshPeriphery.address).send({ from: owner });
    console.log("Add Service Nativecoin: ", tx);

    console.log("Services #####################");
    console.log(await bmcManagement.methods.getServices().call());

    tx = await bshCore.methods.register("ICX", "ICX", 18).send({ from: owner });
    console.log("Register Native coin ICX : ", tx);

    coinAddress = await bshCore.methods.coinId("ICX").call();
    console.log("ICX address: ", coinAddress);

    console.log("Registered Coin Names #####################");
    console.log(await bshCore.methods.coinNames().call())

    tx = await bmcManagement.methods.addLink(link).send({ from: owner, gas: gasAmount });
    console.log("Add Link: ", tx);

    let links = await bmcManagement.methods.getLinks().call();
    console.log("Links #####################");
    console.log(links);

    tx = await bmcManagement.methods.addRelay(link, [owner]).send({ from: owner, gas: gasAmount });
    console.log("Add Relay: ", tx);

    let relays = await bmcManagement.methods.getRelays(link).call();
    console.log("Relays #####################");
    console.log(relays);

    bmcPeriphery.methods.getStatus(link).call().then(console.log)

    tx = await bshProxy.methods.register("ETH", "ETH", 18, 100, bep20tkn.address).send({ from: owner, gas: gasAmount })
    console.log("BSH Register Token: ", tx);

    console.log("Registered Tokens #####################");
    let tokens = await bshProxy.methods.tokenNames().call();
    console.log(tokens);

    //Fund BSH with some tokens
    var balanceBeforeTransfer = await bep20tkn.methods.balanceOf(bshProxy.address).call();
    console.log("BSHProxy Balance Before before funding:", balanceBeforeTransfer)
    await bep20tkn.methods.transfer(bshProxy.address, 10000).send({ from: owner, gas: gasAmount });
    var balanceAfterTransfer = await bep20tkn.methods.balanceOf(bshProxy.address).call();
    console.log("BSHProxy Balance Before after funding:", balanceAfterTransfer)

  } catch (ex) {
    console.log(ex)
    let error = errorJson(ex)
    console.log("Transaction failed, ", error.transactionHash);
    reason(error.transactionHash)
  }
}

init().then(async () => {
  console.log("Provision done")
});


function errorJson(ex) {
  let lines = ex.toString().split('\n');
  lines.splice(0, 1);
  return JSON.parse(lines.join('\n'));
}

function hex_to_ascii(str1) {
  var hex = str1.toString();
  var str = '';
  for (var n = 0; n < hex.length; n += 2) {
    str += String.fromCharCode(parseInt(hex.substr(n, 2), 16));
  }
  return str;
}

async function reason(hash) {
  console.log('tx hash:', hash)
  let tx = await provider.getTransaction(hash)
  if (!tx) {
    console.log('tx not found')
  } else {
    let code = await provider.call(tx, tx.blockNumber)
    console.log(code)
    let reason = hex_to_ascii(code.substr(138))
    console.log('revert reason:', reason)
  }
}



/* const tx = await bmcManagement.methods.addService("TokenBSH", bshImpl.address);
const gas = await tx.estimateGas({ from: owner, gas: gasAmount });
const gasPrice = await web3.eth.getGasPrice();
const data = tx.encodeABI();
const nonce = await web3.eth.getTransactionCount(owner);
const signedTx = await web3.eth.accounts.signTransaction({
  nonce: nonce,
  gas: gas,
  gasPrice: gasPrice,
  data: data,
  to: bmcManagement.address,
},
process.env.PRIVATE_KEY);
const txHash = await web3.eth.sendSignedTransaction(signedTx.rawTransaction);
console.log(txHash); */