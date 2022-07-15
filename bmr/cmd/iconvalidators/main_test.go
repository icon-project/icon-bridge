package main

import (
	"fmt"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/goloop/client"
)

func init() {
	cl = client.NewClientV3("http://localnets:9080/api/v3")
}

func TestBlockHashChain(t *testing.T) {
	var previd []byte
	for i := int64(0); i < 10; i++ {
		h, err := getBlockHeaderByHeight(0x0d12 + i)
		if err != nil {
			t.Fatal(err)
		}
		if i > 0 {
			assert.Equal(t, previd, h.PrevID)
		}
		previd = h.Hash()
		fmt.Println("hash:", common.BytesToHash(h.Hash()), "; prev:", common.BytesToHash(h.PrevID))
	}
}
