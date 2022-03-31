# `ixh.sh`

This is a helper tool to auto-deploy harmony and icon localnets on local/remote host using docker. It has commands to deploy BTP smartcontracts to respective blockchain (`deploysc`), and run tests (`test`) using several subcommands. It also includes a command to run a demo (`demo`).

Please go through the `ixh.sh` file and see other commands at the bottom of the script and their syntax.

## Requirements:

Ensure following tools are installed.

    gradle, jdk@11.x, sdkman, goloop, docker, truffle@5.3.0, node@15.12.0, ethkey

## Run a demo:

1. **Build images**

   `./ixh.sh build`

2. **Publish images**

   `./ixh.sh publish`

3. **Deploy blockchains**

   `./ixh.sh start -d`

   NOTE:

   If you're using a remote docker host, please allow ssh connection either password less or with key added to local ssh-agent. It needs the ablility to ssh into the remote host non-interactively. You also need a local docker registry on the remote host running at 0.0.0.0:5000, to publish necessary docker images of respective blockchains, so that it could be used on that host. Please add remote host to `remotehosts` variable in `ixh.sh` file. If its a hostname and not an IP, and you need to map the hostname to IP in `/etc/hosts` in a linux system, corresponding files in windows/macos system.

   If you're using a local docker host, you cannot use `localhost` as a hostname as the `deploysc` script uses a docker container to create wallets and transfer balances. It will not be able to resolve to your system's host, and will instead to docker container's internal host. So the balance transfer will fail.
   What you can do instead is to map `localdckr` to 127.0.0.1 in your `/etc/hosts` file and use it as a local hostname.

   Sample output:

   ```
   Wallet:
       icon: ✔
       hmny: [src/hmny.wallet.json]  ✔

   Build:
       javascores: ✔

   Deploy:
   icon
       bmc: [tx=0xe10b0835]... ✔ 4s cx88c30f9c8fa3a73ea95e8946d123eb09573b8718
       irc31: [tx=0x6bfeb227]... ✔ 5s cx549e2ba448845431ec0613eed14640e5584d177d
       nativecoin_bsh: [tx=0x9663bfa0]... ✔ 4s cx2b817726ddc4fa92fb14dd9c4a55ab46184c3b59
       btp: btp://0x5b9a77.icon/cx88c30f9c8fa3a73ea95e8946d123eb09573b8718
   hmny
       bmc:  ✔ 110s m=0x77549beEa2e2342e8a7E1689ed644547479Cc7FC, p=0xDd334a2E6DAfa23e45B2FF8035C7DeE87F29375B
       bsh:  ✔ 111s c=0xdc16d7140009A16FDa7AcEB5b928B7aC7Cc2829d, p=0x95671B83c1958204647954502D2fB3Eb0b210001
       btp: btp://0x2.hmny/0xDd334a2E6DAfa23e45B2FF8035C7DeE87F29375B

   Configuring:
   icon
       create_wallet: [_ixh/bmc.owner.json] [tx=0x10e9c7b9]... ✔ 4s
       bmc_add_owner: [hx523609cf] [tx=0x7da58c9b].. ✔ 3s
       bmc_link_hmny_bmc:
           addLink: [tx=0x4c699fac]... ✔ 4s
           setLinkRxHeight: [tx=0x717f4945]... ✔ 4s
           getLinkStatus: rxHeight=7928
       bmc_add_nativecoin_bsh: [tx=0x67e04d66].. ✔ 3s
       create_wallet: [_ixh/nativecoin.icon.owner.json] [tx=0x03477b23].. ✔ 3s
       nativecoin_bsh_add_owner: [hx6840cbc3] [tx=0x1e9139f0].. ✔ 3s
       nativecoin_bsh_register_irc31: [tx=0x21b2160b].. ✔ 3s
       create_wallet: [_ixh/bmr.icon.json] [tx=0xc79ea9d2]... ✔ 4s
       bmc_add_relay: [tx=0x6a731b19]... ✔ 4s
       irc31_add_owner: [cx2b817726] [tx=0x45223e18].. ✔ 3s
   hmny
       bmc_add_bsh:  ✔
       bmc_link_to_icon_bmc:  ✔ ✔
       create_wallet: [_ixh/bmr.hmny.json]  ✔
       bmc_add_relay:  ✔
       bsh_register_coin:  ✔

   deploysc completed in 420s.
   important variables have been written to ./_ixh/ixh.env
   ```

4. **Deploy smartcontracts**

   `./ixh.sh deploysc reset`

   It deploys necessary smartcontracts and generates configurations for Relayers (`_ixh/i2h.config.json` and `_ixh/h2i.config.json`) along with an environment file (`_ixh/ixh.env`) that has all necessary environment variables.

   NOTE: _Wait for 1 minute or more after the first step to do this._

5. **Start Relayers**

   Open two separate command prompts, and go to the `cmd/btpsimple` directory on each. And run following commands:

   ```
   # First shell: harmony -> icon
   cd cmd/btpsimple
   go run . -c ../../devnet/docker/icon-hmny/_ixh/h2i.config.json start

   ...
   D|23:55:44.335355|a6a4|-|hmny|receiver.go:77 receive loop: block notification: height=8033
   D|23:55:44.337270|a6a4|-|hmny|receiver.go:77 receive loop: block notification: height=8034
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
   go run . -c ../../devnet/docker/icon-hmny/_ixh/i2h.config.json start

   ...
   D|23:55:45.182028|531d|-|icon|receiver.go:191 onBlockOfSrc icon: 13146
   D|23:55:45.713797|531d|-|hmny|sender.go:121 final relay message string: ���ظ��������ʸ9btp://0x2.hmny/0xDd334a2E6DAfa23e45B2FF8035C7DeE87F29375B�����>btp://0x5b9a77.icon/cx88c30f9c8fa3a73ea95e8946d123eb09573b8718�9btp://0x2.hmny/0xDd334a2E6DAfa23e45B2FF8035C7DeE87F29375B�bmc�ȄInit����2#, prev: btp://0x5b9a77.icon/cx88c30f9c8fa3a73ea95e8946d123eb09573b8718
   D|23:55:45.894896|531d|-|hmny|sender.go:182 monitor loop: block notification: height=8202
   D|23:55:46.208687|531d|0x2|chain|chain.go:82 after relay rm:0 [i:0,h:0,seq:0,evt:0,txh:&{0x279a807a07dcdfcf19e858c4c99b42820babf04df5dcfac80705e9ea84331024}]
   D|23:55:46.446369|531d|-|icon|receiver.go:191 onBlockOfSrc icon: 13147
   ...
   ```

6. **Run demo**

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

### Disclaimer:

This is not a comprehensive description and is not complete. There are other requirements that are necessary to run the above demo. This document will be updated if somebody reports that they're unable to run the demo using above command.
