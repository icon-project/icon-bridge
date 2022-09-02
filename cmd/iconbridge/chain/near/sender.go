package near

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/errors"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/codec"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const (
	txMaxDataSize        = 64 * 1024
	functionCallMethod   = "handle_relay_message"
	gas                  = uint64(300000000000000)
	defaultSendTxTimeout = 15 * time.Second
	defaultGetTxTimeout  = 15 * time.Second
)

type Sender struct {
	clients     []*Client
	source      chain.BTPAddress
	destination chain.BTPAddress
	wallet      Wallet
	logger      log.Logger
	options     struct {
		BalanceThreshold types.BigInt `json:"balance_threshold"`
	}
}

func NewSender(source, destination chain.BTPAddress, urls []string, wallet wallet.Wallet, options json.RawMessage, logger log.Logger) (chain.Sender, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("empty urls: %v", urls)
	}

	sender := &Sender{
		clients:     NewClients(urls, logger),
		source:      source,
		destination: destination,
		wallet:      wallet,
		logger:      logger,
	}

	if err := json.Unmarshal(options, &sender.options); err != nil {
		logger.Panicf("fail to unmarshal options:%#v err:%+v", options, err)
		return nil, err
	}

	return sender, nil
}

func newMockSender(source, destination chain.BTPAddress, client *Client, wallet wallet.Wallet, _ map[string]interface{}, logger log.Logger) (*Sender, error) {
	clients := make([]*Client, 0)
	clients = append(clients, client)
	sender := &Sender{
		clients:     clients,
		source:      source,
		destination: destination,
		wallet:      wallet,
		logger:      logger,
	}

	return sender, nil
}

func (s *Sender) client() *Client {
	return s.clients[rand.Intn(len(s.clients))]
}

func (s *Sender) Segment(ctx context.Context, msg *chain.Message) (tx chain.RelayTx, newMsg *chain.Message, err error) {
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}

	if len(msg.Receipts) == 0 {
		return nil, msg, nil
	}

	rm := &chain.RelayMessage{
		Receipts: make([][]byte, 0),
	}

	var msgSize uint64

	newMsg = &chain.Message{
		From:     msg.From,
		Receipts: msg.Receipts,
	}
	for i, receipt := range msg.Receipts {
		rlpEvents, err := codec.RLP.MarshalToBytes(receipt.Events)
		if err != nil {
			return nil, nil, err
		}
		rlpReceipt, err := codec.RLP.MarshalToBytes(&chain.RelayReceipt{
			Index:  receipt.Index,
			Height: receipt.Height,
			Events: rlpEvents,
		})
		if err != nil {
			return nil, nil, err
		}
		newMsgSize := msgSize + uint64(len(rlpReceipt))
		if newMsgSize > txMaxDataSize {
			newMsg.Receipts = msg.Receipts[i:]
			break
		}
		msgSize = newMsgSize
		rm.Receipts = append(rm.Receipts, rlpReceipt)
	}

	message, err := codec.RLP.MarshalToBytes(rm)
	if err != nil {
		return nil, nil, err
	}

	tx, err = s.newRelayTransaction(ctx, msg.From.String(), message)
	if err != nil {
		return nil, nil, err
	}

	return tx, newMsg, nil
}

func (s *Sender) newRelayTransaction(ctx context.Context, prev string, message []byte) (*RelayTransaction, error) {
	if nearWallet, Ok := (s.wallet).(*wallet.NearWallet); Ok {
		accountId := nearWallet.Address()

		relayMessage := struct {
			Source  string `json:"source"`
			Message string `json:"message"`
		}{
			Source:  prev,
			Message: base64.URLEncoding.EncodeToString(message),
		}
		data, err := json.Marshal(relayMessage)
		if err != nil {
			return nil, err
		}

		actions := []types.Action{
			{
				Enum: 2,
				FunctionCall: types.FunctionCall{
					MethodName: functionCallMethod,
					Args:       data,
					Gas:        gas,
					Deposit:    *big.NewInt(0),
				},
			},
		}

		transaction := types.Transaction{
			SignerId:   types.AccountId(accountId),
			ReceiverId: types.AccountId(s.destination.ContractAddress()),
			PublicKey:  types.NewPublicKeyFromED25519(*nearWallet.Pkey),
			Actions:    actions,
		}

		return &RelayTransaction{
			Source:      prev,
			Message:     message,
			Transaction: transaction,
			client:      s.client(),
			wallet:      nearWallet,
		}, nil
	}
	return nil, fmt.Errorf("failed to cast wallet")
}

type RelayTransaction struct {
	Source      string `json:"source"`
	Message     []byte `json:"message"`
	Transaction types.Transaction
	client      *Client
	wallet      *wallet.NearWallet
	context     context.Context
}

func (relayTx *RelayTransaction) ID() interface{} {
	if relayTx.Transaction.Txid != [32]byte{} {
		return relayTx.Transaction.Txid
	}
	return nil
}

func (relayTx *RelayTransaction) Receipt(ctx context.Context) (blockHeight uint64, err error) {
	var txStatus types.TransactionResult
	if relayTx.Transaction.Txid == [32]byte{} {
		return 0, fmt.Errorf("no pending tx")
	}

	for i, isPending := 0, true; i < 5 && (isPending || err == errors.ErrUnknownTransaction); i++ {
		time.Sleep(time.Second)
		_, cancel := context.WithTimeout(ctx, defaultGetTxTimeout)
		defer cancel()

		txStatus, err = relayTx.client.api.getTransactionResult(relayTx.Transaction.Txid.Base58Encode(), string(relayTx.Transaction.SignerId))
		if err != nil {
			return blockHeight, err
		}
	}

	if txStatus.TransactionOutcome.BlockHash != [32]byte{} {
		block, err := relayTx.client.api.getBlockByHash(txStatus.TransactionOutcome.BlockHash.Base58Encode())
		if err != nil {
			return 0, err
		}

		blockHeight = uint64(block.Height())
	}

	//TODO: Handle errors
	return blockHeight, err
}

func (relayTx *RelayTransaction) Send(ctx context.Context) (err error) {
	relayTx.client.logger.WithFields(log.Fields{"prev": relayTx.Source}).Debug("handleRelayMessage: send tx")
	_ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeout)
	defer cancel()

	relayTx.context = _ctx
	publicKey := types.NewPublicKeyFromED25519(*relayTx.wallet.Pkey)
	nonce, err := relayTx.client.GetNonce(publicKey, string(relayTx.Transaction.SignerId))
	if nonce == -1 || err != nil {
		return err
	}

	relayTx.Transaction.Nonce = int(nonce) + 1
	blockHash, err := relayTx.client.api.getLatestBlockHash()
	if err != nil {
		return err
	}

	relayTx.Transaction.BlockHash = types.NewCryptoHash(blockHash)

	payload, err := relayTx.Transaction.Payload(relayTx.wallet)
	if err != nil {
		return err
	}

	txId, err := relayTx.client.SendTransaction(payload)
	if err != nil {
		return err
	}

	relayTx.Transaction.Txid = *txId
	relayTx.client.logger.WithFields(log.Fields{"tx": txId.Base58Encode()}).Debug("handleRelayMessage: tx sent")

	return nil
}

func (s *Sender) Status(ctx context.Context) (*chain.BMCLinkStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	status, err := s.client().GetBMCLinkStatus(s.destination, s.source)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (s *Sender) Balance(ctx context.Context) (balance, threshold *big.Int, err error) {
	balance, err = s.client().api.getBalance(s.wallet.Address())
	t := big.Int(s.options.BalanceThreshold)

	return balance, &t, err
}
