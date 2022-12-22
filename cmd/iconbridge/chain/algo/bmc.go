package algo

import (
	"context"
	"fmt"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
)

func (s *sender) HandleRelayMessage(ctx context.Context, _prev []byte, _msg []byte) (
	string, error) {
	res, err := s.callAbi(ctx, "HandleRelayMessage", []interface{}{_prev, _msg})
	if err != nil {
		return "", fmt.Errorf("Error calling Bmc Handle Relay Message: %w", err)
	}
	return res.TxIDs[0], nil
}

func (s *sender) GetBmcStatus(ctx context.Context) (*chain.BMCLinkStatus, error) {
	res, err := s.callAbi(ctx, "GetStatus", []interface{}{})
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
