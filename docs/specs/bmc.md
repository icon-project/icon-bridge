# BTP Message Center(BMC)

## Introduction

BTP Message Center is a smart contract which builds BTP Message and
sends it to BTP Message Relay and handles Relay Message from the other.

## Setup

1. Register BSHs for the services. [BTP Token Service (BTS)](bts.md) is a BSH connected to BMC. (BSH should be deployed before the registration)
2. Add links, BMCs of directly connected blockchains
3. Add routes to other BMCs of in-directly connected blockchains

## BMC Interfaces
The BMC contract has to implement the following methods.
- [sendMessage](#send-a-message-sendmessage)
- [handleRelayMessage](#handlerelaymessage)

The BMC contract has to implement the following events.
- [Message event](#message)

## sendMessage
* #### Method
    - sendMessage
* #### Params
    | Name | Type    | Description                                   |
    |:-----|:--------|:----------------------------------------------|
    | _to  | String  | Network Address of the destination blockchain |
    | _svc | String  | Name of the service.                          |
    | _sn  | Integer | Serial number of the message.                 |
    | _msg | Bytes   | Service message to be delivered.              |
* #### Description

    A service handler sends a message through [BMC.sendMessage](bmc.md#sendmessage).
    It accepts the requests from the registered BTP Service Handler(BSH).
    If service name of those requests is different from
    registration, then they will be rejected.

    Then it builds a BTP Message from the request.
    1. Decide destination BMC from given Network Address
    2. Fill in other information from parameters.
    3. Serialize them for sending.

    Then it tries to send the BTP Message.
    1. Decide next BMC from the destination referring routing information.
    2. Get sequence number corresponding to the next.
    3. Emit the event, [Message](#message) including the information.

    The event will be monitored by the Relay, it will build Relay Message
    for next BMC of destination chain.
* #### Implementation
    1. [Java](/javascore/bmc/src/main/java/foundation/icon/btp/bmc/BTPMessageCenter.java)
    2. [Solidity](/solidity/bmc/contracts/BMCPeriphery.sol)

## handleRelayMessage
* #### Method
    - handleRelayMessage
* #### Params 
    | Name  | Type   | Description                                     |
    |:------|:-------|:------------------------------------------------|
    | _prev | String | BTP Address of the previous BMC                 |
    | _msg  | Bytes  | serialized Relay Message including BTP Messages |
* #### Description

    This method receives the Relay Message from the relay, then it tries to decode it. It may contain multiple BTP Messages.
    It dispatches received BTP Messages one-by-one in the sequence.

    If it is the destination, then it tries to find the BSH for the
    service, and then calls [BSH.handleBTPMessage](bts.md#handlebtpmessage).
    It calls [BSH.handleBTPError](bts.md#handlebtperror) if it's an error.

    If it's not the destination, then it tries to send the message to
    the next route.

    If it fails, then it replies an error.
* #### Implementation
    1. [Java](/javascore/bmc/src/main/java/foundation/icon/btp/bmc/BTPMessageCenter.java)
    2. [Solidity](/solidity/bmc/contracts/BMCPeriphery.sol)


## Message
* #### Event
    - Message
* #### Params
    | Name  | Type    | Description                                |
    |:------|:--------|:-------------------------------------------|
    | _next | String  | BTP Address of the next BMC                |
    | _seq  | Integer | Sequence number of the msg to the next BMC |
    | _msg  | Bytes   | Serialized BTP Message                     |
* #### Description
    The [sendMessage](#sendmessage) method in BMC Contract generates a `Message` eventlog for the relay to pick up and perform transaction based on this event on the destination chain.

* #### Implementation
    1. [Java](/javascore/bmc/src/main/java/foundation/icon/btp/bmc/BTPMessageCenter.java)
    2. [Solidity](/solidity/bmc/contracts/BMCPeriphery.sol)
