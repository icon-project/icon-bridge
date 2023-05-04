echo "Compiling golang tools..."
cd ../../../cmd/tools/algorand
./install-tools.sh

echo "Building ICON smart contracts..."
cd ../../../javascore
./gradlew dummyBSH:optimizedJar
./gradlew bmc:optimizedJar
./gradlew wrapped-token:optimizedJar
./gradlew test-token:optimizedJar
./gradlew escrow:optimizedJar

echo "Building Algorand smart contracts..."
cd ../pyteal
./build.sh bmc.bmc bmc
./build.sh bsh.bsh bsh
./build.sh escrow.escrow escrow
./build.sh reserve.reserve reserve

cd ../devnet/docker/icon-algorand
echo "Exporting path variable"
export PATH=$PATH:~/go/bin

echo "Creating the cache directory"
if [ -d "cache" ]; then
    rm -r cache
fi
mkdir cache

echo $ALGO_TEST_ADR >cache/algod_address
echo $ALGO_TEST_TOK >cache/algo_token

echo "Generating the icon keystore wallet"
goloop ks gen --out icon.keystore.json
ICON_KS_ADDRESS=$(cat icon.keystore.json | jq -r '.address')

echo "Fund wallet $ICON_KS_ADDRESS on Icon faucet and press enter to continue..."
read

echo "Generating the test token minter wallet"
goloop ks gen --out test_token_minter.keystore.json
TTM_KS_ADDRESS=$(cat test_token_minter.keystore.json | jq -r '.address')

echo "Fund wallet $TTM_KS_ADDRESS on Icon faucet and press enter to continue..."
read

echo "Generating the sender keystore wallet"
goloop ks gen --out sender.keystore.json
SENDER_KS_ADDRESS=$(cat sender.keystore.json | jq -r '.address')

echo "Fund wallet $SENDER_KS_ADDRESS on Icon faucet and press enter to continue..."
read

echo "Start kmd locally"
goal kmd start -t 0 -d /tmp/testnet/Node
sleep 5

KMD_ADDRESS="http://$(cat /tmp/testnet/Node/kmd-v0.5/kmd.net)"
KMD_TOKEN="$(cat /tmp/testnet/Node/kmd-v0.5/kmd.token)"
echo "Creating minter KMD wallet"
kmd_output=$(go run kmd.go $KMD_ADDRESS $KMD_TOKEN "Minter")
sleep 5

ALGORAND_MINTER_PRIVATE_KEY=$(echo $kmd_output | jq -r '.private_key')
MINTER_ADDRESS=$(echo $kmd_output | jq -r '.address')

echo Fund minter wallet $MINTER_ADDRESS on Algo faucet and press enter to continue...
read
sleep 5

echo "Deploy Algorand Test Asset"
ASA_ID=$(
    PRIVATE_KEY=$ALGORAND_MINTER_PRIVATE_KEY ALGOD_ADDRESS=$ALGO_TEST_ADR ALGOD_TOKEN=$ALGO_TEST_TOK deploy-asset 1000000000000 6 TABC Test AB Coin http://example.com/ abcd
)
echo $ALGORAND_MINTER_PRIVATE_KEY >cache/algo_minter_private_key
echo $ASA_ID >cache/algo_test_asset_id

echo "Deploy Algorand Wrapped Token"
WTKN_ID=$(
    PRIVATE_KEY=$(cat cache/algo_minter_private_key) ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) deploy-asset 1000000000000 6 WTKN "Wrapped Test Token" http://example.com/ wtkn
)
echo $WTKN_ID > cache/algo_wrapped_token_id

echo "Deploying BMC contract to the ICON network"
CONTRACT=../../../javascore/bmc/build/libs/bmc-optimized.jar
DEPLOY_TXN_ID=$(goloop rpc sendtx deploy $CONTRACT \
    --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex \
    --key_store ./icon.keystore.json --key_password gochain \
    --nid 0x2 \
    --content_type application/java \
    --step_limit 3000000000 \
    --param _net="0x2.icon")
./../../algorand/scripts/wait_for_testnet_txn.sh $DEPLOY_TXN_ID scoreAddress >cache/icon_bmc_addr

echo "Update the config file with the new keystore"
jq --slurpfile new_contents icon.keystore.json '.relays[0].dst.key_store = $new_contents[0]' algo-testnet-config.json >tmpfile && mv tmpfile algo-testnet-config.json

