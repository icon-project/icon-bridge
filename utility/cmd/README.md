# Utility

## Requirements

* Bazel

  **Mac OSX**
    ```
    brew install bazelisk
    ```

  **Windows**
    ```
    choco install bazelisk
    ```

## Build


* Goloop
    ```
    bazel build @com_github_icon_project_goloop//cmd/goloop:goloop
    ```
* NEAR CLI
    ```
    bazel build @near//cli:near_binary --define near_network=testnet
    ```

## Config
```json
[
    {
        "name": "near",
        "key_path": "<PATH>",
        "key_secret": "",
        "sender": "btp-16.testnet",
        "native_coin": "<NATIVE COIN>",
        "bts": "bts.iconbridge.testnet",
        "step_limit": 0,
        "network_id": "0x1",
        "uri": "https://rpc.testnet.near.org",
        "cli": "<PATH>"
    },
    {
        "name": "icon",
        "key_path": "<PATH>",
        "key_secret": "<PATH>",
        "sender": "btp-16.testnet",
        "native_coin": "<NATIVE COIN>",
        "bts": "cx95882bb6a0fda402afc09a52a0141738de8fa133",
        "step_limit": 13610920001,
        "network_id": "0x2",
        "uri": "https://lisbon.net.solidwallet.io/api/v3/icon_dex",
        "cli": "<PATH>"
    }
]
```
## Usage
```
transferToken (near|icon) <AMOUNT> <RECEIVER> --from (near|icon) --config <PATH TO CONFIG>
```