# Solidity BTS Contract

## Prerequisities
- Truffle
    ```sh
    npm install -g truffle
    ```
- Set up truffle console to connect to mainnet/testnet/local. Check truffle-config.js for details.
    ```sh
    truffle console --network bsc
    ```

## Transfer from BSC to other chains

* ### Transfer ERC20 tokens and wrapped tokens
    - The token should be registered on both BSC and destination chain
    - Approve amount to transfer to BTS Core contract.
        ```js
        const erc20 = await ERC20.at(token)
        await erc20.approve(btsCore, amount)
        
        // token     : Address of the token
        // btsCore   : Contract address of BTS Core
        // amount    : Amount to approve to BTS Core for transfer
        ```
    - Call transfer method of BTS Core contract
        ```js
        const btsCore = await BTSCore.deployed()
        await btsCore.transfer(coinName, value, to)

        // coinName  : Registered name for token
        // value     : Amount to transfer
        // to        : BTP Address of destination    
        ```
* ### Transfer BNB
    - BNB should be registered on destination chain
    - Call transferNativeCoin payable method and send amount to transfer in it
        ```js
        const btsCore = await BTSCore.deployed()
        await btsCore.transferNativeCoin(to)

        // to        : BTP Address of destination
        ```
* ### Transfer multiple types of coins at once (TransferBatch)
    - For tokens (ERC20 and wrapped), approve amount to transfer
    - For nativecoin, send amount to transfer
       ```js
        const btsCore = await BTSCore.deployed()
        coinNames = ["btp-0x1.icon-sICX", "btp-0x37.bsc-BUSD"]
        coinAmounts = [100 * 10 ** 18, 100 * 10 ** 18]
        await btsCore.transferBatch(coinNames, coinAmounts)

        // coinNames  : Names of coins to set token limit
        // coinAmounts: Amount to transfer
        ```