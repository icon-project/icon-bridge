package algo

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
)

func (s *sender) HandleRelayMessage(ctx context.Context, receipts [][]byte) (
	[]string, error) {
	abiFuncs := make([]AbiFunc, atomicTxnLimit)
	for _, receipt := range receipts {
		svcName, svcArgs, err := DecodeRelayMessage(hex.EncodeToString(receipt))
		if err != nil {
			return nil, fmt.Errorf("Error Decoding Relay Message: %w", err)
		}
		abiFuncs = append(abiFuncs, AbiFunc{svcName, []interface{}{svcArgs}})
	}

	res, err := s.callAbi(ctx, abiFuncs...)
	if err != nil {
		return nil, fmt.Errorf("Error calling Bmc Handle Relay Message: %w", err)
	}
	return res.TxIDs, nil
}

func (s *sender) GetBmcStatus(ctx context.Context) (*chain.BMCLinkStatus, error) {
	res, err := s.callAbi(ctx, AbiFunc{"GetStatus", []interface{}{}})
	if err != nil {
		return nil, fmt.Errorf("Error calling Bmc Handle Relay Message: %w", err)
	}
	bmcStatus := res.MethodResults[0].ReturnValue

	switch bmcStatus := bmcStatus.(type) {
	case [4]uint64:
		ls := &chain.BMCLinkStatus{
			TxSeq:         bmcStatus[0],
			RxSeq:         bmcStatus[1],
			RxHeight:      bmcStatus[2],
			CurrentHeight: bmcStatus[3],
		}
		return ls, nil
	}
	return nil, fmt.Errorf("BmcStatus - Couldnt parse abi's return interface")
}
