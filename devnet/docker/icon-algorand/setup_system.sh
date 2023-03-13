echo "Exporting path variable"
export PATH=$PATH:~/go/bin

echo "Creating the cache directory"
if [ -d "cache" ]; then
rm -r cache
fi
mkdir cache


echo "Getting goloop NID"
docker exec goloop goloop chain ls | jq -r '.[0] | .nid' >cache/nid
sleep 2

echo "Setting up environment variables for the Algorand node and key management daemon"

ALGOD_ADDRESS=http://localhost:4001
ALGOD_TOKEN=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
KMD_ADDRESS=http://localhost:4002
KMD_TOKEN=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa

echo $ALGOD_ADDRESS >cache/algod_address
echo $KMD_ADDRESS >cache/kmd_address
echo $ALGOD_TOKEN >cache/algo_token

kmd -d /tmp/testnet/Node/kmd-v0.5 &
sleep 5

echo "Generating the icon keystore and transferring 2001 to it"
goloop ks gen --out icon.keystore.json
KS_ADDRESS=$(cat icon.keystore.json | jq -r '.address')
PASSWORD=$(docker exec goloop cat goloop.keysecret)
docker exec goloop goloop rpc sendtx transfer \
    --uri http://localhost:9080/api/v3/icon \
    --nid $(cat cache/nid) \
    --step_limit=3000000000 \
    --key_store goloop.keystore.json --key_password $PASSWORD \
    --to $KS_ADDRESS --value=2001

echo "Update the config file with the new keystore"
jq --slurpfile new_contents icon.keystore.json '.relays[0].dst.key_store = $new_contents[0]' algo-config.json >tmpfile && mv tmpfile algo-config.json

echo "Deploy Algorand Test Asset"
MINTER_PRIVATE_KEY=$(KMD_ADDRESS=$(cat cache/kmd_address) KMD_TOKEN=$(cat cache/algo_token) kmd-extract-private-key 2)
ASA_ID=$(
    PRIVATE_KEY=$MINTER_PRIVATE_KEY ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) deploy-asset 1000000000000 6 TABC Test AB Coin http://example.com/ abcd
)
echo $MINTER_PRIVATE_KEY >cache/algo_minter_private_key
echo $ASA_ID >cache/algo_test_asset_id

echo "Deploying BMC contract to the ICON network"
CONTRACT=../../../javascore/bmc/build/libs/bmc-optimized.jar
DEPLOY_TXN_ID=$(goloop rpc sendtx deploy $CONTRACT \
    --uri http://localhost:9080/api/v3/icon \
    --key_store ./icon.keystore.json --key_password gochain \
    --nid $(cat cache/nid) \
    --content_type application/java \
    --step_limit 3000000000 \
    --param _net="$(cat cache/nid).icon")
./../../algorand/scripts/wait_for_transaction.sh $DEPLOY_TXN_ID scoreAddress >cache/icon_bmc_addr

echo "Deploying dummyBSH contract to the ICON network"
CONTRACT=../../../javascore/dummyBSH/build/libs/dummyBSH-optimized.jar
TXN_ID=$(
    goloop rpc sendtx deploy $CONTRACT \
        --uri http://localhost:9080/api/v3/icon \
        --key_store icon.keystore.json --key_password gochain \
        --nid $(cat cache/nid) --step_limit 10000000000 \
        --content_type application/java \
        --param _bmc=$(cat cache/icon_bmc_addr) \
        --param _to=0x14.algo
)
./../../algorand/scripts/wait_for_transaction.sh $TXN_ID scoreAddress >cache/icon_dbsh_addr

echo "Deploying WrappedToken contract to the ICON network"
CONTRACT=../../../javascore/wrapped-token/build/libs/WrappedToken-optimized.jar
TXN_ID=$(
    goloop rpc sendtx deploy $CONTRACT \
        --uri http://localhost:9080/api/v3/icon \
        --key_store icon.keystore.json --key_password gochain \
        --nid $(cat cache/nid) --step_limit 10000000000 \
        --content_type application/java \
        --param _bmc=$(cat cache/icon_bmc_addr) \
        --param _to=0x14.algo \
        --param _asaId=$ASA_ID \
        --param _name="Wrapped Test Token" \
        --param _symbol="WTT" \
        --param _decimals=6
)
./../../algorand/scripts/wait_for_transaction.sh $TXN_ID scoreAddress >cache/icon_wtt_addr

