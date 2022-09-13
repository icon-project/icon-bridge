#!/bin/bash
set -e 

########################################################

build_all_images_for_local_deployment() {
    if [ -f local/goloop ] || [ -f local/log ] || [ -f local/artifacts ]; then 
        echo "Previous deployment artifacts exist on local/{goloop,log,artifacts}. Clean using sudo make cleaimglocal. Then runimglocal"
        exit 0
    fi
    echo "Start building docker containers for iconbridge, bsc and icon nodes"
    #make -C ../../../ dist-javascore dist-sol iconbridge-image
	mkdir -p ./data/bsc/node1
    
    export DOCKER_DEFAULT_PLATFORM=linux/amd64
    echo "Build BMR"
    cd ../../..
    docker build -f ./devnet/docker/icon-bsc/Dockerfile -t iconbridge_bsc:latest .
    echo "Build BSC"
    docker build --tag bsc-node ./devnet/docker/bsc-node --build-arg KEYSTORE_PASS=Perlia0
    echo "Build ICON"
    cd ./devnet/docker/goloop
    docker build --tag icon-bsc_goloop:latest .
    echo "Build Process Complete"

    echo "To run the containers. Use command: make runimglocal"
    echo "This will  deploy the smart contracts on the locally deployed containers and run relay"
}

run_all_images_for_local_deployment() {
    if [ -f local/goloop ] || [ -f local/log ] || [ -f local/artifacts ]; then 
        echo "Previous deployment artifacts exist on local/{goloop,log,artifacts}. Clean using sudo make cleaimglocal. Then runimglocal"
        exit 0
    fi
    docker-compose -f docker-compose.yml up --remove-orphans -d
    sleep 5
    echo "Containers are running. Check progress with docker logs command"
}

clean_artifacts_of_local_deployment() {
    echo "Cleaning Artifacts. Might require sudo privilege if not used i.e. sudo make cleanimglocal"
    sleep 2
    docker-compose down -v --remove-orphans
	rm -rf local/artifacts
	rm -rf local/goloop
	rm -rf local/log
	rm -rf data/bsc/node1
}

remove_artifacts_of_local_deployment() {
    echo "Removing Artifacts. Might require sudo privilege if not used i.e. sudo make removeimglocal"
    sleep 2
	docker-compose down -v --remove-orphans	
	docker rmi -f icon-bsc_btp
	docker rm -f javascore-dist
	docker rmi -f iconbridge
	docker rmi -f btp/javascore
	rm -rf local/artifacts 
	rm -rf local/log 
	rm -rf local/goloop
	rm -rf data/bsc/node1
}

########################################################

build_smart_contracts() {
    echo "Creating build artifacts on PC"
    #Run the script from icon-bridge/devnet/docker/icon-bsc
    ICON_BSC_DIR=${PWD}
    ROOT_DIR=$(echo "$(cd "$(dirname "../../../../")"; pwd)")
    if [ -d $ICON_BSC_DIR/build ]; then
        echo "Save Previous Build Artifacts"
        local suffix=$(date +%s)
        mv build build_${suffix}
    fi
    CONTRACTS_DIR="$ICON_BSC_DIR/build/contracts"
    JAVASCORE_DIR="$CONTRACTS_DIR/javascore"
    SOLIDITY_DIR="$CONTRACTS_DIR/solidity"

    mkdir -p "$JAVASCORE_DIR"
    mkdir -p "$SOLIDITY_DIR"
    mkdir -p "$ICON_BSC_DIR/_ixh/keystore"

    echo "Creating go build artifacts"
    cd $ROOT_DIR/cmd/iconvalidators/
    go build .
    mv iconvalidators $ICON_BSC_DIR/

    echo "Creating java build artifacts"
    # build contracts
    cd $ROOT_DIR/javascore
    gradle clean
    gradle bmc:optimizedJar
    gradle bts:optimizedJar
    cp bmc/build/libs/bmc-optimized.jar $JAVASCORE_DIR/bmc.jar
    cp bts/build/libs/bts-optimized.jar $JAVASCORE_DIR/bts.jar
    cp lib/irc2Tradeable-0.1.0-optimized.jar $JAVASCORE_DIR/irc2Tradeable.jar

    # irc2-token
    cd $ROOT_DIR
    git clone https://github.com/icon-project/java-score-examples.git
    cd java-score-examples
    gradle irc2-token:clean
    gradle irc2-token:optimizedJar
    cp irc2-token/build/libs/irc2-token-0.9.1-optimized.jar $JAVASCORE_DIR/irc2.jar
    rm -rf $ROOT_DIR/java-score-examples

    echo "Creating solidity build artifacts"
    # copy solidity
    cp -r $ROOT_DIR/solidity/{bmc,bts} $SOLIDITY_DIR/
    cd $SOLIDITY_DIR
    cd bmc && yarn install
    cd ..
    cd bts && yarn install

    echo "Build Artifacts have been created on path "$ICON_BSC_DIR/build
}

build_chain_nodes() {
    echo "Start building docker containers for bsc and icon nodes"
    #make -C ../../../ dist-javascore dist-sol iconbridge-image
	mkdir -p ./data/bsc/node1
    
    export DOCKER_DEFAULT_PLATFORM=linux/amd64
    cd ../../..
    echo "Build BSC"
    docker build --tag bsc-node ./devnet/docker/bsc-node --build-arg KEYSTORE_PASS=Perlia0
    echo "Build ICON"
    cd ./devnet/docker/goloop
    docker build --tag icon-bsc_goloop:latest .

    echo "Build Process Complete"
    echo "To run the blockchain nodes. Use command: make runnodes"
}

