package near

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/icon-project/icon-bridge/common/log"
)

type ReceiverCore struct {
	Log      log.Logger
	Opts     ReceiverOptions
	Cls      []*ethclient.Client
	BlockReq ethereum.FilterQuery
}

type ReceiverOptions struct {
	Verifier        *VerifierOptions `json:"verifier"`
	SyncConcurrency uint64           `json:"syncConcurrency"`
}

type VerifierOptions struct {
}
