# Javascore BMC

## External methods

#### addOwner
* Description
  - Add owners to owner list
  - Owner can call owner guarded methods
* Params
    | Parameters | Type    | Description                          |
    |:-----------|:--------|:-------------------------------------|
    | _owner     | Address | Admin/Owner address to access on BMC |

#### removeOwner

* Description
    - Remove owners to owners list
    - Cannot remove the last owner from the list
* Params
    | Parameters | Type    | Description       |
    |:-----------|:--------|:------------------|
    | _owner     | Address | Address to remove |

  
#### addService
* Description:
    - It registers the smart contract for the service.
    - It's called by the owner/admin to manage the BTP network.
* Params
    | Parameters | Type    | Description                                            |
    |:-----------|:--------|:-------------------------------------------------------|
    | _svc       | String  | the name of the service                                |
    | _addr      | Address | the address of the smart contract handling the service |

#### addServiceCandidate
* Description:
  - Add Service Candidate on the given service
* Params
    | Parameters | Type    | Description                        |
    |:-----------|:--------|:-----------------------------------|
    | _svc       | String  | the name of the service            |
    | _addr      | Address | Address to be added on the service |

#### removeServiceCandidate
* Description:
  * Remove Service Candidate on the given service
* Params
    | Parameters | Type    | Description                        |
    |:-----------|:--------|:-----------------------------------|
    | _svc       | String  | the name of the service            |
    | _addr      | Address | Address to be added on the service |

#### removeService
* Description:
    - It de-registers the smart contract for the service.
    - It's called by the operator to manage the BTP network.
* Params
    | Parameters | Type   | Description             |
    |:-----------|:-------|:------------------------|
    | _svc       | String | the name of the service |

#### handleFragment
* Params
    | Parameters | Type    | Description |
    |:-----------|:--------|:------------|
    | _prev      | String  |             |
    | _msg       | String  |             |
    | _idx       | Integer |             |

#### addLink
* Description
    - If it generates the event related with the link, the relay shall
      handle the event to deliver BTP Message to the BMC.
    - If the link is already registered, or its network is already
      registered then it fails.
    - It initializes status information for the link.
    - It's called by the operator to manage the BTP network.
* Params
    | Parameters | Type   | Description                  |
    |:-----------|:-------|:-----------------------------|
    | _link      | String | BTP Address of connected BMC |

#### setLinkRxHeight
* Description
  - Modify the blockHeight on given link
  - Only owner can manage this on BTP Network
* Params
    | Parameters | Type   | Description          |
    |:-----------|:-------|:---------------------|
    | _link      | String | Added Link           |
    | _height    | Long   | Starting BlockHeight |

#### setLinkRotateTerm
* Description
  - Update the blockHeight, blockInterval and maxAggreation on given link
  - Only owner can manage this on BTP Network
* Params
  | Parameters      | Type    | Description           |
  |:----------------|:--------|:----------------------|
  | _link           | String  | Added Link            |
  | _block_interval | Integer | Block interval        |
  | _max_agg        | Integer | Max Aggregation Value |
  
#### setLinkDelayLimit
* Description
  - Update the `delayLimit` value on given link
  - Only owner can manage this on BTP Network
* Params
    | Parameters | Type    | Description    |
    |:-----------|:--------|:---------------|
    | _link      | String  | Added Link     |
    | _value     | Integer | Block interval |

#### setLinkSackTerm
* Description
  - Update the `sackTerm` value on given link
  - Only owner can manage this on BTP Network
* Params
    | Parameters | Type    | Description    |
    |:-----------|:--------|:---------------|
    | _link      | String  | Added Link     |
    | _value     | Integer | Block interval |


#### removeLink
* Description
    - It removes the link and status information.
    - It's called by the operator to manage the BTP network.
* Params
    | Parameters | Type   | Description                  |
    |:-----------|:-------|:-----------------------------|
    | _link      | String | BTP Address of connected BMC |

#### addRoute
* Description:
    - Add route to the BMC.
    - It may fail if there are more than one BMC for the network.
    - It's called by the operator to manage the BTP network.
* Params
    | Parameters | Type   | Description                                     |
    |:-----------|:-------|:------------------------------------------------|
    | _dst       | String | BTP Address of the destination BMC              |
    | _link      | String | BTP Address of the next BMC for the destination |

#### removeRoute
* Description:
    - Remove route to the BMC.
    - It's called by the operator to manage the BTP network.
