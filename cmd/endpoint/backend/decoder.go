package backend

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/icon"
	"github.com/icon-project/icon-bridge/common/errors"
)

const (
	TransferEndSign      = "TransferEnd(Address,int,int,bytes)"
	TransferStartSign    = "TransferStart(Address,str,int,bytes)"
	TransferReceivedSign = "TransferReceived(str,Address,int,bytes)"
	ICXTransferSign      = "ICXTransfer(Address,Address,int)"
	MessageSign          = "Message(str,int,bytes)"
	ApprovalSign         = "Approval(Address,Address,int)"
)

func decodeIconLog(txn *icon.TxnLog) (*DecodedEvent, error) {

	// for _, l := range txn.EventLogs {

	// 	if l.Indexed[0] == TransferStartSign {
	// 		_, _, _, err := rlpDecodeData(l.Data[len(l.Data)-1], len(l.Data))
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	} else if l.Indexed[0] == TransferReceivedSign {
	// 		_, _, _, err := rlpDecodeData(l.Data[len(l.Data)-1], len(l.Data))
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }
	// newEventLog := make([]chain.EventLog, len(txn.EventLogs))
	// newEventLog[li] = chain.EventLog{Addr: string(l.Addr), Indexed: l.Indexed, Data: l.Data}

	// status := int64(1)
	// if s, err := txn.Status.Value(); err != nil {
	// 	status = s
	// }
	// return &DecodedEvent{
	// 	TxHash:      string(txn.TxHash),
	// 	From:        string(txn.From),
	// 	To:          string(txn.To),
	// 	EventLogs:   newEventLog,
	// 	Status:      status,
	// 	BlockHeight: txn.BlockHeight,
	// }
	return nil, nil
}

type AssetDetailToken struct {
	CoinName string
	Value    *big.Int
	Fee      *big.Int
}
type AssetDetailNative struct {
	CoinName string
	Value    *big.Int
}

func rlpDecodeHex(str string, out interface{}) (interface{}, error) {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}
	input, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	err = rlp.Decode(bytes.NewReader(input), out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func rlpDecodeData(data string, num int) (*string, *big.Int, *big.Int, error) {
	if num == 3 {
		in := &[]AssetDetailToken{}
		res, _ := rlpDecodeHex(data, in)
		w, _ := res.(*[]AssetDetailToken)
		coin := (*w)[0].CoinName
		value := (*w)[0].Value
		fee := (*w)[0].Fee
		return &coin, value, fee, nil
	} else if num == 2 {
		in := &[]AssetDetailNative{}
		res, _ := rlpDecodeHex(data, in)
		w, _ := res.(*[]AssetDetailNative)
		coin := (*w)[0].CoinName
		value := (*w)[0].Value
		return &coin, value, nil, nil
	}
	return nil, nil, nil, errors.New("No Match")
}
