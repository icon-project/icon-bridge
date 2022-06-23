package icon

import (
	"context"
	"fmt"
	"testing"

	"github.com/icon-project/icon-bridge/cmd/btpsimple/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestReceiver(t *testing.T) {
	srcAddress := "btp://0x5b9a77.icon/cx8015df5623344958af75ba1598b4ee14b8574bee"
	dstAddress := "btp://0x6357d2e0.hmny/0x0169AE3f21b67e798fd4AdF50d0FA9FB83d72651"
	srcEndpoint := []string{"http://localhost:9080/api/v3/default"}
	var height uint64 = 1
	var seq uint64 = 0
	opts := map[string]interface{}{
		"verifier": map[string]interface{}{
			"blockHeight":    height,
			"validatorsHash": "0xe4d1a21c32bd86bbcca8fa4dd96ea5f4fe80cc44b990ba1b1ae6b7808f8b7883",
		},
	}
	l := log.New()
	log.SetGlobalLogger(l)
	recv, err := NewReceiver(chain.BTPAddress(srcAddress), chain.BTPAddress(dstAddress), srcEndpoint, opts, l)
	if err != nil {
		panic(err)
	}

	msgCh := make(chan []*TxnLog)
	if errCh, err := recv.Subscribe(
		context.Background(), msgCh, chain.SubscribeOptions{Height: height, Seq: seq}); err != nil {
		panic(err)
	} else {
		for {
			select {
			case err := <-errCh:
				panic(err)
			case msgs := <-msgCh:
				for _, msg := range msgs {
					fmt.Println(msg.BlockHeight)
				}

			}
		}
	}
}
