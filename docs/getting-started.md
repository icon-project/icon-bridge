# Getting Started

This section explains the repo folder structure and how to setup iconbridge and its dependencies, which in this case are the smart contracts that needs to be deployed on the respective chains you want to bridge. There are multiple alternatives to go about this process depending upon which networks you want to bridge, whether testnet, mainnet or locally deployed blockchains. 


Get familiar with some terminologies used in icon-bridge [here.](terminologies.md)

## Introduction To Repo 
| Directory | Description |
|:----------|:------------|
| /cmd | Includes different executables like iconbridge, iconvalidators, e2etest|
| /cmd/iconbridge | Implementation of iconbridge for different chains|
| /cmd/iconbridge/chain | Includes read/write API per blockchain|
| /cmd/iconbridge/relay | Uses chain API to relay message between chains|
| /cmd/iconbridge/stat | System Metrics Collector Service|
| /cmd/iconvalidators | Module that retrieves validator info of a block in ICON chain|
| /cmd/e2etest | End-to-End Testing Module|
| /common | Common Code|
| /devnet | Includes scripts used to deploy iconbridge and blockchain nodes|
| /devnet/docker/bsc-node | Sripts to build BNB Smart Chain docker image|
| /devnet/docker/goloop | Scripts to build ICON Chain docker image|
| /devnet/docker/icon-bsc | Scripts to build iconbridge between icon and bsc chains|
| /devnet/docker/icon-hmny | Scripts to build iconbridge between icon and harmony chains|
| /doc | Documentation(Obsolete) |
| /docs | Documentation |
| /docker and /docker-compose | Scripts to create docker containers (Some of it Obsolete) |
| /javascore | javascore smart contracts |
| /solidity | solidity smart contracts |
| /pyscore | Python Smart Contracts (Obsolete) |
| /testnet |-- (Obsolete) |