* Params
    | Parameters | Type   | Description                        |
    |:-----------|:-------|:-----------------------------------|
    | _dst       | String | BTP Address of the destination BMC |

#### addRelay
* Description
  - Add relayer address on given link ( BTPAddress )
  - Only owner can call this method
* Params
    | Parameters | Type    | Description                 |
    |:-----------|:--------|:----------------------------|
    | _link      | String  | BTP Address of the next BMC |
    | _addr      | Address | Relayer Address             |

#### removeRelay
* Description
  * Remove relayer address on given link (BTPAddress)
  * Only owner can call this method
* Params
    | Parameters | Type    | Description                 |
    |:-----------|:--------|:----------------------------|
    | _link      | String  | BTP Address of the next BMC |
    | _addr      | Address | Relayer Address             |

#### distributeRelayerReward
* Description
  - Distribute relayers' reward to each relayer of registered services

#### claimRelayerReward
* Description
  - Claim the distributed relayer rewards from relayer address

#### setRelayerMinBond
* Description
  - Set Minimum relayer bond while registering as relayer
  - Only Owner can call this method
* Params
    | Parameters | Type       | Description                                  |
    |:-----------|:-----------|:---------------------------------------------|
    | _value     | BigInteger | Value to be set as Minimum Bond as a Relayer |

#### setFeeGatheringTerm
  * Description
    - Set fee feeCollecting BlockHeight
    - Only owner can set the value
  * Params
    | Parameters | Type | Description |
    |:-----------|:-----|:------------|
    | _value     | Long |             |

#### setFeeAggregator
  * Description
    - Set fee feeCollecting address
    - Only owner can set the value
  * Params
    | Parameters | Type    | Description |
    |:-----------|:--------|:------------|
    | _addr      | Address |             |

#### sendFeeGathering
* Description
  - Owner need to call this method.
  - Send Accumulated Fees to `feeAggregator` wallet


#### handleRelayMessage
* Description:
  - It verify and decode RelayMessage with BMV, and dispatch BTP Messages
    to registered BSHs
  - It's allowed to be called by registered Relay.
* Params
    | Parameters | Type   | Description                                                |
    |:-----------|:-------|:-----------------------------------------------------------|
    | _prev      | String | BTP Address of the BMC generates the message               |
    | _msg       | String | base64 encoded string of serialized bytes of Relay Message |

#### registerRelayer
* Description
  - Set relayer to the BMC 
  - Pay Minimum Bond Relayer Fee to be registered
* Params
    | Parameters | Type   | Description         |
    |:-----------|:-------|:--------------------|
    | _desc      | String | Relayer Description |

#### removeRelayer
* Description
  - Only owner can remove the relayer to the BMC 
* Params
    | Parameters | Type    | Description     |
    |:-----------|:--------|:----------------|
    | _addr      | Address | Relayer Address |
    | _refund    | Address | Refund Address  |


#### unregisterRelayer
* Description
  - Remove relayer from the BMC 


## Read-only methods

#### name
* Description
  - Returns name of the Contract
    ```json
    "BTP Message Center"
    ```

#### getBtpAddress
* Description
  - Return BTP address of the BMC
    ```json
    "btp://0x1.icon/cx23a91ee3dd290486a9113a6a42429825d813de53"
    ```
#### getOwners
* Description 
  - Return list of Owners/Admins of BTP
    ```json
    [
      "hxdd6b62ab563ff0cad4ef7248ed2f9458059a18d2",
      "hxc745bee97c35183f9f83b13990a9ae1def328ac9"
    ]
    ```

### isOwner
* Description
  * Returns bool value, `0x1` if owner and `0x0` if not Owner
* Params
    | Parameters | Type    | Description                    |
    |:-----------|:--------|:-------------------------------|
    | _addr      | Address | Address to be checked as Owner |

#### getServices
* Description
    - Get registered services.
* Returns
    - A dictionary with the name of the service as key and address of the BSH
      related with the service as value.
      ```json
         {"bts":"cxcef70e92b89f2d8191a0582de966280358713c32"}
      ```

#### getLinks
* Description
    - Get registered links.
* Returns
    -  A list of links  (BTP Addresses of the BMCs) 
        ```json
        ["btp://0x38.bsc/0x034AaDE86BF402F023Aa17E5725fABC4ab9E9798"]
        ```

#### getFeeGatheringTerm
* Description
  - Returns `feeGatheringTerm` 
    ```string
    43200
    ```

