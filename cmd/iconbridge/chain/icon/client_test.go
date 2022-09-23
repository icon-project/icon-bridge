package icon

import (
	"context"
	"fmt"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon/types"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/stretchr/testify/require"
)

func NewTestClient() *Client {
	uri := "https://ctz.solidwallet.io/api/v3"
	l := log.New().WithFields(log.Fields{"uri": uri})
	return NewClient(uri, l)
}

func TestContextCancel(t *testing.T) {
	urls := []string{
		"https://ctz.solidwallet.io/api/v3/icon_dex",
		"http://138.197.69.76:9000/api/v3/icon_dex",
	}
	l := log.New()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	height, seq := 0x306d1ac, 0

	dstAddr := "btp://0x63564c40.hmny/0xa69712a3813d0505bbD55AeD3fd8471Bc2f722DD"
	blockReq := &types.BlockRequest{
		EventFilters: []*types.EventFilter{{
			Addr:      types.Address("cx997849d3920d338ed81800833fbb270c785e743d"),
			Signature: EventSignature,
			Indexed:   []*string{&dstAddr},
		}},
		Height: types.NewHexInt(int64(height)),
	}

	for i, url := range urls {
		go func(i int, url string) {
			l := l.WithFields(log.Fields{"i": i, "url": url})

			cl := NewClient(url, l)

			h, s := height, seq
			err := cl.MonitorBlock(ctx, blockReq,
				func(conn *websocket.Conn, v *types.BlockNotification) error {
					_h, _ := v.Height.Int()
					if _h != h {
						err := fmt.Errorf("invalid block height: %d, expected: %d", _h, h+1)
						l.Info(err)
						return err
					}
					h++
					s++
					return nil
				},
				func(conn *websocket.Conn) {
					l.WithFields(log.Fields{"local": conn.LocalAddr().String()}).Debug("connected")
				},
				func(conn *websocket.Conn, err error) {
					l.WithFields(log.Fields{"error": err, "local": conn.LocalAddr().String()}).Warn("disconnected")
					_ = conn.Close()
				})
			if err.Error() == "context deadline exceeded" {
				return
			}
			require.NoError(t, err)

		}(i, url)
	}
	time.Sleep(time.Second * 15)
}

// func TestGetHeaderByHeight(t *testing.T) {
// 	Client := NewTestClient()

// 	var height int64 = 50000000

// 	enc := json.NewEncoder(os.Stdout)
// 	enc.SetIndent("", "  ")

// 	h, err := Client.GetBlockHeaderByHeight(height)
// 	require.NoError(t, err)
// 	enc.Encode(h)

// 	vals, err := Client.GetValidatorsByHash(h.NextValidatorsHash)
// 	require.NoError(t, err)

// 	fmt.Println(common.HexBytes(h.NextValidatorsHash))
// 	enc.Encode(vals)

// 	h, err = Client.GetBlockHeaderByHeight(height + 1)
// 	require.NoError(t, err)
// 	enc.Encode(h)

// 	vl, err := Client.GetCommitVoteListByHeight(h.Height)
// 	require.NoError(t, err)
// 	enc.Encode(vl)

// 	p := &BlockHeightParam{Height: NewHexInt(h.Height)}
// 	votes, err := Client.GetVotesByHeight(p)
// 	require.NoError(t, err)
// 	fmt.Println(common.HexBytes(votes))
// }
