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

## Build

    ./ixh.sh build

## Publish images

    ./ixh.sh publish

## Deploy Blockchains

    ./ixh.sh start -d

NOTE:

If you're using a remote docker host, please allow ssh connection either password less or with key added to local ssh-agent. It needs the ablility to ssh into the remote host non-interactively. You also need a local docker registry on the remote host running at 0.0.0.0:5000, to publish necessary docker images of respective blockchains, so that it could be used on that host. If hostname is used instead of an IP address, you need to map the hostname to IP in `/etc/hosts` in a linux system, corresponding files in windows/macos system.

## Deploy Smart Contracts

    ./ixh.sh deploysc reset

It deploys necessary smart contracts and generates configurations for Relayers (`_ixh/i2h.config.json` and `_ixh/h2i.config.json`) along with an environment file (`_ixh/ixh.env`) that has all necessary environment variables.

NOTE: _Wait for 1 minute or more after the first step to do this._

## Start Relayers

Open two separate command prompts, and go to the `cmd/btpsimple` directory on each. And run following commands:

```

# First shell: harmony -> icon

cd cmd/btpsimple
go run . -c ../../devnet/docker/icon-hmny/\_ixh/h2i.config.json start

...
I|12:58:50.253553|----|-|main|main.go:228   ____ _____ ____    ____      _
I|12:58:50.253617|----|-|main|main.go:228  | __ )_   _|  _ \  |  _ \ ___| | __ _ _   _
I|12:58:50.253624|----|-|main|main.go:228  |  _ \ | | | |_) | | |_) / _ \ |/ _` | | | |
I|12:58:50.253629|----|-|main|main.go:228  | |_) || | |  __/  |  _ <  __/ | (_| | |_| |
I|12:58:50.253635|----|-|main|main.go:228  |____/ |_| |_|     |_| \_\___|_|\__,_|\__, |
I|12:58:50.253640|----|-|main|main.go:228                                        |___/
I|12:58:50.253646|----|-|main|main.go:230 Version : unknown
I|12:58:50.253652|----|-|main|main.go:231 Build   : unknown
D|12:58:50.378114|a8e2|-|main|main.go:316 LogForwarderConfig vendor and address is empty string, will be ignore
D|12:58:50.378203|a8e2|-|main|main.go:242 /home/bbist/works/ibriz/code/github.com/icon-project/icon-bridge/devnet/docker/icon-hmny/_ixh/h2i.config.json run/h2i
D|12:58:51.079655|a8e2|0x7|chain|chain.go:321 _init height:0, dst(btp://0x7.icon/cxa2cc386e9db2a72ea6724cbfd12f936a90ba63d2, seq:10), receive:8033
D|12:58:51.079807|a8e2|0x7|chain|chain.go:306 start relayLoop
D|12:58:51.906155|a8e2|-|icon|client.go:221 MonitorBlock WSEvent 192.168.207.251:35580 WSEventInit
D|12:58:51.906322|a8e2|-|icon|sender.go:244 MonitorLoop connected 192.168.207.251:35580
D|12:58:51.906371|a8e2|0x7|chain|chain.go:348 Connect MonitorLoop
D|12:58:55.118300|a8e2|-|hmny|receiver.go:76 receive loop: block notification: height=8033
D|23:55:44.337342|a6a4|-|hmny|receiver.go:67 found event in block 8033: sc=0xDd334a2E6DAfa23e45B2FF8035C7DeE87F29375B
D|23:55:44.337387|a6a4|0x5b9a77|chain|chain.go:187 addRelayMessage rms:1 rps:1 HeightOfDst:0
D|23:55:44.337523|a6a4|-|icon|sender.go:77 HandleRelayMessage prev btp://0x2.hmny/0xDd334a2E6DAfa23e45B2FF8035C7DeE87F29375B, msg: -N_43bjb-NkAuNP40fjPuD5idHA6Ly8weDViOWE3Ny5pY29uL2N4ODhjMzBmOWM4ZmEzYTczZWE5NWU4OTQ2ZDEyM2ViMDk1NzNiODcxOAG4jPiKuDlidHA6Ly8weDIuaG1ueS8weERkMzM0YTJFNkRBZmEyM2U0NUIyRkY4MDM1QzdEZUU4N0YyOTM3NUK4PmJ0cDovLzB4NWI5YTc3Lmljb24vY3g4OGMzMGY5YzhmYTNhNzNlYTk1ZTg5NDZkMTIzZWIwOTU3M2I4NzE4g2JtYwCJyIRJbml0gsHAgh9h
D|23:55:44.337634|a6a4|0x5b9a77|chain|chain.go:82 Going to relay now rm:0 [i:0,h:0,seq:0,evt:0,txh:<nil>]
D|23:55:44.339240|a6a4|-|hmny|receiver.go:77 receive loop: block notification: height=8035
D|23:55:44.341086|a6a4|-|hmny|receiver.go:77 receive loop: block notification: height=8036
...

```

```

# Second shell: icon -> harmony

cd cmd/btpsimple
go run . -c ../../devnet/docker/icon-hmny/\_ixh/i2h.config.json start

...
I|12:58:31.276334|----|-|main|main.go:228   ____ _____ ____    ____      _
I|12:58:31.276385|----|-|main|main.go:228  | __ )_   _|  _ \  |  _ \ ___| | __ _ _   _
I|12:58:31.276392|----|-|main|main.go:228  |  _ \ | | | |_) | | |_) / _ \ |/ _` | | | |
I|12:58:31.276397|----|-|main|main.go:228  | |_) || | |  __/  |  _ <  __/ | (_| | |_| |
I|12:58:31.276403|----|-|main|main.go:228  |____/ |_| |_|     |_| \_\___|_|\__,_|\__, |
I|12:58:31.276408|----|-|main|main.go:228                                        |___/
I|12:58:31.276413|----|-|main|main.go:230 Version : unknown
I|12:58:31.276420|----|-|main|main.go:231 Build   : unknown
D|12:58:31.810368|dA9E|-|main|main.go:316 LogForwarderConfig vendor and address is empty string, will be ignore
D|12:58:31.810531|dA9E|-|main|main.go:242 /home/bbist/works/ibriz/code/github.com/icon-project/icon-bridge/devnet/docker/icon-hmny/_ixh/i2h.config.json run/i2h
D|12:58:34.246307|dA9E|0x6357d2e0|chain|chain.go:321 _init height:0, dst(btp://0x6357d2e0.hmny/0xeA4039A61C7de7057428F8512CddBB9BDb519278, seq:10), receive:13146
D|12:58:34.246453|dA9E|0x6357d2e0|chain|chain.go:306 start relayLoop
D|12:58:35.776128|dA9E|-|icon|client.go:221 MonitorBlock WSEvent 192.168.207.251:34106 WSEventInit
D|12:58:35.776245|dA9E|-|icon|receiver.go:202 ReceiveLoop connected 192.168.207.251:34106
D|12:58:35.776259|dA9E|0x6357d2e0|chain|chain.go:362 Connect ReceiveLoop
D|12:58:35.953389|dA9E|-|icon|receiver.go:191 onBlockOfSrc icon: 13146
D|23:55:45.713797|531d|-|hmny|sender.go:121 final relay message string: ���ظ��������ʸ9btp://0x2.hmny/0xDd334a2E6DAfa23e45B2FF8035C7DeE87F29375B�����>btp://0x5b9a77.icon/cx88c30f9c8fa3a73ea95e8946d123eb09573b8718�9btp://0x2.hmny/0xDd334a2E6DAfa23e45B2FF8035C7DeE87F29375B�bmc�ȄInit����2#, prev: btp://0x5b9a77.icon/cx88c30f9c8fa3a73ea95e8946d123eb09573b8718
D|23:55:45.894896|531d|-|hmny|sender.go:182 monitor loop: block notification: height=8202
D|23:55:46.208687|531d|0x2|chain|chain.go:82 after relay rm:0 [i:0,h:0,seq:0,evt:0,txh:&{0x279a807a07dcdfcf19e858c4c99b42820babf04df5dcfac80705e9ea84331024}]
D|23:55:46.446369|531d|-|icon|receiver.go:191 onBlockOfSrc icon: 13147
...

```

