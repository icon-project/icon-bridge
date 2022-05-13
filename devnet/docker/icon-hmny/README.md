# **`ixh.sh`**

This is a helper tool to auto-deploy harmony and icon localnets on local/remote host using docker. It has commands to deploy BTP smartcontracts to respective blockchain (`deploysc`), and run tests (`test`) using several subcommands. It also includes a command to run a demo (`demo`).

Please go through the `ixh.sh` file and see other commands at the bottom of the script and their syntax.

## Requirements

Make sure that following tools are installed on your system for `ixh.sh` to work without any issue.

1.  ### Docker

    To build, publish, run blockchains (icon/hmny) locally or remote docker host. Download and install docker from https://docs.docker.com/engine/install/ubuntu/

    After installing, make sure that the user account used to run docker (_default is ubuntu_) is added to `docker` group.

        $ sudo groupadd docker
        $ sudo usermod -aG docker $USER
        $ newgrp docker

    Fully logout, and log back in to be able apply the changes.

2.  ### SdkMan

    To install gradle and java.

    1. _`fish`_

       https://github.com/reitzig/sdkman-for-fish

    2. _`bash`_

       https://sdkman.io/install

3.  ### Java and Gradle

    To build javascores.

    1. _`Java`_

       `sdk install java 11.0.11.hs-adpt`

    2. _`gradle`_

       `sdk install gradle 6.7.1`

4.  ### Goloop

    https://github.com/icon-project/goloop

    To interact with icon blockchain using RPC calls and generate keystores.

    `go install github.com/icon-project/goloop/cmd/goloop`

    If `go install` doesn't work use `go get` instead.

5.  ### NodeJS

    To build and deploy solidity smartcontracts.

    1. _`fish`_

       `nvm`: https://github.com/jorgebucaran/nvm.fish

       ```
       $ fisher install jorgebucaran/nvm.fish
       $ nvm install v15.12.0
       $ set --universal nvm_default_version v15.12.0
       $ nvm use v15.12.0
       $ node --version > ~/.nvmrc
       ```

    2. _`bash`_

       https://github.com/nvm-sh/nvm

       ```
       $ curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
       $ nvm install v15.12.0
       $ nvm use v15.12.0
       $ node --version > ~/.nvmrc
       ```

6.  ### Truffle

    https://trufflesuite.com/docs/truffle/getting-started/installation.html

    `npm install -g truffle@5.5.5`

7.  ### Ethkey

    `go get github.com/ethereum/go-ethereum/cmd/ethkey`

8. ### Local BLS Dependencies
    You may need additional build tools to sucessfully build following libraries.

    ```
    git clone https://github.com/harmony-one/bls.git
    git clone https://github.com/harmony-one/mcl.git

    cd bls && make -j8 BLS_SWAP_G=1 && make install && cd ..
    cd mcl && make install && cd ..
    ```

## Build

    ./ixh.sh build

## Publish images

    ./ixh.sh publish

## Deploy Blockchains

    ./ixh.sh start nodes -d

NOTE:

If you're using a remote docker host, please allow ssh connection either password less or with key added to local ssh-agent. It needs the ablility to ssh into the remote host non-interactively. You also need a local docker registry on the remote host running at 0.0.0.0:5000, to publish necessary docker images of respective blockchains, so that it could be used on that host. If hostname is used instead of an IP address, you need to map the hostname to IP in `/etc/hosts` in a linux system, corresponding files in windows/macos system.

## Deploy Smart Contracts

    ./ixh.sh deploysc reset

It deploys necessary smart contracts and generates configuration for Relayer (`_ixh/bmr.config.json`) along with an environment file (`_ixh/ixh.env`) that has all necessary environment variables.

NOTE: _Wait for 1 minute or more before doing this after deploying blockchains to ensure both chains have started properly._

## Start Relayer
If you have `bls` dependencies installed in local system, you can run following commands from the project root directory to start relayer locally.

```
$ cd cmd/btpsimple
$ go run . -config ../../devnet/docker/icon-hmny/_ixh/bmr.config.json
```

If you have docker, you can chose to use docker image: `bmr` instead. Please use following command to start `bmr` container.
```
./ixh.sh start bmr -d && ./ixh.sh docker_compose bmr logs -f --tail=1000 | grep "^bmr*"
```

