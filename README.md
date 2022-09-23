# Icon Bridge
[![codecov](https://codecov.io/gh/icon-project/icon-bridge/branch/main/graph/badge.svg?token=YXV6EE5KB5)](https://codecov.io/gh/icon-project/icon-bridge)
[![OpenSSF
Scorecard](https://api.securityscorecards.dev/projects/github.com/icon-project/icon-bridge/badge)](https://api.securityscorecards.dev/projects/github.com/icon-project/icon-bridge)
## Introduction

We need to build a usable centralized bridge for [BTP](doc/btp.md) Relay System which can deliver digital tokens between multiple chains.

Target chains

- ICON (goloop)
- Polkadot parachain
- Binance Smart Chain
- Harmony One
- NEAR Protocol

Terminologies

| Word            | Description                                                                                                                             |
| :-------------- | :-------------------------------------------------------------------------------------------------------------------------------------- |
| BTP             | Blockchain Transmission Protocol, [ICON BTP Standard](https://github.com/icon-project/IIPs/blob/master/IIPS/iip-25.md) defined by ICON. |
| BTP Message     | A verified message which is delivered by the relay                                                                                      |
| Service Message | A payload in a BTP message                                                                                                              |
| Relay Message   | A message including BTPMessages with proofs for that, and other block update messages.                                                  |
| NetworkAddress  | Network Type and Network ID <br/> _0x1.icon_ <br/> _0x1.icon_                                                                           |
| ContractAddress | Addressing contract in the network <br/> _btp://0x1.icon/cx87ed9048b594b95199f326fc76e76a9d33dd665b_                                    |

> BTP Standard

### Components

- [BTP Message Center(BMC)](doc/bmv.md) - smart contract

  - Receive BTP messages through transactions.
  - Send BTP messages through events.

- [BTP Service Handler(BSH)](doc/bsh.md) - smart contract

  - Handle service messages related to the service.
  - Send service messages through the BMC

- [BTP Message Relay(BMR)](doc/bmr.md) - external software
  - Monitor BTP events
  - Send BTP Relay Message

### Blockchain specifics

- [ICON](doc/icon.md)

## BTP Project

### Documents

- [Build Guide](doc/build.md)
- [Tutorial](doc/tutorial.md)
- [iconbridge command line](doc/iconbridge_cli.md)
- [Binance Smart Chain Guide](doc/bsc-guide.md)
- [Harmony Guide](doc/hmny-guide.md)

### Layout

| Directory                            | Description                                                                                       |
| :----------------------------------- | :------------------------------------------------------------------------------------------------ |
| /cmd                                 | Root of implement of BMR                                                                          |
| /cmd/iconbridge                       | Reference implement of BMR. only provide unidirectional relay. (golang)                           |
| /cmd/iconbridge/relay                 | Implement of common logic of BMR, uses chain package                                              |
| /cmd/iconbridge/chain                 | BMR module interface, common code, and chain specific packages                                    |
| /cmd/iconbridge/chain/`<blockchain>`  | Implement of BMR module (`Sender`,`Receiver`), `<blockchain>` is name of blockchain               |
| /common                              | Common code (golang)                                                                              |
| /doc                                 | Documents                                                                                         |
| /docker                              | Docker related resources                                                                          |
| /`<env>`                             | Root of implement of BTP smart contracts, `<env>` is name of smart contract execution environment |
| /`<env>`/bmc                         | Implement of BMC smart contract                                                                   |
| /`<env>`/lib                         | Library for execution environment                                                                 |
| /`<env>`/`<svc>`                     | Root of implement of BSH smart contract, `<svc>` is name of BTP service                           |
| /`<env>`/token_bsh                   | Reference implement of BSH smart contract for Interchain-Token transfer service                   |
| /`<env>`/token_bsh/sample/irc2_token | Implement of IRC-2.0 smart contract, example for support legacy smart contract                    |

### BMR Modules

| Directory                  | Description                       |
| :------------------------- | :-------------------------------- |
| /cmd/iconbridge/module/icon | BMR module for ICON blockchain    |
| /cmd/iconbridge/module/hmny | BMR module for Harmony blockchain |
