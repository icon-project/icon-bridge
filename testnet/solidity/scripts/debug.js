const Web3 = require('web3');
const HDWalletProvider = require("@truffle/hdwallet-provider");
require('dotenv').config();
var Personal = require('web3-eth-personal')
const ethers = require('ethers')

const BMCManagement = require('../../../solidity/bmc/build/contracts/BMCManagement.json');
const BMCPeriphery = require('../../../solidity/bmc/build/contracts/BMCPeriphery.json');
const BSHImpl = require('../../../solidity/TokenBSH/build/contracts/BSHImpl.json');
const BEP20TKN = require('../../../solidity/TokenBSH/build/contracts/BEP20TKN.json');
const BSHProxy = require('../../../solidity/TokenBSH/build/contracts/BSHProxy.json');
const BSHCore = require('../../../solidity/bsh/build/contracts/BSHCore.json');
const BSHPeriphery = require('../../../solidity/bsh/build/contracts/BSHPeriphery.json');
const address = require('./addresses.json');

//To change for different network
const WS_HTTP = process.env.RPC_HTTP_URL
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
  let _prev = "btp://0x7.icon/cxf9d385353805e4112da8d9d926e7eefd0f44fdb5"
  let _relayMsg = "-QE--QE7uQE4-QE1AbkBLfkBKvkBJ7g5YnRwOi8vMHg2MS5ic2MvMHhGODU5ODUyODgwZDczMEY3RjU3NTNFMDY2MGUzMUNjOTkwN2Y1QjlFCbjp-Oe4OWJ0cDovLzB4Ny5pY29uL2N4NzgwNGQyMzc2YjhiYzg2MWUxZjI1NzA5NTcyNTQ4NWVkZTY4ZTY3Nbg5YnRwOi8vMHg2MS5ic2MvMHhGODU5ODUyODgwZDczMEY3RjU3NTNFMDY2MGUzMUNjOTkwN2Y1QjlFiFRva2VuQlNICLhl-GMAuGD4XqpoeGRmODE3ZDFlMjEzODg0MDgzMzE3YzRhOTc4Y2M1ZTIzOTM2ZThiODCqMHg1MDRjZjJkQWUwQUIzNGFDZWFBRUIyOTU0NjczNjE3ZjVCMUEyMkEwx8aDRVRIBQCDRCRx"
  
 try {

    let buff = new Buffer(_relayMsg, 'base64');
    let data = buff.toString('utf-8');
    
    let tx =await bmcPeriphery.methods.getStatus(link).call()
    console.log(tx)
    //tx =await bmcPeriphery.methods.handleRelayMessage(link, _relayMsg).send({ from: owner, gas: "4712388" })
    //console.log(tx);
    //console.log(await web3.utils.toUtf8("0xf2006d926b3aa2ed82aa34ecccd5ab84791fc6e895144725ef0f72dd478549e4"))
    
  } catch (ex) {
    console.log(ex)
    let error = errorJson(ex)
    console.log("Transaction failed, ", error.transactionHash);
    reason(error.transactionHash)
  }
}

init().then(async () => {
  console.log("Debug done")
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
  const provider1 = new ethers.providers.JsonRpcProvider(http_url)
  let tx = await provider1.getTransaction(hash)
  if (!tx) {
    console.log('tx not found')
  } else {
    let code = await provider1.call(tx, tx.blockNumber)
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