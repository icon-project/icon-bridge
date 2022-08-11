package icon

import (
	"context"
	"net/http"
	"testing"

	vlcodec "github.com/icon-project/goloop/common/codec"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestReceiver(t *testing.T) {
	srcAddress := "btp://0x1.icon/cx997849d3920d338ed81800833fbb270c785e743d"
	dstAddress := "btp://0x63564c40.hmny/0xa69712a3813d0505bbD55AeD3fd8471Bc2f722DD"
	srcEndpoint := []string{"https://ctz.solidwallet.io/api/v3/icon_dex"}
	var height uint64 = 0x307f54a
	var seq uint64 = 628
	opts := map[string]interface{}{
		"verifier": map[string]interface{}{
			"blockHeight":    0x307f540,
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
				if len(msg.Receipts) > 0 && msg.Receipts[0].Height == 50853195 {
					// found event
					return
				}
			}
		}
	}
}

func TestNextValidatorHashFetch(t *testing.T) {

	var conUrl string = "https://ctz.solidwallet.io/api/v3/icon_dex" //devnet
	height := 50852431
	con := jsonrpc.NewJsonRpcClient(&http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1000}}, conUrl)
	getBlockHeaderByHeight := func(height int64, con *jsonrpc.Client) (*BlockHeader, error) {
		var header BlockHeader
		var result []byte
		_, err := con.Do("icx_getBlockHeaderByHeight", &BlockHeightParam{
			Height: NewHexInt(int64(height)),
		}, &result)
		require.NoError(t, err)

		_, err = vlcodec.RLP.UnmarshalFromBytes(result, &header)
		require.NoError(t, err)
		return &header, nil
	}

	getDatabyHash := func(req interface{}, resp interface{}, con *jsonrpc.Client) (interface{}, error) {
		_, err := con.Do("icx_getDataByHash", req, resp)
		require.NoError(t, err)
		return resp, nil
	}

	header, err := getBlockHeaderByHeight(int64(height), con)
	require.NoError(t, err)

	var validatorDataBytes []byte
	_, err = getDatabyHash(&DataHashParam{Hash: NewHexBytes(header.NextValidatorsHash)}, &validatorDataBytes, con)
	require.NoError(t, err)

	var validators [][]byte
	_, err = vlcodec.BC.UnmarshalFromBytes(validatorDataBytes, &validators)
	require.NoError(t, err)

	if common.HexBytes(header.NextValidatorsHash).String() != "0xa6760c547c3f76b7071658ef383d69ec01e11ea71d695600788695b50659e409" {
		err := errors.New("Invalid Validator Hash")
		require.NoError(t, err)
	}
}
