const Wallet = require('ethereumjs-wallet').default
const util = require('ethereumjs-util')
const fs = require('fs')

const run = async () => {    
    let wl=new Wallet(util.toBuffer("PVT_KEY"))
    const publicKey = wl.getPublicKeyString();
    console.log(publicKey);
    const address = wl.getAddressString();
    console.log(address);
    let fn=wl.getV3Filename();
    const data =await wl.toV3("");  //password = ""
    console.log(data)
    fs.writeFileSync(address, JSON.stringify(data)); 
  };

  run();
