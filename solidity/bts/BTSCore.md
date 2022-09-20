# BTS Core

## External Methods 

### Owner guarded methods

#### setFeeRatio
* Description
    - Set fixed fee and fee numerator for registered coins
    - Owner guarded
    - Fee Numerator can range from 0 to 10000. (10000 means 100%)
    - BTP Fee = fixedFee + Amount * feeNumerator / 1000 
* Params
    - _name: string (Name of the coin)
    - _feeNumerator: uint (Fee numerator for coin)
    - _fixedFee: uint (Fixed fee for coin)

#### register
* Description
    - Registers a coin to BTP Token Service
    - Deploys a ERC20Tradable Contract if token does not exist
    - Owner guarded
    - Can't register tokens with same name
* Params
    - _name: string (name of the coin)
    - _symbol: string (symbol of the coin)
    - _decimals: uint8 (number of decimals)
    - _feeNumerator: uint (feeNumerator, same as in setFeeRatio)
    - _fixedFee: uint (fixedFee, same as in setFeeRatio)
    - _address: address (contract address of token) 

#### updateBTSPeriphery
* Description
    - Update BTSPeriphery Contract in BTSCore
    - Owner guarded
    - Must be different than existing BTSPeriphery
    - Can be updated if no request is pending
* Params
    - _btsPeriphery: address (address of new BTSPeriphery)

#### addOwner
* Description
    - Add owners to owner list
* Params
    - _addr: address (address of owner)

#### removeOwner
* Description
    - Remove address from owner list
* Params
    - _addr: address (address of owner to remove)

### Contract guarded methods

#### handleResponseService
* Description
    - Handle response of requested service
    - Only BTSPeriphery
    - Add to refundable if error during transfer
    - Burn wrapped tokens if successfully minted on destination
* Params
    - _requester : address (An address of originator of a requested service)
    - _coinName : string (A name of requested coin)
    - _value : uint (An amount to receive on a destination chain)
    - _fee : uint (An amount of charged fee)
    - _rspCode: uint (Response code)

#### transferFees
* Description
    - Only BTSPeriphery contract
    - Sends BTP Message to send fees to fee aggregator
* Params
    - _fa: string (BTP Address of fee aggregator)

### Methods for users 

#### transferNativeCoin
* Description
    - To transfer nativecoin (BNB) to destination chains
    - Checks for blacklist and token limit 
    - Sends BTP Message for cross-chain transfer of BNB
    - lock BNB on BTS, mint equivalent on destination chain
* Params
    - _to: string (BTP address of destination)

#### transfer
* Description
    - Applicable for all tokens
    - Checks for blacklist and token limit
    - Called after tokens are transferrd to BTS contract
    - Sends BTP Message for cross-chain transfer of tokens
    - subtract value from usable balance
    - lock tokens on BTS, mint equivalent on destination chain
* Params
    - _coinName: string (name of coin)
    - _value: uint (Amount to transfer)
    - _to: string (BTP address of reciever)

#### transferBatch
* Description
    - Transfer multiple tokens/nativecoin in batch
    - Checks for blacklist and token limit
    - locks value of all tokens/nativecoin on BTS, mints on destination
* Params
    - _coinName: string[] (name of coins)
    - _value: uint[] (amount to transfer for each coin in order)
    - _to: string (BTP address of reciever)

#### reclaim
* Description
    - Reclaim tokens/nativecoin if failed to mint on destination chain
    - Reclaim tokens sent to BTS, which are not yet transfered to destination chain
    - Transfers amount of that coin back to user
    - If amount in BTS > _value, remaining amount set to refundable
    - Refundable amount can be claimed anytime
* Params
    - _coinName: string (name of coin)
    - _value: uint (amount to reclaim)

### Readonly methods

#### getNativeCoinName
* Description
    - Get registered name of nativecoin
* Returns
    - Name of nativecoin
    ```json
    "btp-0x38.bsc-BNB"
    ```
 #### getOwners
