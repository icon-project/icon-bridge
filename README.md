# Icon Bridge
![release](https://img.shields.io/github/v/release/icon-project/icon-bridge)
[![codecov](https://codecov.io/gh/icon-project/icon-bridge/branch/main/graph/badge.svg?token=YXV6EE5KB5)](https://codecov.io/gh/icon-project/icon-bridge)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![OpenSSF
Scorecard](https://api.securityscorecards.dev/projects/github.com/icon-project/icon-bridge/badge)](https://api.securityscorecards.dev/projects/github.com/icon-project/icon-bridge)

This repository contains the smart contracts source code and relay source code for Icon bridge. 
## Project Overview

Icon Bridge is a centralized bridge for Blockchain Transmission Protocol(BTP) Relay System which can be used to transfer tokens across multiple chains. Currently, it supports cross chain transfer from ICON and Binance Smart Chain (BSC).

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

If you are a developer, check this out: [Developer guidelines](./docs/developer-guidelines.md)

If you want to contribute to this repository, read the [Contributor Guidelines](CONTRIBUTING.md) for more info. 

The documentation for this project is in the [docs](./docs/) directory.

For the latest mainnet contract addresses, please check [Mainnet Contract Addresses](./docs/mainnet_deployment.json)

For the testnet contract addresses, please check [Testnet Contract Addresses](./docs/testnet_deployment.json)


## RoadMap

## Contributors
