## Requirements

* **NEAR CLI**  
    The NEAR [Command Line Interface (CLI)](https://github.com/near/near-cli) is a tool that enables to interact with the NEAR network directly from the shell. Under the hood, NEAR CLI utilizes the NEAR [JavaScript API](https://github.com/near/near-api-js)


    #### Installation
    ```console
    npm install -g near-cli
    ```
* **NEAR Wallet**  
    - Testnet: https://wallet.testnet.near.org/
    - Mainnet: https://wallet.near.org/

* **Authorize CLI**  
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
2. Transfer [here](#transfer)
3. Check locked balance for the transfered amount [here](#balance)

### Receiving NEAR from Cross-Chain
1. Check usable balance amount to withdraw [here](#balance)
2. Withdraw [here](#withdraw)
3. Check balance
```console
NEAR_ENV=testnet near state <ACCOUNT ID>
```

### Transfer Cross-Chain Native Coins to Cross-Chain 
1. Deposit Cross-Chain Native Coin to BTS [here](#deposit-cross-chain-native-coin-to-bts)
2. Transfer [here](#transfer)
3. Check locked Balance for the transfered amount [here](#locked-balance)

### Receiving Cross-Chain from Cross-Chain
1. Check usable balance amount to withdraw [here](#usable-balance)
2. Withdraw [here](#withdraw)
3. Check balance
```console
NEAR_ENV=testnet near view <NEP141 Contract> ft_balance_of '{"account_id": "<ACCOUNT ID>"'
```

### Reclaiming Failed Transfer
1. Check refundable balance to reclaim [here](#refundable-balance)
2. Reclaim  
```console
NEAR_ENV=testnet near call <BTS> reclaim '{"coin_name": "<Coin Name>", "amount":"<Amount in lowest Denomination>"}' --amount 0.000000000000000000000001 --gas 300000000000000 --accountId <ACCOUNT ID>
```

## Environment

- [Testnet](./testnet.md)


