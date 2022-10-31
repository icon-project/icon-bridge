# ICON Bridge
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-3-orange.svg?style=flat-square)](#contributors-)
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
      <td align="center"><a href="https://github.com/izyak"><img src="https://avatars.githubusercontent.com/u/76203436?v=4?s=100" width="100px;" alt="izyak"/><br /><sub><b>izyak</b></sub></a><br /><a href="https://github.com/icon-project/icon-bridge/commits?author=izyak" title="Tests">⚠️</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=izyak" title="Documentation">📖</a></td>
      <td align="center"><a href="https://github.com/andrii-kl"><img src="https://avatars.githubusercontent.com/u/18900364?v=4?s=100" width="100px;" alt="Andrii"/><br /><sub><b>Andrii</b></sub></a><br /><a href="https://github.com/icon-project/icon-bridge/commits?author=andrii-kl" title="Documentation">📖</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=andrii-kl" title="Tests">⚠️</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=andrii-kl" title="Code">💻</a></td>
      <td align="center"><a href="https://github.com/manishbista28"><img src="https://avatars.githubusercontent.com/u/66529584?v=4?s=100" width="100px;" alt="manishbista28"/><br /><sub><b>manishbista28</b></sub></a><br /><a href="https://github.com/icon-project/icon-bridge/commits?author=manishbista28" title="Documentation">📖</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=manishbista28" title="Tests">⚠️</a> <a href="https://github.com/icon-project/icon-bridge/commits?author=manishbista28" title="Code">💻</a></td>
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