run_chain_nodes() {
    echo "Running block chain nodes locally"
    docker-compose -f docker-compose-nodes.yml up -d
    sleep 5
    local ICON_BSC_DIR=${PWD}
    if [ -d $ICON_BSC_DIR/_ixh ]; then 
        echo "Save Previous Deployment Artifacts"
        local suffix=$(date +%s)
        mv $ICON_BSC_DIR/_ixh $ICON_BSC_DIR/_ixh_${suffix}
    fi
    mkdir -p $ICON_BSC_DIR/_ixh/keystore
    echo "Fetching GodKeys from blockChain nodes"
    docker cp goloop:/goloop/config/goloop.keystore.json $ICON_BSC_DIR/_ixh/keystore/icon.god.wallet.json
    docker cp goloop:/goloop/config/goloop.keysecret $ICON_BSC_DIR/_ixh/keystore/icon.god.wallet.secret
    docker cp goloop:/goloop/config/nid.icon $ICON_BSC_DIR/_ixh/nid.icon
    docker cp binancesmartchain:/bsc/keystore/UTC--2021-07-14T19-55-36.108252000Z--70e789d2f5d469ea30e0525dbfdd5515d6ead30d $ICON_BSC_DIR/_ixh/keystore/bsc.god.wallet.json
    echo -n "Perlia0" > $ICON_BSC_DIR/_ixh/keystore/bsc.god.wallet.secret
    ethkey inspect --json --private --passwordfile $ICON_BSC_DIR/_ixh/keystore/bsc.god.wallet.secret $ICON_BSC_DIR/_ixh/keystore/bsc.god.wallet.json | jq -r .PrivateKey > $ICON_BSC_DIR/_ixh/keystore/bsc.god.wallet.json.priv
    echo "Nodes running"
    docker ps --filter name=goloop --filter name=binancesmartchain
}

clean_deployment_artifacts() {
    echo "cleaning deployment artifacts. Move _ixh folder but retain _ixh/keystore"
    ICON_BSC_DIR=${PWD}
    if [ -d $ICON_BSC_DIR/_ixh ]; then 
        echo "Save Previous Deployment Artifacts"
        local suffix=$(date +%s)
        mv $ICON_BSC_DIR/_ixh $ICON_BSC_DIR/_ixh_${suffix}
        if [ -d $ICON_BSC_DIR/_ixh_${suffix}/keystore ]; then
            mkdir -p $ICON_BSC_DIR/_ixh
            echo "Backup "$ICON_BSC_DIR/_ixh_${suffix}/keystore " inside "$ICON_BSC_DIR/_ixh/
            cp -r $ICON_BSC_DIR/_ixh_${suffix}/keystore $ICON_BSC_DIR/_ixh/keystore
        fi
    fi
    echo "Finished cleaning. Deploy new set of contracts with make runsc"
}

deploy_smart_contracts_on_testnet() {
    echo "Deploying contracts. ./deploysc.sh"
	cd scripts
	./deploysc.sh
}

deploy_smart_contracts_on_localnet() {
    echo "Deploying contracts. ./deploysc.sh"
	cd scripts
    if [ ! -f config_testnet.sh ]; then
        cp config.sh config_testnet.sh
        cp config_local.sh config.sh
    fi
	./deploysc.sh
}

########################################################

build_relay_img() {
    echo "Build BMR"
    cd ../../..
    docker build -t localnets:5000/bmr:latest -f devnet/docker/icon-bsc/bmr.Dockerfile .
}

run_relay_img() {
    local relayConfigPath=${PWD}/_ixh/bmr.config.json
    if [ ! -f relayConfigPath ]; then
        echo "relay config does not exist on path "$relayConfigPath
        exit 0
    fi
    cp $relayConfigPath ./bmr.config.json
    docker-compose -f docker-compose-bmr.yml up -d 
    sleep 5 
    rm bmr.config.json
    echo "Relay Running"
    docker ps --filter name=bmr
}

run_relay_from_source() {
    local relayConfigPath=${PWD}/_ixh/bmr.config.json
    if [ ! -f $relayConfigPath ]; then
        echo "relay config does not exist on path "$relayConfigPath
        exit 0
    fi
    cd ../../../cmd/iconbridge/
    go run . -config $relayConfigPath
}

stop_relay_img() {
    docker-compose -f docker-compose-bmr.yml down
}

stop_chain_nodes() {
    docker-compose -f docker-compose-nodes.yml down
}
########################################################



if [ $# -eq 0 ]; then
    echo "No arguments supplied. Check README for details"
elif [ $1 == "buildimglocal" ]; then
    build_all_images_for_local_deployment
elif [ $1 == "runimglocal" ]; then
    run_all_images_for_local_deployment
elif [ $1 == "cleanimglocal" ]; then
    clean_artifacts_of_local_deployment
elif [ $1 == "removeimglocal" ]; then
    remove_artifacts_of_local_deployment
elif [ $1 == "buildsc" ]; then
    build_smart_contracts
elif [ $1 == "buildnodes" ]; then 
    build_chain_nodes
elif [ $1 == "runnodes" ]; then 
    run_chain_nodes
elif [ $1 == "stopnodes" ]; then 
    stop_chain_nodes
elif [ $1 == "deploysctestnet" ]; then 
    deploy_smart_contracts_on_testnet
elif [ $1 == "deploysclocalnet" ]; then 
    deploy_smart_contracts_on_localnet
elif [ $1 == "runrelaysrc" ]; then 
    run_relay_from_source
elif [ $1 == "cleanartifacts" ]; then 
    clean_deployment_artifacts
elif [ $1 == "buildrelayimg" ]; then 
    build_relay_img
elif [ $1 == "runrelayimg" ]; then 
    run_relay_img
elif [ $1 == "stoprelayimg" ]; then 
    stop_relay_img
else
    echo "To build on docker container: make buildimg. To run on local PC: make buildsc  Check README.md for more"
fi
