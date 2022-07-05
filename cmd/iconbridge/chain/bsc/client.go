/*
 * Copyright 2021 ICON Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bsc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/common/wallet"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/bsc/binding"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	DefaultTimeout  = 50 * time.Second //
	DefaultGasLimit = 8000000
)

var (
	BlockRetryInterval        = time.Second * 3
	BlockRetryLimit           = 5
	ConnectionSleepInterval   = time.Second * 40
	ConnectionSleepRetryLimit = 4
)

type Wallet interface {
	Sign(data []byte) ([]byte, error)
	Address() string
}

type client struct {
	log     log.Logger
	eth     []*ethclient.Client
	rpc     []*rpc.Client
	chainID *big.Int
	stop    <-chan bool
	bmc     []*binding.BMC
}

func (c *client) newTransactOpts(w Wallet) (*bind.TransactOpts, error) {
	txo, err := bind.NewKeyedTransactorWithChainID(w.(*wallet.EvmWallet).Skey, c.chainID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	txo.GasPrice, _ = c.ethClient().SuggestGasPrice(ctx)
	txo.GasLimit = uint64(DefaultGasLimit)
	return txo, nil
}

func (c *client) SignTransaction(signerKey *ecdsa.PrivateKey, tx *types.Transaction) error {
	signer := types.LatestSignerForChainID(c.chainID)
	tx, err := types.SignTx(tx, signer, signerKey)
	if err != nil {
		c.log.Errorf("could not sign tx: %v", err)
		return err
	}
	return nil
}

func (c *client) SendTransaction(tx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	err := c.ethClient().SendTransaction(ctx, tx)

	if err != nil {
		c.log.Errorf("could not send tx: %v", err)
		return nil
	}

	return nil
}

func (c *client) CallContract(callData ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	data, err := c.ethClient().CallContract(ctx, callData, blockNumber)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *client) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	tr, err := c.ethClient().TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

func (c *client) GetTransaction(hash common.Hash) (*types.Transaction, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	tx, pending, err := c.ethClient().TransactionByHash(ctx, hash)
	if err != nil {
		return nil, pending, err
	}
	return tx, pending, err
}

func (c *client) GetBlockByHeight(height *big.Int) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	block, err := c.ethClient().BlockByNumber(ctx, height)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *client) GetHeaderByHeight(height *big.Int) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	block, err := c.ethClient().BlockByNumber(ctx, height)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *client) GetBlockReceipts(block *types.Block) ([]*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	var receipts []*types.Receipt
	for _, tx := range block.Transactions() {
		receipt, err := c.ethClient().TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}
	return receipts, nil
}

func (c *client) FilterLogs(query ethereum.FilterQuery) ([]types.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	logs, err := c.ethClient().FilterLogs(ctx, query)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (c *client) GetChainID() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return c.ethClient().ChainID(ctx)
}

func (c *client) MonitorBlock(ctx context.Context, p *BlockRequest, cb func(b *BlockNotification) error) error {
	return c.Poll(ctx, p, cb)
}

func (c *client) Poll(ctx context.Context, p *BlockRequest, cb func(b *BlockNotification) error) error {

	c.log.Infof("Polling start")
	current := p.Height
	var retry = BlockRetryLimit
	var sleepRetry = ConnectionSleepRetryLimit
	for {
		select {
		case <-ctx.Done():
			fmt.Errorf("Context Closed")
		case <-c.stop:
			return nil
		default:
			// Exhausted all error retries
			if retry == 0 {
				c.log.Error("Polling failed, retries exceeded")
				//l.sysErr <- ErrFatalPolling
				if sleepRetry == 0 {
					c.log.Errorf("Cannot connect even after sleeping for %d retries each for %d Seconds", ConnectionSleepRetryLimit, ConnectionSleepInterval.Seconds())
					return nil
					//todo: stop relay panic here
				}
				c.log.Errorf("Going to sleep for %d seconds", ConnectionSleepInterval.Seconds())
				sleepRetry--
				<-time.After(ConnectionSleepInterval)
				retry = BlockRetryLimit
				continue
			}
			latestHeader, err := c.ethClient().HeaderByNumber(ctx, current) // c.GetHeaderByHeight(current)
			if err != nil {
				//c.log.Error("Unable to get latest block ", current, err)
				retry--
				<-time.After(BlockRetryInterval)
				continue
			}

			if latestHeader.Number.Cmp(current) < 0 {
				c.log.Info("Block not ready, will retry", "target:", current, "latest:", latestHeader.Number)
				<-time.After(BlockRetryInterval)
				continue
			}

			query := ethereum.FilterQuery{
				FromBlock: current,
				ToBlock:   current,
				Addresses: []common.Address{
					p.SrcContractAddress,
				},
			}

			logs, err := c.FilterLogs(query)
			if err != nil {
				c.log.Info("Unable to get logs ", err)
				continue
			}

			v := &BlockNotification{
				Height: current,
				Hash:   latestHeader.Hash(),
				Header: latestHeader,
				Logs:   logs,
			}
			if err := cb(v); err != nil {
				c.log.Errorf(err.Error())
			}

			current.Add(current, big.NewInt(1))
			retry = BlockRetryLimit
			sleepRetry = ConnectionSleepRetryLimit
		}
	}
}

func (c *client) Monitor(cb func(b *BlockNotification) error) error {
	subch := make(chan *types.Header)
	sub, err := c.ethClient().SubscribeNewHead(context.Background(), subch)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				c.log.Fatal(err)
			case header := <-subch:
				b := &BlockNotification{Hash: header.Hash(), Height: header.Number, Header: header}
				err := cb(b)
				if err != nil {
					return
				}
				c.log.Debugf("MonitorBlock %v", header.Number.Int64())
			}
		}
	}()

	return nil
}

func (c *client) CloseAllMonitor() {
	c.log.Debugf("CloseMonitor")
	for _, eth := range c.eth {
		eth.Close()
	}
	for _, rpc := range c.rpc {
		rpc.Close()
	}

}

func NewClient(urls []string, bmc string, log log.Logger) (cl *client, err error) {
	if len(urls) == 0 {
		log.Errorf("invalid client urls: %v", urls)
	}
	c := &client{
		log: log,
	}
	for _, url := range urls {
		rpcCl, err := rpc.Dial(url)
		if err != nil {
			log.Errorf("failed to create BSC rpc client: %v", err)
			return nil, err
		}
		ethCl := ethclient.NewClient(rpcCl)
		c.rpc = append(c.rpc, rpcCl)
		c.eth = append(c.eth, ethCl)
		bmc, err := binding.NewBMC(common.HexToAddress(bmc), ethCl)
		if err != nil {
			log.Errorf("failed to create bmc binding to bsc ethclient: , %v", err)
			return nil, err
		}
		c.bmc = append(c.bmc, bmc)
		c.chainID, _ = c.GetChainID()
		log.Tracef("Client Connected Chain ID: ", c.chainID)
	}
	return c, nil
}

func (c *client) ethClient() *ethclient.Client {
	id := rand.Intn(len(c.eth))
	return c.eth[id]
}

func (c *client) bmcCl() *binding.BMC {
	id := rand.Intn(len(c.bmc))
	log.Tracef("bmc ID: ", id)
	return c.bmc[id]
}
