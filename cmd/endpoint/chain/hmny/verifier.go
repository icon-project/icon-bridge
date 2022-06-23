package hmny

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"

	libbls "github.com/harmony-one/bls/ffi/go/bls"
	"github.com/harmony-one/harmony/crypto/bls"
	"github.com/harmony-one/harmony/numeric"
)

var (
	bigOne  = big.NewInt(1)
	bigZero = big.NewInt(0)
)

type VerifierOptions struct {
	BlockHeight     uint64        `json:"blockHeight"`
	CommitBitmap    hexutil.Bytes `json:"commitBitmap"`
	CommitSignature hexutil.Bytes `json:"commitSignature"`
}

type Verifier interface {
	Epoch() uint64
	Verify(h *Header, bitmap, signature []byte) (ok bool, err error)
	Update(h *Header) (err error)
}

type dumbVerifier struct{}

func (*dumbVerifier) Epoch() uint64 {
	return 0
}

func (*dumbVerifier) Update(
	h *Header) (err error) {
	return nil
}

func (*dumbVerifier) Verify(h *Header,
	bitmap, signature []byte) (ok bool, err error) {
	return true, nil
}

type verifier struct {
	epoch uint64
	mu    sync.RWMutex
	smsk  map[uint32]*bls.Mask
}

func (vr *verifier) Epoch() uint64 { return vr.epoch }

func (vr *verifier) Verify(h *Header, bitmap, signature []byte) (bool, error) {
	vr.mu.RLock()
	msk, ok := vr.smsk[h.ShardID]
	vr.mu.RUnlock()
	if !ok {
		return false, fmt.Errorf("invalid shard id: %d", h.ShardID)
	}
	mask := *msk
	mask.Clear()
	if err := mask.SetMask(bitmap); err != nil {
		return false, err
	}
	asig := &libbls.Sign{}
	if err := asig.Deserialize(signature); err != nil {
		return false, err
	}
	return asig.VerifyHash(mask.AggregatePublic, vr.payload(h)), nil
}

func (vr *verifier) Update(h *Header) (err error) {
	vr.mu.Lock()
	defer vr.mu.Unlock()

	if len(h.ShardState) == 0 {
		return nil
	}

	var epoch uint64

	spks := make(map[uint32][]bls.SerializedPublicKey)

	ss := ShardState{}
	if err = rlp.DecodeBytes(h.ShardState, &ss); err == nil {
		for _, cmt := range ss.Shards {
			pkws := make([]bls.SerializedPublicKey, 0, len(cmt.Slots))
			for _, slt := range cmt.Slots {
				pkws = append(pkws, slt.BLSPublicKey)
			}
			spks[cmt.ShardID] = pkws
		}
		epoch = ss.Epoch.Uint64()
	} else {
		lss := LegacyShardState{}
		if err = rlp.DecodeBytes(h.ShardState, &lss); err != nil {
			return err
		}
		for _, cmt := range lss {
			pkws := make([]bls.SerializedPublicKey, 0, len(cmt.Slots))
			for _, slt := range cmt.Slots {
				pkws = append(pkws, slt.BLSPublicKey)
			}
			spks[cmt.ShardID] = pkws
		}
		epoch = h.Epoch.Uint64() + 1
	}

	smsk := make(map[uint32]*bls.Mask)

	for sid, pks := range spks {
		pubs := make([]bls.PublicKeyWrapper, len(pks))
		for i, pk := range pks {
			pubs[i].Bytes = pk
			pubs[i].Object, err = bls.BytesToBLSPublicKey(pubs[i].Bytes[:])
			if err != nil {
				return err
			}
		}
		mask, err := bls.NewMask(pubs, nil)
		if err != nil {
			return err
		}
		smsk[sid] = mask
	}

	vr.epoch, vr.smsk = epoch, smsk
	return nil
}

func (vl *verifier) payload(h *Header) []byte {
	hash := h.Hash().Bytes()
	payload := make([]byte, 8+len(hash)+8)
	binary.LittleEndian.PutUint64(payload, h.Number.Uint64())
	copy(payload[8:], hash)
	binary.LittleEndian.PutUint64(payload[8+len(hash):], h.ViewID.Uint64()) // after staking epoch
	return payload
}

type LegacyShardState []struct {
	ShardID uint32
	Slots   []struct {
		EcdsaAddress common.Address
		BLSPublicKey bls.SerializedPublicKey
	}
}

type ShardState struct {
	Epoch  *big.Int
	Shards []struct {
		ShardID uint32
		Slots   []struct {
			EcdsaAddress   common.Address
			BLSPublicKey   bls.SerializedPublicKey
			EffectiveStake *numeric.Dec `rlp:"nil"`
		}
	}
}
