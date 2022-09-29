# BMR Test Flow
## Introduction
A complete Github Actions Workflow for running Relay.

## Starting the Flow
The main workflow Github Actions file is [deploy-testnet-workflow.yml](https://github.com/icon-project/icon-bridge/blob/main/.github/workflows/deploy-testnet-workflow.yml) which is set to trigger either manually or runs on schedule. 

## Workflow
The workflow triggers multiple other workflows:
1. [deploy-javascore](https://github.com/icon-project/icon-bridge/blob/main/.github/workflows/deploy-javascore-testnet.yml): This workflow involves jobs following tasks to Icon Testnet:
	1. Deploy BMC Javascore
	2. Deploy BTS Javascore
	3. Register BTS to BMC
	4. Sync BMC and BTS address to S3 bucket
2. [deploy-solidity](https://github.com/icon-project/icon-bridge/blob/main/.github/workflows/deploy-solidity-testnet.yml): This workflow runs following jobs on BSC Testnet:
	1. Deploy BMC Solidity
	2. Deploy BTS Solidity
	3. Register BTS to BMC
3. [deploy-token-javascore](https://github.com/icon-project/icon-bridge/blob/main/.github/workflows/deploy-token-javascore.yml): This deploy token job involves:
	1. Deploy irc2 token from [icon-project/java-score-examples repo]
(https://github.com/icon-project/java-score-examples)
	2. Deploy Javascore token
	3. Add BMC owner
	4. Set BMC fee Aggregator
	5. Javascore add BTS and owner
	6. Set BTS ICX Fee
4. [link-register-token](https://github.com/icon-project/icon-bridge/blob/main/.github/workflows/link-register-coin.yml):
	1. Javascore register native(sICX) and wrapped token (BNB)
	2. Javascore add link ICON-BSC and Set link height(ICON)
	3. BMC add relay
	4. Get BTP coin ID: native and wrapped
	5. Solidity add link BSC-ICON
	6. Set link height(BSC)
	7. Add icon relay
	8. Solidity add BMC and BTS owner
	9. Solidity set fee ratio
	10. BSC register native and wrapped token
	11. BSC get coin ID(BUSD and ICX)
	12. Sync address files to S3
5. [deploy-bmr](https://github.com/icon-project/icon-bridge/blob/main/.github/workflows/deploy-bmr.yml):
	1. Build BMR docker image
	2. Generate config file for BMR deployment
	3. Start BMR using docker-compose file
