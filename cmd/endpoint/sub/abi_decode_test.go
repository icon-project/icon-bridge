package sub

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/icon-bridge/common"
)

func TestTopic(t *testing.T) {
	const (
		TransferStartSignature         = "TransferStart(address,string,uint256,(string,uint256,uint256)[])"    //0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a
		TransferEndSignature           = "TransferEnd(address,uint256,uint256,string)"                         //0x9b4c002cf17443998e01f132ae99b7392665eec5422a33a1d2dc47308c59b6e2
		TransferReceivedSignature      = "TransferReceived(string,address,uint256,(string,uint256)[])"         //0x78e3e55e26c08e043fbd9cc0282f53e2caab096d30594cb476fcdfbbe7ce8680
		TransferReceivedSignatureToken = "TransferReceived(string,address,uint256,(string,uint256,uint256)[])" //0xd2221859bf6855d034602a0388473f88313afe64aa63f26788e51caa087ed15c
	)

	fmt.Println(crypto.Keccak256Hash([]byte(TransferStartSignature)))
	fmt.Println(crypto.Keccak256Hash([]byte(TransferEndSignature)))
	fmt.Println(crypto.Keccak256Hash([]byte(TransferReceivedSignature)))
	fmt.Println(crypto.Keccak256Hash([]byte(TransferReceivedSignatureToken)))
}

func TestAbiDecode(t *testing.T) {
	abi, err := abi.JSON(strings.NewReader(bshPeripherABI))
	if err != nil {
		t.Fatal(err)
	}
	const transferEndHex = "000000000000000000000000000000000000000000000000000000000000002e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000105472616e73666572205375636365737300000000000000000000000000000000"
	const transferReceivedHex = "000000000000000000000000000000000000000000000000000000000000002e00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000dbd2fc137a3000000000000000000000000000000000000000000000000000000000000000000034f4e450000000000000000000000000000000000000000000000000000000000"
	data, err := hex.DecodeString(transferReceivedHex)
	if err != nil {
		t.Fatal(err)
	}

	type TransferEnd struct {
		From     common.Address
		Sn       *big.Int
		Code     *big.Int
		Response string
	}

	type TransferReceived struct {
		From         string
		To           common.Address
		Sn           *big.Int
		AssetDetails []struct {
			CoinName string
			Value    *big.Int
		}
	}

	var ev TransferReceived
	err = abi.UnpackIntoInterface(&ev, "TransferReceived", data)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(ev)
}

func TestAbiDecodeTransferStart(t *testing.T) {
	abi, err := abi.JSON(strings.NewReader(bshPeripherABI))
	if err != nil {
		t.Fatal(err)
	}
	const transferStartHex = "0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000002f00000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000003e6274703a2f2f30783562396137372e69636f6e2f68783734356634333239636235336166313238376662343561613437636330383635646636366161313600000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000001b7a5f826f46000000000000000000000000000000000000000000000000000000470de4df82000000000000000000000000000000000000000000000000000000000000000000034f4e450000000000000000000000000000000000000000000000000000000000"

	data, err := hex.DecodeString(transferStartHex)
	if err != nil {
		t.Fatal(err)
	}
	type AssetTransferDetail struct {
		CoinName string
		Value    *big.Int
		Fee      *big.Int
	}
	type TransferStart struct {
		From         common.Address
		To           string
		Sn           *big.Int
		AssetDetails []AssetTransferDetail
	}

	var ev TransferStart

	err = abi.UnpackIntoInterface(&ev, "TransferStart", data)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(ev)
	// ev := map[string]interface{}{}
	// err = abi.UnpackIntoMap(ev, "TransferStart", data)
	// if err != nil {
	// 	t.Error(err)
	// }
	// fmt.Println(ev)
}

const bshPeripherABI = `[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "_from",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "_sn",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "_code",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "string",
          "name": "_response",
          "type": "string"
        }
      ],
      "name": "TransferEnd",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "string",
          "name": "_from",
          "type": "string"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "_to",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "_sn",
          "type": "uint256"
        },
        {
          "components": [
            {
              "internalType": "string",
              "name": "coinName",
              "type": "string"
            },
            {
              "internalType": "uint256",
              "name": "value",
              "type": "uint256"
            }
          ],
          "indexed": false,
          "internalType": "struct Types.Asset[]",
          "name": "_assetDetails",
          "type": "tuple[]"
        }
      ],
      "name": "TransferReceived",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "_from",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "string",
          "name": "_to",
          "type": "string"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "_sn",
          "type": "uint256"
        },
        {
          "components": [
            {
              "internalType": "string",
              "name": "coinName",
              "type": "string"
            },
            {
              "internalType": "uint256",
              "name": "value",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "fee",
              "type": "uint256"
            }
          ],
          "indexed": false,
          "internalType": "struct Types.AssetTransferDetail[]",
          "name": "_assetDetails",
          "type": "tuple[]"
        }
      ],
      "name": "TransferStart",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "string",
          "name": "_from",
          "type": "string"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "_sn",
          "type": "uint256"
        }
      ],
      "name": "UnknownResponse",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "requests",
      "outputs": [
        {
          "internalType": "string",
          "name": "from",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "to",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [],
      "name": "serviceName",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_bmc",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "_bshCore",
          "type": "address"
        },
        {
          "internalType": "string",
          "name": "_serviceName",
          "type": "string"
        }
      ],
      "name": "initialize",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "hasPendingRequest",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_from",
          "type": "address"
        },
        {
          "internalType": "string",
          "name": "_to",
          "type": "string"
        },
        {
          "internalType": "string[]",
          "name": "_coinNames",
          "type": "string[]"
        },
        {
          "internalType": "uint256[]",
          "name": "_values",
          "type": "uint256[]"
        },
        {
          "internalType": "uint256[]",
          "name": "_fees",
          "type": "uint256[]"
        }
      ],
      "name": "sendServiceMessage",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_from",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "_svc",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "_sn",
          "type": "uint256"
        },
        {
          "internalType": "bytes",
          "name": "_msg",
          "type": "bytes"
        }
      ],
      "name": "handleBTPMessage",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "_svc",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "_sn",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_code",
          "type": "uint256"
        },
        {
          "internalType": "string",
          "name": "_msg",
          "type": "string"
        }
      ],
      "name": "handleBTPError",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_to",
          "type": "string"
        },
        {
          "components": [
            {
              "internalType": "string",
              "name": "coinName",
              "type": "string"
            },
            {
              "internalType": "uint256",
              "name": "value",
              "type": "uint256"
            }
          ],
          "internalType": "struct Types.Asset[]",
          "name": "_assets",
          "type": "tuple[]"
        }
      ],
      "name": "handleRequestService",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_fa",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "_svc",
          "type": "string"
        }
      ],
      "name": "handleFeeGathering",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_to",
          "type": "string"
        }
      ],
      "name": "checkParseAddress",
      "outputs": [],
      "stateMutability": "pure",
      "type": "function",
      "constant": true
    }
  ]`

