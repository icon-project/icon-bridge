package algo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
)

func Test_Abi(t *testing.T) {
	s, err := createTestSender(testnetAccess)
	if err != nil {
		t.Logf("Failed creting new sender:%v", err)
		t.FailNow()
	}
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	abiCall, err := s.(*sender).callAbi(ctx, AbiFunc{"sendMessage",
		[]interface{}{"this", "string", 19}})

	if err != nil {
		t.Logf("Failed to call abi:%v", err)
		t.FailNow()
	}
	fmt.Println(abiCall)
}

func Test_Segment(t *testing.T) {
	s, err := createTestSender(testnetAccess)
	if err != nil {
		t.Logf("Failed creting new sender:%v", err)
		t.FailNow()
	}
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	msg := &chain.Message{
		From: chain.BTPAddress(iconBmc),
		Receipts: []*chain.Receipt{{
			Index:  10,
			Height: 1,
			Events: []*chain.Event{
				{Next: "algobmc", Sequence: 19, Message: []byte{97, 98, 99, 100, 101, 102}},
			}}, {
			Index:  20,
			Height: 1,
			Events: []*chain.Event{
				{Next: "algobmc", Sequence: 20, Message: []byte{64, 2, 4, 111, 55, 23}},
			}},
		},
	}
	_, _, err = s.Segment(ctx, msg)

	if err != nil {
		t.Logf("Couldn't segment message:%v", err)
		t.FailNow()
	}
}