echo "Registering dBSH to the BMC contract on the ICON network"
TXN_ID=$(goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) \
    --method addService \
    --value 0 \
    --param _addr=$(cat cache/icon_dbsh_addr) \
    --param _svc="dbsh" \
    --step_limit=3000000000 \
    --uri http://localhost:9080/api/v3/icon \
    --key_store icon.keystore.json --key_password gochain \
    --nid=$(cat cache/nid))
./../../algorand/scripts/wait_for_transaction.sh $TXN_ID

echo "Registering Wrapped Token to the BMC contract on the ICON network"
TXN_ID=$(
    goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) \
        --method addService \
        --value 0 \
        --param _addr=$(cat cache/icon_wtt_addr) \
        --param _svc="wtt" \
        --step_limit=3000000000 \
        --uri http://localhost:9080/api/v3/icon \
        --key_store icon.keystore.json --key_password gochain \
        --nid=$(cat cache/nid)
)
./../../algorand/scripts/wait_for_transaction.sh $TXN_ID

echo "Extracting the private key for KMD deployer"
PRIVATE_KEY=$(KMD_ADDRESS=$KMD_ADDRESS KMD_TOKEN=$KMD_TOKEN kmd-extract-private-key 1)

echo "Deploying Algorand contracts"
BMC_TX_ID=$(PRIVATE_KEY=$PRIVATE_KEY ALGOD_ADDRESS=$ALGOD_ADDRESS ALGOD_TOKEN=$ALGOD_TOKEN deploy-contract ../../../pyteal/teal/bmc)
DUMMY_BSH_TX_ID=$(PRIVATE_KEY=$PRIVATE_KEY ALGOD_ADDRESS=$ALGOD_ADDRESS ALGOD_TOKEN=$ALGOD_TOKEN deploy-contract ../../../pyteal/teal/bsh)
ESCROW_TX_ID=$(PRIVATE_KEY=$PRIVATE_KEY ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) deploy-contract ../../../pyteal/teal/escrow)

echo "Setting BMC_APP_ID"
BMC_APP_ID=$(ALGOD_ADDRESS=$ALGOD_ADDRESS ALGOD_TOKEN=$ALGOD_TOKEN get-app-id $BMC_TX_ID)

echo "Setting DUMMY_BSH_APP_ID"
DUMMY_BSH_APP_ID=$(ALGOD_ADDRESS=$ALGOD_ADDRESS ALGOD_TOKEN=$ALGOD_TOKEN get-app-id $DUMMY_BSH_TX_ID)

echo "Setting ESCROW_APP_ID"

ESCROW_APP_ID=$(ALGOD_ADDRESS=$(cat cache/algod_address) ALGOD_TOKEN=$(cat cache/algo_token) get-app-id $ESCROW_TX_ID)

echo "Updating config file with the new key"
jq --arg PRIVATE_KEY "$PRIVATE_KEY" '.relays[1].dst.key_store.id=$PRIVATE_KEY' algo-config.json >config_tmp.json
mv config_tmp.json algo-config.json

echo "Updating cache files"
echo $PRIVATE_KEY >cache/algo_private_key
printf '%s' "$BMC_APP_ID" >cache/bmc_app_id
printf '%s' "$DUMMY_BSH_APP_ID" >cache/dbsh_app_id
printf '%s' "$ESCROW_APP_ID" >cache/escrow_app_id

echo "Getting algo_btp_addr"
goal app info --app-id $BMC_APP_ID -d /tmp/testnet/Node |
    awk -F ':' '/Application account:/ {gsub(/^[[:space:]]+|[[:space:]]+$/,"",$2); \
          print "btp://0x14.algo/" $2}' >cache/algo_btp_addr

echo "Setting link to algo btp on icon bmc"
LINK_TXN_ID=$(goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) \
    --method addLink --param _link=$(cat cache/algo_btp_addr) \
    --key_store ./icon.keystore.json --key_password gochain \
    --nid $(cat cache/nid) --step_limit 3000000000 --uri http://localhost:9080/api/v3/icon)