* Description
    - Get list of all the owners of BTS Core contract
* Returns
    - List of owners of BTS Core Contract
    ```json
    [
        "0xdd6b62ab563ff0cad4af7248ed2f9458059a18d2",
        "0x61fa24aab5dc30d645daadeb996def021661d2a1"
    ]
    ```
 
#### isOwner
* Description
    - Check if an address if a BTS owner
* Params
    - _addr: address (user address)
* Returns
    - Boolean
    ```json
    "0x1"
    ```
#### coinNames
* Description
    - Returns all registered coins on BTS (including ICX)
* Returns
    - List of all registered coins
    ```json
    [
    "btp-0x1.icon-ICX",
    "btp-0x1.icon-sICX",
    "btp-0x1.icon-bnUSD",
    "btp-0x38.bsc-BNB",
    "btp-0x38.bsc-BUSD",
    "btp-0x38.bsc-USDT",
    "btp-0x38.bsc-USDC",
    "btp-0x38.bsc-BTCB",
    "btp-0x38.bsc-ETH"
    ]
    ```
 
#### coinId
* Description
    - Returns contract address for coin
* Params
    - _coinName: (name of coin)
* Returns
    - Contract address of token
    ```json
    "cxeb14f67eeb6742c7a1ff474308ce82d874469703"
    ```

#### isValidCoin
* Description
    - Check if a coin is registered on solidity BTSCore
* Params
    - _coinName: (name of coin)
* Returns
    - Boolean
    - If coin is registered
    ```json
    "0x1"
    ```
#### feeRatio
* Description
    - Get fixed fee and fee numerator for a coin
* Params
    - _name: string (coinname)
* Returns 
    - Mapping of fixedFee and feeNumerator as keys and respective values
    ```json
    (
        "fixedFee": "0x100",
        "feeNumerator": "0x64"
    )
    ```

#### balanceOf
* Description
    - Returns usable, locked, refundable and user balance
    - Usable      : Amount allowed for BTS Contract to use
    - Locked      : Amount locked in ICON, but exists on other destination chains
    - Refundable  : Amount that can be refunded back to user
    - UserBalance : Current balance of the user
* Params
    - _owner: address (user to check balance of)
    - _coinName: string (name of coin to check balance of) 
* Returns
    - Map of usable, locked, refundable and userBalance and corresponding value
    ```json
    {
        "usable" : "0x0",
        "locked" : "0x12345",
        "refundable" : "0x0",
        "userBalance" : "0x56a76b623167000"
    }
    ```
 
#### balanceOfBatch
* Description
    - balanceOf, but for multiple coins at once
* Params
    - _owner: address (user to check balance of)
    - _coinNames: string[] (list of coins to check balance of for user)
* Returns 
    - List of map retuned by balanceOf
    ```json
    [
        ( 
            "usable" : "0x0",
            "locked" : "0x12345",
            "refundable" : "0x0",
            "userBalance" : "0x56a76b623167000"
        ),
        (
            "usable" : "0x213aab8721000",
            "locked" : "0x0",
            "refundable" : "0x0",
            "userBalance" : "0x56a7623ab267000" 
        )
    ]
    ```
 
#### getAccumulatedFees
* Description
    - Returns accumulated fees for all coins, not yet sent to fee gathering address
* Returns
    - Map of coinName as key, and fee accumulated for that coin as value
    ```json
    (
        "btp-0x1.icon-ICX": "0xb30601a7228a0000",
        "btp-0x1.icon-bnUSD": "0x14d1120d7b160000",
        "btp-0x1.icon-sICX": "0xd87e555900180000",
        "btp-0x38.bsc-BNB": "0x0",
        "btp-0x38.bsc-BTCB": "0x0",
        "btp-0x38.bsc-BUSD": "0x0",
        "btp-0x38.bsc-ETH": "0x0",
        "btp-0x38.bsc-USDC": "0x0",
        "btp-0x38.bsc-USDT": "0x0"
    )
    ```