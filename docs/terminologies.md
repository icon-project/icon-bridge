# Terminologies

* [Network Address](#network-address)

  A string to identify blockchain network

* [BTP Address](#btp-address)

  A string of URL for locating an account of the blockchain network

* [BTP Message](#btp-message)

  Standardized messages delivered between different blockchains

* [BTP Message Center(BMC)](./specs/bmc.md)

  BMC sends message to relay and as well as accepts messages from a relay (Relay Messages).

* [BTP Service Handler(BSH)](./specs/bsh.md)

  BSH handles BTP Messages of the service. It also sends messages according to different service scenarios.

* [BTP Token Service(BTS)](./specs/bts.md)

  BTP is a service handler which handles BTP Messages of token transfers. It is responsible for cross chain token transfers.

---

### Network Address

A string to identify blockchain network

```
<NID>.<Network Name>
```

**NID**: ID of the network in the blockchain network system.

> Example

| Network Address | Description                                  |
|:----------------|:---------------------------------------------|
| `0x1.icon`      | ICON Network with nid="0x1" <- ICON Main-Net |
| `0x38.bsc`      | BSC Network with nid="0x38" <- BSC Main-Net  |

### BTP Address

A string of URL for locating an account of the blockchain network

> Example
```
btp://<Network Address>/<Account Identifier>
```
**Account Identifier**:
Identifier of the account including smart contract.
It should be composed of URL safe characters except "."(dot).

> Example
```
btp://0x1.icon/hxc0007b426f8880f9afbab72fd8c7817f0d3fd5c0
btp://0x38.bsc/0x429731644462ebcfd22185df38727273f16f9b87
```

It could be expanded to other resources.

### BTP Message

A message delivered across blockchains. It has the following fields.

| Name | Type    | Description                                          |
|:-----|:--------|:-----------------------------------------------------|
| src  | String  | BTP Address of source BMC                            |
| dst  | String  | BTP Address of destination BMC                       |
| svc  | String  | name of the service                                  |
| sn   | Integer | serial number of the message                         |
| msg  | Bytes   | serialized bytes of Service Message or Error Message |

if **sn** is negative, **msg** should be Error Message.
It would be serialized in [RLP serialization](#rlp-serialization).


### Error Message

A message for delivering error information.

| Name | Type    | Description   |
|:-----|:--------|:--------------|
| code | Integer | error code    |
| msg  | String  | error message |

It would be serialized in [RLP serialization](#rlp-serialization).


### RLP serialization

For encoding [BTP Message](#btp-message) and [Error Message](#error-message), it uses Recursive Length Prefix (RLP).
RLP supports bytes and list naturally.
Here are some descriptions about other types.
