# BTSPeriphery

## External Methods

### Contract guarded

#### addToBlacklist
* Description
    - Add users to blacklist
    - Only BTSPeriphery, made external to use try/catch
    - converts string into address format before saving
* Params
    | Parameter | Type     | Description                      |
    |:----------|:---------|:---------------------------------|
    | _address  | string[] | list of addresses in string form |

#### removeFromBlacklist
* Description
    - Remove users from blacklist
    - Only BTSPeriphery, made external to use try/catch
    - converts string into address format
* Params
    | Parameter | Type     | Description                      |
    |:----------|:---------|:---------------------------------|
    | _address  | string[] | list of addresses in string form |

#### setTokenLimit
* Description
    - Set token limit on tokens, registered or not
    - Only BTSPeriphery, made external to use try/catch
* Params
    | Parameter    | Type     | Description                |
    |:-------------|:---------|:---------------------------|
    | _coinNames   | string[] | Array of names of the coin |
    | _tokenLimits | uint[]   | Token limit for coins      |

#### sendServiceMessage 
* Description
    - Only BTS Core
    - Generates required message and call sendMessage of BMC Contract to send message cross chain
    - Updates serial number (sn)
* Params
    | Parameter  | Type      | Description                |
    |:-----------|:----------|:---------------------------|
    | _to        | string    | BTP Address of destination |
    | _from      | address   | from address               |
    | _coinNames | string[]  | array of names of coin     |
    | _values    | uint256[] | array of value to send     |
    | _fees      | uint256[] | array of fee for each coin |

#### handleBTPMessage
* Description
    - Handles token transfer request/response, blacklist, token limit
    - Token transfer from other chain (handleRequestService)
    - Token transfer from BSC (Handle successful/unsccessful response incoming from destination chain)
* Params
    | Parameter | Type   | Description           |
    |:----------|:-------|:----------------------|
    | _from     | string | BTP Address of source |
    | _svc      | string | Service Type          |
    | _sn       | uint   | Service Number        |
    | _msg      | bytes  | BTS Message in bytes  |

#### handleRequestService
* Description
    - Only BTSPeriphery, made external to use try/catch
    - Request token transfer/mint on this chain
* Params
    | Parameter | Type          | Description                                 |
    |:----------|:--------------|:--------------------------------------------|
    | _to       | string        | BTP Address of this chain to recieve tokens |
    | _assets   | Types.Asset[] | Asset details, coinNames and amount         |

#### handleResponseService
* Description
    - Only BTSPeriphery, made external to use try/catch
    - Calls handleResponseService on btsCore
* Params
    | Parameter | Type  | Description           |
    |:----------|:------|:----------------------|
    | _sn       | uint  | Service Number        |
    | _code     | uint  | Code for error/sucess |
    | _msg      | bytes | Message in bytes      |

#### handleBTPError
* Description
    - For handling BTP Error Messages
    - Called if BTP Message couldn't be resolved on destination chain
    - Handled for token transfers
    - Only BMC
* Params
    | Parameter | Type   | Description           |
    |:----------|:-------|:----------------------|
    | _src      | string | BTP Address of source |
    | _svc      | string | Service Type Name     |
    | _sn       | uint   | Service number        |
    | _code     | long   | BTP Transaction code  |
    | _msg      | string | error message if any  |

#### handleFeeGathering
* Description
    - Collects fee accumulated for all native tokens / nativecoins on this chain
    - Only BMC
* Params
    | Parameter | Type   | Description             |
    |:----------|:-------|:------------------------|
    | _fa       | string | fee accumulator address |
    | _svc      | string | service type            |


### General Methods
#### checkTransferRestrictions
* Description
    - Check for transfer restrictions, blacklist and token limit
    - Check if user is blacklisted on this chain
    - Check if that amount can be transferred
* Params
    | Parameter | Type    | Description                                              |
    |:----------|:--------|:---------------------------------------------------------|
    | _coinName | string  | Name of coin                                             |
    | _user     | address | Address of user to check if it is blacklisted            |
    | _value    | uint    | Amount to check if it exceeded token limit for that coin |

## Readonly Methods
#### hasPendingRequest
* Description
    - Check if any request is pending
* Returns
    - If any request is pending
        ```json
        "0x0"
        ```

## Eventlogs

#### TransferStart
* Params
    | Parameter     | Type             | Description                               |
    |:--------------|:-----------------|:------------------------------------------|
    | _from         | address(indexed) | Address of user who is to transfer tokens |
    | _to           | string           | BTP address of destination                |
    | _sn           | uint             | Service Number                            |
    | _assetDetails | bytes            | TransferTransactionData in bytes          |
* Description
    - When token is to be transferred from source to destination, it is generated in source chain.


#### TransferRecieved
* Params
    | Parameter     | Type             | Description                      |
    |:--------------|:-----------------|:---------------------------------|
    | _from         | string(indexed)  | BTP address of BMC               |
    | _to           | address(indexed) | Address of reciever              |
    | _sn           | uint             | Service Number                   |
    | _assetDetails | bytes            | TransferTransactionData in bytes |
* Description
    - When token is recieved in destination chain, this event is thrown.
    - To indicate token transfer was successful in destination chain, and tokens has been minted 

#### TransferEnd
* Params
    | Parameter | Type             | Description                    |
    |:----------|:-----------------|:-------------------------------|
    | _from     | address(indexed) | Address of transfer originator |
    | _sn       | uint             | Service Number                 |
    | _code     | uint             | Successful or unsuccessful     |
    | _msg      | bytes            | Message of any in bytes        |
* Description
    - After successful response of token transfer is recieved in the destination chain, it is generated in source chain.

#### UnknownResponse
* Params
    | Parameter | Type    | Description   |
    |:----------|:--------|:--------------|
    | _from     | address | User Address  |
    | _sn       | uint    | Serial Number |
* Description
    - Generated if unknown service type messae is recieved