package near

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/near/borsh-go"
	"github.com/reactivex/rxgo/v2"
	"math/big"
	"net/http"
	url_pkg "net/url"
	"strconv"
	"time"
)

const BmcContractMessageStateKey = "bWVzc2FnZQ=="

type Client struct {
	subClients []*Client
	api  Api
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

func NewClients(urls []string, logger log.Logger) []*Client {
	transport := &http.Transport{MaxIdleConnsPerHost: 1000}
	clients := make([]*Client, 0)

	for _, url := range urls {
		url, err := url_pkg.Parse(url)
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			logger:          logger,
			isMonitorClosed: false,
			api: &api{
				host:   url.Host,
				Client: jsonrpc.NewJsonRpcClient(&http.Client{Transport: transport}, url.String()),
				logger: logger,
			},
		}
		clients = append(clients, client)
	}

	return clients
}

func (c *Client) Call(method string, args interface{}, res interface{}) (*jsonrpc.Response, error) {
	return c.Client.Do(method, args, res)
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

func (c *Client) MonitorBlocks(height uint64, source string, concurrency uint, callback func(rxgo.Observable) error, subClient func() *Client) error {
	return callback(c.MonitorBlockHeight(int64(height)).Map(func(_ context.Context, offset interface{}) (interface{}, error) {
		if offset, Ok := (offset).(int64); Ok {
			block, err := subClient().api.getBlockByHeight(offset)
			bn := types.NewBlockNotification(offset)

			if err != nil {
				return bn, nil // TODO: Handle Error
			}

			bn.SetBlock(block)
			receipts, _ := subClient().GetReceipts(&block, source) // TODO: Handle Error

			bn.SetReceipts(receipts)

			return bn, nil
		}
		return nil, fmt.Errorf("error casting offset to int64")
	}, rxgo.WithPool(int(concurrency))).Serialize(int(height), func(_bn interface{}) int {
		bn := _bn.(*types.BlockNotification)

		return int(bn.Offset())
	}, rxgo.WithPool(int(concurrency)), rxgo.WithErrorStrategy(rxgo.ContinueOnError)).TakeUntil(func(i interface{}) bool {
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
			Index: uint64(i),
			Events: []*chain.Event{
				{
					Next:     event.Next,
					Sequence: uint64(sequence),
					Message:  event.Message,
				},
			},
			Height: uint64(block.Height()),
		})

	}

	return receipts, nil
}
