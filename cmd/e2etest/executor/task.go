package executor

import (
	"context"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	PRIVKEYPOS = 0
	PUBKEYPOS  = 1
)

type evt struct {
	msg       *chain.EventLogInfo
	chainType chain.ChainType
}

type args struct {
	watchRequestID uint64
	log            log.Logger
	src            chain.SrcAPI
	dst            chain.DstAPI
	srcKey         string
	srcAddr        string
	dstAddr        string
	coinName       string
	sinkChan       <-chan *evt
}

type callBackFunc func(ctx context.Context, args *args) error

type Script struct {
	Name        string
	Description string
	Callback    callBackFunc
}