#### getFeeAggregator
* Description
  - Returns `feeAggregator` wallet for BTP
      ```json
      "hxcef70e92b89f2d8191a0582de966280358713c32"
      ```

#### getRoutes
* Description:
    - Get routing information.
* Return
    - A dictionary with the BTP Address of the destination BMC as key and
      the BTP Address of the next as value.
      ```json
      {
        "btp://0x2.iconee/cx1d6e4decae8160386f4ecbfc7e97a1bc5f74d35b", 
        "btp://0x1.iconee/cx9f8a75111fd611710702e76440ba9adaffef8656"
      }
      ```

#### getRelayers
* Description:
    - Get relayers' information.
* Return
    - A dictionary with the BTP Address of the destination BMC as key and
      the BTP Address of the next as value.
      ```json
      {
        "btp://0x2.iconee/cx1d6e4decae8160386f4ecbfc7e97a1bc5f74d35b",
        "btp://0x1.iconee/cx9f8a75111fd611710702e76440ba9adaffef8656"
      }
      ```
#### getRelays
  * Description:
      - Get relay address of given BTP link
  * Params
      | Parameters | Type | Description                            |
      |:-----------|:-----|:---------------------------------------|
      | _link      | str  | BTP Address of BMC for the destination |
  * Return
      - A list with the BTP relay Address of the destination.
        ```string
        ["hx04dc2d402fdc0b31ce1044c459ec3283c8d04aed"]
        ```

#### getStatus
  * Description:
      - Get status of BMC.
      - It's used by the relay to resolve next BTP Message to send.
      - If target is not registered, it will fail.
  * Params
      | Parameters | Type   | Description                      |
      |:-----------|:-------|:---------------------------------|
      | _link      | String | BTP Address of the connected BMC |
  * Return 
    * The object contains followings fields.

      | Field              | Type       | Description                                      |
      |--------------------|------------|--------------------------------------------------|
      | block_interval_dst | Integer    | next sequence number of the next sending message |
      | block_interval_src | Integer    | next sequence number of the message to receive   |
      | cur_height         | Long       | Current block height                             |
      | delay_limit        | Integer    |                                                  |
      | max_agg            | Integer    |                                                  |
      | relay_idx          | Integer    |                                                  |
      | relays             | Array      |                                                  |
      | rotate_height      | Long       |                                                  |
      | rotate_term        | Integer    |                                                  |
      | rx_height          | Long       |                                                  |
      | rx_height_src      | Long       |                                                  |
      | rx_seq             | BigInteger |                                                  |
      | sack_height        | Long       |                                                  |
      | sack_next          | Long       |                                                  |
      | sack_seq           | BigInteger |                                                  |
      | sack_term          | Integer    |                                                  |
      | tx_seq             | BigInteger |                                                  |
      | verifier           | BMCStatus  |                                                  |

        ```json
        {"block_interval_dst":"0x0","block_interval_src":"0x0","cur_height":"0x35104f2","delay_limit":"0x0","max_agg":"0x0","relay_idx":"0x0","relays":null,"rotate_height":"0x0","rotate_term":"0x0","rx_height":"0x1485c3b","rx_height_src":"0x0","rx_seq":"0x51","sack_height":"0x0","sack_next":"0x0","sack_seq":null,"sack_term":"0x0","tx_seq":"0x5d","verifier":null}
        ```  
  
#### getServiceCandidates
* Description 
  * Returns the list of serviceCandidate address list of all services registered to the BTP

#### getRelayerTerm
* Description
  * Returns value of current `relayerTerm`
    ```json
    43200
    ```

#### getRelayerMinBond
* Description
  * Returns value of minimum bond to register as relayer
    ```json
    100
    ```


### Events

#### Message
* Indexed(2)
* Params
    | Parameters | Type       | Description                                                 |
    |:-----------|:-----------|:------------------------------------------------------------|
    | _next      | String     | BTP Address of the BMC to handle the message                |
    | _seq       | BigInteger | sequence number of the message from current BMC to the next |
    | _msg       | Bytes      | serialized bytes of BTP Message                             |
* Description
    - It sends the message to the next BMC.
    - The relay monitors this event.

#### ErrorOnBTPError
* Indexed(2)
* Params
    | Parameters | Type       | Description |
    |:-----------|:-----------|:------------|
    | _svc       | String     |             |
    | _seq       | BigInteger |             |
    | _code      | Long       |             |
    | _msg       | String     |             |
    | _ecode     | Long       |             |
    | _emsg      | String     |             |