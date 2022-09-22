## API Documents

### Writable methods

### Owner guarded methods

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


#### setBMCPeriphery
* Params
    - _addr : Address (BMCPeriphery Address to handle the messages)

* Description
  - Set BMCPeriphery Contract address on the BMC contract which contains the core methods of BMC.
  - Only Owner can call this method.


#### addService
* Params
    - _svc: String (the name of the service)
    - _addr: Address (the address of the smart contract handling the service)
* Description:
    - It registers the smart contract for the service.
    - It's called by the owner/admin to manage the BTP network.

#### removeService
* Params
    - _svc: String (the name of the service)
* Description:
    - It de-registers the smart contract for the service.
    - It's called by the operator to manage the BTP network.


#### handleRelayMessage
* Params
    - _prev: String ( BTP Address of the BMC generates the message )
    - _msg: String ( base64 encoded string of serialized bytes of Relay Message )
* Description:
    - It verify and decode RelayMessage with BMV, and dispatch BTP Messages
      to registered BSHs
    - It's allowed to be called by registered Relay.

#### sendMessage
* Params
    - _to: String ( Network Address of destination network )
    - _svc: String ( name of the service )
    - _sn: Integer ( serial number of the message, it should be positive )
    - _msg: Bytes ( serialized bytes of Service Message )
* Description:
    - It sends the message to specific network.
    - It's allowed to be called by registered BSHs.


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

#### setLink
* Params 
  - _link: String (BTP Address of Connected BMC)
  - _blockInterval: Uint256 (Sync Interval (By Default 1000))
  - _maxAggregation: Uint256
  - _delayLimit: Uint256 (Delay Seconds (3))
* Description 
  - Edit already added Link parameters
  - Only Owner can call this method

#### setLinkRxHeight
* Params
  - _link: String (Added Link)
  - _height: Uint256 (Starting BlockHeight)
* Description
  * Modify the blockHeight on given link
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
    - _link: String ( BTP Address of the next BMC for the destination )
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

#### rotateRelay
* Params 
  * string memory _link: String  (Link)
  * uint256 _currentHeight: Uint256
  * _relayMsgHeight: Uint256
  * _hasMsg: Bool
* Description
  * Only BMCPeriphery contract can call this method.

### Read-only methods

#### getServices
* Description
    - Get registered services.
* Returns
    - A dictionary with the name of the service as key and address of the BSH
      related with the service as value.
      ```string
        bsc,0x0a7792fe75548b26b287871081Aa6b05f48D9e89,sicx,0xc0c1aA22F99bb6724dC4159C256A5989D90A659C
      ```

#### getLinks
* Description
    - Get registered links.
* Returns
    -  A list of links ( BTP Addresses of the BMCs )
  ```json
  [ "btp://0x1.iconee/cx9f8a75111fd611710702e76440ba9adaffef8656" ]
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

#### getStatus
* Params
    - _link: String ( BTP Address of the connected BMC )
  * Description:
      - Get status of BMC.
      - It's used by the relay to resolve next BTP Message to send.
      - If target is not registered, it will fail.
  * Return 
    * The object contains followings fields.
      ```string
      90,79,55521830,21494158
      ```
    | Field       | Type    | Description                                           |
    |-------------|---------|-------------------------------------------------------|
    | tx_seq      | Integer | next sequence number of the next sending message      |
    | rx_seq      | Integer | next sequence number of the message to receive        |
    | rxHeight    | Integer | status information of the link connection blockheight |
    | blockNumber | Integer | Current block height                                  |  



### Events

#### Message
* Indexed: 1
* Params
    - _next: String ( BTP Address of the BMC to handle the message )
    - _seq: Integer ( sequence number of the message from current BMC to the next )
    - _msg: Bytes ( serialized bytes of BTP Message )
* Description
    - It sends the message to the next BMC.
    - The relay monitors this event.