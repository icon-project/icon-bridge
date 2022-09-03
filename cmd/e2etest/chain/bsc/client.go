package bsc

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/common/log"
)

func NewClients(urls []string, l log.Logger) (cls []*Client, err error) {
	for _, url := range urls {
		clrpc, err := rpc.Dial(url)
		if err != nil {
			l.Errorf("failed to create bsc rpc client: url=%v, %v", url, err)
			return nil, err
		}
		cleth := ethclient.NewClient(clrpc)
		cls = append(cls, &Client{
			log: l,
			rpc: clrpc,
			eth: cleth,
			//bmc: clbmc,
		})
	}
	return cls, nil
}

// grouped rpc api clients
type Client struct {
	log log.Logger
	rpc *rpc.Client
	eth *ethclient.Client
	//bmc *BMC
}

type BlockNotification struct {
	Hash   common.Hash
	Height *big.Int
	Header *types.Header
	Logs   []types.Log
}
