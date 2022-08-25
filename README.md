# ICON Bridge
<p>
  <img alt="GitHub release (latest SemVer including pre-releases)" src="https://img.shields.io/github/v/release/icon-project/icon-bridge?include_prereleases">
  <img src="https://goreportcard.com/badge/github.com/icon-project/icon-bridge">
  <a href="https://discord.gg/ZkgByPn92j">
    <img src="https://img.shields.io/discord/880651922682560582?label=discord" alt="Discord">
  </a>
  <a href="https://opensource.org/licenses/Apache-2.0">
    <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License">
  </a>
</p>
  
The ICON Bridge is an open‑source blockchain interoperability application. It is a fork from the [Blockchain Transmission Protocol](https://github.com/icon-project/btp) and differs in that message verification is done off-chain.

> **_NOTE:_**  The ICON Bridge is currently in beta and has not completed auditing. Use at your own risk.

# Getting Started

We are currently refactoring our process for getting started and will update this file when it is in a better state.

### Documents

- [iconbridge command line](doc/iconbridge_cli.md)
- TODO: add more docs


### Layout

| Directory                            | Description                                                                                       |
| :----------------------------------- | :------------------------------------------------------------------------------------------------ |
| /cmd                                 | Root of implement of the relay                                                                          |
| /cmd/iconbridge/relay                 | Implement of common logic of relay, uses chain package                                              |
| /cmd/iconbridge/chain                 | Relay module interface, common code, and chain specific packages                                    |
| /cmd/iconbridge/chain/`<blockchain>`  | Implement of relay module (`Sender`,`Receiver`), `<blockchain>` is name of blockchain               |
| /common                              | Common code (golang)                                                                              |
| /doc                                 | Documents                                                                                         |
| /docker                              | Docker related resources                                                                          |
| /`<env>`                             | Root of implement of ICON Bridge smart contracts, `<env>` is name of smart contract execution environment |
| /`<env>`/lib                         | Library for execution environment                                                                 |

### BMR Modules

| Directory                  | Description                       |
| :------------------------- | :-------------------------------- |
| /cmd/iconbridge/module/icon | Relay module for ICON blockchain    |
| /cmd/iconbridge/module/hmny | Relay module for Harmony blockchain |
| /cmd/iconbridge/module/hmny | Relay module for Harmony blockchain |
