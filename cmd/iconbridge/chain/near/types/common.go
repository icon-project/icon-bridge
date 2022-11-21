package types

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/near/borsh-go"
	"github.com/shopspring/decimal"
)

const (
	ED25519   = 0
	SECP256K1 = 1
)

type Height uint64

func (h *Height) Int() int {
	if *h > math.MaxInt {
		panic("overflow")
	}

	return int(*h)
}

type AccountId string

type CryptoHash [32]byte

func (c *CryptoHash) UnmarshalJSON(p []byte) error {
	var cryptoHash string
	var data [32]byte

	err := json.Unmarshal(p, &cryptoHash)
	if err != nil {
		return err
	}
	if cryptoHash == "" {
		c = nil
		return nil
	}

	copy(data[:], base58.Decode(cryptoHash))
	*c = CryptoHash(data)

	return nil
}

func NewCryptoHash(hash string) CryptoHash {
	var data [32]byte
	copy(data[:], base58.Decode(hash))

	return CryptoHash(data)
}

func (c *CryptoHash) Base58Encode() string {
	return base58.Encode(c[:])
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
	Data    [32]byte
}

func (pk *PublicKey) UnmarshalJSON(p []byte) error {
	var s string

	err := json.Unmarshal(p, &s)
	if err != nil {
		return err
	}

	if s == "" {
		pk = nil
		return nil
	}

	*pk = NewPublicKeyFromString(s)

	return nil
}

func (pk *PublicKey) Base58Encode() string {
	if pk.KeyType == ED25519 {
		return "ed25519:" + base58.Encode(pk.Data[:])
	} else {
		return "secp256k1:" + base58.Encode(pk.Data[:])
	}
}

func NewPublicKeyFromString(s string) PublicKey {
	var data [32]byte
	var keyType uint8
	if strings.Contains(s, "ed25519:") {
		copy(data[:], base58.Decode(s[len("ed25519:"):]))
		keyType = ED25519
	} else if strings.Contains(s, "secp256k1:") {
		copy(data[:], base58.Decode(s[len("secp256k1:"):]))
		keyType = SECP256K1
	}

	return PublicKey{
		KeyType: keyType,
		Data:    data,
	}
}

func NewPublicKeyFromED25519(pk ed25519.PublicKey) PublicKey {
	var data [32]byte
	copy(data[:], pk)

	return PublicKey{
		KeyType: ED25519,
		Data:    data,
	}
}

type Signature struct {
	KeyType uint8
	Data    [64]byte
}

func NewSignatureFromString(s string) Signature {
	var data [64]byte
	var keyType uint8

	if strings.Contains(s, "ed25519:") {
		copy(data[:], base58.Decode(s[len("ed25519:"):]))
		keyType = ED25519
	} else if strings.Contains(s, "secp256k1:") {
		copy(data[:], base58.Decode(s[len("secp256k1:"):]))
		keyType = SECP256K1
	}

	return Signature{
		KeyType: keyType,
		Data:    data,
	}
}

func (s Signature) Base58Encode() string {
	if s.KeyType == ED25519 {
		return "ed25519:" + base58.Encode(s.Data[:])
	} else {
		return "secp256k1:" + base58.Encode(s.Data[:])
	}
}

func (s *Signature) Bytes() []byte {
	return append([]byte{s.KeyType}, s.Data[:]...)
}

func (sig *Signature) UnmarshalJSON(p []byte) error {
	var s string

	err := json.Unmarshal(p, &s)
	if err != nil {
		return err
	}

	if s == "" {
		sig = nil
		return nil
	}

	*sig = NewSignatureFromString(s)

	return nil
}

type BigInt big.Int

func NewBigInt(bigInt string) BigInt {
	dec, err := decimal.NewFromString(bigInt)
	if err != nil {
		return BigInt(big.Int{})
	}

	return BigInt(*dec.BigInt())
}

func (b *BigInt) UnmarshalJSON(p []byte) error {
	var bigInt string
	err := json.Unmarshal(p, &bigInt)
	if err != nil {
		return err
	}

	if bigInt == "" {
		return nil
	}

	dec, err := decimal.NewFromString(bigInt)
	if err != nil {
		return err
	}

	*b = BigInt(*dec.BigInt())

	return nil
}

