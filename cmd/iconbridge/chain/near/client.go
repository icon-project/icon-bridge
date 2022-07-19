package near

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/account"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
	"github.com/icon-project/icon-bridge/common/log"
	"net/http"
	"strings"
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
	broadcastTransactionAsync(string) (string, error)
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

func (c *Client) GetNonce(publicKey string, accountId string) (int64, error) {
	var err error
	var publicKeyString string

	if !strings.HasPrefix(publicKey, "ed25519:") {
		var publicKeyBytes []byte

		if len(publicKey) == 64 {
			publicKeyBytes, err = hex.DecodeString(publicKey)
			if err != nil {
				return -1, err
			}

			publicKeyString = account.PublicKeyToString(publicKeyBytes)

		} else {
			publicKeyBytes = base58.Decode(publicKey)
			if len(publicKeyBytes) == 0 {
				return -1, fmt.Errorf("b58 decode public key error, %s", publicKey)
			}

			publicKeyString = "ed25519:" + publicKey
		}
	} else {
		publicKeyString = publicKey
	}

	nonce, err := c.api.getNonce(accountId, publicKeyString)
	if err != nil {
		return -1, err
	}
	return nonce, nil
}
