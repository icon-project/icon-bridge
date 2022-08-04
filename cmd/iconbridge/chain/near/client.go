package near

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"net/http"
)

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
		return  nil, err
	}

	return &txId, nil

}