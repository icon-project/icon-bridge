# Javascore BTS Unit Tests

## Introduction

These tests are designed to ensure that the proper functioning of features on the javascore BTS Contract.

## How to run 
```sh
cd javascore/bts
./gradlew :bts:test
```

## Covered test cases
- Prerequisties
  - ICON: Primary Network 
  - METACHAIN: Connected Network
  - PARA: A token on METACHAIN network
  - TEST_TOKEN: A token on ICON network
  - UINT_CAP: Value equal to (2 ** 256 - 1)


- Add to locklist
    - Add user to blacklist on ICON
    - Add user to blacklist on other chain
    - Successful blacklist response from other chain
    - Unsuccessful blacklist response from other chain
    - handleBTPError for add to blacklist
    

- Remove from locklist
  - Remove user from blacklist on ICON
  - Remove user from blacklist on other chain
  - Successful remove from blacklist response from other chain
  - Unsuccessful remove from blacklist response from other chain
  - handleBTPError for remove from blacklist


- Registered wrapped token
  - Register a wrapped token PARA
  - Only owner should be able to register
  - Deploys IRC2 Tradeable contract
  - Added to registered coin array
  - Token limit should be set to UINT_CAP


- Register IRC2 Token
  - Does not deploy a IRC2 Tradeable Contract
  - Owner should be able to register
  - Non owner should not 
  - Should not be able to register token that already exists
  - Token added to coin array
  - Token limit set to UINT_CAP
  

- Operation on unregistered token
  - Cannot do any operations on unregistered token
  - BTS Contract should not accept those tokens
  

- Reclaim tokens transferred to BTS
  - Transfer tokens to BTS and reclaim w/o cross chain transfer
  - No fee deducted
  - Partial reclaim
  - Full reclaim


- Transfer nativecoin
  - Transfer nativecoin to METACHAIN
  - Blacklist sender on ICON, transfer fail
  - Blacklist reciever on METACHAIN, transfer fail
  - Amount exceed token limit, transfer fail
  - handle unsuccessful transfer on destination chain
  - relay fee deducted on cross chain transfer
  - handleBTPError on transferNativecoin
    - Transfer nativecoin back to user successful
    - Transfer nativecoin back to user unsuccessful, set to refundable
    - Reclaim refundable amount


- Handle Fee Gathering
  - Gather fees after transfer/transfer batch
  - Send fees to fee aggregator wallet


- Transfer tokens
  - Transfer IRC2 token to destination
  - Transfer invalid amount should fail
  - Transfer unregistered token should fail
  - Once transfer originates, it's in locked state
  - TransferTransaction for that sn not null
  - Response recieved through handleBTPMessage
    - Successful response
      - Set amount to refundable if refund fails
  - Response recieved through handleBTPError
    - Locked balance = 0
    - Refund (amount-fee) back to user
    - If refund fails, set that amount to refundable
  - User can reclaim refundable amount
  - Locked amount deducted if successful on destination


- TransferBatch
  - Transfer nativecoin, IRC2 and wrapped coins all at once
  - Fee deducted for each token transfered


- Unknown response
  - Generates UnknownResponse eventlog


- Token limit check and blacklisted users check
  - Only owner can change token limit
  - Only owners can blacklist users
  - Limit can be set for tokens not registered as well
  - Transfer from/to blacklisted user fails
  - Amount greater than token limit means txn reverted
  - Token limit status set to false before response is recieved


- Owner tests 
  - Non-Owner tries to add a new Owner
  - Owner tries to add themselves
  - Current Owner adds a new Owner
  - After adding a new Owner, owner registers a new coin
  - New Owner registers a new coin
  - Newly added owner tries to add owner
  - Current Owner removes another Owner
  - Owner tries to add itself again
  - The last Owner removes him/herself


- Set Fee Ratio
  - Only owner can update fee ratio
  - Fee_numerator should not be set higher than fee_denominator
  - Try to set negative fee
