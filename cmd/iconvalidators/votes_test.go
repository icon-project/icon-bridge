package main

import (
	"fmt"
	"testing"

	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/goloop/common/crypto"
)

type commitVoteItem struct {
	Timestamp int64
	Signature common.Signature
}

type commitVoteList struct {
	Round          int32
	BlockPartSetID *PartSetID
	Items          []commitVoteItem
}

type PartSetID struct {
	Count uint16
	Hash  []byte
}

type VoteType byte

type _HR struct {
	Height int64
	Round  int32
}

type voteBase struct {
	_HR
	Type           VoteType
	BlockID        []byte
	BlockPartSetID *PartSetID
}

type vote struct {
	voteBase
	Timestamp int64
}

func TestVotes(t *testing.T) {
	vl := new(commitVoteList)
	msg := "f87300e201a003d1b1c67c5d6a806904f372e7d75037dfd69402bc367ac139fc2718f59ea0e6f84df84b8705da40a369f87eb84125dc8de24456decf5067956135a89c4067934a9f4e154465a6122d1fe434ff4a24e8d2b1e9a7c418a5b382f1aebf4a1ccd9aad9f7d81ea23070cc7e99ec1089601"
	codec.BC.MustUnmarshalFromBytes(ecommon.Hex2Bytes(msg), vl)

	v := &vote{
		voteBase: voteBase{
			_HR: _HR{
				Height: 0x0d13,
				Round:  vl.Round,
			},
			Type:           1,
			BlockID:        ecommon.Hex2Bytes("214f49db3a82b5a0ad6f062d09daed8a230bedac6e0aa06b57c624be3f76f94a"),
			BlockPartSetID: vl.BlockPartSetID,
		},
		Timestamp: vl.Items[0].Timestamp,
	}

	vby := codec.BC.MustMarshalToBytes(v)

	vh := crypto.SHA3Sum256(vby)
	fmt.Println(ecommon.Bytes2Hex(vh))
	fmt.Println(ecommon.Bytes2Hex(vby))

}
