package executor

import (
	"context"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
)

const (
	PRIVKEYPOS = 0
	PUBKEYPOS  = 1
)

type evt struct {
	msg       *chain.EventLogInfo
	chainType chain.ChainType
}

type callBackFunc func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) error

type Script struct {
	Name        string
	Description string
	Callback    callBackFunc
}
