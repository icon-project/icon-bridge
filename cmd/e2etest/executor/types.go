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

type callBackFunc func(ctx context.Context, srcChain, dstChain chain.ChainType, coinNames []string, ts *testSuite) (txnRec *txnRecord, err error)

type Script struct {
	Name        string
	Type        string
	Description string
	Callback    callBackFunc
}

type configureCallBack func(ctx context.Context, conf *configPoint, ts *testSuite) (txnRec *txnRecord, err error)
type ConfigureScript struct {
	Name        string
	Type        string
	Description string
	Callback    configureCallBack
}

type keypair struct {
	PrivKey string
	PubKey  string
}

type txnRecord struct {
	feeRecords []*feeRecord
	addresses  map[chain.ChainType][]keypair
}

type feeRecord struct {
	ChainName chain.ChainType
	Sn        *big.Int
	Fee       map[string]*big.Int
}

type pointGenerator struct {
	cfgPerChain    map[chain.ChainType]*chain.Config
	maxBatchSize   *int
	transferFilter func([]*transferPoint) []*transferPoint
	configFilter   func([]*configPoint) []*configPoint
}

type transferPoint struct {
	SrcChain  chain.ChainType
	DstChain  chain.ChainType
	CoinNames []string
	Amounts   []*big.Int
}

type configPoint struct {
	chainName   chain.ChainType
	TokenLimits map[string]*big.Int
	Fee         map[string][2]*big.Int
}

var (
	ZeroEvents               = errors.New("Got zero event logs, expected at least one")
	StatusCodeZero           = errors.New("Got status code zero(failed)")
	ExternalContextCancelled = errors.New("External Context Cancelled")
	NilEventReceived         = errors.New("Nil Event Received")
	InsufficientNativeToken  = errors.New("Insufficient Native Token")
	InsufficientWrappedCoin  = errors.New("Insufficient Wrapped Coin")
	InsufficientUnknownCoin  = errors.New("Insufficient Unknown Coin")
	UnsupportedCoinArgs      = errors.New("Unsupported Coin Args")
	IgnoreableError          = errors.New("Ignoreable Error")
)

type Config struct {
	Chains               []*chain.Config `json:"chains"`
	FeeAggregatorAddress string          `json:"fee_aggregator"`
}

// common
