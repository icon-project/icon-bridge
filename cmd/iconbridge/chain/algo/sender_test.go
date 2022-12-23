package algo

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/bmizerany/assert"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
)

func Test_Abi(t *testing.T) {
	s, err := createTestSender(sandboxAccess)
	if err != nil {
		t.Logf("Failed creting new sender:%v", err)
		t.FailNow()
	}
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	_, err = s.(*sender).callAbi(ctx, "sendMessage",
		[]interface{}{"this", "string", 19})

	/* if err != nil {
		t.Logf("Failed calling abi:%v", err)
		t.FailNow()
	}
	concatString := ret.MethodResults[0].ReturnValue.(string)
	assert.Equal(t, concatString, "thisstringisjoined") */
}

func Test_Segment(t *testing.T) {
	s, err := createTestSender(sandboxAccess)
	if err != nil {
		t.Logf("Failed creting new sender:%v", err)
		t.FailNow()
	}
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	msg := &chain.Message{
		From: chain.BTPAddress(icon_bmc),
		Receipts: []*chain.Receipt{{
			Index:  0,
			Height: 1,
			Events: []*chain.Event{
				{
					Message: []byte{97, 98, 99, 100, 101, 102},
				},
				{
					Message: []byte{44, 32, 33, 4, 101, 255},
				},
			},
		},
			{
				Index:  0,
				Height: 2,
				Events: []*chain.Event{
					{
						Message: []byte{55, 56, 222, 34, 6, 3},
					},
					{
						Message: []byte{64, 2, 4, 111, 55, 23},
					},
				},
			}},
	}
	tx, _, err := s.Segment(ctx, msg)

	if err != nil {
		t.Logf("Couldn't segment message:%v", err)
		t.FailNow()
	}

	sss := tx.(*relayTx).msg

	recovered_pay_bytes := make([]byte, 1000000)
	base64.StdEncoding.Decode(recovered_pay_bytes, sss)
	rm := &chain.RelayMessage{}
	msgpack.Decode(recovered_pay_bytes, &rm)

	recSli := make([]chain.Receipt, 0)
	for _, r := range rm.Receipts {
		decodedReceipt := make([]byte, 1000)
		base64.StdEncoding.Decode(decodedReceipt, r)

		var finalRcp chain.Receipt
		msgpack.Decode(decodedReceipt, &finalRcp)

		recSli = append(recSli, finalRcp)
	}

	for i := range recSli {
		assert.Equal(t, *msg.Receipts[i], recSli[i])
	}
}

func Test_file(t *testing.T) {

	absPath, _ := filepath.Abs("../../../../pyteal/build/bmc")
	fmt.Println(absPath)
}
