# E2E Tests

- [Design Document](https://github.com/icon-project/icon-bridge/discussions/141)
## Introduction

These tests are designed to ensure that the flow of data between chains completes as expected. A separate module named e2etest provides the interface necessary to write the script tailored for this task. 

It mainly covers Transfer and Configuration APIs. Transfer API involves methods like GetCoinBalance, Approve, Transfer, TransferBatch while Configuration AP involves methods like SetTokenLimit, AddBlackListAddress, RemoveBlackListAddress, SetFeeGatheringTerm, etc. These and other APIs enable a tester to write a script that executes a set of operations, wait for inter-chain events and assert that the result is as expected.

## Steps to run e2e tests
### Prerequisites
* #### Relay: 
    To run an e2e test, the deployment step should have been completed and the relay should have initialized i.e. caught up to the latest blocks.
* #### Config: 
    During the deployment, a file name e2e.config.json is auto-generated that includes the necessary configuration to run the tests. The file is generated on path `icon-bridge/devnet/docker/icon-bsc/_ixh/e2e.config.json`
* #### GodWallet: 
    To execute transfers, we require a wallet with sufficient funds and it does not need to have any additional privilege for the deployed smart contracts. The keystore file is present on path `icon-bridge/devnet/docker/icon-bsc/_ixh/chain.god.wallet.json` for god-account of each chain. 
* #### Smart Contract Owner Wallet: 
    To execute configuration api calls, we need a privileged user with sufficient funds. Since the deployer account is a privileged user of all smart contracts, funding it and then specifying it as the contract owner will enable calling configuration APIs for all deployed smart contracts. The keystore file for the deployer account is the god wallet path mentioned above.
    
Alternatively, you can fund bts and bmc smart contract owner wallets separately. Doing so is necessary especially if the god wallet you’ve used is a non-privileged user.

Path to relevant keystore files are present in e2e.config.json. By default, the config file uses privileged deployer accounts for both transfer and configuration calls to avoid having to fund multiple wallets and ease the testing procedure. If you’d rather use separate wallets, please fund relevant ones and specify their path in the e2e.config.json file.

### Commands
After you’ve finished deployment procedure and have run relay, run the following commands to start e2etests
```sh
cd icon-bridge/devnet/docker/icon-bsc
make rune2etests
```

> Note: By default e2e.config.json uses the deployed smart contract address and deployer account (chain.god.wallet.json) as an owner of all smart contracts

## Tests and Scripts
_Following are the scripts that are run by default:_
* TransferUniDirection
    - Transfer a coin from source chain to destination chain
* TransferBiDirection,
    - Transfer a coin from source chain to destination chain
    - Retrieve the same coin from destination chain to source chain
* TransferBatchBiDirection,
	- Transfer batch of coins from source chain to destination chain
	- Retrieve the same batch of coins from destination chain to source chain
* TransferFromBlackListedSrcAddress
	- Transfer a coin from source to destination chain
	- Blacklist an address on source chain
	- Try transferring again and notice failure due to blacklisting
	- Remove from blacklisted address
	- Try transferring angin and notice success due to removal from blacklist
* TransferToBlackListedDstAddress,
	- Transfer a coin from source to destination chain
	- Blacklist an address on destination chain
	- Try transferring again and notice failure due to blacklisting
	- Remove from blacklisted address
	- Try transferring angin and notice success due to removal from blacklist
* TransferEqualToFee,
	- Transfer a coin such that the fee charged is equal to user provided amount and net transferable amount is zero
* TransferLessThanFee
	- Transfer a coin such that the fee charged is equal to user provided amount and net transferable amount is negative
* TransferToZeroAddress
	- Transfer to zero address on destination chain. 
* TransferToUnknownNetwork,
	- Transfer to an unsupported block-chain network.
* TransferWithoutApprove
	- Transfer token without Approve step

Each of these scripts runs checks to ensure that the change in amount held in different accounts all adds up. This includes
- checking for the correct fee charged
- considering native coin used for gas fees 
- checking the correct amount transferred and received
- checking the correct amount increase or decrease in wallets

## Adding custom scripts
Please read the [design document](https://github.com/icon-project/icon-bridge/discussions/141) to get an understanding on how the e2etest module is designed. You may find default scripts on path `iconbridge/cmd/e2etest/executor/scripts*.go`
Please go through them to get an idea on how the scripts can be written.

