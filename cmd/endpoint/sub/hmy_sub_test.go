package sub

import (
	"context"
	"fmt"
	"testing"

	"github.com/icon-project/icon-bridge/cmd/btpsimple/chain"
	"github.com/icon-project/icon-bridge/cmd/btpsimple/chain/hmny"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestHmnyReceiver(t *testing.T) {
	const (
		src = "btp://0x6357d2e0.hmny/0x0169AE3f21b67e798fd4AdF50d0FA9FB83d72651"
		dst = "btp://0x5b9a77.icon/cx8015df5623344958af75ba1598b4ee14b8574bee"
		url = "http://localhost:9500"
	)

	var opts map[string]interface{} = map[string]interface{}{
		// "options": map[string]interface{}{
		"syncConcurrency": 100,
		"verifier": map[string]interface{}{
			"blockHeight":     1171, //1171,
			"commitBitmap":    "0xff",
			"commitSignature": "0xf102d2c65b9012f9e02f2066b381bcd91dfcaa0b41ce76671838dce44d5aff4837768982fd7d5e0c544318aa560af502a28de49329c992c59e69277828362b7dabecae1b462fb2ed43d2c4d010e31f3ae099689ea8224ab6e65051cf5e6db992",
		},
		// },
	}

	l := log.New()
	log.SetGlobalLogger(l)

	rx, err := hmny.NewReceiver(chain.BTPAddress(src), chain.BTPAddress(dst), []string{url}, opts, l)
	if err != nil {
		log.Fatal((err))
	}

	srcMsgCh := make(chan *chain.Message)
	srcErrCh, err := rx.Subscribe(context.TODO(),
		srcMsgCh,
		chain.SubscribeOptions{
			Seq:    1,
			Height: 1172,
		})
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-srcErrCh:
			log.Fatal(err)

		case msg := <-srcMsgCh:
			fmt.Println("msg", msg.From, "  ", len(msg.Receipts))
			for _, r := range msg.Receipts {
				fmt.Println(r.Height, r.Index)
				for _, er := range r.Events {
					fmt.Println(er.Next, er.Sequence, er.Message)
				}
			}
		}
	}
}
