# BTP Message Center(BMC)

## Introduction

BTP Message Center is a smart contract which builds BTP Message and
sends it to BTP Message Relay and handles Relay Message from the other.

## Setup

1. Registers [BSH](bsh.md)s for the services.
   (BSH should be deployed before the registration)
2. Add links, BMCs of directly connected blockchains
3. Add routes to other BMCs of in-directly connected blockchains

## Send a message

BSH sends a message through [BMC.sendMessage](bmc.md#sendmessage).
It accepts the requests from the registered BTP Service Handler(BSH).
Of course, if service name of those requests is different from
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
for next BMC.

## Receive a message

It receives the Relay Message, then it tries to decode it with registered
BMV. It may contain multiple BTP Messages.
It dispatches received BTP Messages one-by-one in the sequence.

If it is the destination, then it tries to find the BSH for the
service, and then calls [BSH.handleBTPMessage](bsh.md#handlebtpmessage).
It calls [BSH.handleBTPError](bsh.md#handlebtperror) if it's an error.

If it's not the destination, then it tries to send the message to
the next route.

If it fails, then it replies an error.
BTP Message for error reply is composed of followings.
* sn : negated serial number of the message.
* dst : BTP Address of the source.
* src : BTP Address of the BMC.
* msg : Error Message including error code and message.
