<p align="center">
    <img src="./docs/img/iconbridge_350x80.png" alt="ICON Bridge Logo" width="350px"></img>
</p>
<div align="center">
    <a align="center" href='https://icon.community/learn/icon-bridge/'><button type='button' style='font-weight:semibold; background:#30AAAE; border-radius:5px; border:0px; box-shadow:1px; padding:4px 6px; color:white; cursor:pointer; '>What is ICON Bridge?</button></a>
</div>
<h3 align="center">
    ICON Bridge is an early iteration of ICON's cutting-edge interoperability product, BTP, which allows cross-chain transfers and integration with any blockchain that suppots smart contracts.    
</h3>
<p align="center">
    <a href="https://twitter.com/iconfoundation_"><img src="https://img.shields.io/twitter/follow/iconfoundation_?style=social"></a>
    &nbsp;
   <a href="https://twitter.com/helloiconworld"><img src="https://img.shields.io/twitter/follow/helloiconworld?style=social"></a>
</p>

## ICON Bridge
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-4-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->
![release](https://img.shields.io/github/v/release/icon-project/icon-bridge)
[![codecov](https://codecov.io/gh/icon-project/icon-bridge/branch/main/graph/badge.svg?token=YXV6EE5KB5)](https://codecov.io/gh/icon-project/icon-bridge)

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![OpenSSF
Scorecard](https://api.securityscorecards.dev/projects/github.com/icon-project/icon-bridge/badge)](https://api.securityscorecards.dev/projects/github.com/icon-project/icon-bridge)


This repository contains the smart contracts source code and relay source code for ICON bridge. 
## Project Overview

ICON Bridge is a centralized bridge for Blockchain Transmission Protocol(BTP) Relay System which can be used to transfer tokens across multiple chains. Currently, it supports cross chain transfer from ICON and Binance Smart Chain (BSC).

The main components of icon bridge are:
* ### BTP Message Relay (BMR)
    - It serves to relay BTP Message across connected chains and monitor BTP events
* ### Contracts
    * #### BTP Message Center (BMC)
        - Receive BTP messages through transactions.
        - Send BTP messages through events.


    * #### BTP Service Handler (BSH)
        - Services that can be serviced by ICON-Bridge
        - BTP Token Service (BTS) is a BSH that is responsible for token transfers cross chain.
        - Currently, BTS is the only service handler for icon bridge
        - Handle service messages related to the service.
        - Send service messages through the BMC


## Getting Started
[Terminologies](./docs/terminologies.md) used in ICON Bridge.

Getting started section can be found [here](./docs/getting-started.md). It contains information about folder structure of the repo, how to build ICON Bridge on local or testnet/mainnet and how to run the tests from scratch.

If you want to contribute to this repository, read the [Contributor Guidelines](CONTRIBUTING.md) for more info. 

The documentation for this project is in the [docs](./docs/) directory.

For the latest mainnet contract addresses, please check [Mainnet Contract Addresses](./docs/mainnet_deployment.json)

For the testnet contract addresses, please check [Testnet Contract Addresses](./docs/testnet_deployment.json)


## Roadmap

Please see our quarterly roadmap [here](https://github.com/orgs/icon-project/projects/4).

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center"><a href="https://github.com/izyak"><img src="https://avatars.githubusercontent.com/u/76203436?v=4?s=100" width="100px;" alt="izyak"/><br /><sub><b>izyak</b></sub></a><br /><a href="https://github.com/icon-project/icon-bridge/commits?author=izyak" title="Tests">⚠️</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=izyak" title="Documentation">📖</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=izyak" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/andrii-kl"><img src="https://avatars.githubusercontent.com/u/18900364?v=4?s=100" width="100px;" alt="Andrii"/><br /><sub><b>Andrii</b></sub></a><br /><a href="https://github.com/icon-project/icon-bridge/commits?author=andrii-kl" title="Documentation">📖</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=andrii-kl" title="Tests">⚠️</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=andrii-kl" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/CyrusVorwald-ICON"><img src="https://avatars.githubusercontent.com/u/90732384?v=4?s=100" width="100px;" alt="CyrusVorwald-ICON"/><br /><sub><b>CyrusVorwald-ICON</b></sub></a><br /><a href="#business-CyrusVorwald-ICON" title="Business development">💼</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=CyrusVorwald-ICON" title="Documentation">📖</a> <a href="#ideas-CyrusVorwald-ICON" title="Ideas, Planning, & Feedback">🤔</a> <a href="#maintenance-CyrusVorwald-ICON" title="Maintenance">🚧</a> <a href="#projectManagement-CyrusVorwald-ICON" title="Project Management">📆</a></td>
      <td align="center"><a href="http://themissingsemicolon.in"><img src="https://avatars.githubusercontent.com/u/35568964?v=4?s=100" width="100px;" alt="Shreyas S Bhat"/><br /><sub><b>Shreyas S Bhat</b></sub></a><br /><a href="https://github.com/icon-project/icon-bridge/commits?author=shreyasbhat0" title="Documentation">📖</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=shreyasbhat0" title="Tests">⚠️</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=shreyasbhat0" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://allcontributors.org) specification.
Contributions of any kind are welcome!
