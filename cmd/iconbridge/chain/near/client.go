package near

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/near/borsh-go"
	"github.com/reactivex/rxgo/v2"
)

const BmcContractMessageStateKey = "bWVzc2FnZQ=="

type Client struct {
	api Api
	*jsonrpc.Client
	logger          log.Logger
	isMonitorClosed bool
}

type Wallet interface {
	Address() string
	Sign(data []byte) ([]byte, error)
}

type Api interface {
	broadcastTransaction(string) (string, error)
	broadcastTransactionAsync(string) (types.CryptoHash, error)
	getBlockByHash(string) (types.Block, error)
	getBlockByHeight(int64) (types.Block, error)
	getBmcLinkStatus(accountId string, link *chain.BTPAddress) (types.BmcStatus, error)
	getBmvStatus(accountId string) (types.BmvStatus, error)
	getContractStateChange(height int64, accountId string, keyPrefix string) (types.ContractStateChange, error)
	getLatestBlockHash() (string, error)
	getLatestBlockHeight() (int64, error)
	getNextBlockProducers(*types.CryptoHash) (types.NextBlockProducers, error)
	getNonce(string, string) (int64, error)
	getReceiptProof(blockHash, receiptId *types.CryptoHash, receiverId string) (types.ReceiptProof, error)
	getTransactionResult(string, string) (types.TransactionResult, error)
	getBalance(string) (*big.Int, error)
}

func newClients(urls []string, logger log.Logger) []*Client {
	transport := &http.Transport{MaxIdleConnsPerHost: 1000}
	clients := make([]*Client, 0)
	for _, url := range urls {
		client := &Client{
			logger:          logger,
			isMonitorClosed: false,
			api: &api{
				Client: jsonrpc.NewJsonRpcClient(&http.Client{Transport: transport}, url),
				logger: logger,
			},
		}
		clients = append(clients, client)
	}

	return clients
}

func (c *Client) GetBMCLinkStatus(destination, source chain.BTPAddress) (*chain.BMCLinkStatus, error) {
	bmcStatus, err := c.api.getBmcLinkStatus(destination.ContractAddress(), &source)
	if err != nil {
		return nil, err
	}

	linkstatus := &chain.BMCLinkStatus{}
	linkstatus.TxSeq = bmcStatus.TxSeq
	linkstatus.RxSeq = bmcStatus.RxSeq
	linkstatus.BMRIndex = bmcStatus.BMRIndex
	linkstatus.RotateHeight = bmcStatus.RotateHeight
	linkstatus.RotateTerm = bmcStatus.RotateTerm
	linkstatus.DelayLimit = bmcStatus.DelayLimit
	linkstatus.MaxAggregation = bmcStatus.MaxAggregation
	linkstatus.CurrentHeight = bmcStatus.CurrentHeight
	linkstatus.RxHeight = bmcStatus.RxHeight
	linkstatus.RxHeightSrc = bmcStatus.RxHeightSrc
	linkstatus.BlockIntervalSrc = bmcStatus.BlockIntervalSrc
	linkstatus.BlockIntervalDst = bmcStatus.BlockIntervalDst

	return linkstatus, nil
}

func (c *Client) GetNonce(publicKey types.PublicKey, accountId string) (int64, error) {
	nonce, err := c.api.getNonce(accountId, publicKey.Base58Encode())
	if err != nil {
		return -1, err
	}
	return nonce, nil
}

func (c *Client) SendTransaction(payload string) (*types.CryptoHash, error) {
	txId, err := c.api.broadcastTransactionAsync(payload)

	if err != nil {
		return nil, err
	}

	return &txId, nil

}

func (c *Client) MonitorBlockHeight(offset int64) rxgo.Observable {
	channel := make(chan rxgo.Item)
	go func(offset int64) {
		defer close(channel)

		lastestBlockHeight, err := c.api.getLatestBlockHeight()
		if err != nil {
			// TODO: Handle Error
			channel <- rxgo.Error(err)
			return
		}

		if lastestBlockHeight < 1 {
			channel <- rxgo.Error(errors.New("invalid block height"))
			return
		}

		for {
			rangeHeight := lastestBlockHeight - offset
			if rangeHeight < 5 {
				lastestBlockHeight, err = c.api.getLatestBlockHeight()
				if err != nil {
					// TODO: Handle Error
					fmt.Println(err)
				}

				rangeHeight = lastestBlockHeight - offset
				if rangeHeight < 3 {
					time.Sleep(time.Second * 2)
					continue
				}
			}

			channel <- rxgo.Of(offset)
			offset += 1
		}
	}(offset)

	return rxgo.FromChannel(channel, rxgo.WithCPUPool())
}

func (c *Client) MonitorBlocks(height uint64, callback func(rxgo.Observable) error) error {
	return callback(c.MonitorBlockHeight(int64(height)).FlatMap(func(i rxgo.Item) rxgo.Observable {
		if i.E != nil {
			return rxgo.Just(errors.New(i.E.Error()))()
		}

		block, err := c.api.getBlockByHeight(i.V.(int64))
		if err != nil {
			c.logger.Error(err)
			return rxgo.Empty()
		}

		return rxgo.Just(block)()
	}).TakeUntil(func(i interface{}) bool {
		return c.isMonitorClosed
	}))
}

func (c *Client) CloseMonitor() {
	c.isMonitorClosed = true
}

func (c *Client) GetReceipts(block *types.Block, accountId string) ([]*chain.Receipt, error) {
	receipts := make([]*chain.Receipt, 0)

	response, err := c.api.getContractStateChange(block.Height(), accountId, BmcContractMessageStateKey)
	if err != nil {
		return nil, err
	}

	for i, change := range response.Changes {
		var event struct {
			Next     chain.BTPAddress
			Sequence string
			Message  []byte
		}

		var eventsData string

		eventDataBytes, err := base64.URLEncoding.Strict().DecodeString(change.Data.ValueBase64)
		if err != nil {
			return nil, err
		}

		err = borsh.Deserialize(&eventsData, eventDataBytes)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(eventsData), &event)
		if err != nil {
			return nil, err
		}

		sequence, err := strconv.ParseInt(event.Sequence, 10, 64)
		if err != nil {
			return nil, err
		}

		receipts = append(receipts, &chain.Receipt{
			Index:  uint64(i),
			Events: []*chain.Event{
				{
					Next: event.Next,
					Sequence: uint64(sequence),
					Message: event.Message,
				},
			},
			Height: uint64(block.Height()),
		})

	}

	return receipts, nil
}
