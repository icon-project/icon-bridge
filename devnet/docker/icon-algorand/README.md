# Algorand - Icon Bridge Integration

## Intro
The icon-algorand bridge, aims to integrate the Algorand blockchain into the icon-bridge system.
It should allow users on both ends to exchange messages between the two blockchains.
The primary use case built, is to provide a token transfer bridge between the two chains,
powered by the messaging exchange system.

## How to run the bridge
There are three ways we developed to run a new instance of the icon-algorand bridge, which will be describded
throughout the next sections.
These are a github actions workflow, running automated builds and tests, with a local chains setup, which should be able run on CI and validate the wellfunctioning of the system.
Additionally, we have development scripts to launch the systems on local machines, which aim the localnets 
and testnets respectively.


### Workflow
In order to showcase a token transfer, the [algorand-integration](/.github/workflows/algorand-integration.yml) workflow was created.
It aims to set up local instances of both chains and send messages in both directions asking for
token transfers.
It starts by installing the required dependencies, running a container from the icon
goloop image and setting up a new algorand local network.
Once these are running, the smart contracts for both are built, deployed and registered on the
opposite chains, creating a valid link between them.
Throughout the execution of the previous steps, a set of environmental variables will be created,
which will then be copied to the relayer [config file](/devnet/docker/icon-algorand/algo-config.json), allowing it to run accordinglly.


### Local execution
It's also possible to replicate the same steps locally, which can be much more helpful for debugging:
1. Install Algorand and Pyteal, using the same cmds provided on the workflow file.
2. Run  [prepare_local_env.sh](devnet/docker/icon-algorand/prepare_local_env.sh) to setup local chain nodes and build the smart contracts.
3. Run [setup_system.sh](devnet/docker/icon-algorand/setup_system.sh) to deploy the contracts and setup the relayer config file.
4. Go to ``./cmd/iconbridge`` and run ``go run . -config=../../devnet/docker/icon-algorand/algo-config.json``
to start the relayer.
5. To execute any of the integration tests, go to ``/devnet/docker/icon-algorand`` and run the
respective script.


### Testnet execution
At last, the relayer can also be ran connected to the respective testnets. Note that first you need
to obtain a valid testnet node access from algorand and set the env vars ``$ALGO_TEST_ADR`` and ``$ALGO_TEST_TOK``, to the given
testnet address and token values.
1. Go to ``./devnet/docker/icon-algorand`` and run [testnet_start_relay.sh](/devnet/docker/icon-algorand/testnet_start_relay.sh) - Beware that this script will create wallet accounts that need to be funded.
on the respective testnet faucets. The algorand one can be accessed [here](https://bank.testnet.algorand.network/).
2. To execute any of the integration tests, go to ``/devnet/docker/icon-algorand`` and run the
respective script.
