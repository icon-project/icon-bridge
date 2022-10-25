package substrate_eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"sync"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/common/eth/types"
	ethSubVerifier "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/common/substrate-eth"
	"github.com/icon-project/icon-bridge/common"
	"github.com/pkg/errors"
)

const (
	extraVanity          = 32          // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal            = 65          // Fixed number of extra-data suffix bytes reserved for signer seal
	defaultEpochLength   = uint64(200) // Default number of blocks of checkpoint to update validatorSet from contract
	validatorBytesLength = ethCommon.AddressLength
)

var (
	big1      = big.NewInt(1)
)

var (
	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")

	// errUnauthorizedValidator is returned if a header is signed by a non-authorized entity.
	errUnauthorizedValidator = errors.New("unauthorized validator")
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

type IVerifier interface {
	Next() *big.Int
	Verify(header *ethTypes.Header, nextHeader *ethTypes.Header, receipts types.Receipts) error
	Update(header *ethTypes.Header) (err error)
	ParentHash() ethCommon.Hash
	IsValidator(addr ethCommon.Address) bool
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

func (vr *Verifier) Verify(header *ethTypes.Header, nextHeader *ethTypes.Header, receipts types.Receipts) error {
	if nextHeader.Number.Cmp((&big.Int{}).Add(header.Number, big1)) != 0 {
		return fmt.Errorf("Different height between successive header: Prev %v New %v", header.Number, nextHeader.Number)
	}
	if header.Hash != nextHeader.ParentHash {
		return fmt.Errorf("Different hash between successive header: (%v): Prev %v New %v", header.Number, header.Hash, nextHeader.ParentHash)
	}
	if vr.Next().Cmp(header.Number) != 0 {
		return fmt.Errorf("Unexpected height: Got %v Expected %v", header.Number, vr.Next())
	}
	if header.ParentHash != vr.ParentHash() {
		return fmt.Errorf("Unexpected Hash(%v): Got %v Expected %v", header.Number, header.ParentHash, vr.ParentHash())
	}

	if err := ethSubVerifier.VerifyHeader(nextHeader); err != nil {
		return errors.Wrapf(err, "verifyHeader %v", err)
	}

	if err := ethSubVerifier.VerifyCascadingFields(nextHeader, header); err != nil {
		return errors.Wrapf(err, "verifyCascadingFields %v", err)
	}
	signer, err := ethSubVerifier.VerifySeal(nextHeader, vr.ChainID())

	if err != nil {
		return errors.Wrapf(err, "verifySeal %v", err)
	}

	if ok := vr.IsValidator(signer); !ok {
		return errors.Wrapf(errUnauthorizedValidator, "Signer %v", signer)
	}
	// TODO: check if signer is a recent Validator; avoid recent validators for spam protection

	if len(receipts) > 0 {
		if err := ethSubVerifier.ValidateState(nextHeader, receipts); err != nil {
			return errors.Wrapf(err, "validateState %v", err)
		}
	}
	return nil

}

func (vr *Verifier) Update(header *ethTypes.Header) (err error) {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	vr.parentHash = header.Hash
	vr.next.Add(header.Number, big1)
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