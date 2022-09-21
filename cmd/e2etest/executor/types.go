package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/config"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

type evt struct {
	msg       *chain.EventLogInfo
	chainType chain.ChainType
}

type callBackFunc func(ctx context.Context, tp *transferPoint, ts *testSuite) (txnRec *txnRecord, err error)

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

type eventTs struct {
	ChainName     chain.ChainType
	Sn            *big.Int
	EventType     chain.EventLogType
	BlockNumber   uint64
	TransactionID uint64
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
	config.FileConfig          `json:",squash"`
	LogLevel                   string            `json:"log_level"`
	ConsoleLevel               string            `json:"console_level"`
	LogWriter                  *log.WriterConfig `json:"log_writer,omitempty"`
	EnableExperimentalFeatures bool              `json:"enable_experimental_features"`
	Chains                     []*chain.Config   `json:"chains"`
	FeeAggregatorAddress       string            `json:"fee_aggregator"`
}

// common

func LoadConfig(file string) (*Config, error) {
	if len(file) == 0 {
		return nil, fmt.Errorf("Should specify config file using -config flag. Empty provided")
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open file %v", file)
	}
	cfg := &Config{}
	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "json.Decode file %v", file)
	}
	return cfg, nil
}

func SetLogger(cfg *Config) log.Logger {
	l := log.New()
	log.SetGlobalLogger(l)

	if lv, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.Panicf("Invalid log_level=%s", cfg.LogLevel)
	} else {
		l.SetConsoleLevel(lv)
	}

	if cfg.LogWriter != nil {
		if cfg.LogWriter.Filename == "" {
			log.Fatalln("Empty LogWriterConfig filename!")
		}
		var lwCfg log.WriterConfig
		lwCfg = *cfg.LogWriter
		lwCfg.Filename = cfg.ResolveAbsolute(lwCfg.Filename + "_" + strconv.Itoa(int(time.Now().Unix())))
		w, err := log.NewWriter(&lwCfg)
		if err != nil {
			log.Panicf("Fail to make writer err=%+v", err)
		}
		err = l.SetFileWriter(w)
		if err != nil {
			log.Panicf("Fail to set file l err=%+v", err)
		}
	}
	if lv, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.Panicf("Invalid log_level=%s", cfg.LogLevel)
	} else {
		l.SetLevel(lv)
	}
	if lv, err := log.ParseLevel(cfg.ConsoleLevel); err != nil {
		log.Panicf("Invalid console_level=%s", cfg.ConsoleLevel)
	} else {
		l.SetConsoleLevel(lv)
	}
	return l
}