>TLDR version of deployment is available [here](#tdlr)

## Local Deployment
We’ll start with a description on setting up iconbridge on a locally deployed blockchain. A local deployment involves the following steps:

- Set up & run block-chain nodes that you want to bridge
- Deploy smart contracts on those chains
- Run relay to exchange messages given these deployed contracts

We’ll briefly touch upon these steps so that the commands to run these processes can add up easily when we provide them later.

* ### Setting up Blockchain Nodes
    This requires blockchain nodes, usually available as docker images, to be run. Each node has some configuration parameters, which we require later in the process. For example, rpc endpoint, network-id and god-wallet are required to communicate and execute transactions on the network.
    ```sh
    make cleanimglocal          # Clean Previous build & deployment artifacts if present
    make buildimglocal          # Build BMR,BSC,ICON images
    ```

* ### Deploying smart contracts
    Alongside blockchain configuration, we also need to prepare artifacts needed to deploy our smart contracts to the networks. Setting up these artifacts includes building jar files for javascore smart contracts and installing node modules for solidity smart contracts. 
    Once these build artifacts are generated, we can use commands to interact with and deploy smart contracts to the block-chains. Contract addresses of the deployed contracts are needed to interact with them later.

* ### Run relay
    We now use the contract addresses output from the deployment process as an input parameter to the relay. The relay itself can be run either as an executable or a docker container on any system.

    ```sh
    make runimglocal            # Run containers. deploys smart contract, run relay
    # Optional:
        make removeimglocal     # Cleans artifacts and also removes previously built images
    ```


## Alternatives to Local Deployment:
Depending on whether you want the contract deployment plus relay-run to be done inside a docker container or on your local system, we have a couple of alternatives you can follow.

* ### Single-step process of Local Deployment
    The former method (i.e. running a docker container) is more suited for someone new to the process as it takes care of setting up most of the dependencies and the user is presented with a usable version of locally deployed iconbridge.

    > Note: This process is susceptible to problems on M1 MacOS if you are trying to build the docker-container yourself and not use the one available in docker-hub. While creating docker images, the build-smart contracts stage is prone to get stuck on M1 MacOS specifically. To get past this problem, you may cancel and rerun until it completes or use the alternate method.
    
    _After you’ve cloned the icon-bridge repository, go through the following steps_
    - Change directory to where icon-bsc bridge deployment scripts are placed.
        ```sh
        cd icon-bridge/devnet/docker/icon-bsc/
        ```
    - Clean previous build and deployment artifacts, if any
        ```sh
        make cleanimglocal
        ```
    - Create docker images of ICON, BSC and BMR (a.k.a iconbridge)
        ```sh
        make buildimglocal
        ```
    - Run docker containers of the above images. Then, deploy smart contracts on the ICON and BSC nodes. Afterwards, run the relay image using the deployed smart contract addresses through a config file.
        ```sh
        make runimglocal
        ```
    - ( Optional ) Clean build and deployment artifacts and remove previously built docker images
        ```sh
        make removeimglocal
        ```
    > Note: The removal process can require you to use sudo (root user privilege) because some files generated by ICON nodes get stored with root privilege.

    The demerits of the process is that since the source code is packaged inside the docker container and multiple tasks get auto-started at the start of the docker container, it will be difficult to intervene and change the packaged container if one requires changing source code. This makes the process less flexible for a developer to work with.

* ### Multi-Step Process of Local Deployment
    Here, each of the steps are separately run with respective commands on your PC. The dependencies that need to be installed are provided [here](#dependencies)
    
    _After you’ve cloned repository and installed dependencies, run the following commands_

    - Change directory to where icon-bsc bridge deployment scripts are placed.
        ```sh
        cd icon-bridge/devnet/docker/icon-bsc/
        ```
    - Builds smart contracts i.e jar and node modules
        ```sh
        make buildsc                    
        ```
    - Build ICON and BSC docker images
        ```sh
        make buildnodes
        ```
    - Run ICON and BSC docker containers
        ```sh
        make runnodes
        ```
    - Deploy smart contracts on ICON and BSC Nodes. Be sure to run this step only after it’s been at least a minute since nodes were run as they need time to initialize.
        ```sh
        make deploysclocalnet
        ```
    - Run relay from soure code. Optionally you can create docker image of relay and run it using the following
        ```sh
        make runrelaysrc
        ```
        #### _Optional Commands_
        - Stop relay containers if running
            ```sh
            make stoprelayimg
            ```
        - Build docker image of relay
            ```sh
            make buildrelayimg
            ```
        - Run relay docker containers
            ```sh
            make runrelayimg
            ```
    - (Optional) These optional commands can be used if needed
        - Stop ICON & BSC nodes
            ```sh
            make stopnodes
            ```
        - Remove deployment artifacts but reuses pre-existing keystore files.
            ```sh
            make cleanartifacts
            ```
        - Run e2etests using the deployment artifacts.
            ```sh
            make rune2etests
            ```

## Mainnet/Testnet Deployment
Mainnet/Testnet deployment is similar to the Multi-Step local deployment process mentioned above. The differences are:

- Block-chain nodes are already running and we only need their parameters
- Since we do not have god wallet, we require an account in the network that has sufficient balance to do the contract deployment

The steps for testnet/mainnet deployment can be listed as follows:
- Builds smart contracts i.e jar and node modules
    ```sh
    make buildsc
    ```
- Deploys smart contract. Will prompt you to fund deployer wallet if not already present in the required path. Blockchain’s parameters are used from a config file (scripts/config.sh)
    ```sh
    make deploysctestnet
    ```
- Run relay from source
    ```sh
    make runrelaysrc
    ```
- (Optional) Run e2e tests
    ```sh
    make rune2etests
    ```

## Things to remember about deployment
- Keystore files with appropriate names should be present inside icon-bsc/_ixh/keystore. A set of keystore files (for deployer/god, bmc, bts, bmr and fee-aggregator) are auto-generated if not already present during the smart contract deployment process. 
- Fund god wallets (named *.god.wallet.json) before deploying smart contracts
- Fund bmr (named *.bmr.wallet.json) wallet before running relay
- Fund bts (named *.bts.wallet.json) wallet before running e2etests
- For multi-step deployments, running make cleanartifacts will not delete the keystore files if previously present. This is done to ease the redeployment process.
- If you’re switching from multi-step testnet deployment to multi-step local deployment, make sure to remove keystore files because the same keystores can not be used for localnet.
- Though not advisable, you can use the same wallet for god, bts and bmr accounts for testnet and localnet deployments. This eliminates having to fund necessary accounts for different tasks.
- For multi-step deployment, if the process terminates or needs to be terminated, then it can be resumed by rerunning the smart contract deployment command
- Cleaning build artifacts for single-step deployment can require sudo privilege
- Relay takes some time to synchronize upto the latest block. Once it does, relay’s logs will include “block notification” messages of recent block numbers
- You can run e2etests with the command given above only after relay has synchronized. More about e2etests on following sections.

## How to run tests
* Running unit tests
    - For go code
        ```sh
        go test
        ```
    - For solidity
        ```sh
        yarn test
        ```
    - For javascore
        ```sh
        ./gradlew :<project name>:test
        ```
* Running e2e tests
    ```sh
    cd devnet/docker/icon-bsc && make rune2etests
    ```
    Details about e2e tests [here](./tests/e2e-test.md)


##  Dependencies:
The following must be installed on your PC for piecewise deployment and testnet/mainnet deployment.
- Docker 
- java 11 and gradle 6.7.1
- goloop 
- NodeJS
- Truffle  v5.5.5
- EthKey 
- Go >= 1.13 (ref: go.mod)
 
   
1. ### Docker
 
    To build, publish, run blockchains (icon/bsc) locally or remote docker host. Download and install docker from https://docs.docker.com/engine/install/ubuntu/
 
    After installing, make sure that the user account used to run docker (_default is ubuntu_) is added to `docker` group.

    ```sh
    sudo groupadd docker

    sudo usermod -aG docker $USER

    newgrp docker
    ```
 

 
    Fully logout, and log back in to be able apply the changes.
 

 
2. ### SdkMan
 

 
    To install gradle and java.
 

 
    1. _`fish`_
 

 
       https://github.com/reitzig/sdkman-for-fish
 

 
    2. _`bash`_
 

 
       https://sdkman.io/install
 

 
3. ### Java and Gradle

    To build javascores.

    1. _`Java`_
  
       ```sh
       sdk install java 11.0.11.hs-adpt
       ```
 

 
    2. _`gradle`_
       ```sh
       sdk install gradle 6.7.1
       ```
 

 
4.  ### Goloop

    To interact with icon blockchain using RPC calls and generate keystores.
 
 
    https://github.com/icon-project/goloop
 
 
    ```sh
    go install github.com/icon-project/goloop/cmd/goloop@latest
    ```

    Note: If `go install` doesn't work use `go get` instead.
 

 
5.  ### NodeJS 
    To build and deploy solidity smart contracts.

    1. _`fish`_
  
       `nvm`: https://github.com/jorgebucaran/nvm.fish
  
       ```sh 
       fisher install jorgebucaran/nvm.fish
 
       nvm install v15.12.0
 
       set --universal nvm_default_version v15.12.0
 
       nvm use v15.12.0
 
       node --version > ~/.nvmrc 
       ```

    2. _`bash`_
 
       https://github.com/nvm-sh/nvm
 
        ```sh
        curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash

        nvm install v15.12.0

        nvm use v15.12.0

        node --version > ~/.nvmrc 
        ```
 
6.  ### Truffle

    https://trufflesuite.com/docs/truffle/getting-started/installation.html
  
    `npm install -g truffle@5.5.5`

 
7.  ### Ethkey 
    `go get github.com/ethereum/go-ethereum/cmd/ethkey`


<br/>

# TDLR
 
Deployment Involves the following steps:
|Task|Description|
|:-|:-|
|1. RUN_NODES              | There should be icon and bsc nodes running (either as testnet/mainnet) or as local docker containers|
|2. BUILD_ARTIFACTS        | Build artifacts (jar & node_modules) need to be created|
|3. DEPLOY_SMART_CONTRACTS | Smart Contracts should be deployed on these icon and bsc nodes using build artifacts |
|4. RELAY_CONFIG           | After smart contract deployment completes, the deployment artifacts includes relay config bmr.config.json|
|5. RELAY_RUN              | Relay (also called bmr, iconbridge) should be run|
|6. USE_CASES              | Deployment artifacts include addresses.json, Use the contract addresses mentioned there for any use cases (e.g. token transfer)|
 


 
## Deploying ICON_BRIDGE
 
1. Deploy on local machine
    - [Full Deployment](#deploy-on-local-machine)        (Useful for newcomers. All the steps done inside a docker container. )
    - [Piecewise Deployment](#piecewise-deployment)    (Useful for developers. Steps run on docker + PC)
2. [Deploy on Mainnet/Testnet](#deploy-on-mainnettestnet)



##  Description for Deployment Processes:

 
* ###  Deploy on local machine
 
    - #### Full Deployment 
        ```sh 
        make cleanimglocal          # Clean Previous build & deployment artifacts if present
 
        make buildimglocal          # Build BMR,BSC,ICON images
 
        make runimglocal            # Run containers. deploys smart contract, run relay
 
        # Optional:
 
            make removeimglocal     # Cleans artifacts and also removes previously built images
        ```
 

 
    - #### PieceWise Deployment
    
        First, install required [dependencies](#dependencies)

        ```sh 
        make buildsc                    # Builds smart contracts
 
        make buildnodes                 # Build ICON & BSC nodes
  
        make runnodes                   # RUN ICON & BSC nodes locally
 
        make deploysclocalnet           # Deploy smart contracts on ICON & BSC nodes
 
        make runrelaysrc                # Run relay from source
 
        # Optional:                     # Run relay from docker container instead of source
 
                make stoprelayimg       # Stop relay docker container (if present)
 
                make buildrelayimg      # Build relay docker image 
 
                make runrelayimg        # Run relay docker container
 
        # Optional:
 
                make stopnodes          # Stop icon & bsc nodes
 
                make cleanartifacts     # Clean artifacts generated by smart contract deployment
        ```

 

 
* ### Deploy on Mainnet/Testnet

    First, install required [dependencies](#dependencies)
    
    ```sh
    make buildsc                    # Builds smart contracts
 
    make deploysctestnet            # Deploy smart contracts on Testnet. Using scripts/config.sh
 
    make runrelaysrc                # Run relay from source
 
    # Optional:                     # Run relay from docker container instead of source
 
            make stoprelayimg       # Stop relay docker container (if present)
 
            make buildrelayimg      # Build relay docker image 
 
            make runrelayimg        # Run relay docker container
 
    # Optional:
 
            make cleanartifacts     # Clean artifacts generated by smart contract deployment
 
    ```
 

 