echo "Deploying dummyBSH contract to the ICON network"
CONTRACT=../../../javascore/dummyBSH/build/libs/dummyBSH-optimized.jar
TXN_ID=$(
    goloop rpc sendtx deploy $CONTRACT \
        --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex \
        --key_store icon.keystore.json --key_password gochain \
        --nid 0x2 --step_limit 10000000000 \
        --content_type application/java \
        --param _bmc=$(cat cache/icon_bmc_addr) \
        --param _to=0x14.algo
)
./../../algorand/scripts/wait_for_testnet_txn.sh $TXN_ID scoreAddress >cache/icon_dbsh_addr

echo "Deploying WrappedToken contract to the ICON network"
CONTRACT=../../../javascore/wrapped-token/build/libs/WrappedToken-optimized.jar
TXN_ID=$(
    goloop rpc sendtx deploy $CONTRACT \
        --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex \
        --key_store icon.keystore.json --key_password gochain \
        --nid 0x2 --step_limit 10000000000 \
        --content_type application/java \
        --param _bmc=$(cat cache/icon_bmc_addr) \
        --param _to=0x14.algo \
        --param _asaId=$ASA_ID \
        --param _name="Wrapped Test Token" \
        --param _symbol="WTT" \
        --param _decimals=6
)
./../../algorand/scripts/wait_for_testnet_txn.sh $TXN_ID scoreAddress >cache/icon_wtt_addr

echo "Deploying TestToken contract to the ICON network"
CONTRACT=../../../javascore/test-token/build/libs/TestToken-optimized.jar
TXN_ID=$(
    goloop rpc sendtx deploy $CONTRACT \
        --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex \
        --key_store test_token_minter.keystore.json --key_password gochain \
        --nid 0x2 --step_limit 10000000000 \
        --content_type application/java \
        --param _name="Test Token" \
        --param _symbol="TKN" \
        --param _decimals=6 \
        --param _amount=100000000000
)
./../../algorand/scripts/wait_for_testnet_txn.sh $TXN_ID scoreAddress > cache/icon_test_token_addr

echo "Deploying Escrow contract to the ICON network"
CONTRACT=../../../javascore/escrow/build/libs/Escrow-optimized.jar
TXN_ID=$(
    goloop rpc sendtx deploy $CONTRACT \
        --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex \
        --key_store icon.keystore.json --key_password gochain \
        --nid 0x2 --step_limit 10000000000 \
        --content_type application/java \
        --param _bmc=$(cat cache/icon_bmc_addr) \
        --param _to=0x14.algo \
        --param _asaId=$WTKN_ID \
)
./../../algorand/scripts/wait_for_testnet_txn.sh $TXN_ID scoreAddress > cache/icon_escrow_addr

echo "Registering dBSH to the BMC contract on the ICON network"
TXN_ID=$(goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) \
    --method addService \
    --value 0 \
    --param _addr=$(cat cache/icon_dbsh_addr) \
    --param _svc="dbsh" \
    --step_limit=3000000000 \
    --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex \
    --key_store icon.keystore.json --key_password gochain \
    --nid=0x2)
./../../algorand/scripts/wait_for_testnet_txn.sh $TXN_ID

echo "Registering Wrapped Token to the BMC contract on the ICON network"
TXN_ID=$(
    goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) \
        --method addService \
        --value 0 \
        --param _addr=$(cat cache/icon_wtt_addr) \
        --param _svc="wtt" \
        --step_limit=3000000000 \
        --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex \
        --key_store icon.keystore.json --key_password gochain \
        --nid=0x2
)
./../../algorand/scripts/wait_for_testnet_txn.sh $TXN_ID

echo "Registering Escrow to the BMC contract on the ICON network"
TXN_ID=$(
    goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) \
        --method addService \
        --value 0 \
        --param _addr=$(cat cache/icon_escrow_addr) \
        --param _svc="i2a" \
        --step_limit=3000000000 \
        --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex \
        --key_store icon.keystore.json --key_password gochain \
        --nid=0x2
)
./../../algorand/scripts/wait_for_testnet_txn.sh $TXN_ID

echo "Creating creator KMD wallet"
kmd_output=$(go run kmd.go $KMD_ADDRESS $KMD_TOKEN "Creator")
sleep 5

PRIVATE_KEY=$(echo $kmd_output | jq -r '.private_key')
ADDRESS=$(echo $kmd_output | jq -r '.address')

echo Fund creator wallet $ADDRESS on Algo faucet and press enter to continue...
read
sleep 5

echo "Creating algo receiver KMD wallet"
kmd_output=$(go run kmd.go $KMD_ADDRESS $KMD_TOKEN "AlgoReceiver")
sleep 5

ALGO_RECEIVER_PRIVATE_KEY=$(echo $kmd_output | jq -r '.private_key')
ALGO_RECEIVER_ADDRESS=$(echo $kmd_output | jq -r '.address')

