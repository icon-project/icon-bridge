package hmny

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/harmony/core/types"
)

type BlockNotification struct {
	Hash     common.Hash
	Height   *big.Int
	Header   *Header
	Receipts types.Receipts
}

// harmony types
type BlockWithTxHash struct {
	Number           *hexutil.Big   `json:"number"`
	ViewID           *hexutil.Big   `json:"viewID"`
	Epoch            *hexutil.Big   `json:"epoch"`
	Hash             common.Hash    `json:"hash"`
	ParentHash       common.Hash    `json:"parentHash"`
	Nonce            uint64         `json:"nonce"`
	MixHash          common.Hash    `json:"mixHash"`
	LogsBloom        ethtypes.Bloom `json:"logsBloom"`
	StateRoot        common.Hash    `json:"stateRoot"`
	Miner            string         `json:"miner"`
	Difficulty       uint64         `json:"difficulty"`
	ExtraData        hexutil.Bytes  `json:"extraData"`
	Size             hexutil.Uint64 `json:"size"`
	GasLimit         hexutil.Uint64 `json:"gasLimit"`
	GasUsed          hexutil.Uint64 `json:"gasUsed"`
	VRF              common.Hash    `json:"vrf"`
	VRFProof         hexutil.Bytes  `json:"vrfProof"`
	Timestamp        hexutil.Uint64 `json:"timestamp"`
	TransactionsRoot common.Hash    `json:"transactionsRoot"`
	ReceiptsRoot     common.Hash    `json:"receiptsRoot"`
	Uncles           []common.Hash  `json:"uncles"`
	Transactions     []common.Hash  `json:"transactions"`
	EthTransactions  []common.Hash  `json:"transactionsInEthHash"`
	StakingTxs       []common.Hash  `json:"stakingTransactions"`
	Signers          []string       `json:"signers,omitempty"`
}

type Header struct {
	ParentHash           common.Hash    `json:"parentHash"`
	Miner                common.Address `json:"miner"`
	StateRoot            common.Hash    `json:"stateRoot"`
	TransactionsRoot     common.Hash    `json:"transactionsRoot"`
	ReceiptsRoot         common.Hash    `json:"receiptsRoot"`
	OutgoingReceiptsRoot common.Hash    `json:"outgoingReceiptsRoot"`
	IncomingReceiptsRoot common.Hash    `json:"incomingReceiptsRoot"`
	LogsBloom            ethtypes.Bloom `json:"logsBloom"`
	Number               *big.Int       `json:"number"`
	GasLimit             uint64         `json:"gasLimit"`
	GasUsed              uint64         `json:"gasUsed"`
	Timestamp            *big.Int       `json:"timestamp"`
	ExtraData            hexutil.Bytes  `json:"extraData"`
	MixHash              common.Hash    `json:"mixHash"`
	ViewID               *big.Int       `json:"viewID"`
	Epoch                *big.Int       `json:"epoch"`
	ShardID              uint32         `json:"shardID"`
	LastCommitSignature  hexutil.Bytes  `json:"lastCommitSignature"`
	LastCommitBitmap     hexutil.Bytes  `json:"lastCommitBitmap"`
	Vrf                  hexutil.Bytes  `json:"vrf"`
	Vdf                  hexutil.Bytes  `json:"vdf"`
	ShardState           hexutil.Bytes  `json:"shardState"`
	CrossLink            hexutil.Bytes  `json:"crossLink"`
	Slashes              hexutil.Bytes  `json:"slashes"`
}

func (h *Header) Hash() common.Hash {
	return common.BytesToHash(crypto.Keccak256(h.RLPMarshalToBytes()))
}

func (h *Header) RLPMarshalToBytes() []byte {
	rlph, _ := rlp.EncodeToBytes([]interface{}{"HmnyTgd", "v3", h})
	return rlph
}

// const (
// 	TransferStartSignature    = "TransferStart"
// 	TransferEndSignature      = "TransferEnd"
// 	TransferReceivedSignature = "TransferReceived"
// )

type LogResult struct {
	TxHash   common.Hash
	LogIndex uint
	Address  common.Address
	Topic    string
	Logs     interface{}
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

type TransferStart struct {
	From         common.Address
	To           string
	Sn           *big.Int
	AssetDetails []struct {
		CoinName string
		Value    *big.Int
		Fee      *big.Int
	}
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
