# Getting Started 
This section is for deployment of ICON-Tezos bridge in the testnet. This file also documents all the dependencies and prerequisits that are necessary for the bridge. 

## Setting up on ICON
Install `goloop` cli for interacting with the ICON chain. 

```
    https://github.com/icon-project/goloop/blob/master/doc/build.md
```

## Setting up on Tezos 
Install `octez-client`, cli tool in tezos for interacting with the Tezos chain.

```
https://opentezos.com/tezos-basics/cli-and-rpc/
```
Once the `octez-client` is all set, you will have to change the network that your `octez-client` is connected to.

```sh
octez-client --endpoint https://ghostnet.tezos.marigold.dev config update
```

## ICON Tezos Bridge Configuration(Testnet)
For the complete deployment, build and running the relay navigate to 
```sh
$ cd $BASE_DIR/devnet/docker/icon-tezos/scripts
$ bash testnet.i2t.bmr.sh
```
When running the script, first wallets are created and we wil have to fund the wallets using the faucet of the respective chains.

For example you will get a message like this.
```
Fund the recently created wallet and run the script once again
icon bmc wallet: hxbc77026b7c3823744d5746507ab5bdb570ef9ca3
icon bts wallet: hx05256b36068144ec4523fd3943966ea018208ddb
icon bmr wallet: hx416d86b01f48ffe425a8a8a5f61182b1c17a339d
icon fa wallet : hxb7d5d680576de816459ad7b4af659886b2b0e4e3
tz bmr wallet : tz1dxhHuEcZNXoFyX3PX5A8NNpWnJ3MKHDY2
```
Fund the icon wallets using the faucet
```
https://faucet.iconosphere.io/
```
Fund the tezos wallets using the faucet
```
https://faucet.marigold.dev/
```

Rerun the script again and wait for the relay to start. This script will deploy the ICON smart contracts, Tezos smart contracts and register the native coin of the destination chains in its own chain. 

```sh
$ bash testnet.i2t.bmr.sh
```