# ICON BTS Contract

## Requirements
* #### goloop-cli
    - This tool will be used to interact with the ICON BTS Contract from shell.
    - Alternative: [ICON sdk](https://docs.icon.community/icon-stack/client-apis)
    - Installation
        ```sh
        go install github.com/icon-project/goloop/cmd/goloop@latest 
        ```
* #### ICON Wallet
    - [Hana Wallet](https://hanawallet.io/)
    - Get keystore and password for wallet and load funds
    - [Testnet Faucet](https://faucet.iconosphere.io/)

* #### NID
    - NID for Mainnet: 0x1
    - NID for Testnets
        - Lisbon: 0x2
        - Berlin: 0x7
* #### ENDPOINTS
    - Mainnet: https://ctz.solidwallet.io/api/v3
    - Testnet
        - Lisbon: https://lisbon.net.solidwallet.io/api/v3
        - Berlin: https://berlin.net.solidwallet.io/api/v3

* #### BTP Address
    - BTP Address is a format for address used in icon-bridge
    
    

## API

* ### Balance
    - #### Method
        - balanceOf
    - #### Parameters
        | Parameters | Type    | Info               |
        |:-----------|:--------|:-------------------|
        | _owner     | Address | Valid ICON Address |
        | _coinName  | string  | Name of coin       |

    -   #### CLI Command
        ```sh
        goloop rpc call --uri <ENDPOINT> \
            --to <BTS> \
            --method balanceOf \
            --param _owner=<Address> \
            --param _coinName=<Name of coin>

        ```
    - #### Returns
        ```json
        {
            "usable" : "0x0",
            "locked" : "0x12345",
            "refundable" : "0x0",
            "userBalance" : "0x56a76b623167000"
        }
        ```

* ### Fees
    - #### Method
        - feeRatio
        
    - #### Parameters
        | Parameters | Type   | Info         |
        |:-----------|:-------|:-------------|
        | _name      | string | Name of coin |

    - #### CLI Command
        ```sh
        goloop rpc call --uri <ENDPOINT> \
            --to <BTS> \
            --method feeRatio \
            --param _name=<Name of coin>
        ```
    - #### Description
        - When `value` amount is to be transferred cross chain
        - `fee = (value * feeNumerator / FEE_DENOMINATOR) + fixedFee`
        - where, FEE_DENOMINATOR = 10000
        - The amount that gets transferred is: `value - fee`

    - #### Returns
        ```json
        {
            "fixedFee": "0x56a76b623167000",
            "feeNumerator": "0x0"
        }
        ```
        
* ### Token Transfer to BTS
    - #### Method
        - transfer
    - #### Parameters
        | Parameters | Type       | Info                     |
        |:-----------|:-----------|:-------------------------|
        | _to        | Address    | ICON Address of reciever |
        | _value     | BigInteger | Amount to transfer       |


    - #### CLI Command
        ```sh
        goloop rpc sendtx call --uri <ENDPOINT> \
            --to <IRC2 Token Address> \
            --key_store <KEYSTORE_FILE> \
            --key_secret <KEYSTORE_SECRET> \
            --nid <NID> \
            --method transfer \
            --param _to=<BTS> \
            --param _value=<Amount>
        ```

* ### Wrapped Tokens Approvals
    - #### Method
        - approve
    - #### Parameters
        | Parameters | Type       | Info                    |
        |:-----------|:-----------|:------------------------|
        | spender    | Address    | ICON Address of spender |
        | amount     | BigInteger | Amount to approve       |
    - #### CLI Command
        ```sh
        goloop rpc sendtx call --uri <ENDPOINT> \
            --to <Wrapped Token Address> \
            --key_store <KEYSTORE_FILE> \
            --key_secret <KEYSTORE_SECRET> \
            --nid <NID> \
            --method approve \
            --param spender=<BTS> \
            --param amount=<Amount>
        ```

* ### Transfer Tokens Cross Chain
    - #### Method
        - transfer
    
    - #### Parameters
        | Parameters | Type       | Info                       |
        |:-----------|:-----------|:---------------------------|
        | _coinName  | String     | Name of coin to transfer   |
        | _to        | String     | BTP Address of destination |
        | _value     | BigInteger | Amount to transfer         |

    - #### CLI Command
        ```sh
        goloop rpc sendtx call --uri <ENDPOINT> \
            --to <BTS> \
            --key_store <KEYSTORE_FILE> \
            --key_secret <KEYSTORE_SECRET> \
            --nid <NID> \
            --method transfer \
            --param _coinName=<Name of Coin> \
            --param _to=<BTP Address of destination> \
            --param _value=<Amount>
        ```

* ### Transfer Nativecoin Cross Chain
    - #### Method
        - transferNativeCoin
    - #### Parameters
        | Parameters | Type   | Info                       |
        |:-----------|:-------|:---------------------------|
        | _to        | String | BTP Address of destination |
    - #### CLI Command

        ```sh
        goloop rpc sendtx call --uri <ENDPOINT> \
            --to <BTS> \
            --key_store <KEYSTORE_FILE> \
            --key_secret <KEYSTORE_SECRET> \
            --nid <NID> \
            --method transferNativeCoin \
            --value <Amount> \
            --param _to=<BTP Address of destination>
        ```

* ### TransferBatch Cross Chain
    - #### Method
        - transferBatch
    - #### Parameters
        | Parameters  | Type         | Info                                     |
        |:------------|:-------------|:-----------------------------------------|
        | _coinName[] | String[]     | Array of names of coin to transfer       |
        | _value[]    | BigInteger[] | Amount to transfer on order of coinNames |
        | _to         | String       | BTP Address of destination               |

    - #### CLI Command

        ```sh
        goloop rpc sendtx call --uri <ENDPOINT> \
            --to <BTS> \
            --key_store <KEYSTORE_FILE> \
            --key_secret <KEYSTORE_SECRET> \
            --nid <NID> \
            --method transferBatch \
            --value <Amount> \ # 0 if nativecoin is not to be transferred
            --raw "{\"params\":{\"_coinNames\":\"<Array of coinNames>\",\"_values\":\"<Array of values>\"}}"
        ```

## Usage

* ### To transfer ICON native IRC2 Tokens
    - The token must be registered in ICON BTS and BTS of destination chain
    - Transfer the token to ICON BTS Contract [here](#token-transfer-to-bts)
    - Call transfer method of BTS Contract with following parameters [here](#transfer-tokens-cross-chain)
        
* ### To transfer wrapped tokens of other chains on ICON to other chains
    - The token must be registered in BTS of ICON and destination chain
    - Approve amount to transfer to ICON BTS Contract [here](#token-transfer-to-bts)
    - Call transfer method of BTS Contract with following parameters [here](#transfer-tokens-cross-chain)
    
* ### To transfer nativecoin(ICX) to other chains
    - Nativecoin is registered by default on ICON, should be registered on destination chain
    - Call transferNativeCoin method of BTS contract and send required amount of ICX to transfer [here](#transfer-nativecoin-cross-chain)

* ### Transfer multiple types of coins at once (TransferBatch)
    - For IRC2 Tokens, transfer to BTS Contract as above.
    - For wrapped tokens, approve to BTS Contract as above.
    - For nativecoin, send required amount of ICX to transfer in transferBatch method
    - Call transferBatch method with array of coinNames and amount to transfer [here](#transferbatch-cross-chain)
* ### Reclaim
    - To reclaim balance of user that BTS holds
    - Check refundable balance of user [here](#balance)
    - Method
        - reclaim
    - Parameters
        | Parameters | Type       | Info             |
        |:-----------|:-----------|:-----------------|
        | _coinName  | String     | Name of coin     |
        | _value     | BigInteger | Amount to redeem |
    - CLI Command
        ```sh
        goloop rpc sendtx call --uri <ENDPOINT> \
            --to <BTS> \
            --key_store <KEYSTORE_FILE> \
            --key_secret <KEYSTORE_SECRET> \
            --nid <NID> \
            --method reclaim \
            --params _coinName=<Name of Coin> \
            --params _value=<Amount>
        ``` 

* ### Blacklist users on any chain
    - Can blacklist a batch of addresses on any chain at once
    - Owner guarded
    - #### Method
        - addBlacklistAddress
    - #### Parameters
        | Parameters | Type     | Info                                    |
        |:-----------|:---------|:----------------------------------------|
        | _net       | String   | Network to blacklist on                 |
        | _addresses | String[] | Array of addresses to blacklist on _net |
    - #### CLI Command
        ```sh
        goloop rpc sendtx call --uri <ENDPOINT> \
            --to <BTS> \
            --key_store <KEYSTORE_FILE> \
            --key_secret <KEYSTORE_SECRET> \
            --nid <NID> \
            --method addBlacklistAddress \
            --raw "{\"params\":{\"_net\":\"<network name>\",\"_addresses\":\"<Address Array>\"}}"
        ``` 
    - #### Implementation
        - [here](/devnet/docker/icon-bsc/scripts/blacklist.sh)


* ### Remove users from blacklist on any chain
    - Can remove a batch of addresses from blacklist on any chain at once
    - Owner guarded
    - #### Method
        - removeBlacklistAddress
    - #### Parameters
        | Parameters | Type     | Info                                    |
        |:-----------|:---------|:----------------------------------------|
        | _net       | String   | Network to blacklist on                 |
        | _addresses | String[] | Array of addresses to blacklist on _net |
    - #### CLI Command
        ```sh
        goloop rpc sendtx call --uri <ENDPOINT> \
            --to <BTS> \
            --key_store <KEYSTORE_FILE> \
            --key_secret <KEYSTORE_SECRET> \
            --nid <NID> \
            --method removeBlacklistAddress \
            --raw "{\"params\":{\"_net\":\"${net}\",\"_addresses\":\"<Address Array>\"}}"
        ``` 
    - #### Implementation
        - [here](/devnet/docker/icon-bsc/scripts/blacklist.sh)


* ### Set token limit
    - To set token limit for tokens (registered or not) owner calls setTokenLimit method
    - Can set limit for multiple coins across all connected networks 
    - Owner guarded
    - #### Method
        - setTokenLimit
    - #### Parameters
        | Parameters   | Type         | Info                               |
        |:-------------|:-------------|:-----------------------------------|
        | _coinNames   | String[]     | Array of coinNames                 |
        | _tokenLimits | BigInteger[] | Array of tokenLimits for coinNames |
    - #### CLI Command
         ```sh
        goloop rpc sendtx call --uri <ENDPOINT> \
            --to <BTS> \
            --key_store <KEYSTORE_FILE> \
            --key_secret <KEYSTORE_SECRET> \
            --nid <NID> \
            --method setTokenLimit \
            --raw "{\"params\":{\"_coinNames\":\"<Array of coinNames>\",\"_addresses\":\"<Array of tokenLimits>\"}}"
        ``` 
    - #### Implementation
        - [here](/devnet/docker/icon-bsc/scripts/tokenLimit.sh)


