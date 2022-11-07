package near

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/errors"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/near/borsh-go"
	"github.com/reactivex/rxgo/v2"
)

const BmcContractMessageStateKey = "bWVzc2FnZQ=="

type Client struct {
	api             IApi
	logger          log.Logger
	isMonitorClosed bool
}

type Wallet interface {
	Address() string
	Sign(data []byte) ([]byte, error)
}

type IApi interface {
	Block(param interface{}) (response types.Block, err error)
	BroadcastTxCommit(param interface{}) (response types.TransactionResult, err error)
	BroadcastTxAsync(param interface{}) (response types.CryptoHash, err error)
	CallFunction(param interface{}) (response types.CallFunctionResponse, err error)
	Changes(param interface{}) (response types.ContractStateChange, err error)
	Chunk(param interface{}) (response types.ChunkHeader, err error)
	LightClientProof(param interface{}) (response types.ReceiptProof, err error)
	Status(param interface{}) (response types.ChainStatus, err error)
	Transaction(param interface{}) (response types.TransactionResult, err error)
	ViewAccessKey(param interface{}) (response types.AccessKeyResponse, err error)
	ViewAccount(param interface{}) (response types.Account, err error)
}

type IClient interface {
	Api() IApi
	CloseMonitor()
	GetBalance(types.AccountId) (*big.Int, error)
	GetBlockByHash(types.CryptoHash) (types.Block, error)
	GetBlockByHeight(int64) (types.Block, error)
	GetBmcLinkStatus(destination, source chain.BTPAddress) (*chain.BMCLinkStatus, error)
	GetLatestBlockHash() (types.CryptoHash, error)
	GetNonce(publicKey types.PublicKey, accountId string) (int64, error)
	GetTransactionResult(types.CryptoHash, types.AccountId) (types.TransactionResult, error)
	GetReceipts(block *types.Block, accountId string) ([]*chain.Receipt, error)
	Logger() log.Logger
	MonitorBlocks(height uint64, source string, concurrency uint, callback func(rxgo.Observable) error, subClient func() IClient) error
	SendTransaction(payload string) (*types.CryptoHash, error)
	GetLatestBlockHeight() (int64, error)
	MonitorBlockHeight(offset int64) rxgo.Observable
	IsMonitorClosed() bool
}

func (c *Client) CloseMonitor() {
	c.isMonitorClosed = true
}

func (c *Client) Api() IApi {
	return c.api
}

func (c *Client) GetBalance(accountId types.AccountId) (balance *big.Int, err error) {
	param := struct {
		AccountId    types.AccountId `json:"account_id"`
		Finality     string          `json:"finality"`
		Request_type string          `json:"request_type"`
	}{
		AccountId:    accountId,
		Finality:     "final",
		Request_type: "view_account",
	}
	response, err := c.api.ViewAccount(param)
	if err != nil {
		return nil, err
	}

	return (*big.Int)(&response.Amount), nil
}

func (c *Client) GetBlock(param interface{}) (types.Block, error) {
	block, err := c.api.Block(param)
	if err != nil {
		return types.Block{}, err
	}

	return block, nil
}

func (c *Client) GetBlockByHash(blockHash types.CryptoHash) (types.Block, error) {
	param := struct {
		BlockId string `json:"block_id"`
	}{
		BlockId: blockHash.Base58Encode(),
	}

	return c.GetBlock(param)
}

func (c *Client) GetBlockByHeight(height int64) (types.Block, error) {
	param := struct {
		BlockId int64 `json:"block_id"`
	}{
		BlockId: height,
	}
	return c.GetBlock(param)
}

func (c *Client) GetBmcLinkStatus(destination, source chain.BTPAddress) (*chain.BMCLinkStatus, error) {
	var bmcStatus types.BmcStatus

	methodParam, err := json.Marshal(struct {
		Link string `json:"link"`
	}{
		Link: source.String(),
	})

	if err != nil {
		return nil, err
	}

	param := types.CallFunction{
		RequestType:  "call_function",
		Finality:     "final",
		AccountId:    types.AccountId(destination.ContractAddress()),
		MethodName:   "get_status",
		ArgumentsB64: base64.URLEncoding.EncodeToString(methodParam),
	}

	response, err := c.api.CallFunction(param)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response.Result, &bmcStatus)
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

func (c *Client) GetChainStatus() (types.ChainStatus, error) {
	chainStatus, err := c.api.Status([]interface{}{})
	if err != nil {
		return types.ChainStatus{}, err
	}

	return chainStatus, nil
}

