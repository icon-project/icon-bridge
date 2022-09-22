# API Documents

## Writable methods

#### addOwner
* Params
  - _owner : Address (Admin/Owner address to access on BMC)
* Description
  - Add owners to owner list
  - Owner can call owner guarded methods

#### removeOwner
* Params
    - _owner : Address
* Description
    - Remove owners to owners list
    - Cannot remove the last owner from the list

  
#### addService
* Params
    - _svc: String (the name of the service)
    - _addr: Address (the address of the smart contract handling the service)
* Description:
    - It registers the smart contract for the service.
    - It's called by the owner/admin to manage the BTP network.

#### addServiceCandidate
* Params:
  * _svc: String (the name of the service)
  * _addr: Address (Address to be added on the service)
* Description:
  * Add Service Candidate on the given service

#### removeServiceCandidate
* Params:
  * _svc: String (the name of the service)
  * _addr: Address (Address to be added on the service)
* Description:
  * Remove Service Candidate on the given service

#### removeService
* Params
    - _svc: String (the name of the service)
* Description:
    - It de-registers the smart contract for the service.
    - It's called by the operator to manage the BTP network.

#### handleFragment
* Params
  * _prev: String
  * _msg: String
  * _idx: Integer

#### addLink
* Params
    - _link: String (BTP Address of connected BMC)
* Description
    - If it generates the event related with the link, the relay shall
      handle the event to deliver BTP Message to the BMC.
    - If the link is already registered, or its network is already
      registered then it fails.
    - It initializes status information for the link.
    - It's called by the operator to manage the BTP network.

#### setLinkRxHeight
* Params
  - _link: String (Added Link)
  - _height: Long (Starting BlockHeight)
* Description
  * Modify the blockHeight on given link
  * Only owner can manage this on BTP Network

#### setLinkRotateTerm
* Params
  - _link: String (Added Link)
  - _block_interval: Integer (Block interval)
  - _max_agg: Integer (Max Aggregation Value)
* Description
  * Update the blockHeight, blockInterval and maxAggreation on given link
  * Only owner can manage this on BTP Network
  
#### setLinkDelayLimit
* Params
  - _link: String (Added Link)
  - _value: Integer (Block interval)
* Description
  * Update the `delayLimit` value on given link
  * Only owner can manage this on BTP Network

#### setLinkSackTerm
* Params
  - _link: String (Added Link)
  - _value: Integer (Block interval)
* Description
  * Update the `sackTerm` value on given link
  * Only owner can manage this on BTP Network


#### removeLink
* Params
    - _link: String (BTP Address of connected BMC)
* Description
    - It removes the link and status information.
    - It's called by the operator to manage the BTP network.

#### addRoute
* Params
    - _dst: String ( BTP Address of the destination BMC )
    - _link: String (BTP Address of the next BMC for the destination )
* Description:
    - Add route to the BMC.
    - It may fail if there are more than one BMC for the network.
    - It's called by the operator to manage the BTP network.

#### removeRoute
* Params
    - _dst: String ( BTP Address of the destination BMC )
* Description:
    - Remove route to the BMC.
    - It's called by the operator to manage the BTP network.

#### addRelay
* Params
  * _link: String (BTP Address of the next BMC)
  * _addr: Address (Relayer Address)
* Description
  * Add relayer address on given link (BTPAddress)
  * Only owner can call this method

#### removeRelay
* Params
  * _link: String (BTP Address of the next BMC)
  * _addr: Address (Relayer Address)
* Description
  * Remove relayer address on given link (BTPAddress)
  * Only owner can call this method

#### distributeRelayerReward
* Params
* Description
  * Distribute relayers' reward to each relayer of registered services

#### claimRelayerReward
* Params
* Description
  * Claim the distributed relayer rewards from relayer address

#### setRelayerMinBond
* Params 
  * _value: BigInteger (Value to be set as Minimum Bond as a Relayer)
* Description
  * Set Minimum relayer bond while registering as relayer
  * Only Owner can call this method

#### setFeeGatheringTerm
  * Params
    * _value : Long
  * Description
    * Set fee feeCollecting BlockHeight
    * Only owner can set the value

