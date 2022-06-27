package decoder_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/icon"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder"
	ctr "github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts/nativeHmy"
	"github.com/icon-project/icon-bridge/cmd/endpoint/decoder/contracts/tokenIcon"
)

func TestBshEvent(t *testing.T) {
	urlPerChain := map[chain.ChainType]string{chain.HMNY: "http://127.0.0.1:9500"}
	btp_hmny_token_bsh_impl := "0xfAC8B63F77d8056A9BB45175b3DEd7D316D868D4"
	btp_hmny_nativecoin_bsh_periphery := "0xfEe5c5B2bc2f927335C60879d78304e4305CdBaC"
	contractUsed := btp_hmny_nativecoin_bsh_periphery
	m := map[string]ctr.ContractName{
		btp_hmny_token_bsh_impl:           ctr.TokenHmy,
		btp_hmny_nativecoin_bsh_periphery: ctr.NativeHmy,
	}
	dec, err := decoder.New(urlPerChain, m)
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

	if out, err := dec.DecodeEventLogData(log, log.Address.String()); err != nil {
		t.Fatal(err)
	} else {
		for k, v := range out {
			if k == "TransferStart" && contractUsed == btp_hmny_nativecoin_bsh_periphery {
				res, ok := v.(*nativeHmy.NativeHmyTransferStart)
				if !ok {
					t.Fatal(errors.New("Problem"))
				} else {
					fmt.Println("First ", res.From, res.To, res.Sn, res.AssetDetails)
				}
			}
		}
	}
}

func TestIconEvent(t *testing.T) {
	btp_icon_token_bsh := "cx5924a147ae30091ed9c6fe0c153ef77de4132902"
	m := map[string]ctr.ContractName{
		btp_icon_token_bsh: ctr.TokenIcon,
	}
	urlPerChain := map[chain.ChainType]string{chain.HMNY: "http://127.0.0.1:9500"}
	dec, err := decoder.New(urlPerChain, m)
	if err != nil {
		t.Fatal(err)
	}
	log := icon.TxnEventLog{
		Addr:    icon.Address("cx5924a147ae30091ed9c6fe0c153ef77de4132902"),
		Indexed: []string{"TransferStart(Address,str,int,bytes)", "hx4a707b2ecbb5f40a8d761976d99244f53575eeb6"},
		Data:    []string{"btp://0x6357d2e0.hmny/0x8BE8641225CC0Afdb24499409863E8E3f6557C32", "0x25", "0xd6d583455448880dbd2fc137a30000872386f26fc10000"},
	}
	if out, err := dec.DecodeEventLogData(log, btp_icon_token_bsh); err != nil {
		t.Fatalf("%+v", err)
	} else {
		for k, v := range out {
			if k == "TransferStart" && log.Addr == icon.Address(btp_icon_token_bsh) {
				res, ok := v.(*tokenIcon.TokenIconTransferStart)
				if !ok {
					t.Fatal(errors.New("Problem"))
				} else {
					fmt.Println("First ", res.From, res.To, res.Sn, res.Assets)
				}
			}
		}
	}

}
