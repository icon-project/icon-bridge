## Requirements

* **NEAR CLI**  
    The NEAR [Command Line Interface (CLI)](https://github.com/near/near-cli) is a tool that enables to interact with the NEAR network directly from the shell. Under the hood, NEAR CLI utilizes the NEAR [JavaScript API](https://github.com/near/near-api-js)


    #### Installation
    ```npm install -g near-cli```
* **NEAR Wallet**  
    - Testnet: https://wallet.testnet.near.org/
    - Mainnet: https://wallet.near.org/

* **Authorize CLI**  
    ```near login```
## API
*Path to BTS Rust ReadMe*

## Usage

**Transfer NEAR**  
- Deposit  
```console
NEAR_ENV=testnet near call <BTS> deposit --amount <AMOUNT in NEAR> --accountId <ACCOUNT ID>
```
- Transfer  
```console
NEAR_ENV=testnet near call <BTS> transfer '{"coin_name": "<Registered Coin Name>", "destination": "btp://<Network>/<Address>", "amount": "<Amount in lowest Denomination, For NEAR in yoctoNEAR ie 1 NEAR = 1^24 yoctoNEAR>"}' --gas 300000000000000 --accountId <ACCOUNT ID>
```

**Receiving Cross-Chain Native Coins**
- Withdraw  
```console
NEAR_ENV=testnet near call <BTS> withdraw '{"coin_name": "<Registered Coin Name>", "amount":"<Amount in lowest Denomination>"}' --amount 0.000000000000000000000001 --gas 300000000000000 --accountId <ACCOUNT ID>
```

**Reclaiming Failed Transfer**
- Reclaim  
```console
NEAR_ENV=testnet near call <BTS> reclaim '{"coin_name": "<Registered Coin Name>", "amount":"<Amount in lowest Denomination>"}' --amount 0.000000000000000000000001 --gas 300000000000000 --accountId <ACCOUNT ID>
```

## Environment

- [Testnet](./testnet.md)


