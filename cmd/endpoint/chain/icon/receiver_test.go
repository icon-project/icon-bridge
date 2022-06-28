package icon

import (
	"context"
	"fmt"
	"testing"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestReceiver(t *testing.T) {
	srcAddress := "btp://0x5b9a77.icon/cx0f011b8b10f2c0d850d5135ef57ea42120452003"
	dstAddress := "btp://0x6357d2e0.hmny/0x7a6DF2a2CC67B38E52d2340BF2BDC7c9a32AaE91"
	srcEndpoint := []string{"http://localhost:9080/api/v3/default"}
	var height uint64 = 15000
	btp_icon_token_bsh := "cx3e836c763af780392a00a9ac2fc6e0471c95cb50"
	btp_icon_nativecoin_bsh := "cxe4b60a773c63961aa2303961483c3c95b9de3360"
	addrToName := map[string]chain.ContractName{
		btp_icon_token_bsh:      chain.TokenIcon,
		btp_icon_nativecoin_bsh: chain.NativeIcon,
	}
	l := log.New()
	log.SetGlobalLogger(l)
	recv, err := NewReceiver(chain.BTPAddress(srcAddress), chain.BTPAddress(dstAddress), srcEndpoint, l, addrToName)
	if err != nil {
		panic(err)
	}

	if err := recv.Subscribe(context.Background(), height); err != nil {
		panic(err)
	} else {
		for {
			select {
			case err := <-recv.errChan:
				panic(err)
			case msgs := <-recv.sinkChan:
				res, ok := msgs.Res.([]*TxnEventLog)
				if !ok {
					panic(err)
				}
				for _, msg := range res {
					fmt.Println(msg)
				}

			}
		}
	}

}

func TestIconEvent(t *testing.T) {
	btp_icon_token_bsh := "cx5924a147ae30091ed9c6fe0c153ef77de4132902"
	m := map[string]chain.ContractName{
		btp_icon_token_bsh: chain.TokenIcon,
	}
	parser, err := NewParser(m)
	if err != nil {
		t.Fatal(err)
	}
	log := &TxnEventLog{
		Addr:    Address("cx5924a147ae30091ed9c6fe0c153ef77de4132902"),
		Indexed: []string{"TransferStart(Address,str,int,bytes)", "hx4a707b2ecbb5f40a8d761976d99244f53575eeb6"},
		Data:    []string{"btp://0x6357d2e0.hmny/0x8BE8641225CC0Afdb24499409863E8E3f6557C32", "0x25", "0xd6d583455448880dbd2fc137a30000872386f26fc10000"},
	}
	res, eventType, err := parser.Parse(log)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("EventType %v  Res %+v", eventType, res)
}