Whichever approach you prefer from above options, you should see console logs similar to the following.
```
...
bmr   | I|06:16:53.751186|0xEC|relay|i2h|relay.go:61 init: link.rxSeq=0, link.rxHeight=1658
bmr   | D|06:16:53.767850|0xEC|icon|i2h|rx_client.go:230 MonitorBlock WSEvent 127.0.0.1:41872 WSEventInit
bmr   | D|06:16:53.767912|0xEC|icon|i2h|rx_receiver.go:252 connected local=127.0.0.1:41872
bmr   | I|06:16:53.778151|hx44|relay|h2i|relay.go:61 init: link.rxSeq=0, link.rxHeight=1641
...
bmr   | D|06:16:53.808659|0xEC|icon|i2h|rx_receiver.go:241 block notification height=1936
bmr   | I|06:16:53.808712|0xEC|relay|i2h|relay.go:128 srcMsg added seq=[1 1]
bmr   | D|06:16:53.808755|0xEC|icon|i2h|rx_receiver.go:241 block notification height=1937
...
bmr   | D|06:16:56.057912|hx44|hmny|h2i|rx_receiver.go:369 block notification height=1930
bmr   | I|06:16:56.058026|hx44|relay|h2i|relay.go:128 srcMsg added seq=[1 1]
bmr   | D|06:16:56.061405|hx44|hmny|h2i|rx_receiver.go:369 block notification height=1931
...
bmr   | D|06:16:58.251836|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2163
bmr   | D|06:16:58.752245|0xEC|relay|i2h|relay.go:100 relaySignal
bmr   | D|06:16:58.755332|0xEC|hmny|i2h|tx_sender.go:211 handleRelayMessage: send tx prev=btp://0x5b9a77.icon/cxbaf6c209178820d6969316ea5b1dd4f3a91c463a
bmr   | D|06:16:58.760923|0xEC|hmny|i2h|tx_sender.go:230 handleRelayMessage: tx sent txh=0x51f3a6286b02a0afb7bb9f2f44fc49892a68959c9258fb4a4fa7cc85db6ee937 msg=0xf8e8f8e6b8e4f8e200b8dcf8daf8d8b8406274703a2f2f307836333537643265302e686d6e792f30783761364446326132434336374233384535326432333430424632424443376339613332416145393101b893f891b83e6274703a2f2f30783562396137372e69636f6e2f637862616636633230393137383832306436393639333136656135623164643466336139316334363361b8406274703a2f2f307836333537643265302e686d6e792f30783761364446326132434336374233384535326432333430424632424443376339613332416145393183626d630089c884496e697482c1c082078f
bmr   | D|06:16:58.779202|hx44|relay|h2i|relay.go:100 relaySignal
bmr   | D|06:16:58.789691|hx44|icon|h2i|tx_sender.go:239 handleRelayMessage: send tx prev=btp://0x6357d2e0.hmny/0x7a6DF2a2CC67B38E52d2340BF2BDC7c9a32AaE91
bmr   | D|06:16:58.799490|hx44|icon|h2i|tx_sender.go:258 handleRelayMessage: tx sent msg=0xf8e6f8e4b8e2f8e000b8daf8d8f8d6b83e6274703a2f2f30783562396137372e69636f6e2f63786261663663323039313738383230643639363933313665613562316464346633613931633436336101b893f891b8406274703a2f2f307836333537643265302e686d6e792f307837613644463261324343363742333845353264323334304246324244433763396133324161453931b83e6274703a2f2f30783562396137372e69636f6e2f63786261663663323039313738383230643639363933313665613562316464346633613931633436336183626d630089c884496e697482c1c082078a txh=0x9fb05de0c09b196beb7431e34ad6407e932a9129fdcd1e2d439b06b98c37b04e
bmr   | D|06:17:00.262092|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2164
bmr   | D|06:17:00.830769|hx44|hmny|h2i|rx_receiver.go:369 block notification height=2107
bmr   | D|06:17:01.769243|0xEC|hmny|i2h|tx_sender.go:285 handleRelayMessage: success txh=0x51f3a6286b02a0afb7bb9f2f44fc49892a68959c9258fb4a4fa7cc85db6ee937
bmr   | D|06:17:02.256921|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2165
bmr   | D|06:17:02.805919|hx44|icon|h2i|tx_sender.go:307 handleRelayMessage: success txh=0x9fb05de0c09b196beb7431e34ad6407e932a9129fdcd1e2d439b06b98c37b04e
bmr   | D|06:17:02.811400|hx44|hmny|h2i|rx_receiver.go:369 block notification height=2108
...
bmr   | D|06:19:28.308031|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2238
bmr   | D|06:19:28.771589|0xEC|relay|i2h|relay.go:100 relaySignal
bmr   | D|06:19:28.774908|0xEC|hmny|i2h|tx_sender.go:211 handleRelayMessage: send tx prev=btp://0x5b9a77.icon/cxbaf6c209178820d6969316ea5b1dd4f3a91c463a
bmr   | D|06:19:28.780688|0xEC|hmny|i2h|tx_sender.go:230 handleRelayMessage: tx sent txh=0xc4c3ef7b164f0061ff26d0f6ad47ca3390574a5c7b2661b504a293380b422976 msg=0xf90159f90156b90153f9015000b90149f90146f90143b8406274703a2f2f307836333537643265302e686d6e792f30783761364446326132434336374233384535326432333430424632424443376339613332416145393102b8fef8fcb83e6274703a2f2f30783562396137372e69636f6e2f637862616636633230393137383832306436393639333136656135623164643466336139316334363361b8406274703a2f2f307836333537643265302e686d6e792f3078376136444632613243433637423338453532643233343042463242444337633961333241614539318a6e6174697665636f696e01b86cf86a00b867f865aa687836393165616438386264353934356134336338613164613333316666366464383065323933366565aa307838666336363832373562346661303332333432656133303339363533643834316630363961383362cecd83494358881b7a5f826f4600008208bc
bmr   | D|06:19:28.794036|hx44|relay|h2i|relay.go:100 relaySignal
bmr   | D|06:19:29.830622|hx44|hmny|h2i|rx_receiver.go:369 block notification height=2182
bmr   | D|06:19:30.308766|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2239
bmr   | D|06:19:31.789246|0xEC|hmny|i2h|tx_sender.go:285 handleRelayMessage: success txh=0xc4c3ef7b164f0061ff26d0f6ad47ca3390574a5c7b2661b504a293380b422976
bmr   | D|06:19:31.831117|hx44|hmny|h2i|rx_receiver.go:369 block notification height=2183
bmr   | D|06:19:32.311219|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2240
bmr   | D|06:19:33.772493|0xEC|relay|i2h|relay.go:100 relaySignal
bmr   | D|06:19:33.794611|hx44|relay|h2i|relay.go:100 relaySignal
bmr   | D|06:19:33.832965|hx44|hmny|h2i|rx_receiver.go:369 block notification height=2184
bmr   | D|06:19:34.310398|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2241
bmr   | D|06:19:35.831364|hx44|hmny|h2i|rx_receiver.go:369 block notification height=2185
bmr   | I|06:19:35.831503|hx44|relay|h2i|relay.go:128 srcMsg added seq=[2 2]
bmr   | D|06:19:36.310873|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2242
bmr   | D|06:19:37.830611|hx44|hmny|h2i|rx_receiver.go:369 block notification height=2186
bmr   | D|06:19:38.313039|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2243
bmr   | D|06:19:38.772564|0xEC|relay|i2h|relay.go:100 relaySignal
bmr   | D|06:19:38.795434|hx44|relay|h2i|relay.go:100 relaySignal
bmr   | D|06:19:38.802847|hx44|icon|h2i|tx_sender.go:239 handleRelayMessage: send tx prev=btp://0x6357d2e0.hmny/0x7a6DF2a2CC67B38E52d2340BF2BDC7c9a32AaE91
bmr   | D|06:19:38.804125|hx44|icon|h2i|tx_sender.go:258 handleRelayMessage: tx sent txh=0xd409884e72a74da714aef0bce98308f9cf61df7f535d334afa9789b481f75dbc msg=0xf8eaf8e8b8e6f8e400b8def8dcf8dab83e6274703a2f2f30783562396137372e69636f6e2f63786261663663323039313738383230643639363933313665613562316464346633613931633436336102b897f895b8406274703a2f2f307836333537643265302e686d6e792f307837613644463261324343363742333845353264323334304246324244433763396133324161453931b83e6274703a2f2f30783562396137372e69636f6e2f6378626166366332303931373838323064363936393331366561356231646434663361393163343633618a6e6174697665636f696e0186c50283c20080820889
bmr   | D|06:19:40.316202|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2244
bmr   | D|06:19:41.831356|hx44|hmny|h2i|rx_receiver.go:369 block notification height=2187
bmr   | D|06:19:42.313364|0xEC|icon|i2h|rx_receiver.go:241 block notification height=2245
bmr   | D|06:19:42.809590|hx44|icon|h2i|tx_sender.go:307 handleRelayMessage: success txh=0xd409884e72a74da714aef0bce98308f9cf61df7f535d334afa9789b481f75dbc
bmr   | D|06:19:42.813254|hx44|hmny|h2i|rx_receiver.go:369 block notification height=2188
...
```

