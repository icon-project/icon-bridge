# BTP Token Service (BTS)

## External Methods

### Owner guarded methods

#### setFeeRatio
* Description
    - Set fixed fee and fee numerator for registered coins
    - Owner guarded
    - Fee Numerator can range from 0 to 10000. (10000 means 100%)
    - BTP Fee = fixedFee + Amount * feeNumerator / 1000 
* Params
    - _name: String (Name of the coin)
    - _feeNumerator: BigInteger (Fee numerator for coin)
    - _fixedFee: BigInteger (Fixed fee for coin)

#### register
* Description
    - Registers a coin to BTP Token Service
    - Deploys a IRC2Tradeable Contract if token does not exist
    - Owner guarded
    - Can't register tokens with same name
* Params
    - _name: String (name of the coin)
    - _symbol: String (symbol of the coin)
    - _decimals: int (number of decimals)
    - _feeNumerator: BigInteger (feeNumerator, same as in setFeeRatio)
    - _fixedFee: BigInteger (fixedFee, same as in setFeeRatio)
    - _address: @Optional Address (contract address of token) 

#### setTokenLimit
* Description
    - Set token limit for tokens
    - Can set limit for not registered tokens as well
    - Initially set limit on BTS javascore
    - Sends Token Limit BTP Message to connected chains
    - Other chains pick up BTP Message, and set on their smart contract
    - Send response back to BTS if operation was successful or not
    - If response unsuccessful: reset token limit on javascore BTS
    - Owner guarded
    - Size of array should be same
    - Order of coin name and token limit in param is same
* Params
    - _coinNames: String[] (Array of coin names)
    - _tokenLimits: BigInteger[] (Array of token limits)

#### addBlacklistAddress
* Description
    - Add users to blacklist on certain networks
    - Owner guarded
    - Maintains information of all addresses blacklisted on all chains
    - Sends Blacklist BTP Message to connected chains
    - No BTP Message if blacklisted only on ICON
    - Save address in lowercase string 
* Params
    - _net: String (chain to blacklist on)
    - _addresses: String[] (list of addresses to blacklist)


#### removeBlacklistAddress
* Description
    - Remove users from blacklist on certain networks
    - Owner guarded
    - Sends Blacklist BTP Message to connected chains
* Params
    - _net: String (chain name)
    - _addresses: String[] (list of addresses remove from blacklist)

#### addOwner
* Description
    - Add owners to owner list
    - Owner can call previously mentioned owner guarded methods
* Params
    - _addr: Address (icon address of owner)

#### removeOwner
* Description
    - Remove address from owner list
    - Cannot remove deployer (score owner) from list
* Params
    - _addr: Address (icon address of owner to remove)

#### addRestrictions
* Description
    - Enables restrictions
    - Token limit and blacklist are checked of restriction is enabled

#### disableRestrictions
* Description
    - Disable restrictions
    - Token limit and blacklist checks skipped if restriction is disabled


### Contract guarded methods (only BTP Message Center)

#### handleBTPMessage
* Description
    - Handles token transfer request/response, blacklist, token limit
    - Token transfer from other chain: Handle request
    - Token transfer from ICON: Handle successful/unsccessful response incoming from destination chain
    - locking/unlocking of tokens/ICX
    - set to refundable on unsuccessful response
    - handles add/remove from blacklist based on response
    - checks networks where token limit has been added 
* Params
    - _from: String (BTP Address of source)
    - _svc: String (Service Type)
    - _sn: BigInteger (Service Number)
    - _msg: byte[] (BTS Message in bytes)

#### handleFeeGathering
* Description
    - Sends fee accumulated on BTS Contract to fee gathering address defined on BMC
    - Called once after a certain block height
* Params
    - _fa : String (fee accumulator address)
    - _svc: String (service type)

#### handleBTPError
* Description
    - For handling BTP Error Messages
    - Called if BTP Message couldn't be resolved on destination chain
    - Handled for token transfers, token limit and blacklist error response 
* Params
    - _src: String (BTP Address of source)
    - _svc: String (Service Type Name)
    - _sn: BigInteger (Service number)
    - _code: long (BTP Transaction code)
    - _msg: String (error message if any)

