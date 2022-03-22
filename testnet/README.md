
Javascore Deployment:
---------------------------------------------------------------------------------------------------------------------------------------------------------------

1. Copy the javascore jar files from the build folder after the `make build` to artifacts folder(Skip this step if there are no javascore code changes)
2. update the proper keystore file with funds 
3. refer to [javascore_deployment_commands](https://github.com/icon-project/btp/blob/btp_nobmv/testnet/javascore/javascore_deploy_commands.md) file for java deployment commands & configurations
4. update the keystore file path&name on `keystore` variable & update the `password` variable with the keystore password on the commands env setup step
5. paste the env export commands in the terminal
6. proceed with the deployment after BMC deployment, deploy BMC on the solidity side to get the link address


Solidity deployment:
---------------------------------------------------------------------------------------------------------------------------------------------------------------

1. Follow the [solidity_deploy_commands.md](https://github.com/icon-project/btp/blob/btp_nobmv/testnet/solidity/solidity_deploy_commands.md) for deploying sol contracts
2. Remember to clean up local artifacts folder before running the deploy command
2. `npm install` in the solidity/scripts folder
3. rename `.env_example` to `.env` & update the private key


Relay:
---------------------------------------------------------------------------------------------------------------------------------------------------------------

1. Edit both direction config files in the config folder with proper URLs & wallet & data folder & Link addresses in source and destination
2. change offset to the block in the config to BMC deployment block number/ before add Link txn block number
3. add bsc.secret & icon.secret with password to the respective wallet
4. create a data folder if doesnt exist
5. start the relay by passing the proper config files
`./btpsimple start --config ./icon.config.json`