## Run Demo

`./ixh.sh demo > _ixh/demo.log`

After both Relayers have started and forwarded the first BTP messages to other chains, run the demo in separate shell using above command. It will dump all the transaction details into the `_ixh/demo.log` file for debugging purpose. And it will print the transfer stats on the console via stderr.

Here is the sample output:

```
ICON:
    NativeCoins: ["ICX","ONE"]
    IRC2 Tokens: ["ETH"]
HMNY:
    NativeCoins: ["ONE","ICX"]
    ERC20 Tokens: ["ETH"]

Funding demo wallets...
    ICON (hx691ead88bd5945a43c8a1da331ff6dd80e2936ee): 250.00 ICX, 10.00 ETH
    HMNY (0x8fc668275b4fa032342ea3039653d841f069a83b): 10.00 ONE, 10.00 ETH


Balance:
    ICON: hx691ead88bd5945a43c8a1da331ff6dd80e2936ee
        ICX: 250.00
        ONE (Wrapped): 0
        ETH (IRC2): 10.00
    HMNY: 0x8fc668275b4fa032342ea3039653d841f069a83b
        ONE: 10.00
        ICX (Wrapped): 0
        ETH (ERC20): 10.00

Transfer Native ICX (ICON -> HMNY):
    amount=2.00


Balance:
    ICON: hx691ead88bd5945a43c8a1da331ff6dd80e2936ee
        ICX: 247.99
        ONE (Wrapped): 0
        ETH (IRC2): 10.00
    HMNY: 0x8fc668275b4fa032342ea3039653d841f069a83b
        ONE: 10.00
        ICX (Wrapped): 1.98
        ETH (ERC20): 10.00

Transfer Native ONE (HMNY -> ICON):
    amount=2.00


Balance:
    ICON: hx691ead88bd5945a43c8a1da331ff6dd80e2936ee
        ICX: 247.99
        ONE (Wrapped): 1.98
        ETH (IRC2): 10.00
    HMNY: 0x8fc668275b4fa032342ea3039653d841f069a83b
        ONE: 7.98
        ICX (Wrapped): 1.98
        ETH (ERC20): 10.00

Approve ICON NativeCoinBSH to access ONE
    Allowance: 100000.00
Approve HMNY BSHCore to access ICX
    Allowance: 100000.00

Transfer Wrapped ICX (HMNY -> ICON):
    amount=1.00

Balance:
    ICON: hx691ead88bd5945a43c8a1da331ff6dd80e2936ee
        ICX: 248.98
        ONE (Wrapped): 1.98
        ETH (IRC2): 10.00
    HMNY: 0x8fc668275b4fa032342ea3039653d841f069a83b
        ONE: 7.96
        ICX (Wrapped): .98
        ETH (ERC20): 10.00

Transfer Wrapped ONE (ICON -> HMNY):
    amount=1.00


Balance:
    ICON: hx691ead88bd5945a43c8a1da331ff6dd80e2936ee
        ICX: 248.98
        ONE (Wrapped): .98
        ETH (IRC2): 10.00
    HMNY: 0x8fc668275b4fa032342ea3039653d841f069a83b
        ONE: 8.95
        ICX (Wrapped): .98
        ETH (ERC20): 10.00

Transfer irc2.ETH (ICON -> HMNY):
    amount=1.00


Balance:
    ICON: hx691ead88bd5945a43c8a1da331ff6dd80e2936ee
        ICX: 248.98
        ONE (Wrapped): .98
        ETH (IRC2): 9.00
    HMNY: 0x8fc668275b4fa032342ea3039653d841f069a83b
        ONE: 8.95
        ICX (Wrapped): .98
        ETH (ERC20): 10.99

Transfer erc20.ETH (HMNY -> ICON):
    amount=1.00

Balance:
    ICON: hx691ead88bd5945a43c8a1da331ff6dd80e2936ee
        ICX: 248.98
        ONE (Wrapped): .98
        ETH (IRC2): 9.99
    HMNY: 0x8fc668275b4fa032342ea3039653d841f069a83b
        ONE: 8.93
        ICX (Wrapped): .98
        ETH (ERC20): 9.99
```

## Disclaimer:

This is not a comprehensive description and is not complete. There are other requirements that are necessary to run the above demo. This document will be updated if somebody reports that they're unable to run the demo using above command.
