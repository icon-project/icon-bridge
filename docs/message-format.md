## Message Formats

1. BTS sends an Service Message through BMC.

   * BTS calls [BMC.sendMessage](./specs/bmc.md#sendmessage) with followings.

        | Name    | Type    | Description                                   |
        |:--------|:--------|:----------------------------------------------|
        | _to     | String  | Network Address of the destination blockchain |
        | _svc    | String  | Name of the service.                          |
        | _sn     | Integer | Serial number of the message.                 |
        | _msg    | Bytes   | Service message to be delivered.              |

   * BMC lookup the destination BMC belonging to *_to*.
     If there is no known BMC to the network, then it fails.

   * BMC builds an BTP Message.

        | Name | Type    | Description                                   |
        |:-----|:--------|:----------------------------------------------|
        | src  | String  | BTP Address of current BMC                    |
        | dst  | String  | BTP Address of destination BMC in the network |
        | svc  | String  | Given service name                            |
        | sn   | Integer | Given serial number                           |
        | msg  | Bytes   | Given service message                         |

   * BMC decide the next BMC according to the destination.
     If there is no route to the destination BMC.

   * BMC generates a `Message` event with BTP Message.
   
        | Name  | Type    | Description                                |
        |:------|:--------|:-------------------------------------------|
        | _next | String  | BTP Address of the next BMC                |
        | _seq  | Integer | Sequence number of the msg to the next BMC |
        | _msg  | Bytes   | Serialized BTP Message                     |

2. The BTP Message Relay(BMR) detects Message event.
   * The relay detects [BMC.Message](./specs/bmc.md#message) through various ways.
   * The relay can confirm that it occurs and it's finalized.
   
3. BMR gathers proofs
    * Relay gathers proofs of the event(POE)s
        - Proof for the new block
        - Proof for the event in the block
    * Relay builds Relay Message including followings.
        - Proof of the new events
        - New events including the BTP Message.
    * Relay calls [BMC.handleRelayMessage](./specs/bmc.md#handlerelaymessage)
        with built Relay Message.
     
        | Name  | Type   | Description                                     |
        |:------|:-------|:------------------------------------------------|
        | _prev | String | BTP Address of the previous BMC                 |
        | _msg  | Bytes  | serialized Relay Message including BTP Messages |
     
4. BSH handles service message

    * BMC dispatches BTP Messages.
    * If the destination BMC isn't current one, then it locates
     the next BMC and generates the event.
    * If the destination BMC is the current one, then it locates BSH
     for the service of the BTP Message.
    * Calls [BSH.handleBTPMessage](./specs/bsh.md#handlebtpmessage) if
     the message has a positive value as *_sn*.

        | Name  | Type    | Description                           |
        |:------|:--------|:--------------------------------------|
        | _from | String  | Network Address of the source network |
        | _svc  | String  | Given service name                    |
        | _sn   | Integer | Given serial number                   |
        | _msg  | Bytes   | Given service message                 |
     
   * Otherwise, it calls [BSH.handleBTPError](./specs/bsh.md#handlebtperror).
   
        | Name  | Type    | Description                                    |
        |:------|:--------|:-----------------------------------------------|
        | _src  | String  | BTP Address of the BMC that generated the error|
        | _svc  | String  | Given service name                             |
        | _sn   | Integer | Given serial number                            |
        | _code | Integer | Given error code                               |
        | _msg  | String  | Given error message                            |
