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

type SenderConfig struct {
	source      chain.BTPAddress
	destination chain.BTPAddress
	options     json.RawMessage
	wallet      wallet.Wallet
}

type Sender struct {
	clients     []IClient
	source      chain.BTPAddress
	destination chain.BTPAddress
	wallet      Wallet
	logger      log.Logger
	options     struct {
		BalanceThreshold types.BigInt `json:"balance_threshold"`
	}
}

func senderFactory(source, destination chain.BTPAddress, urls []string, wallet wallet.Wallet, options json.RawMessage, logger log.Logger) (chain.Sender, error) {
	clients, err := newClients(urls, logger)
	if err != nil {
		return nil, err
	}

	return NewSender(SenderConfig{source, destination, options, wallet}, logger, clients...)
}

func NewSender(config SenderConfig, logger log.Logger, clients ...IClient) (*Sender, error) {
	if len(clients) == 0 {
		return nil, fmt.Errorf("nil clients")
	}

	s := &Sender{
		clients:     clients,
		wallet:      config.wallet,
		logger:      logger,
		source:      config.source,
		destination: config.destination,
	}

	if err := json.Unmarshal(config.options, &s.options); err != nil && config.options != nil {
		logger.Panicf("fail to unmarshal options:%#v err:%+v", config.options, err)
		return nil, err
	}

	return s, nil
}

func (s *Sender) client() IClient {
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

	receipt := msg.Receipts[0]

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

	rm.Receipts = append(rm.Receipts, rlpReceipt)

	message, err := codec.RLP.MarshalToBytes(rm)
	if err != nil {
		return nil, nil, err
	}

	tx, err = s.newRelayTransaction(ctx, msg.From.String(), message)
	if err != nil {
		return nil, nil, err
	}

	return tx, &chain.Message{
		From:     msg.From,
		Receipts: msg.Receipts[1:],
	}, nil
}

func (s *Sender) newRelayTransaction(ctx context.Context, prev string, message []byte) (*RelayTransaction, error) {
	if nearWallet, Ok := (s.wallet).(*wallet.NearWallet); Ok {
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

		return NewRelayTransaction(ctx, nearWallet, s.destination.ContractAddress(), s.client(), actions), nil
	}

	return nil, fmt.Errorf("failed to cast wallet")
}

type RelayTransaction struct {
	Transaction types.Transaction
	client      IClient
	wallet      *wallet.NearWallet
	context     context.Context
}

func NewRelayTransaction(context context.Context, wallet *wallet.NearWallet, destination string, client IClient, actions []types.Action) *RelayTransaction {
	transaction := types.Transaction{
		SignerId:   types.AccountId(wallet.Address()),
		ReceiverId: types.AccountId(destination),
		PublicKey:  types.NewPublicKeyFromED25519(*wallet.Pkey),
		Actions:    actions,
	}

	return &RelayTransaction{
		Transaction: transaction,
		client:      client,
		wallet:      wallet,
		context:     context,
	}
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

		txStatus, err = relayTx.client.GetTransactionResult(relayTx.Transaction.Txid, relayTx.Transaction.SignerId)
		if err != nil {
			return blockHeight, err
		}
	}

	if txStatus.TransactionOutcome.BlockHash != [32]byte{} {
		block, err := relayTx.client.GetBlockByHash(txStatus.TransactionOutcome.BlockHash)
		if err != nil {
			return 0, err
		}

		blockHeight = uint64(block.Height())
	}

	//TODO: Handle errors
	return blockHeight, err
}

func (relayTx *RelayTransaction) Send(ctx context.Context) (err error) {
	relayTx.client.Logger().WithFields(log.Fields{"signer": relayTx.Transaction.SignerId}).Debug("prepare tx")
	_ctx, cancel := context.WithTimeout(ctx, defaultSendTxTimeout)
	defer cancel()

	relayTx.context = _ctx
	publicKey := types.NewPublicKeyFromED25519(*relayTx.wallet.Pkey)
	nonce, err := relayTx.client.GetNonce(publicKey, string(relayTx.Transaction.SignerId))
	if nonce == -1 || err != nil {
		return err
	}

	relayTx.Transaction.Nonce = int(nonce) + 1
	relayTx.Transaction.BlockHash, err = relayTx.client.GetLatestBlockHash()
	if err != nil {
		return err
	}

	payload, err := relayTx.Transaction.Payload(relayTx.wallet)
	if err != nil {
		return err
	}

	txId, err := relayTx.client.SendTransaction(payload)
	if err != nil {
		return err
	}

	relayTx.Transaction.Txid = *txId
	relayTx.client.Logger().WithFields(log.Fields{"tx": txId.Base58Encode()}).Debug("tx sent")

	return nil
}

func (s *Sender) Status(ctx context.Context) (*chain.BMCLinkStatus, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	status, err := s.client().GetBmcLinkStatus(s.destination, s.source)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (s *Sender) Balance(ctx context.Context) (balance, threshold *big.Int, err error) {
	balance, err = s.client().GetBalance(types.AccountId(s.wallet.Address()))
	t := big.Int(s.options.BalanceThreshold)

	return balance, &t, err
}