func (b BigInt) Serialize() ([]byte, error) {
	return borsh.Serialize(big.Int(b))
}

type SuccessValue []byte

func (s *SuccessValue) UnmarshalJSON(p []byte) error {
	var v string

	err := json.Unmarshal(p, &v)
	if err != nil {
		return err
	}

	if v == "" {
		s = nil
		return nil
	}

	data, err := base64.URLEncoding.Strict().DecodeString(v)
	if err != nil {
		return err
	}

	*s = data

	return nil
}

func CombineHash(hash1 [32]byte, hash2 [32]byte) [32]byte {
	combined := new(bytes.Buffer)
	combined.Write(hash1[:])
	combined.Write(hash2[:])

	hash := sha256.Sum256(combined.Bytes())

	return hash
}

type MerklePathItem struct {
	Hash      CryptoHash `json:"hash"`
	Direction string     `json:"direction"`
}

type MerklePath []MerklePathItem

type ExecutionStatus struct {
	Unknown          string       `json:"Unknown"`
	Failure          Failure      `json:"Failure"`
	SuccessValue     SuccessValue `json:"SuccessValue"`
	SuccessReceiptId CryptoHash   `json:"SuccessReceiptId"`
}

func (e ExecutionStatus) Serialize() ([]byte, error) {
	if e.Failure != *new(Failure) {
		return borsh.Serialize(struct {
			Enum    uint8
			Failure Failure
		}{
			Enum:    1,
			Failure: e.Failure,
		})
	} else if e.SuccessValue != nil {
		return borsh.Serialize(struct {
			Enum         uint8
			SuccessValue SuccessValue
		}{
			Enum:         2,
			SuccessValue: e.SuccessValue,
		})
	} else if e.SuccessReceiptId != [32]byte{} {
		return borsh.Serialize(struct {
			Enum             uint8
			SuccessReceiptId CryptoHash
		}{
			Enum:             3,
			SuccessReceiptId: e.SuccessReceiptId,
		})
	} else {
		return borsh.Serialize(struct {
			Enum    uint8
			Unknown string
		}{
			Enum:    0,
			Unknown: e.Unknown,
		})
	}
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

type ExecutionOutcomeWithIdView struct {
	Proofs    MerklePath           `json:"proof"`
	BlockHash CryptoHash           `json:"block_hash"`
	ReceiptId CryptoHash           `json:"id"`
	Outcome   ExecutionOutcomeView `json:"outcome"`
}

type ExecutionOutcomeView struct {
	ReceiptIds  []CryptoHash    `json:"receipt_ids"`
	GasBurnt    uint64          `json:"gas_burnt"`
	TokensBurnt BigInt          `json:"tokens_burnt"`
	ExecutorId  AccountId       `json:"executor_id"`
	Status      ExecutionStatus `json:"status"`
	Logs        []string        `json:"logs"`
}

func (e ExecutionOutcomeView) Serialize() ([]byte, error) {
	return borsh.Serialize(struct {
		ReceiptIds  []CryptoHash
		GasBurnt    uint64
		TokensBurnt BigInt
		ExecutorId  AccountId
		Status      ExecutionStatus
	}{
		ReceiptIds:  e.ReceiptIds,
		GasBurnt:    e.GasBurnt,
		TokensBurnt: e.TokensBurnt,
		ExecutorId:  e.ExecutorId,
		Status:      e.Status,
	})
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

type ActionView struct {
	Enum           borsh.Enum `borsh_enum:"true"` // treat struct as complex enum when serializing/deserializing
	CreateAccount  borsh.Enum
	DeployContract DeployContract
	FunctionCall   FunctionCallView
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
	Deposit    BigInt
}

type FunctionCallView struct {
	MethodName string
	Args       []byte
	Gas        uint64
	Deposit    BigInt
}

type Transfer struct {
	Deposit BigInt
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

type SlashedValidator struct {
	AccountId    AccountId
	IsDoubleSign bool
}