echo Fund receiver wallet $ALGO_RECEIVER_ADDRESS on Algo faucet and press enter to continue...
read
sleep 5

echo "Deploying BMC contract to the Algorand network"
BMC_TX_ID=$(PRIVATE_KEY=$PRIVATE_KEY ALGOD_ADDRESS=$ALGO_TEST_ADR ALGOD_TOKEN=$ALGO_TEST_TOK deploy-contract ../../../pyteal/teal/bmc)
echo "Deploying dBSH contract to the Algorand network"
DUMMY_BSH_TX_ID=$(PRIVATE_KEY=$PRIVATE_KEY ALGOD_ADDRESS=$ALGO_TEST_ADR ALGOD_TOKEN=$ALGO_TEST_TOK deploy-contract ../../../pyteal/teal/bsh)
echo "Deploying Escrow contract to the Algorand network"
ESCROW_TX_ID=$(PRIVATE_KEY=$PRIVATE_KEY ALGOD_ADDRESS=$ALGO_TEST_ADR ALGOD_TOKEN=$ALGO_TEST_TOK deploy-contract ../../../pyteal/teal/escrow)
echo "Deploying Reserve contract to the Algorand network"
RESERVE_TX_ID=$(PRIVATE_KEY=$PRIVATE_KEY ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) deploy-contract ../../../pyteal/teal/reserve)

echo "Setting BMC_APP_ID"
BMC_APP_ID=$(ALGOD_ADDRESS=$ALGO_TEST_ADR ALGOD_TOKEN=$ALGO_TEST_TOK get-app-id $BMC_TX_ID)

echo "Setting DUMMY_BSH_APP_ID"
DUMMY_BSH_APP_ID=$(ALGOD_ADDRESS=$ALGO_TEST_ADR ALGOD_TOKEN=$ALGO_TEST_TOK get-app-id $DUMMY_BSH_TX_ID)

echo "Setting ESCROW_APP_ID"
ESCROW_APP_ID=$(ALGOD_ADDRESS=$ALGO_TEST_ADR ALGOD_TOKEN=$ALGO_TEST_TOK get-app-id $ESCROW_TX_ID)

echo "Setting RESERVE_APP_ID"
RESERVE_APP_ID=$(ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-app-id $RESERVE_TX_ID)

echo "Updating config file with the new key"
jq --arg PRIVATE_KEY "$PRIVATE_KEY" '.relays[1].dst.key_store.id=$PRIVATE_KEY' algo-testnet-config.json >config_tmp.json
mv config_tmp.json algo-testnet-config.json

echo "Updating cache files"
echo $PRIVATE_KEY >cache/algo_private_key
echo $ALGO_RECEIVER_PRIVATE_KEY >cache/algo_receiver_private_key
printf '%s' "$BMC_APP_ID" >cache/bmc_app_id
printf '%s' "$DUMMY_BSH_APP_ID" >cache/dbsh_app_id
printf '%s' "$ESCROW_APP_ID" >cache/escrow_app_id
printf '%s' "$RESERVE_APP_ID" >cache/reserve_app_id

echo "Setting link to algo btp on icon bmc"
LINK_TXN_ID=$(goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) \
    --method addLink --param _link=$(cat cache/algo_btp_addr) \
    --key_store ./icon.keystore.json --key_password gochain \
    --nid 0x2 --step_limit 3000000000 --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex)

./../../algorand/scripts/wait_for_testnet_txn.sh $LINK_TXN_ID

echo "Getting latest algo round"
get-last-round $ALGO_TEST_ADR $ALGO_TEST_TOK >cache/algo_last_round

echo "Setting round on icon bmc"
ADD_ROUND=$(goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) --method setLinkRxHeight \
    --param _link=$(cat cache/algo_btp_addr) --param _height=$(cat cache/algo_last_round) \
    --key_store ./icon.keystore.json --key_password gochain \
    --nid 0x2 --step_limit 3000000000 --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex)

./../../algorand/scripts/wait_for_testnet_txn.sh $ADD_ROUND

echo "Adding relay to icon bmc"
ADD_RELAY=$(goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) --method addRelay \
    --param _link=$(cat cache/algo_btp_addr) \
    --param _addr=$(cat icon.keystore.json | jq -r .address) \
    --key_store ./icon.keystore.json --key_password gochain \
    --nid 0x2 --step_limit 3000000000 --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex)

./../../algorand/scripts/wait_for_testnet_txn.sh $ADD_RELAY

cd ../../../cmd/iconbridge

ICON_ALGO="../../devnet/docker/icon-algorand"

