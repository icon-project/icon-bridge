# BTP Service Handler (BSH)

## Introduction

BTP Service Handler (BSH) may send messages through BTP Message Center(BMC) from any user request. Of course, the request may come from other smart contracts. It also have responsibility to handle the message from other BSHs. BSH can communicate other BSHs with same service name. If there is already the service using same name, then it should choose
other name for the service when it registers a new service.

The BTS contract is a BSH contract, which originates cross chain token transfer and gets response back from BTS of destination chain and services based on the response received. BTS contract maintains information about blacklist and token transfer limit as well.

BTS, being a BSH, needs to be registered to the BMC, before being able to send a BTP Message or handle incoming response messages.

For a contract to be BSH, the following are required:

1. Implement [BSH interface](#bsh-interface).
2. Registered to the BMC through [BMC.addService](bmc.md#addservice)

After the registration, it may send messages through
[BMC.sendMessage](bmc.md#sendmessage).

If there is an error while it delivers the message, then it will
return error information though [handleBTPError](#handlebtperror).

If it's successfully delivered, then BMC will call
[handleBTPMessage](#handlebtpmessage) of the target BTS.

While it processes the message, it may reply though
[BMC.sendMessage](bmc.md#sendmessage).

## Setup

1. Registers [BSH](bsh.md) for the services in BMC.
   (BSH should be deployed before the registration)
2. Add links, BMCs of directly connected blockchains
3. Add routes to other BMCs of in-directly connected blockchains

## BSH Interface
The BTP Service Handler (BSH) should implement the following 3 methods:
- [handleBTPMessage](#handlebtpmessage)
- [handleBTPError](#handlebtperror)
- [handleFeeGathering](#handlefeegathering)

*  ### handleBTPMessage
   - #### Method
      - handleBTPMessage
   - #### Parameters
        | Parameters | Type    | Info                                              |
        |:-----------|:--------|:--------------------------------------------------|
        | _from      | string  | An originated network address of a request        |
        | _svc       | string  | Service name of BSH contract                      |
        | _sn        | integer | Serial number of a service request                |
        | _msg       | bytes   | RLP message of a service request/service response |
   - #### Implementation
      - [Java](/javascore/bts/src/main/java/foundation/icon/btp/bts/BTPTokenService.java)
      - [Solidity](/solidity/bts/contracts/)

* ### handleBTPError
   - #### Method
      - handleBTPError
   - #### Parameters

        | Parameters | Type    | Info                                        |
        |:-----------|:--------|:--------------------------------------------|
        | _src       | string  | An originated network address of a request  |
        | _svc       | string  | Service name of BSH contract                |
        | _sn        | integer | Serial number of a service request          |
        | _code      | integer | Response code of a message (RC_OK / RC_ERR) |
        | _msg       | string  | Response message                            |

   - #### Implementation
      - [Java](/javascore/bts/src/main/java/foundation/icon/btp/bts/BTPTokenService.java)
      - [Solidity](/solidity/bts/contracts/)

* ###  handleFeeGathering
   - #### Method
      - handleFeeGathering
   - #### Parameters

        | Parameters | Type   | Info                          |
        |:-----------|:-------|:------------------------------|
        | _fa        | string | BTP Address of fee aggregator |
        | _svc       | string | Service name of BSH contract  |

   - #### Implementation
      - [Java](/javascore/bts/src/main/java/foundation/icon/btp/bts/BTPTokenService.java)
      - [Solidity](/solidity/bts/contracts/)