### Methods for users 

#### tokenFallback
* Description
    - Applicable for regular IRC2 Tokens
    - Called on IRC2 Transfer to BTS
    - Adds transferred amount to usable balance
    - Only for registered tokens
* Params
    - _from: Address (From address)
    - _value: BigInteger (Amount to send to BTS)
    - _data: byte[] (data if any)

#### transferNativeCoin
* Description
    - To transfer nativecoin (ICX) to destination chains
    - Checks for blacklist and token limit
    - Sends BTP Message for cross-chain transfer of ICX
    - lock ICX on BTS, mint equivalent on destination chain
* Params
    - _to: String (BTP Address of destination)

#### transfer
* Description
    - Applicable for IRC2 Tradable and regular IRC2 Tokens
    - Checks for blacklist and token limit
    - Called after IRC2 Tradable Tokens are approved to BTS contract
    - Called after IRC2 Tokens are transferrd to BTS contract
    - Sends BTP Message for cross-chain transfer of tokens
    - subtract value from usable balance
    - lock tokens on BTS, mint equivalent on destination chain
* Params
    - _coinName: String (name of coin)
    - _value: BigInteger (Amount to transfer)
    - _to: String (BTP Address of reciever)

#### transferBatch
* Description
    - Transfer multiple tokens/nativecoin in batch
    - Checks for blacklist and token limit
    - locks value of all tokens/nativecoin on BTS, mints on destination
* Params
    - _coinName: String[] (name of coins)
    - _value: BigInteger[] (amount to transfer for each coin in order)
    - _to: String (BTP Address of reciever)

#### reclaim
* Description
    - Reclaim tokens/nativecoin if failed to mint on destination chain
    - Reclaim IRC2 Tokens sent to BTS, which are not yet transfered to destination chain
    - Transfers amount of that coin back to user
    - If amount in BTS > _value, remaining amount set to refundable
    - Refundable amount can be claimed anytime
* Params
    - _coinName: String (name of coin)
    - _value: BigInteger (amount to reclaim)

### Readonly methods

#### name
* Description
    - Name of contract "BTP Token Service"
* Returns
    ```json
    "BTP Token Service"
    ```

#### feeRatio
* Description
    - Get fixed fee and fee numerator for a coin
* Params
    - _name: String (coinname)
* Returns 
    - Mapping of fixedFee and feeNumerator as keys and respective values
    ```json
    {
        "fixedFee": "0x100",
        "feeNumerator": "0x64"
    }
    ```
 
#### getTokenLimit
* Description
    - Get token limit for coins
    - uint(256) - 1 if not set
* Params
    - _name: String (coinname)
* Returns 
    - Token limit for coin in BigInteger
    ```json
    "0x64"
    ```
 
#### getSn
* Description
    - Get total number of cross-chain BTP Messages sent via BTS Javascore
* Returns
    - BigInteger value 
    ```json
    "0x2"
    ```

#### isUserBlackListed
* Description
    - Check if user is blacklisted on a network/chain
* Params
    - _net: String (name of network/chain)
    - _address: String (address of user)
* Returns
    - Boolean
    ```json
    "0x0"
    ```
 
#### getBlackListedUsers
* Description
    - Get all users who are blacklisted on a network within a range
* Params
    - _net: String (name of network/chain)
    - _start: int (start index)
    - _end: int (end index)
* Returns
    - List of users
    ```json
    ["hx5ed806a9a612ff7e2aa1ba67ed80a438bcdd11fa", "hx3023cfd7cc407b67fea8de216d409758f5bb10dd"]
    ```
 
#### getRegisteredTokensCount
* Description
    - Get number of registered coins on BTS (does not include nativecoin ICX)