ICON_BTP="btp://0x2.icon/$(cat "$ICON_ALGO/cache/icon_bmc_addr")"
echo $ICON_BTP >$ICON_ALGO/cache/icon_btp_addr

ALGO_BTP=$(cat "$ICON_ALGO/cache/algo_btp_addr")
ALGO_ROUND=$(cat "$ICON_ALGO/cache/algo_last_round" | tr -d '\n')

echo "Getting icon last block height"
ICON_HEIGHT=$(goloop rpc lastblock --uri https://lisbon.net.solidwallet.io/api/v3/icon_dex | jq .height)
ICON_HEIGHT_HEX=$(printf "%x" $((ICON_HEIGHT - 1)))

echo "Building icon validators"
go build -o ../iconvalidators ./../iconvalidators
VALIDATORS=$(URI=https://lisbon.net.solidwallet.io/api/v3/icon_dex HEIGHT=0x$ICON_HEIGHT_HEX ./../iconvalidators/iconvalidators | jq -r .hash)

echo "Update config file with multiple variables"
jq --arg ICON_BTP "$ICON_BTP" \
    --arg ALGO_BTP "$ALGO_BTP" \
    --argjson ALGO_ROUND "$ALGO_ROUND" \
    --argjson BMC_ID "$BMC_APP_ID" \
    --arg VALIDATORS "$VALIDATORS" \
    --argjson ICON_HEIGHT "$((ICON_HEIGHT - 1))" \
    '.relays[0].dst.address=$ICON_BTP |
    .relays[0].src.address=$ALGO_BTP |
    .relays[0].src.options.verifier.round=$ALGO_ROUND |
    .relays[0].src.options.appID=$BMC_ID |
    .relays[1].dst.options.bmc_id=$BMC_ID |
    .relays[1].src.address=$ICON_BTP |
    .relays[1].dst.address=$ALGO_BTP |
    .relays[1].src.options.verifier.validatorsHash=$VALIDATORS |
    .relays[1].src.options.verifier.blockHeight=$ICON_HEIGHT' $ICON_ALGO/algo-testnet-config.json >config_tmp.json
mv config_tmp.json $ICON_ALGO/algo-testnet-config.json

echo "Register Algorand BSHs"
ALGOD_ADDRESS=$ALGO_TEST_ADR ALGOD_TOKEN=$ALGO_TEST_TOK PRIVATE_KEY=$(cat $ICON_ALGO/cache/algo_private_key) init-and-register-dbsh ../../pyteal/teal bsh $(cat $ICON_ALGO/cache/bmc_app_id) $(cat $ICON_ALGO/cache/dbsh_app_id) $ICON_BTP dbsh
ALGOD_ADDRESS=$ALGO_TEST_ADR ALGOD_TOKEN=$ALGO_TEST_TOK PRIVATE_KEY=$(cat $ICON_ALGO/cache/algo_private_key) init-and-register-wtt ../../pyteal/teal escrow $(cat $ICON_ALGO/cache/bmc_app_id) $(cat $ICON_ALGO/cache/escrow_app_id) $(cat $ICON_ALGO/cache/icon_btp_addr) wtt $(cat $ICON_ALGO/cache/algo_test_asset_id)
echo "Register new contractt"
ALGOD_ADDRESS=$(cat $ICON_ALGO/cache/algod_address) ALGOD_TOKEN=$(cat $ICON_ALGO/cache/algo_token) PRIVATE_KEY=$(cat $ICON_ALGO/cache/algo_private_key) MINTER_PRIVATE_KEY=$(cat $ICON_ALGO/cache/algo_minter_private_key) RECEIVER_PRIVATE_KEY=$(cat $ICON_ALGO/cache/algo_receiver_private_key) init-and-register-i2a ../../pyteal/teal reserve $(cat $ICON_ALGO/cache/bmc_app_id) $(cat $ICON_ALGO/cache/reserve_app_id) $(cat $ICON_ALGO/cache/icon_btp_addr) i2a $(cat $ICON_ALGO/cache/algo_wrapped_token_id)

echo "Start Algorand Link Status file with latest chain heights"
echo '{"tx_seq":0,"rx_seq":1,"rx_height":'$ICON_HEIGHT',"tx_height":'$ALGO_ROUND'}' >chain/algo/linkStatus.json

echo "Create Algorand Services Map file with dummy BSH app id"
echo '{"dbsh":'$DUMMY_BSH_APP_ID', "wtt":'$ESCROW_APP_ID', "i2a":'$RESERVE_APP_ID'}' > chain/algo/serviceMap.json

go run . -config=../../devnet/docker/icon-algorand/algo-testnet-config.json
