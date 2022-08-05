Run bsc node with
cd icon-bridge/devnet/docker/bsc-node
docker run -d -p 8545:8545 -p 8546:8546 bsc-node

Run goloop node using "icon" docker image from
icon-bridge/devnet/docker/icon-hmny/src/docker-compose.nodes.yml

To build javascore
make buildsc

Provide parameters in scripts/config.sh 

Deploy smart contract with 
cd ./scripts/
./deploysc.sh

Run relay with
export bmr_config_json=$(cat ./_ixh/bmr.config.json)
docker-compose -f docker-compose-bmr.yml up -d