## Run Demo

`./ixh.sh demo > _ixh/demo.log`

After both Relayers have started and forwarded the first BTP messages to other chains, run the demo in separate shell using above command. It will dump all the transaction details into the `_ixh/demo.log` file for debugging purpose. And it will print the transfer stats on the console via stderr.

Here is the sample output:

```

    Icon Wrapped Coins:
        ["ICX","ONE_DEV"]
    Hmny Wrapped Coins:
        ["ONE_DEV","ICX"]

    Balance:
        Icon: hxff0ea998b84ab9955157ab27915a9dc1805edd35
            Native: 162698527232959994832983024666
            Wrapped (ONE_DEV): 0
        Hmny: 0xa5241513da9f4463f1d4874b548dfbac29d91f34
            Native: 8063401632391062000000000000
            Wrapped (ICX): 0

    TransferNativeCoin (Icon -> Hmny):
        amount=54232842410986664944327674888
        [tx=0x76aa8c67bfb4a8284caa10c2698bf9ee23f2145228ddefcada45d287fbb33657].. ✔ 3s

    Balance:
        Icon: hxff0ea998b84ab9955157ab27915a9dc1805edd35
            Native: 108465684821973329888563567278
            Wrapped (ONE_DEV): 0
        Hmny: 0xa5241513da9f4463f1d4874b548dfbac29d91f34
            Native: 8063401632391062000000000000
            Wrapped (ICX): 54178609568575678279383347214

    TransferNativeCoin (Hmny -> Icon):
        amount=2687800544130354000000000000


    Balance:
        Icon: hxff0ea998b84ab9955157ab27915a9dc1805edd35
            Native: 108465684821973329888563567278
            Wrapped (ONE_DEV): 2685112743586223645999500000
        Hmny: 0xa5241513da9f4463f1d4874b548dfbac29d91f34
            Native: 5375601088250117000000000000
            Wrapped (ICX): 54178609568575678279383347214

    Approve Icon NativeCoinBSH
        [tx=0x686fc59fea6ec38fad7e0bd120e8911b603ba40bacda098d51556c16e18126e0]... ✔ 5s
        Status: "0x1"
    Approve Hmny BSHCore
        Status: true

    TransferWrappedCoin ICX (Hmny -> Icon):
        amount=27089304784287839139691673607

    Balance:
        Icon: hxff0ea998b84ab9955157ab27915a9dc1805edd35
            Native: 135527900301476881189099761812
            Wrapped (ONE_DEV): 2685112743586223645999500000
        Hmny: 0xa5241513da9f4463f1d4874b548dfbac29d91f34
            Native: 5375601088237940000000000000
            Wrapped (ICX): 27089304784287839139691673607

    TransferWrapped Coin ONE_DEV (Icon -> Hmny):
        amount=1342556371793111822999750000
        [tx=0x5dd69708f16339cc261ba24d1f234819315092f07e0238768e658422d099ebfc]... ✔ 5s

    Balance:
        Icon: hxff0ea998b84ab9955157ab27915a9dc1805edd35
            Native: 135527900301476881189041244312
            Wrapped (ONE_DEV): 1342556371793111822999750000
        Hmny: 0xa5241513da9f4463f1d4874b548dfbac29d91f34
            Native: 6716814903659260000000000000
            Wrapped (ICX): 27089304784287839139691673607

```

## Disclaimer:

This is not a comprehensive description and is not complete. There are other requirements that are necessary to run the above demo. This document will be updated if somebody reports that they're unable to run the demo using above command.
