package icon

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	vlcodec "github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/icon-bridge/bmr/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/bmr/common"
	"github.com/icon-project/icon-bridge/bmr/common/jsonrpc"
	"github.com/icon-project/icon-bridge/bmr/common/log"
)

func TestReceiver(t *testing.T) {
	srcAddress := "btp://0x1.icon/cx997849d3920d338ed81800833fbb270c785e743d"
	dstAddress := "btp://0x63564c40.hmny/0xa69712a3813d0505bbD55AeD3fd8471Bc2f722DD"
	srcEndpoint := []string{"https://ctz.solidwallet.io/api/v3/icon_dex"}
	// var height uint64 = 0x307f245 // seq 0x0a
	// var height uint64 = 0x307f24f
	var height uint64 = 0x308024f
	// var seq uint64 = 611
	var seq uint64 = 628
	opts := map[string]interface{}{
		"verifier": map[string]interface{}{
			"blockHeight":    0x307f24f,
			"validatorsHash": "0xa6760c547c3f76b7071658ef383d69ec01e11ea71d695600788695b50659e409",
		},
	}
	l := log.New()
	log.SetGlobalLogger(l)
	// log.AddForwarder(&log.ForwarderConfig{Vendor: log.HookVendorSlack, Address: "https://hooks.slack.com/services/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", Level: "info"})
	recv, err := NewReceiver(chain.BTPAddress(srcAddress), chain.BTPAddress(dstAddress), srcEndpoint, opts, l)
	if err != nil {
		panic(err)
	}
	msgCh := make(chan *chain.Message)
	if errCh, err := recv.Subscribe(
		context.Background(), msgCh, chain.SubscribeOptions{Height: height, Seq: seq}); err != nil {
		panic(err)
	} else {
		for {
			select {
			case err := <-errCh:
				panic(err)
			case msg := <-msgCh:
				if len(msg.Receipts) > 0 {
					l.Info(msg.From, len(msg.Receipts), msg.Receipts[0].Height)
				}
			}
		}
	}
}

func TestNextValidatorHashFetch(t *testing.T) {

	//var conUrl string = "https://ctz.solidwallet.io/api/v3/icon_dex" //devnet
	var conUrl string = "http://127.0.0.1:9080/api/v3/default" // mainnet
	height := 9999

	con := jsonrpc.NewJsonRpcClient(&http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1000}}, conUrl)

	getBlockHeaderByHeight := func(height int64, con *jsonrpc.Client) (*BlockHeader, error) {
		var header BlockHeader
		var result []byte
		_, err := con.Do("icx_getBlockHeaderByHeight", &BlockHeightParam{
			Height: NewHexInt(int64(height)),
		}, &result)
		if err != nil {
			return nil, err
		}
		_, err = vlcodec.RLP.UnmarshalFromBytes(result, &header)
		if err != nil {
			return nil, err
		}
		return &header, nil
	}

	getDatabyHash := func(req interface{}, resp interface{}, con *jsonrpc.Client) (interface{}, error) {
		_, err := con.Do("icx_getDataByHash", req, resp)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	header, err := getBlockHeaderByHeight(int64(height), con)
	if err != nil {
		log.Fatal(err)
	}

	var validatorDataBytes []byte
	_, err = getDatabyHash(&DataHashParam{Hash: NewHexBytes(header.NextValidatorsHash)}, &validatorDataBytes, con)
	if err != nil {
		log.Fatal(err)
	}
	var validators [][]byte
	_, err = vlcodec.BC.UnmarshalFromBytes(validatorDataBytes, &validators)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(common.HexBytes(header.NextValidatorsHash), NewHexInt(int64(height)))
}
