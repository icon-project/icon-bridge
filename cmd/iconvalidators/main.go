package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/icon-project/goloop/client"
	"github.com/icon-project/goloop/common"
	"github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/goloop/common/crypto"
	"github.com/icon-project/goloop/server/jsonrpc"
	v3 "github.com/icon-project/goloop/server/v3"
)

var (
	v2c = codec.BC
	cl  *client.ClientV3
)

func init() {
	cl = client.NewClientV3(os.Getenv("URI"))
}

type v2Header struct {
	Version                int
	Height                 int64
	Timestamp              int64
	Proposer               []byte
	PrevID                 []byte
	VotesHash              []byte
	NextValidatorsHash     []byte
	PatchTransactionsHash  []byte
	NormalTransactionsHash []byte
	LogsBloom              []byte
	Result                 []byte
}

func (h *v2Header) Hash() []byte {
	return crypto.SHA3Sum256(v2c.MustMarshalToBytes(h))
}

type validatorList []common.HexBytes

func main() {

	hexInt := os.Getenv("HEIGHT")
	v, err := getBlockHeaderByHeightInHexInt(hexInt)
	if err != nil {
		panic("get block header: HEIGHT=" + hexInt + "; " + err.Error())
	}

	rlpv, err := cl.GetDataByHash(
		&v3.DataHashParam{Hash: jsonrpc.HexBytes(
			common.HexHash(v.NextValidatorsHash).String())})
	if err != nil {
		panic(err)
	}
	var vl validatorList
	_, err = v2c.UnmarshalFromBytes(rlpv, &vl)
	if err != nil {
		panic(err)
	}
	var validators []common.HexBytes
	for _, v := range vl {
		validators = append(validators, v)
	}
	json.NewEncoder(os.Stdout).Encode(map[string]interface{}{
		"hash":       common.HexBytes(v.NextValidatorsHash),
		"validators": validators,
	})
}

func getBlockHeaderByHeightInHexInt(h string) (*v2Header, error) {
	rlpn, err := cl.GetBlockHeaderByHeight(&v3.BlockHeightParam{
		Height: jsonrpc.HexInt(h)})
	if err != nil {
		return nil, err
	}
	var v v2Header
	_, err = codec.BC.UnmarshalFromBytes(rlpn, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil

}

func getBlockHeaderByHeight(n int64) (*v2Header, error) {
	return getBlockHeaderByHeightInHexInt(hexInt(n))
}

func hexInt(v int64) string {
	return fmt.Sprintf("0x%x", v)
}
