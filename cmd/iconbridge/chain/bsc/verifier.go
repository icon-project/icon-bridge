package bsc

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"sync"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/icon-bridge/common"
	"github.com/pkg/errors"
)

const (
	extraVanity          = 32          // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal            = 65          // Fixed number of extra-data suffix bytes reserved for signer seal
	defaultEpochLength   = uint64(200) // Default number of blocks of checkpoint to update validatorSet from contract
	validatorBytesLength = ethCommon.AddressLength

	ParliaGasLimitBoundDivisor uint64 = 256                // The bound divisor of the gas limit, used in update calculations.
	MinGasLimit                uint64 = 5000               // Minimum the gas limit may ever be.
	MaxGasLimit                uint64 = 0x7fffffffffffffff // Maximum the gas limit (2^63-1).
)

var (
	big1      = big.NewInt(1)
	uncleHash = types.CalcUncleHash(nil)
)

var (
	// errUnknownBlock is returned when the list of validators is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	// errMissingVanity is returned if a block's extra-data section is shorter than
	// 32 bytes, which is required to store the signer vanity.
	errMissingVanity = errors.New("extra-data 32 byte vanity prefix missing")

	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")

	errMissingValidators = errors.New("epoch block does not have validators")

	// errExtraValidators is returned if non-sprint-end block contain validator data in
	// their extra-data fields.
	errExtraValidators = errors.New("non-sprint-end block contains extra validator list")

	// errInvalidSpanValidators is returned if a block contains an
	// invalid list of validators (i.e. non divisible by 20 bytes).
	errInvalidSpanValidators = errors.New("invalid validator list on sprint end block")

	// errInvalidMixDigest is returned if a block's mix digest is non-zero.
	errInvalidMixDigest = errors.New("non-zero mix digest")

	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")

	// errInvalidDifficulty is returned if the difficulty of a block is missing.
	errInvalidDifficulty = errors.New("invalid difficulty")

	// errUnauthorizedValidator is returned if a header is signed by a non-authorized entity.
	errUnauthorizedValidator = errors.New("unauthorized validator")

	// errCoinBaseMisMatch is returned if a header's coinbase do not match with signature
	errCoinBaseMisMatch = errors.New("coinbase do not match with signature")
)

type VerifierOptions struct {
	BlockHeight   uint64          `json:"blockHeight"`
	BlockHash     common.HexBytes `json:"parentHash"`
	ValidatorData common.HexBytes `json:"validatorData"`
}

// next points to height whose parentHash is expected
// parentHash of height h is got from next-1's hash
type Verifier struct {
	chainID    *big.Int
	mu         sync.RWMutex
	next       *big.Int
	parentHash ethCommon.Hash
	validators map[ethCommon.Address]bool
}

func (vr *Verifier) Next() *big.Int {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return (&big.Int{}).Set(vr.next)
}

func (vr *Verifier) ChainID() *big.Int {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return (&big.Int{}).Set(vr.chainID)
}

func (vr *Verifier) ParentHash() ethCommon.Hash {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	return ethCommon.BytesToHash(vr.parentHash.Bytes())
}

func (vr *Verifier) IsValidator(addr ethCommon.Address) bool {
	vr.mu.RLock()
	defer vr.mu.RUnlock()
	_, exists := vr.validators[addr]
	return exists
}

// prove that header is linked to verified nextHeader
// only then can header be used for receiver.Callback or vr.Update()
func (vr *Verifier) Verify(header *types.Header, nextHeader *types.Header, receipts ethTypes.Receipts) error {

	if nextHeader.Number.Cmp((&big.Int{}).Add(header.Number, big1)) != 0 {
		return fmt.Errorf("Different height between successive header: Prev %v New %v", header.Number, nextHeader.Number)
	}
	if !bytes.Equal(header.Hash().Bytes(), nextHeader.ParentHash.Bytes()) {
		return fmt.Errorf("Different hash between successive header: (%v): Prev %v New %v", header.Number, header.Hash(), nextHeader.ParentHash)
	}
	if vr.Next().Cmp(header.Number) != 0 {
		return fmt.Errorf("Unexpected height: Got %v Expected %v", header.Number, vr.Next())
	}
	if !bytes.Equal(header.ParentHash.Bytes(), vr.ParentHash().Bytes()) {
		return fmt.Errorf("Unexpected Hash(%v): Got %v Expected %v", header.Number, header.ParentHash, vr.ParentHash())
	}

	if err := vr.verifyHeader(nextHeader); err != nil {
		return errors.Wrapf(err, "verifyHeader %v", err)
	}
	if err := vr.verifyCascadingFields(nextHeader, header); err != nil {
		return errors.Wrapf(err, "verifyCascadingFields %v", err)
	}
	if err := vr.verifySeal(nextHeader, vr.ChainID()); err != nil {
		return errors.Wrapf(err, "verifySeal %v", err)
	}
	if len(receipts) > 0 {
		if err := vr.validateState(nextHeader, receipts); err != nil {
			return errors.Wrapf(err, "validateState %v", err)
		}
	}
	return nil
}

