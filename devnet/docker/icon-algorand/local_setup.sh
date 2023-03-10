echo "Building Goloop container..."
cd ./../goloop
docker build -t icon-algorand_goloop .

echo "Running Goloop container..."
docker run -d \
  --name goloop \
  -p 9080:9080 \
  -e GOLOOP_NODE_DIR=/goloop/data/goloop \
  -e GOLOOP_LOG_WRITER_FILENAME=/goloop/data/log/goloop.log \
  -t icon-algorand_goloop

echo "Running Algod..."
cd ./../../algorand
goal network create -r /tmp/testnet -t ./template.json
cp ./config.json /tmp/testnet/Node
cp ./algod.token /tmp/testnet/Node
cp ./kmd_config.json /tmp/testnet/Node/kmd-v0.5/kmd_config.json
cp ./kmd.token /tmp/testnet/Node/kmd-v0.5/kmd.token
goal network start -r /tmp/testnet

echo "Compiling golang tools..."
cd ../../cmd/tools/algorand
./install-tools.sh

echo "Building ICON smart contracts..."
cd ../../../javascore
./gradlew dummyBSH:optimizedJar
./gradlew bmc:optimizedJar
./gradlew wrapped-token:optimizedJar


echo "Building Algorand smart contracts..."
cd ../pyteal
./build.sh bmc.bmc bmc
./build.sh bsh.bsh bsh
./build.sh escrow.escrow escrow
