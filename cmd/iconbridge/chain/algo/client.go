package algo

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	AlgoBlockRate   = 4
	BlockRetryLimit = 5
)

type Client struct {
	log   log.Logger
	algod *algod.Client
}

type IClient interface {
	Log() log.Logger
	GetBalance(ctx context.Context, hexAddr string) (*big.Int, error)
	WaitForTransaction(ctx context.Context, txId string) (models.PendingTransactionInfoResponse, error)
	GetLatestRound(ctx context.Context) (uint64, error)
	GetBlockbyRound(ctx context.Context, round uint64) (block *types.Block, err error)
	DecodeBtpMsg(log string) (*chain.Event, error)
}

func newClient(algodAccess []string, l log.Logger) (*Client, error) {
	algodClient, err := algod.MakeClient(algodAccess[0], algodAccess[1])
	if err != nil {
		return nil, err
	}
	cli := &Client{
		log:   l,
		algod: algodClient,
	}
	return cli, nil
}

func (cl *Client) WaitForTransaction(ctx context.Context, txId string) (models.PendingTransactionInfoResponse, error) {
	return future.WaitForConfirmation(cl.algod, txId, BlockRetryLimit, ctx)
}

func (cl *Client) GetBalance(ctx context.Context, walletAddr string) (*big.Int, error) {
	accountInfo, err := cl.algod.AccountInformation(walletAddr).Do(ctx)
	if err != nil {
		return nil, err
	} else {
		return new(big.Int).SetUint64(accountInfo.Amount), nil
	}
}

func (c *Client) Log() log.Logger {
	return c.log
}

// get latest block round
func (cl *Client) GetLatestRound(ctx context.Context) (uint64, error) {
	sta, err := cl.algod.Status().Do(ctx)
	return sta.LastRound, err
}

// get latest block number
func (cl *Client) GetBlockbyRound(ctx context.Context, round uint64) (*types.Block, error) {
	for i := 1; i <= BlockRetryLimit; i++ {
		block, err := cl.algod.Block(round).Do(ctx)
		if err != nil {
			time.Sleep(AlgoBlockRate * time.Second)
			continue
		}
		return &block, nil
	}
	err := fmt.Errorf("GetBlock reached retry limit")
	return nil, err
}

func (cl *Client) DecodeBtpMsg(log string) (*chain.Event, error) {
	//TODO this func should use ABI logic to go through the log string and decode it into a BTP message event
	return &chain.Event{}, nil
}