func (c *Client) GetLatestBlockHash() (types.CryptoHash, error) {
	chainStatus, err := c.GetChainStatus()
	if err != nil {
		return types.CryptoHash{}, err
	}

	return chainStatus.SyncInfo.LatestBlockHash, nil
}

func (c *Client) GetLatestBlockHeight() (int64, error) {
	chainStatus, err := c.GetChainStatus()
	if err != nil {
		return 0, err
	}

	return chainStatus.SyncInfo.LatestBlockHeight, nil
}

func (c *Client) GetNonce(publicKey types.PublicKey, accountId string) (int64, error) {
	param := struct {
		AccountId    string `json:"account_id"`
		PublicKey    string `json:"public_key"`
		Finality     string `json:"finality"`
		Request_type string `json:"request_type"`
	}{
		AccountId:    accountId,
		PublicKey:    publicKey.Base58Encode(),
		Finality:     "final",
		Request_type: "view_access_key",
	}

	accessKeyResponse, err := c.api.ViewAccessKey(param)
	if err != nil {
		return -1, err
	}

	return accessKeyResponse.Nonce, nil
}

func (c *Client) GetReceipts(block *types.Block, accountId string) ([]*chain.Receipt, error) {
	receipts := make([]*chain.Receipt, 0)

	param := struct {
		ChangeType string   `json:"changes_type"`
		AccountIds []string `json:"account_ids"`
		KeyPrefix  string   `json:"key_prefix_base64"`
		BlockId    int64    `json:"block_id"`
	}{
		ChangeType: "data_changes",
		AccountIds: []string{accountId},
		KeyPrefix:  BmcContractMessageStateKey,
		BlockId:    block.Height(),
	}

	stateChanges, err := c.api.Changes(param)
	if err != nil {
		return nil, err
	}

	for i, change := range stateChanges.Changes {
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

func (c *Client) GetTransactionResult(transactionId types.CryptoHash, senderId types.AccountId) (types.TransactionResult, error) {
	param := []string{transactionId.Base58Encode(), string(senderId)}

	transactionResult, err := c.api.Transaction(param)
	if err != nil {
		return types.TransactionResult{}, err
	}

	return transactionResult, nil
}

func (c *Client) Logger() log.Logger {
	return c.logger
}

func (c *Client) IsMonitorClosed() bool {
	return c.isMonitorClosed
}

func (c *Client) MonitorBlockHeight(offset int64) rxgo.Observable {
	channel := make(chan rxgo.Item)
	go func(offset int64) {
		defer close(channel)

		lastestBlockHeight, err := c.GetLatestBlockHeight()
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
				lastestBlockHeight, err = c.GetLatestBlockHeight()
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

func (c *Client) MonitorBlocks(height uint64, source string, concurrency uint, callback func(rxgo.Observable) error, subClient func() IClient) error {
	return callback(c.MonitorBlockHeight(int64(height)).Map(func(_ context.Context, offset interface{}) (interface{}, error) {
		if offset, Ok := (offset).(int64); Ok {
			block, err := c.GetBlockByHeight(offset)
			bn := types.NewBlockNotification(offset)

			if err != nil && errors.Is(err, errors.ErrUnknownBlock) {
				return bn, nil
			} else if err != nil && !errors.Is(err, errors.ErrUnknownBlock) {
				return nil, err
			}

			bn.SetBlock(block)

			receipts, err := subClient().GetReceipts(&block, source)
			if err != nil && !errors.Is(err, errors.ErrUnknownBlock) {
				return bn, err
			}

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

func NewClient(endpoint string, logger log.Logger) (IClient, error) {
	transport := &http.Transport{MaxIdleConnsPerHost: 1000}
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %v, err: %v", endpoint, err)
	}

	return &Client{
		logger:          logger,
		isMonitorClosed: false,
		api: &api{
			host: url.Host,
			Client: jsonrpc.NewJsonRpcClient(&http.Client{Transport: transport}, url.String()).SetErrFunc(func(buffer json.RawMessage) error {
				var rpcErr errors.RpcError
				err = json.Unmarshal(buffer, &rpcErr)
				if err != nil {
					return err
				}

				return &rpcErr
			}),
		},
	}, nil
}

func newClients(urls []string, logger log.Logger) ([]IClient, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}

	clients := make([]IClient, 0)

	for _, url := range urls {
		client, err := NewClient(url, logger)
		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	return clients, nil
}

func (c *Client) SendTransaction(payload string) (*types.CryptoHash, error) {
	param := []string{payload}

	txId, err := c.api.BroadcastTxAsync(param)
	if err != nil {
		return nil, err
	}

	return &txId, nil
}
