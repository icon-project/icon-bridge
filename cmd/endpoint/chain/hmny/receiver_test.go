package hmny

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/harmony-one/harmony/core/types"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

func TestHmnyReceiver(t *testing.T) {
	const (
		src = "btp://0x6357d2e0.hmny/0x7a6DF2a2CC67B38E52d2340BF2BDC7c9a32AaE91"
		dst = "btp://0x5b9a77.icon/cx0f011b8b10f2c0d850d5135ef57ea42120452003"
		url = "http://localhost:9500"
	)

	l := log.New()
	log.SetGlobalLogger(l)
	btp_hmny_token_bsh_impl := "0x8283e3bE7ac5f6dB332Df605f20E2B4c9977c662"
	btp_hmny_nativecoin_bsh_periphery := "0xfad748a1063a40FF447B5D766331904d9bedDC26"
	addrToName := map[string]chain.ContractName{
		btp_hmny_token_bsh_impl:           chain.TokenHmy,
		btp_hmny_nativecoin_bsh_periphery: chain.NativeHmy,
	}
	rx, err := NewReceiver(chain.BTPAddress(src), chain.BTPAddress(dst), []string{url}, l, addrToName)
	if err != nil {
		log.Fatal((err))
	}

	err = rx.Subscribe(context.TODO(), 100)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-rx.errChan:
			t.Fatalf("%+v", err)

		case msgs := <-rx.sinkChan:
			res, ok := msgs.Res.([]*types.Log)
			if !ok {
				t.Fatalf("%+v", err)
			}
			for _, msg := range res {
				fmt.Println(msg)
			}
		}
	}
}

func TestBshEvent(t *testing.T) {

	btp_hmny_token_bsh_impl := "0xfAC8B63F77d8056A9BB45175b3DEd7D316D868D4"
	btp_hmny_nativecoin_bsh_periphery := "0xfEe5c5B2bc2f927335C60879d78304e4305CdBaC"
	contractUsed := btp_hmny_nativecoin_bsh_periphery
	m := map[string]chain.ContractName{
		btp_hmny_token_bsh_impl:           chain.TokenHmy,
		btp_hmny_nativecoin_bsh_periphery: chain.NativeHmy,
	}
	p, err := NewParser("http://127.0.0.1:9500", m)
	if err != nil {
		t.Fatal(err)
	}
	const transferStartHex = "0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000001d00000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000003e6274703a2f2f30783562396137372e69636f6e2f68786663303136306133306565373033393861303134393438303531623331383963326563373261306200000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000dbd2fc137a30000000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000000000034554480000000000000000000000000000000000000000000000000000000000"
	txStartBytes, err := hex.DecodeString(transferStartHex)
	if err != nil {
		t.Fatal(err)
	}

	log := &types.Log{
		Address:     common.HexToAddress(contractUsed),
		Topics:      []common.Hash{common.HexToHash("0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a"), common.HexToHash("0x00000000000000000000000052c08a9a3a457e9ec8db545793ab9f0630dec4b4")},
		Data:        txStartBytes,
		BlockNumber: 100,
		TxHash:      common.HexToHash("0x123"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x456"),
		Index:       2,
	}

	res, eventType, err := p.Parse(log)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("EventType %v  Res %+v", eventType, res)
}
