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
	"github.com/icon-project/btp/common/errors"
)

type VerifierOptions struct {
	BlockHeight     int64         `json:"blockHeight"`
	CommitBitmap    hexutil.Bytes `json:"commitBitmap"`
	CommitSignature hexutil.Bytes `json:"commitSignature"`
}

type Verifier interface {
	Update(h *Header) (err error)
	Verify(h *Header, bitmap, signature []byte) (ok bool, err error)
	CatchUp(cl *Client, height int64) (err error)
}

type dumbVerifier struct{}

func (vl *dumbVerifier) Update(
	h *Header) (err error) {
	return nil
}
func (vl *dumbVerifier) Verify(h *Header,
	bitmap, signature []byte) (ok bool, err error) {
	return true, nil
}
func (vr *dumbVerifier) CatchUp(
	cl *Client, height int64) (err error) {
	return nil
}

type verifier struct {
	epoch int64
	mu    sync.RWMutex
	smsk  map[uint32]*bls.Mask
}

func (vr *verifier) Update(h *Header) (err error) {
	vr.mu.Lock()
	defer vr.mu.Unlock()

	if len(h.ShardState) == 0 {
		return nil
	}

	var epoch int64

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
		epoch = ss.Epoch.Int64()
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
		epoch = h.Epoch.Int64() + 1
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

func (vr *verifier) CatchUp(cl *Client, height int64) (err error) {
	h, err := cl.GetHmyHeaderByHeight(big.NewInt(height), 0)
	if err != nil {
		return err
	}
	for epoch := vr.epoch + 1; epoch < h.Epoch.Int64(); epoch++ {
		elb, err := cl.GetEpochLastBlock(big.NewInt(epoch))
		if err != nil {
			return errors.Wrapf(err, "cl.GetEpochLastBlock: %v", err)
		}
		elh, err := cl.GetHmyHeaderByHeight(elb, 0)
		if err != nil {
			return errors.Wrapf(err, "cl.GetHmyHeaderByHeight: %v", err)
		}
		if err = vr.Update(elh); err != nil {
			return errors.Wrapf(err, "vr.Update: %v", err)
		}
		cl.log.Debugf("caught up to epoch: %d, h=%d", epoch, elb)
	}
	return nil
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

func NewVerifier(cl *Client, opts *VerifierOptions) (Verifier, error) {
	if opts == nil {
		return &dumbVerifier{}, nil
	}
	h, err := cl.GetHmyHeaderByHeight(big.NewInt(opts.BlockHeight), 0.75)
	if err != nil {
		return nil, errors.Wrapf(err, "cl.GetHeaderByHeight(%d): %v", opts.BlockHeight, err)
	}
	ssh := h // shard state header
	if ssh.Epoch.Cmp(bigZero) <= 0 {
		if ssh.Number.Cmp(bigZero) > 0 {
			ssh, err = cl.GetHmyHeaderByHeight(bigZero, 0.75)
			if err != nil {
				return nil, errors.Wrapf(err, "cl.GetHeaderByHeight(%d): %v", 0, err)
			}
		}
	} else {
		epoch := new(big.Int).Sub(ssh.Epoch, bigOne)
		elb, err := cl.GetEpochLastBlock(epoch)
		if err != nil {
			return nil, errors.Wrapf(err, "cl.GetEpochLastBlock(%d): %v", epoch, err)
		}
		ssh, err = cl.GetHmyHeaderByHeight(elb, 0.75)
		if err != nil {
			return nil, errors.Wrapf(err, "cl.GetHeaderByHeight(%d): %v", elb, err)
		}
	}
	vr := verifier{}
	if err = vr.Update(ssh); err != nil {
		return nil, errors.Wrapf(err, "verifier.Update: %v", err)
	}
	ok, err := vr.Verify(h, opts.CommitBitmap, opts.CommitSignature)
	if !ok || err != nil {
		return nil, errors.Wrapf(err, "invalid signature: %v", err)
	}
	return &vr, nil
}