#### setFeeAggregator
  * Params
    * _addr : Address
  * Description
    * Set fee feeCollecting address
    * Only owner can set the value

#### sendFeeGathering
* Params 
* Description
  * Owner need to call this method.
  * Send Accumulated Fees to `feeAggregator` wallet


#### handleRelayMessage
* Params
  - _prev: String ( BTP Address of the BMC generates the message )
  - _msg: String ( base64 encoded string of serialized bytes of Relay Message )
* Description:
  - It verify and decode RelayMessage with BMV, and dispatch BTP Messages
    to registered BSHs
  - It's allowed to be called by registered Relay.

#### registerRelayer
* Params
  * _desc: String (Relayer Description)
* Description
  * Set relayer to the BMC 
  * Pay Minimum Bond Relayer Fee to be registered

#### removeRelayer
* Params
  * _addr: Address (Relayer Address)
  * _refund: Address (Refund Address)
* Description
  * Only owner can remove the relayer to the BMC 


#### unregisterRelayer
* Params

* Description
  * Remove relayer from the BMC 


## Read-only methods

#### name
* Description
  * Returns name of the Contract
  ```string
  BTP Message Center 
  ```

#### getBtpAddress
* Description
  * Return BTP address of the BMC
  ```string
   btp://0x1.icon/cx23a91ee3dd290486a9113a6a42429825d813de53 
  ```
#### getOwners
* Description 
  * Return list of Owners/Admins of BTP
  ```string
  ["hxdd6b62ab563ff0cad4ef7248ed2f9458059a18d2","hxc745bee97c35183f9f83b13990a9ae1def328ac9"]
  ```

### isOwner
* Params 
  * _addr: Address (Address to be checked as Owner)
* Description
  * Returns bool value, `0x1` if owner and `0x0` if not Owner

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
    -  A list of links ( BTP Addresses of the BMCs )
  ```json
  ["btp://0x38.bsc/0x034AaDE86BF402F023Aa17E5725fABC4ab9E9798"]
  ```

#### getFeeGatheringTerm
* Description
  * Returns `feeGatheringTerm` 
  ```string
  43200
  ```

#### getFeeAggregator
* Description
  * Returns `feeAggregator` wallet for BTP
  ```string
  hxcef70e92b89f2d8191a0582de966280358713c32
  ```

#### getRoutes
* Description:
    - Get routing information.
* Return
    - A dictionary with the BTP Address of the destination BMC as key and
      the BTP Address of the next as value.
      ```json
      {
        "btp://0x2.iconee/cx1d6e4decae8160386f4ecbfc7e97a1bc5f74d35b": "btp://0x1.iconee/cx9f8a75111fd611710702e76440ba9adaffef8656"
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
        "btp://0x2.iconee/cx1d6e4decae8160386f4ecbfc7e97a1bc5f74d35b": "btp://0x1.iconee/cx9f8a75111fd611710702e76440ba9adaffef8656"
      }
      ```
#### getRelays
* Params 
  * _link: str (BTP Address of BMC for the destination)
* Description:
    - Get relay address of given BTP link
* Return
    - A list with the BTP relay Address of the destination.
      ```string
      ["hx04dc2d402fdc0b31ce1044c459ec3283c8d04aed"]
      ```

#### getStatus
* Params
    - _link: String ( BTP Address of the connected BMC )
  * Description:
      - Get status of BMC.
      - It's used by the relay to resolve next BTP Message to send.
      - If target is not registered, it will fail.
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

```string
43200
```

#### getRelayerMinBond
* Params
* Description
  * Returns value of minimum bond to register as relayer
  ```string
  100
  ```


### Events

#### Message
* Indexed: 2
* Params
    - _next: String ( BTP Address of the BMC to handle the message )
    - _seq: BigInteger ( sequence number of the message from current BMC to the next )
    - _msg: Bytes ( serialized bytes of BTP Message )
* Description
    - It sends the message to the next BMC.
    - The relay monitors this event.

#### ErrorOnBTPError
* Indexed: 2
* Params:
  * _svc: String
  * _seq: BigInteger
  * _code: Long
  * _msg: String
  * _ecode: Long
  * _emsg: String