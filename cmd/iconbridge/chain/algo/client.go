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

// TODO create bmc interface w/ methods to compile, deploy and interact with the smart contract

type IClient interface {
	Log() log.Logger

	//algod
	GetBalance(ctx context.Context, hexAddr string) (*big.Int, error)
	PendingTransactionsByAddress(address string) ([]types.SignedTxn, error)
	WaitForTransaction(ctx context.Context, txId string) (models.PendingTransactionInfoResponse, error)
	GetLatestRound() (uint64, error)
	GetBlockbyRound(round uint64) (block *types.Block, err error)
	GetBlockHash(round uint64) (hash string, err error)
}

func newClient(algodAccess []string, l log.Logger) (cli *Client, err error) {

	algodClient, err := algod.MakeClient(algodAccess[0], algodAccess[1])

	if err != nil {
		l.Fatalf("Algod client could not be created: %s\n", err)
	}

	cli = &Client{
		log:   l,
		algod: algodClient,
	}
	return
}

func (cl *Client) WaitForTransaction(ctx context.Context, txId string) (models.PendingTransactionInfoResponse, error) {
	return future.WaitForConfirmation(cl.algod, txId, BlockRetryLimit, ctx)
}

func (cl *Client) GetBalance(ctx context.Context, walletAddr string) (*big.Int, error) {
	accountInfo, err := cl.algod.AccountInformation(walletAddr).Do(context.Background())
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
func (cl *Client) GetLatestRound() (uint64, error) {
	sta, err := cl.algod.Status().Do(context.Background())

	return sta.LastRound, err
}

// get latest block number
func (cl *Client) GetBlockbyRound(round uint64) (*types.Block, error) {
	for i := 1; i <= BlockRetryLimit; i++ {
		block, err := cl.algod.Block(round).Do(context.Background())
		if err != nil {
			time.Sleep(AlgoBlockRate * time.Second)
			continue
		}
		return &block, nil
	}
	err := fmt.Errorf("GetBlock reached retry limit")
	return nil, err
}

// get latest block number
func (cl *Client) GetBlockHash(round uint64) (hash string, err error) {
	for i := 1; i <= BlockRetryLimit; i++ {
		hashResponse, err := cl.algod.GetBlockHash(round).Do(context.Background())
		if err != nil {
			time.Sleep(AlgoBlockRate * time.Second)
			continue
		}
		hash = hashResponse.Blockhash
		return hash, nil
	}
	return "", err
}

func (cl *Client) DecodeBtpMsg(log string) (*chain.Event, error) {
	//TODO this func should use ABI logic to go through the log string and decode it into a BTP message event
	return &chain.Event{}, nil
}

func (cl *Client) GetBmcStatus(ctx context.Context) (*chain.BMCLinkStatus, error) {
	//TODO replace hardocded struct with value from BMC contract call
	ls := &chain.BMCLinkStatus{}
	ls.TxSeq = uint64(19)
	ls.RxSeq = uint64(17)
	ls.RxHeight = uint64(12341332)
	ls.CurrentHeight = uint64(100432121)
	return ls, nil
}

// not used atm
func (cl *Client) PendingTransactionsByAddress(address string) ([]types.SignedTxn, error) {
	_, tx, err := cl.algod.PendingTransactionsByAddress(address).Do(context.Background())

	return tx, err
}