./../../algorand/scripts/wait_for_transaction.sh $LINK_TXN_ID

echo "Getting latest algo round"
goal node lastround -d /tmp/testnet/Node >cache/algo_last_round

echo "Setting round on icon bmc"
ADD_ROUND=$(goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) --method setLinkRxHeight \
    --param _link=$(cat cache/algo_btp_addr) --param _height=$(cat cache/algo_last_round) \
    --key_store ./icon.keystore.json --key_password gochain \
    --nid $(cat cache/nid) --step_limit 3000000000 --uri http://localhost:9080/api/v3/icon)

./../../algorand/scripts/wait_for_transaction.sh $ADD_ROUND

echo "Adding relay to icon bmc"
ADD_RELAY=$(goloop rpc sendtx call --to $(cat cache/icon_bmc_addr) --method addRelay \
    --param _link=$(cat cache/algo_btp_addr) \
    --param _addr=$(cat icon.keystore.json | jq -r .address) \
    --key_store ./icon.keystore.json --key_password gochain \
    --nid $(cat cache/nid) --step_limit 3000000000 --uri http://localhost:9080/api/v3/icon)

./../../algorand/scripts/wait_for_transaction.sh $ADD_RELAY

cd ../../../cmd/iconbridge

ICON_ALGO="../../devnet/docker/icon-algorand"

ICON_BTP="btp://$(cat "$ICON_ALGO/cache/nid").icon/$(cat "$ICON_ALGO/cache/icon_bmc_addr")"
echo $ICON_BTP >$ICON_ALGO/cache/icon_btp_addr

ALGO_BTP=$(cat "$ICON_ALGO/cache/algo_btp_addr")
ALGO_ROUND=$(cat "$ICON_ALGO/cache/algo_last_round" | tr -d '\n')

echo "Getting icon last block height"
ICON_HEIGHT=$(goloop rpc lastblock --uri http://localhost:9080/api/v3/icon | jq .height)
ICON_HEIGHT_HEX=$(printf "%x" $((ICON_HEIGHT - 1)))

echo "Building icon validators"
go build -o ../iconvalidators ./../iconvalidators
VALIDATORS=$(URI=http://localhost:9080/api/v3/icon HEIGHT=0x$ICON_HEIGHT_HEX ./../iconvalidators/iconvalidators | jq -r .hash)

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
    .relays[1].src.options.verifier.blockHeight=$ICON_HEIGHT' $ICON_ALGO/algo-config.json >config_tmp.json
mv config_tmp.json $ICON_ALGO/algo-config.json

echo "Register Algorand BSHs"
ALGOD_ADDRESS=$(cat $ICON_ALGO/cache/algod_address) ALGOD_TOKEN=$(cat $ICON_ALGO/cache/algo_token) PRIVATE_KEY=$(cat $ICON_ALGO/cache/algo_private_key) init-and-register-dbsh ../../pyteal/teal bsh $(cat $ICON_ALGO/cache/bmc_app_id) $(cat $ICON_ALGO/cache/dbsh_app_id) $ICON_BTP dbsh
ALGOD_ADDRESS=$(cat $ICON_ALGO/cache/algod_address) ALGOD_TOKEN=$(cat $ICON_ALGO/cache/algo_token) PRIVATE_KEY=$(cat $ICON_ALGO/cache/algo_private_key) init-and-register-wtt ../../pyteal/teal escrow $(cat $ICON_ALGO/cache/bmc_app_id) $(cat $ICON_ALGO/cache/escrow_app_id) $(cat $ICON_ALGO/cache/icon_btp_addr) wtt $(cat $ICON_ALGO/cache/algo_test_asset_id)

echo "Start Algorand Link Status file with latest chain heights"
echo '{"tx_seq":0,"rx_seq":1,"rx_height":'$ICON_HEIGHT',"tx_height":'$ALGO_ROUND'}' >chain/algo/linkStatus.json

echo "Create Algorand Services Map file with dummy BSH app id"
echo '{"dbsh":'$DUMMY_BSH_APP_ID', "wtt":'$ESCROW_APP_ID'}' > chain/algo/serviceMap.json


