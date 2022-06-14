const Web3 = require('web3');
const HDWalletProvider = require("@truffle/hdwallet-provider");
require('dotenv').config();
var Personal = require('web3-eth-personal')

const BEP20TKN = require('../../../solidity/TokenBSH/build/contracts/BEP20TKN.json');
const BSHProxy = require('../../../solidity/TokenBSH/build/contracts/BSHProxy.json');
const address = require('./addresses.json');

const WS_HTTP = process.env.RPC_HTTP_URL
const ws_url = process.env.RPC_WS_URL
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

  const bshProxy = await new web3.eth.Contract(
    BSHProxy.abi,
    address.solidity.BSHProxy,
    { from: owner, gas: gasAmount }
  );

  const bep20tkn = await new web3.eth.Contract(
    BEP20TKN.abi,
    address.solidity.BEP20TKN,
    { from: owner, gas: gasAmount }
  );

  bshProxy.address = address.solidity.BSHProxy
  bep20tkn.address = address.solidity.BEP20TKN

  //Approve & initiate Transfer
  try {
    let transferAmount = 10
    var _to = 'btp://0x7.icon/hx275c118617610e65ba572ac0a621ddd13255242b';
    tokenName = "ETH"
    //approve & initiate transfer
    var balanceBefore = await bshProxy.methods.getBalanceOf(owner, tokenName).call()
    console.log("BEP20 token user balance before transfer", balanceBefore)
    await bep20tkn.methods.approve(bshProxy.address, transferAmount).send({ from: owner, gas: gasAmount });
    var txn = await bshProxy.methods.transfer(tokenName, transferAmount, _to).send({ from: owner, gas: gasAmount })
    console.log(txn);
    //check balance after transfer
    var balanceafter = await bshProxy.methods.getBalanceOf(owner, tokenName).call()
    console.log("BEP20 token user balance after transfer", balanceafter)
    let bshBal = await bep20tkn.methods.balanceOf(bshProxy.address).call();
    console.log("Balance of BSH Proxy" + bshProxy.address + " after the transfer:" + bshBal);

    console.log("Transfer init done")
  } catch (ex) {
    console.log(ex)
    let error = errorJson(ex)
    console.log("Init transfer Transaction failed, ", error.transactionHash);
    reason(error.transactionHash)
  }
  //End of Approve & initiate Transfer
}

init().then(async () => {
  console.log("Transfer done.")
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