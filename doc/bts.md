# BTP Token Service (BTS)

## Introduction

BTP Token Service is a service handler smart contract which handles token transfer operations. It sends BTP message to other chains through BTP Message Center(BMC). Similarly, it recieves incoming requests from other chain via BMC, handles response and send a message back to source chain. The BTS contract originates cross chain token transfer and gets response back from BTS of destination chain and services based on the response reiceved.

BTS contract maintains information about blacklist and token transfer limit as well.

BTS needs to be registered to the BMC, before being able to send a BTP Message or handle incoming response messages.

For a contract to be BTS, the following are required:

1. Implement BSH interface.
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

1. Registers [BTS](bts.md) s for the services in BMC.
   (BTS should be deployed before the registration)
2. Add links, BMCs of directly connected blockchains
3. Add routes to other BMCs of in-directly connected blockchains
