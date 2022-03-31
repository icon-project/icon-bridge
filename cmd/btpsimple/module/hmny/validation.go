package hmny

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"

	libbls "github.com/harmony-one/bls/ffi/go/bls"
	"github.com/harmony-one/harmony/crypto/bls"
	"github.com/harmony-one/harmony/numeric"
	"github.com/icon-project/btp/common/errors"
)

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

type Validator struct {
	mu   sync.RWMutex
	smsk map[uint32]*bls.Mask
}

func (vl *Validator) update(h *Header) (err error) {
	vl.mu.Lock()
	defer vl.mu.Unlock()

	if len(h.ShardState) == 0 {
		return nil
	}

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
	}

	vl.smsk = make(map[uint32]*bls.Mask)

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
		vl.smsk[sid] = mask
	}

	return nil
}

func (vl *Validator) verify(h *Header, sig, bitmap []byte) (bool, error) {
	vl.mu.RLock()
	msk, ok := vl.smsk[h.ShardID]
	vl.mu.RUnlock()
	if !ok {
		return false, fmt.Errorf("invalid shard id: %d", h.ShardID)
	}
	mask := *msk
	mask.Clear()
	if err := mask.SetMask(bitmap); err != nil {
		return false, err
	}
	asig := &libbls.Sign{}
	if err := asig.Deserialize(sig); err != nil {
		return false, err
	}
	return asig.VerifyHash(mask.AggregatePublic, vl.payload(h)), nil
}

func (vl *Validator) payload(h *Header) []byte {
	hash := h.Hash().Bytes()
	payload := make([]byte, 8+len(hash)+8)
	binary.LittleEndian.PutUint64(payload, h.Number.Uint64())
	copy(payload[8:], hash)
	binary.LittleEndian.PutUint64(payload[8+len(hash):], h.ViewID.Uint64()) // after staking epoch
	return payload
}

func newValidator(cl *Client, height *big.Int) (*Validator, error) {
	h, err := cl.GetHmyHeaderByHeight(height, 0.75)
	if err != nil {
		return nil, errors.Wrapf(err, "GetHeaderByHeight(%d): %v", height, err)
	}
	if h.Epoch.Cmp(bigZero) <= 0 {
		if h.Number.Cmp(bigZero) > 0 {
			h, err = cl.GetHmyHeaderByHeight(bigZero, 0.75)
			if err != nil {
				return nil, errors.Wrapf(err, "GetHeaderByHeight(%d): %v", 0, err)
			}
		}
	} else {
		epoch := new(big.Int).Sub(h.Epoch, bigOne)
		elb, err := cl.GetEpochLastBlock(epoch)
		if err != nil {
			return nil, errors.Wrapf(err, "GetEpochLastBlock(%d): %v", epoch, err)
		}
		h, err = cl.GetHmyHeaderByHeight(elb, 0.75)
		if err != nil {
			return nil, errors.Wrapf(err, "GetHeaderByHeight(%d): %v", elb, err)
		}
	}
	vl := Validator{}
	return &vl, vl.update(h)
}
