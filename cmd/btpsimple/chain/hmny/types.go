//go:build hmny

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

type BlockV2WithTxHash struct {
	Number           *big.Int       `json:"number"`
	ViewID           *big.Int       `json:"viewID"`
	Epoch            *big.Int       `json:"epoch"`
	Hash             common.Hash    `json:"hash"`
	ParentHash       common.Hash    `json:"parentHash"`
	Nonce            uint64         `json:"nonce"`
	MixHash          common.Hash    `json:"mixHash"`
	LogsBloom        ethtypes.Bloom `json:"logsBloom"`
	StateRoot        common.Hash    `json:"stateRoot"`
	Miner            string         `json:"miner"`
	Difficulty       uint64         `json:"difficulty"`
	ExtraData        hexutil.Bytes  `json:"extraData"`
	Size             uint64         `json:"size"`
	GasLimit         uint64         `json:"gasLimit"`
	GasUsed          uint64         `json:"gasUsed"`
	VRF              common.Hash    `json:"vrf"`
	VRFProof         hexutil.Bytes  `json:"vrfProof"`
	Timestamp        *big.Int       `json:"timestamp"`
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