const bshCoreABI = `[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "remover",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "formerOwner",
          "type": "address"
        }
      ],
      "name": "RemoveOwnership",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "promoter",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "newOwner",
          "type": "address"
        }
      ],
      "name": "SetOwnership",
      "type": "event"
    },
    {
      "inputs": [],
      "name": "feeNumerator",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [],
      "name": "fixedFee",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_nativeCoinName",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "_feeNumerator",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_fixedFee",
          "type": "uint256"
        }
      ],
      "name": "initialize",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_owner",
          "type": "address"
        }
      ],
      "name": "addOwner",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_owner",
          "type": "address"
        }
      ],
      "name": "removeOwner",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_owner",
          "type": "address"
        }
      ],
      "name": "isOwner",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [],
      "name": "getOwners",
      "outputs": [
        {
          "internalType": "address[]",
          "name": "",
          "type": "address[]"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_bshPeriphery",
          "type": "address"
        }
      ],
      "name": "updateBSHPeriphery",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "_feeNumerator",
          "type": "uint256"
        }
      ],
      "name": "setFeeRatio",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "_fixedFee",
          "type": "uint256"
        }
      ],
      "name": "setFixedFee",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_name",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "_symbol",
          "type": "string"
        },
        {
          "internalType": "uint8",
          "name": "_decimals",
          "type": "uint8"
        }
      ],
      "name": "register",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "coinNames",
      "outputs": [
        {
          "internalType": "string[]",
          "name": "_names",
          "type": "string[]"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_coinName",
          "type": "string"
        }
      ],
      "name": "coinId",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_coinName",
          "type": "string"
        }
      ],
      "name": "isValidCoin",
      "outputs": [
        {
          "internalType": "bool",
          "name": "_valid",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_owner",
          "type": "address"
        },
        {
          "internalType": "string",
          "name": "_coinName",
          "type": "string"
        }
      ],
      "name": "getBalanceOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "_usableBalance",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_lockedBalance",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_refundableBalance",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_owner",
          "type": "address"
        },
        {
          "internalType": "string[]",
          "name": "_coinNames",
          "type": "string[]"
        }
      ],
      "name": "getBalanceOfBatch",
      "outputs": [
        {
          "internalType": "uint256[]",
          "name": "_usableBalances",
          "type": "uint256[]"
        },
        {
          "internalType": "uint256[]",
          "name": "_lockedBalances",
          "type": "uint256[]"
        },
        {
          "internalType": "uint256[]",
          "name": "_refundableBalances",
          "type": "uint256[]"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [],
      "name": "getAccumulatedFees",
      "outputs": [
        {
          "components": [
            {
              "internalType": "string",
              "name": "coinName",
              "type": "string"
            },
            {
              "internalType": "uint256",
              "name": "value",
              "type": "uint256"
            }
          ],
          "internalType": "struct Types.Asset[]",
          "name": "_accumulatedFees",
          "type": "tuple[]"
        }
      ],
      "stateMutability": "view",
      "type": "function",
      "constant": true
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_to",
          "type": "string"
        }
      ],
      "name": "transferNativeCoin",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function",
      "payable": true
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_coinName",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "_value",
          "type": "uint256"
        },
        {
          "internalType": "string",
          "name": "_to",
          "type": "string"
        }
      ],
      "name": "transfer",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string[]",
          "name": "_coinNames",
          "type": "string[]"
        },
        {
          "internalType": "uint256[]",
          "name": "_values",
          "type": "uint256[]"
        },
        {
          "internalType": "string",
          "name": "_to",
          "type": "string"
        }
      ],
      "name": "transferBatch",
      "outputs": [],
      "stateMutability": "payable",
      "type": "function",
      "payable": true
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_coinName",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "_value",
          "type": "uint256"
        }
      ],
      "name": "reclaim",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_to",
          "type": "address"
        },
        {
          "internalType": "string",
          "name": "_coinName",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "_value",
          "type": "uint256"
        }
      ],
      "name": "refund",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_to",
          "type": "address"
        },
        {
          "internalType": "string",
          "name": "_coinName",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "_value",
          "type": "uint256"
        }
      ],
      "name": "mint",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "_requester",
          "type": "address"
        },
        {
          "internalType": "string",
          "name": "_coinName",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "_value",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_fee",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "_rspCode",
          "type": "uint256"
        }
      ],
      "name": "handleResponseService",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_fa",
          "type": "string"
        }
      ],
      "name": "transferFees",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`
