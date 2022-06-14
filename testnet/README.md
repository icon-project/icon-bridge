
Javascore Deployment:
---------------------------------------------------------------------------------------------------------------------------------------------------------------

1. Copy the javascore jar files from the build folder after the `make build` to artifacts folder(Skip this step if there are no javascore code changes)
    Note: If BMC has any changes, build it separately and copy the jar file to `lib` folder under javascore folder. before running `make build`
2. Update the proper ICON keystore wallet file with funds 
3. Refer to [javascore_deployment_commands](https://github.com/icon-project/icon-bridge/blob/bridge_bsc/testnet/javascore/javascore_deploy_commands.md) file for java deployment commands & configurations
4. Update the keystore file path&name on `keystore` variable & update the `password` variable with the keystore password on the commands env setup step
5. Paste the env export commands in the terminal
6. Proceed with the deployment: after javascore BMC deployment, deploy BMC on the solidity side to get the link address.


Solidity deployment:
---------------------------------------------------------------------------------------------------------------------------------------------------------------

1. Follow the [solidity_deploy_commands.md](https://github.com/icon-project/icon-bridge/blob/bridge_bsc/testnet/solidity/solidity_deploy_commands.md) for deploying sol contracts
2. Remember to clean up local artifacts folder before running the deploy command
2. `npm install` in the solidity/scripts folder
3. rename `.env_example` to `.env` & update the private key 
Note: please remember to use the same wallet's private key that will be used for both deployment & running Relay


Relay:
---------------------------------------------------------------------------------------------------------------------------------------------------------------

1. Edit both direction config files in the config folder with proper URLs & keystore wallet & data folder & Link addresses in source and destination
2. change offset to the block in the config to BMC deployment block number/ before add Link txn block number
3. add bsc.secret & icon.secret with password to the respective wallet
4. create a data folder if doesnt exist
5. start the relay by passing the proper config files
`./btpsimple start --config ./icon.config.json`