func (vr *Verifier) Update(header *types.Header) (err error) {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	vr.parentHash = header.Hash()
	vr.next.Add(header.Number, big1)
	if header.Number.Uint64()%defaultEpochLength != 0 {
		return nil
	}
	// update validators if epoch block
	vr.validators, err = getValidatorMapFromHex(header.Extra)
	return
}

func getValidatorMapFromHex(headerExtra common.HexBytes) (map[ethCommon.Address]bool, error) {
	if len(headerExtra) < extraVanity+extraSeal {
		return nil, errMissingSignature
	}
	addrs := headerExtra[extraVanity : len(headerExtra)-extraSeal]
	numAddrs := len(addrs) / validatorBytesLength
	newVals := make(map[ethCommon.Address]bool, numAddrs)
	for i := 0; i < numAddrs; i++ {
		newVals[ethCommon.BytesToAddress(addrs[i*validatorBytesLength:(i+1)*validatorBytesLength])] = true
	}
	return newVals, nil
}

func (vr *Verifier) verifyHeader(header *types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	number := header.Number.Uint64()

	// Don't waste time checking blocks from the future
	if header.Time > uint64(time.Now().Unix()) {
		return consensus.ErrFutureBlock
	}
	// Check that the extra-data contains the vanity, validators and signature.
	if len(header.Extra) < extraVanity {
		return errMissingVanity
	}
	if len(header.Extra) < extraVanity+extraSeal {
		return errMissingSignature
	}

	// check extra data
	isEpoch := number%defaultEpochLength == 0

	// Ensure that the extra-data contains a signer list on checkpoint, but none otherwise
	signersBytes := len(header.Extra) - extraVanity - extraSeal
	if !isEpoch && signersBytes != 0 {
		return errExtraValidators
	}

	if isEpoch && signersBytes == 0 {
		return errMissingValidators
	}

	if isEpoch && signersBytes%validatorBytesLength != 0 {
		return errInvalidSpanValidators
	}

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != (ethCommon.Hash{}) {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in PoA
	if header.UncleHash != uncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if number > 0 && header.Difficulty == nil {
		return errInvalidDifficulty
	}
	return nil
}

func (vr *Verifier) verifyCascadingFields(header *types.Header, parent *types.Header) error {
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	// Verify that the gas limit is <= 2^63-1
	capacity := MaxGasLimit
	if header.GasLimit > capacity {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, capacity)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}

	// Verify that the gas limit remains within allowed bounds
	diff := int64(parent.GasLimit) - int64(header.GasLimit)
	if diff < 0 {
		diff *= -1
	}
	limit := parent.GasLimit / ParliaGasLimitBoundDivisor

	if uint64(diff) >= limit || header.GasLimit < MinGasLimit {
		return fmt.Errorf("invalid gas limit: have %d, want %d += %d", header.GasLimit, parent.GasLimit, limit)
	}
	return nil
}

func (vr *Verifier) verifySeal(header *types.Header, chainID *big.Int) error {
	// Resolve the authorization key and check against validators
	signer, err := ecrecover(header, chainID)
	if err != nil {
		return err
	}
	if signer != header.Coinbase {
		return errCoinBaseMisMatch
	}

	if ok := vr.IsValidator(signer); !ok {
		return errUnauthorizedValidator
	}
	// TODO: check if signer is a recent Validator; avoid recent validators for spam protection
	return nil
}

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(header *types.Header, chainId *big.Int) (ethCommon.Address, error) {
	if len(header.Extra) < extraSeal {
		return ethCommon.Address{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(SealHash(header, chainId).Bytes(), signature)
	if err != nil {
		return ethCommon.Address{}, err
	}
	var signer ethCommon.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	return signer, nil
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header *types.Header, chainId *big.Int) (hash ethCommon.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header, chainId)
	hasher.Sum(hash[:0])
	return hash
}

func encodeSigHeader(w io.Writer, header *types.Header, chainId *big.Int) {
	err := rlp.Encode(w, []interface{}{
		chainId,
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-65], // this will panic if extra is too short, should check before calling encodeSigHeader
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}

func (vr *Verifier) validateState(header *types.Header, receipts types.Receipts) error {
	rbloom := types.CreateBloom(receipts)
	if rbloom != header.Bloom {
		return fmt.Errorf("invalid bloom (remote: %x  local: %x)", header.Bloom, rbloom)
	}
	receiptSha := types.DeriveSha(receipts, trie.NewStackTrie(nil))
	if receiptSha != header.ReceiptHash {
		return fmt.Errorf("invalid receipt root hash (remote: %x local: %x)", header.ReceiptHash, receiptSha)
	}
	return nil
}