* Returns
    - Integer
    ```json
    "0x2"
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
 
#### balanceOf
* Description
    - Returns usable, locked, refundable and user balance
    - Usable      : Amount allowed for BTS Contract to use
    - Locked      : Amount locked in ICON, but exists on other destination chains
    - Refundable  : Amount that can be refunded back to user
    - UserBalance : Current balance of the user
* Params
    - _owner: Address (user to check balance of)
    - _coinName: String (name of coin to check balance of) 
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
    - _owner: Address (user to check balance of)
    - _coinNames: String[] (list of coins to check balance of for user)
* Returns 
    - List of map retuned by balanceOf
    ```json
    [
        { 
            "usable" : "0x0",
            "locked" : "0x12345",
            "refundable" : "0x0",
            "userBalance" : "0x56a76b623167000"
        },
        {
            "usable" : "0x213aab8721000",
            "locked" : "0x0",
            "refundable" : "0x0",
            "userBalance" : "0x56a7623ab267000" 
        }
    ]
    ```
 
#### getAccumulatedFees
* Description
    - Returns accumulated fees for all coins, not yet sent to fee gathering address
* Returns
    - Map of coinName as key, and fee accumulated for that coin as value
    ```json
    {
        "btp-0x1.icon-ICX": "0xb30601a7228a0000",
        "btp-0x1.icon-bnUSD": "0x14d1120d7b160000",
        "btp-0x1.icon-sICX": "0xd87e555900180000",
        "btp-0x38.bsc-BNB": "0x0",
        "btp-0x38.bsc-BTCB": "0x0",
        "btp-0x38.bsc-BUSD": "0x0",
        "btp-0x38.bsc-ETH": "0x0",
        "btp-0x38.bsc-USDC": "0x0",
        "btp-0x38.bsc-USDT": "0x0"
    }
    ```
 
#### getOwners
* Description
    - Get list of all the owners of BTS contract
* Returns
    - List of owners of BTS Contract
    ```json
    [
        "hxdd6b62ab563ff0cad4af7248ed2f9458059a18d2",
        "hx61fa24aab5dc30d645daadeb996def021661d2a1"
    ]
    ```
 
#### isOwner
* Description
    - Check if an address if a BTS owner
* Params
    - _addr: Address (user address)
* Returns
    - Boolean
    ```json
    "0x1"
    ```
 
#### isRestrictionEnabled
* Description
    - Flag to check if restriction is enabled
    - If enabled, token limit and blacklist is checked, else skipped
 * Returns
    - Boolean
    ```json
    "0x1"
    ```

### Eventlogs

#### TransferStart
* Indexed 1
* Params
    - _from: Address (Address of user who is to transfer tokens)
    - _to: String (BTP Address of destination)
    - _sn: BigInteger (Service Number)
    - _assetDetails: byte[] (TransferTransactionData in bytes)
* Description
    - When token is to be transferred from source to destination, it is generated in source chain.


#### TransferRecieved
* Indexed 2
* Params
    - _from: String (BTP Address of BMC)
    - _to: Address (Address of reciever)
    - _sn: BigInteger (Service Number)
    - _assetDetails: byte[] (TransferTransactionData in bytes)
* Description
    - When token is recieved in destination chain, this event is thrown.
    - To indicate token transfer was successful in destination chain, and tokens has been minted 

#### TransferEnd
* Indexed 1
* Params
    - _from: Address (Address of transfer originator)
    - _sn: BigInteger (Service Number)
    - _code: BigInteger (Successful or unsuccessful)
    - _msg: bytes[] (Message of any in bytes)
* Description
    - After successful response of token transfer is recieved in the destination chain, it is generated in source chain.

#### UnknownResponse
* Indexed 1
* Params
    - _from: Address (User Address)
    - _sn: BigInteger
* Description
    - Generated if unknown service type messae is recieved

#### AddedToBlacklist
* Indexed 1
* Params
    - sn: BigInteger (Service Number)
    - bytes: byte[] (Response message if any)
* Description
    - Generated when response for successful add to blacklist is recieved

#### RemovedFromBlacklist
* Indexed 1
* Params
    - sn: BigInteger (Service Number)
    - bytes: byte[] (Response message if any)
* Description
    - Generated when a response for succesful remove from blacklist is recieved
 
#### TokenLimitSet
* Indexed 1
* Params
    - sn: BigInteger (Service NUmber)
    - bytes: byte[] (Response message if any)
* Description
    - Generated when successful response of token limit set is recieved 