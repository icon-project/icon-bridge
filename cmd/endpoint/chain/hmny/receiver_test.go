package hmny

import (
	"context"
	"testing"

	"github.com/icon-project/icon-bridge/cmd/btpsimple/chain"
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
			"blockHeight":     29075, //1171,
			"commitBitmap":    "0xff",
			"commitSignature": "0x44e933d799e2e44f5f7d84fc2f8400b429337a4779271798bb45b6e07e94d311e6de89272f3a41ba766316b3a850b20701c207926c263a8d009264629a22412d1b8e6ad7905444b4da81c81362abc5dca42b970d0233e0ed980dd79677a96204",
		},
		// },
	}

	l := log.New()
	log.SetGlobalLogger(l)

	rx, err := NewReceiver(chain.BTPAddress(src), chain.BTPAddress(dst), []string{url}, opts, l)
	if err != nil {
		log.Fatal((err))
	}

	srcMsgCh := make(chan *chain.Message)
	srcErrCh, err := rx.Subscribe(context.TODO(),
		srcMsgCh,
		chain.SubscribeOptions{
			Seq:    138,
			Height: 29076,
		})
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-srcErrCh:
			log.Fatal(err)

		case <-srcMsgCh:
		}
	}
}
