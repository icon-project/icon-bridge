# Solidity BTS Contract

Due to the contract size limitation, the BTS Contract has been split into 2 contracts. BTSCore and BTSPeriphery. The transfer logic originates on BTSCore contract. 

## Transfer from BSC to other chains

* ### Transfer ERC20 tokens and wrapped tokens
    - The token should be registered on both BSC and destination chain
    - Approve amount to transfer to BTS Core contract.
        ```py
        erc20.approve(btsCore, amount)

        # btsCore   : Contract address of BTS Core
        # amount    : Amount to approve to BTS Core for transfer
        ```
    - Call transfer method of BTS Core contract
        ```py
        btsCore.transfer(coinName, value, to)

        # coinName  : Registered name for token
        # value     : Amount to transfer
        # to        : BTP Address of destination    
        ```
* ### Transfer BNB
    - BNB should be registered on destination chain
    - Call transferNativeCoin payable method and send amount to transfer in it
        ```py
        btsCore.transferNativeCoin(to)

        # to        : BTP Address of destination
        ```
* ### Transfer multiple types of coins at once (TransferBatch)
    - For tokens (ERC20 and wrapped), approve amount to transfer
    - For nativecoin, send amount to transfer
       ```py
        coinNames = ["btp-0x1.icon-sICX", "btp-0x37.bsc-BUSD"]
        coinAmounts = [100 * 10 ** 18, 100 * 10 ** 18]
        btsCore.transferBatch(coinNames, coinAmounts)

        # coinNames  : Names of coins to set token limit
        # coinAmounts: Amount to transfer
        ```