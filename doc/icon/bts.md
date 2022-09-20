# ICON BTS Contract

## Transfer from ICON to other chains

* ### To transfer ICON native IRC2 Tokens
    - The token must be registered in ICON BTS and BTS of destination chain
    - Transfer that token to ICON BTS Contract
        ```py
        token.transfer(bts, amount, b"data")
        # bts       : Contract Address of ICON BTS
        # amount    : Amount to transfer
        ```
    - Call transfer method of BTS Contract with following parameters
        ```py
        bts.transfer(coinName, amount, to)
        # coinName  : registered name of the token
        # amount    : amount to transfer (fees is deducted from it)
        # to        : BTP Address of destination
        ```
* ### To transfer wrapped tokens of other chains on ICON to other chains
    - The token must be registered in BTS of ICON and destination chain
    - Approve amount to transfer to ICON BTS Contract
        ```py
        wrappedToken.approve(bts, amount)
        # bts       : Contract Address of ICON BTS
        # amount    : Amount to approve to BTS
        ```
    - Call transfer method of BTS Contract with following parameters
        ```py
        bts.transfer(coinName, amount, to)
        # coinName  : registered name of the token
        # amount    : amount to transfer (fees is deducted from it)
        # to        : BTP Address of destination
        ```
* ### To transfer nativecoin(ICX) to other chains
    - Nativecoin is registered by default on ICON, should be registered on destination chain
    - Call transferNativeCoin method and send required amount of ICX to transfer.
        ```py
        bts.transferNativeCoin(to)
        # to        : BTP Address of destination
        ```

* ### Transfer multiple types of coins at once (TransferBatch)
    - For IRC2 Tokens, transfer to BTS Contract as above.
    - For arapped tokens, approve to BTS Contract as above.
    - For nativecoin, send required amount of ICX to transfer in transferBatch method
    - Call transferBatch method with array of coinNames and amount to transfer
        ```py
        coinNames = ["btp-0x1.icon-sICX", "btp-0x37.bsc-BUSD"]
        coinAmounts = [100 * 10 ** 18, 100 * 10 ** 18]
        bts.transferBatch(coinNames, coinAmounts)
        # coinNames  : Names of coins to set token limit
        # coinAmounts: Amount to transfer
        ```

## Blacklist
* ### Blacklist users on any chain
    - To blacklist users on any chain, owner calls addBlacklistAddress method
    - Can blacklist a batch of addresses on any chain at once
        ```py
        net = "0x1.icon" # "0x61.bsc"
        addresses = [
            "hx1212121212121212121212121212121212121212",
            "hxabababababababababababababababababababab"
        ]
        bts.addBlacklistAddress(net, addresses)
        # net       : Network to blacklist on
        # addresses : Array of addresses of that network
        ```
    - These 2 addresses will now be blacklisted on ICON.

## Token Limit
* ### Set token limit
    - To set token limit for tokens (registered or not) owner calls setTokenLimit method
    - Can set limit for multiple coins across all connected networks
    - Makes token limit consistent across all connected networks
        ```py
        coinNames = ["btp-0x1.icon-sICX", "btp-0x37.bsc-BUSD"]
        tokenLimits = [1000000000000000000000, 3000000000000000000]
        bts.setTokenLimit(coinNames, tokenLimits)
        # coinNames  : Names of coins to set token limit
        # tokenLimits: Token limit to set
        ```
