package decoder_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder"
	ctr "github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts/bshPeriphery"
)

func TestBshEvent(t *testing.T) {
	url := "http://127.0.0.1:9500"
	btp_hmny_token_bsh_impl := "0xfAC8B63F77d8056A9BB45175b3DEd7D316D868D4"
	btp_hmny_nativecoin_bsh_periphery := "0xfEe5c5B2bc2f927335C60879d78304e4305CdBaC"
	contractUsed := btp_hmny_nativecoin_bsh_periphery
	m := map[ctr.ContractName]common.Address{
		ctr.BSHImpl:      common.HexToAddress(btp_hmny_token_bsh_impl),
		ctr.BSHPeriphery: common.HexToAddress(btp_hmny_nativecoin_bsh_periphery),
	}
	dec, err := decoder.New(url, m)
	if err != nil {
		t.Fatal(err)
	}
	const transferStartHex = "0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000001d00000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000003e6274703a2f2f30783562396137372e69636f6e2f68786663303136306133306565373033393861303134393438303531623331383963326563373261306200000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000dbd2fc137a30000000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000000000034554480000000000000000000000000000000000000000000000000000000000"
	txStartBytes, err := hex.DecodeString(transferStartHex)
	if err != nil {
		t.Fatal(err)
	}

	log := types.Log{
		Address:     common.HexToAddress(contractUsed),
		Topics:      []common.Hash{common.HexToHash("0x50d22373bb84ed1f9eeb581c913e6d45d918c05f8b1d90f0be168f06a4e6994a"), common.HexToHash("0x00000000000000000000000052c08a9a3a457e9ec8db545793ab9f0630dec4b4")},
		Data:        txStartBytes,
		BlockNumber: 100,
		TxHash:      common.HexToHash("0x123"),
		TxIndex:     1,
		BlockHash:   common.HexToHash("0x456"),
		Index:       2,
	}

	if out, err := dec.DecodeEventLogData(log); err != nil {
		t.Fatal(err)
	} else {
		for k, v := range out {
			if k == "TransferStart" && contractUsed == btp_hmny_nativecoin_bsh_periphery {
				res, ok := v.(*bshPeriphery.BshPeripheryTransferStart)
				if !ok {
					t.Fatal(errors.New("Problem"))
				} else {
					fmt.Println("First ", res.From, res.To, res.Sn, res.AssetDetails)
				}
			}
		}
	}
}
