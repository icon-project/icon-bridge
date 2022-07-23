package icon

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestMonitorBlockMissingNotification(t *testing.T) {
	urls := []string{
		"https://ctz.solidwallet.io/api/v3/icon_dex",
		"http://138.197.69.76:9000/api/v3/icon_dex",
	}
	l := log.New()
	ctx := context.Background()

	height, seq := 0x306d1ac, 0

	dstAddr := "btp://0x63564c40.hmny/0xa69712a3813d0505bbD55AeD3fd8471Bc2f722DD"
	blockReq := &BlockRequest{
		EventFilters: []*EventFilter{{
			Addr:      Address("cx997849d3920d338ed81800833fbb270c785e743d"),
			Signature: EventSignature,
			Indexed:   []*string{&dstAddr},
		}},
		Height: NewHexInt(int64(height)),
	}

	for i, url := range urls {
		go func(i int, url string) {
			l := l.WithFields(log.Fields{"i": i, "url": url})

			cl := NewClient(url, l)

			h, s := height, seq
			err := cl.MonitorBlock(ctx, blockReq,
				func(conn *websocket.Conn, v *BlockNotification) error {
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
			if err != nil {
				panic(err)
			}

		}(i, url)
	}
	time.Sleep(time.Hour)
}
