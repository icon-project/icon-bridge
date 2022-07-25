package executor

import (
	"context"
	"math/big"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
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

type keypair struct {
	PrivKey string
	PubKey  string
}

type fee struct {
	fixed       *big.Int
	numerator   *big.Int
	denominator *big.Int
}

var (
	ZeroEvents     = errors.New("Got zero event logs, expected at least one")
	StatusCodeZero = errors.New("Got status code zero(failed)")
)
