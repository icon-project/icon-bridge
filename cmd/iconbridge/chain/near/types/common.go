package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/near/borsh-go"
)

const (
	ED25519   = 0
	SECP256K1 = 1
)

type AccountId string

type CryptoHash []byte

func (c *CryptoHash) UnmarshalJSON(p []byte) error {
	var cryptoHash string
	err := json.Unmarshal(p, &cryptoHash)
	if err != nil {
		return err
	}
	if cryptoHash == "" {
		*c = nil
		return nil
	}
	*c = CryptoHash(base58.Decode(cryptoHash))
	return nil
}

func (c *CryptoHash) Base58Encode() string {
	return base58.Encode(*c)
}

type Timestamp uint64

func (t *Timestamp) UnmarshalJSON(p []byte) error {
	var timestamp string
	err := json.Unmarshal(p, &timestamp)
	if err != nil {
		return err
	}
	if timestamp == "" {
		return nil
	}

	n, err := strconv.ParseUint(timestamp, 10, 0)
	if err != nil {
		return fmt.Errorf("not a uint64: %s %s", timestamp, err)
	}
	*t = Timestamp(n)
	return nil
}

type PublicKey struct {
	KeyType uint8
	Data    []byte
}

func (pk *PublicKey) UnmarshalJSON(p []byte) error {
	var publicKey string
	err := json.Unmarshal(p, &publicKey)
	if err != nil {
		return err
	}

	if publicKey == "" {
		pk = nil
		return nil
	}

	if strings.Contains(publicKey, "ed25519:") {
		pk = &PublicKey{
			KeyType: ED25519,
			Data:    base58.Decode(publicKey[len("ed25519:"):]),
		}
	} else if strings.Contains(publicKey, "secp256k1:") {
		pk = &PublicKey{
			KeyType: SECP256K1,
			Data:    base58.Decode(publicKey[len("secp256k1:"):]),
		}
	} else {
		pk = nil
	}
	return nil
}

type Signature struct {
	KeyType uint8
	Data    []byte
}

func (s Signature) Base58Encode() string {
	if s.KeyType == ED25519 {
		return "ed25519:" + base58.Encode(s.Data[:])
	} else {
		return "secp256k1:" + base58.Encode(s.Data[:])
	}
}

func (s *Signature) Bytes() []byte {
	if s.KeyType == ED25519 {
		return append([]byte{ED25519}, s.Data...)
	} else {
		return append([]byte{SECP256K1}, s.Data...)
	}
}

func (s *Signature) UnmarshalJSON(p []byte) error {
	var signature string
	err := json.Unmarshal(p, &signature)
	if err != nil {
		return err
	}

	if signature == "" {
		s = nil
		return nil
	}

	if strings.Contains(signature, "ed25519:") {
		s = &Signature{
			KeyType: ED25519,
			Data:    base58.Decode(signature[len("ed25519:"):]),
		}
	} else if strings.Contains(signature, "secp256k1:") {
		s = &Signature{
			KeyType: SECP256K1,
			Data:    base58.Decode(signature[len("secp256k1:"):]),
		}
	} else {
		s = nil
	}
	return nil
}

type BigInt string

func (b *BigInt) Int() (big.Int, error) {
	n := new(big.Int)
	n, ok := n.SetString(string(*b), 10)
	if !ok {
		return big.Int{}, fmt.Errorf("not a valid bigint: %s", string(*b))
	}
	return *n, nil
}

func CombineHash(hash1 []byte, hash2 []byte) []byte {
	combined := new(bytes.Buffer)
	combined.Write(hash1[:])
	combined.Write(hash2[:])

	hash := sha256.Sum256(combined.Bytes())

	return hash[:]
}

type MerklePathItem struct {
	Hash      CryptoHash `json:"hash"`
	Direction string     `json:"direction"`
}

type MerklePath []MerklePathItem

type ExecutionStatus struct {
	SuccessValue     string     `json:"SuccessValue"`
	SuccessReceiptId CryptoHash `json:"SuccessReceiptId"`
	Failure          Failure    `json:"Failure"`
	Unknown          string     `json:"Unknown"`
}

// TODO: Add More Errors
type Failure struct {
	ActionError ActionError `json:"ActionError"`
}

type ActionError struct {
	Index uint64 `json:"index"`
	Kind  Kind   `json:"kind"`
}

type Kind struct {
	FunctionCallError   FunctionCallError   `json:"FunctionCallError"`
	AccountDoesNotExist AccountDoesNotExist `json:"AccountDoesNotExist"`
}

type FunctionCallError struct {
	ExecutionError string `json:"ExecutionError"`
}

type AccountDoesNotExist struct {
	AccountId AccountId `json:"account_id"`
}

type CallFunction struct {
	RequestType  string    `json:"request_type"`
	Finality     string    `json:"finality"`
	AccountId    AccountId `json:"account_id"`
	MethodName   string    `json:"method_name"`
	ArgumentsB64 string    `json:"args_base64"`
}

type ExecutionOutcomeWithIdView struct {
	Proofs    MerklePath           `json:"proof"`
	BlockHash CryptoHash           `json:"block_hash"`
	ReceiptId CryptoHash           `json:"id"`
	Outcome   ExecutionOutcomeView `json:"outcome"`
}

type ExecutionOutcomeView struct {
	Logs        []string        `json:"logs"`
	ReceiptIds  []CryptoHash    `json:"receipt_ids"`
	GasBurnt    uint64          `json:"gas_burnt"`
	TokensBurnt string          `json:"tokens_burnt"`
	ExecutorId  string          `json:"executor_id"`
	Status      ExecutionStatus `json:"status"`
}

type ApprovalMessage struct {
	Type                []byte
	PreviousBlockHash   CryptoHash
	PreviousBlockHeight int64
	TargetHeight        int64
}

type Action struct {
	Enum           borsh.Enum `borsh_enum:"true"` // treat struct as complex enum when serializing/deserializing
	CreateAccount  borsh.Enum
	DeployContract DeployContract
	FunctionCall   FunctionCall
	Transfer       Transfer
	Stake          Stake
	AddKey         AddKey
	DeleteKey      DeleteKey
	DeleteAccount  DeleteAccount
}

type DeployContract struct {
	Code []byte
}

type FunctionCall struct {
	MethodName string
	Args       []byte
	Gas        uint64
	Deposit    big.Int
}

type Transfer struct {
	Deposit big.Int
}

type Stake struct {
	Stake     big.Int
	PublicKey PublicKey
}

type AddKey struct {
	PublicKey PublicKey
	AccessKey AccessKey
}

type DeleteKey struct {
	PublicKey PublicKey
}

type DeleteAccount struct {
	BeneficiaryID string
}

type AccessKey struct {
	Nonce      uint64
	Permission AccessKeyPermission
}

type AccessKeyPermission struct {
	Enum         borsh.Enum `borsh_enum:"true"`
	FunctionCall FunctionCallPermission
	FullAccess   borsh.Enum
}

type FunctionCallPermission struct {
	Allowance   *big.Int
	ReceiverId  string
	MethodNames []string
}
