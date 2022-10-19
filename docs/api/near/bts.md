# NEAR BTS Contract

## Requirements
* #### NEAR CLI
    - The NEAR [Command Line Interface (CLI)](https://github.com/near/near-cli) is a tool that enables to interact with the NEAR network directly from the shell. Under the hood, NEAR CLI utilizes the NEAR [JavaScript API](https://github.com/near/near-api-js)
    - Installation
    ```console
    npm install -g near-cli
    ```
* #### NEAR Wallet  
    - Testnet: https://wallet.testnet.near.org/
    - Mainnet: https://wallet.near.org/

* #### Authorize CLI
    ```console
    near login
    ```

## API

### Balance

#### Usable Balance
**Method** 
- balance_of

| Parameters | Type | Info |
|:---------|:--------|:--------|
| account_id | string | should be a valid account id |
| coin_name | string |  |

**CLI Command** 
```console
NEAR_ENV=testnet near view <BTS> balance_of '{"account_id": "<ACCOUNT ID>", "coin_name": "<Coin Name>"}'
```

#### Refundable Balance
**Method** 
- refundable_balance_of

| Parameters | Type | Info |
|:---------|:--------|:--------|
| account_id | string | should be a valid account id |
| coin_name | string |  |

**CLI Command** 
```console
NEAR_ENV=testnet near view <BTS> refundable_balance_of '{"account_id": "<ACCOUNT ID>", "coin_name": "<Coin Name>"}'
```

#### Locked Balance
**Method** 
- locked_balance_of

| Parameters | Type | Info |
|:---------|:--------|:--------|
| account_id | string | should be a valid account id |
| coin_name | string |  |

**CLI Command** 
```console
NEAR_ENV=testnet near view <BTS> locked_balance_of '{"account_id": "<ACCOUNT ID>", "coin_name": "<Coin Name>"}'
```

### Deposit

#### Deposit NEAR to BTS
**Method** 
- deposit

**CLI Command** 
```console
NEAR_ENV=testnet near call <BTS> deposit --amount <AMOUNT in NEAR> --accountId <ACCOUNT ID>
```
#### Deposit Cross-Chain Native Coin to BTS
**Method** 
- ft_transfer_call

| Parameters | Type | Info |
|:---------|:--------|:--------|
| receiver_id | string | should be a valid account id |
| amount | string |  |
| msg | string |  |

**CLI Command** 
```console
NEAR_ENV=testnet near call <NEP141 Contract> ft_transfer_call '{"receiver_id": "<BTS>", "amount": "<AMOUNT>", "msg": ""}' --accountId <ACCOUNT ID> --amount <1 yoctoNEAR in highest Denomination ie 0.000000000000000000000001>
```

### Withdraw
**Method** 
- withdraw

| Parameters | Type | Info |
|:---------|:--------|:--------|
| coin_name | string |  |
| amount | string |  |

**CLI Command** 
```console
NEAR_ENV=testnet near call <BTS> withdraw '{"coin_name": "<COIN NAME>", "amount":"<Amount in lowest Denomination>"}' --amount <1 yoctoNEAR in highest Denomination ie 0.000000000000000000000001> --gas 300000000000000 --accountId <ACCOUNT ID>
```

### Fee

#### Get Fee
**Method** 
- get_fee

| Parameters | Type | Info |
|:---------|:--------|:--------|
| coin_name | string |  |
| amount | string |  |

**CLI Command** 
```console
NEAR_ENV=testnet near view <BTS> get_fee '{"coin_name": "<Coin Name>", "amount": "<Amount in lowest Denomination, For NEAR in yoctoNEAR ie 1 NEAR = 1^24 yoctoNEAR>"}'
```

### Transfer
**Method** 
- transfer

| Parameters | Type | Info |
|:---------|:--------|:--------|
| coin_name | string |  |
| destination | string | valid btp address |
| amount | string |  |

**CLI Command** 
```console
NEAR_ENV=testnet near call <BTS> transfer '{"coin_name": "<Coin Name>", "destination": "btp://<Network>/<Address>", "amount": "<Amount in lowest Denomination, For NEAR in yoctoNEAR ie 1 NEAR = 1^24 yoctoNEAR>"}' --gas 300000000000000 --accountId <ACCOUNT ID>
```


## Usage

### Transfer NEAR to Cross-Chain 
1. Deposit NEAR to BTS [here](#deposit-near-to-bts)  
**Example** 
```console
NEAR_ENV=testnet near call bts.iconbridge-6.testnet deposit --amount 10 --accountId dev-20211206025826-24100687319598
```
2. Query Transfer Fee [here](#get-fee)  
**Example** 
```console
NEAR_ENV=testnet near view bts.iconbridge-6.testnet get_fee '{"coin_name": "btp-0x2.near-NEAR", "amount": "10000000000000000000000000"}'
```
3. Transfer [here](#transfer)  
**Example** 
```console
NEAR_ENV=testnet near call bts.iconbridge-6.testnet transfer '{"coin_name": "btp-0x2.near-NEAR", "destination": "btp://0x2.icon/hx54d9ba221fbe8a475a8bf38c7d048675b5d7b85a", "amount": "10000000000000000000000000"}' --gas 300000000000000 --accountId dev-20211206025826-24100687319598
```
4. Check locked balance for the transfered amount [here](#balance)  
**Example** 
```console
NEAR_ENV=testnet near view bts.iconbridge-6.testnet locked_balance_of '{"account_id": "dev-20211206025826-24100687319598", "coin_name": "btp-0x2.near-NEAR"}'
```

### Receiving NEAR from Cross-Chain
1. Check usable balance amount to withdraw [here](#balance)  
**Example** 
```console
NEAR_ENV=testnet near view bts.iconbridge-6.testnet balance_of '{"account_id": "dev-20211206025826-24100687319598", "coin_name": "btp-0x2.near-NEAR"}'
```
2. Withdraw [here](#withdraw)  
**Example** 
```console
NEAR_ENV=testnet near call bts.iconbridge-6.testnet withdraw '{"coin_name": "btp-0x2.near-NEAR", "amount":"10000000000000000000000000"}' --amount 0.000000000000000000000001 --gas 300000000000000 --accountId dev-20211206025826-24100687319598
```
3. Check balance  
**Example** 
```console
NEAR_ENV=testnet near state dev-20211206025826-24100687319598
```

### Transfer Cross-Chain Native Coins to Cross-Chain 
1. Deposit Cross-Chain Native Coin to BTS [here](#deposit-cross-chain-native-coin-to-bts)  
**Example** 
```console
NEAR_ENV=testnet near call btp-icx.bts.iconbridge-6.testnet ft_transfer_call '{"receiver_id": "bts.iconbridge-6.testnet", "amount": "10000000000000000000000000", "msg": ""}' --accountId dev-20211206025826-24100687319598 --amount 0.000000000000000000000001
```
2. Query Transfer Fee [here](#get-fee)  
**Example** 
```console
NEAR_ENV=testnet near view bts.iconbridge-6.testnet get_fee '{"coin_name": "btp-0x2.icon-ICX", "amount": "10000000000000000000000000"}'
```
3. Transfer [here](#transfer)  
**Example** 
```console
NEAR_ENV=testnet near call bts.iconbridge-6.testnet transfer '{"coin_name": "btp-0x2.icon-ICX", "destination": "btp://0x2.icon/hx54d9ba221fbe8a475a8bf38c7d048675b5d7b85a", "amount": "10000000000000000000000000"}' --gas 300000000000000 --accountId dev-20211206025826-24100687319598
```
4. Check locked Balance for the transfered amount [here](#locked-balance)  
**Example** 
```console
NEAR_ENV=testnet near view bts.iconbridge-6.testnet locked_balance_of '{"account_id": "dev-20211206025826-24100687319598", "coin_name": "btp-0x2.icon-ICX"}'
```

### Receiving Cross-Chain from Cross-Chain
1. Check usable balance amount to withdraw [here](#usable-balance)  
**Example** 
```console
NEAR_ENV=testnet near view bts.iconbridge-6.testnet balance_of '{"account_id": "dev-20211206025826-24100687319598", "coin_name": "btp-0x2.icon-ICX"}'
```
2. Withdraw [here](#withdraw)  
**Example** 
```console
NEAR_ENV=testnet near call bts.iconbridge-6.testnet withdraw '{"coin_name": "btp-0x2.icon-ICX", "amount":"10000000000000000000000000"}' --amount 0.000000000000000000000001 --gas 300000000000000 --accountId dev-20211206025826-24100687319598
```
3. Check balance  
**Example** 
```console
NEAR_ENV=testnet near view btp-icx.bts.iconbridge-6.testnet ft_balance_of '{"account_id": "dev-20211206025826-24100687319598"'
```

### Reclaiming Failed Transfer
1. Check refundable balance to reclaim [here](#refundable-balance)  
**Example** 
```console
NEAR_ENV=testnet near view bts.iconbridge-6.testnet refundable_balance_of '{"account_id": "dev-20211206025826-24100687319598", "coin_name": "btp-0x2.icon-ICX"}'
```
2. Reclaim  
**Example** 
```console
NEAR_ENV=testnet near call bts.iconbridge-6.testnet reclaim '{"coin_name": "btp-0x2.icon-ICX", "amount":"10000000000000000000000000"}' --amount 0.000000000000000000000001 --gas 300000000000000 --accountId dev-20211206025826-24100687319598
```

## Environment

- [Testnet](./testnet.md)


