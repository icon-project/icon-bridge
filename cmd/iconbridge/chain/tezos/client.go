package tezos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	// "io"
	"time"

	"github.com/icon-project/icon-bridge/common/log"

	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/micheline"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/signer"
	"blockwatch.cc/tzgo/tezos"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/tezos/types"
)

const (
	DefaultSendTransactionRetryInterval        = 30 * time.Second
	DefaultGetTransactionResultPollingInterval = 15 * time.Second
	DefaultBlockWaitInterval                   = 15 * time.Second
)

type IClient interface {
	// Call(ctx context.Context, callArgs contract.CallArguments, opts *rpc.CallOptions)
	GetBalance(ctx context.Context, connection *rpc.Client, account tezos.Address, blockLevel int64)
	GetBlockByHeight(ctx context.Context, connection *rpc.Client, blockLevel int64) (*rpc.Block, error)
	GetBlockHeightByHash(ctx context.Context, connection *rpc.Client, hash tezos.BlockHash) (int64, error)
	// GetBlockHeaderByHeight(ctx context.Context, connection *rpc.Client, blockLevel int64)
	// GetBlockMetadataByHash(ctx context.Context, connection *rpc.Client, blockHash tezos.Hash)

	MonitorBlock(ctx context.Context, client *rpc.Client, connection *contract.Contract, blockLevel int64, callback func(v *types.BlockNotification) error) (*rpc.Block, error)
	// MonitorEvent(ctx context.Context, connection *rpc.Client, blockLevel int64)

	GetLastBlock(ctx context.Context, connection *rpc.Client) (*rpc.Block, error)
	GetStatus(ctx context.Context, contr *contract.Contract) (TypesLinkStats, error)

	HandleRelayMessage(ctx context.Context, callArgs contract.CallArguments) (*rpc.Receipt, error)
}

// tezos periphery
type TypesLinkStats struct {
	RxSeq         *big.Int
	TxSeq         *big.Int
	RxHeight      *big.Int
	CurrentHeight *big.Int
}

type Client struct {
	Log log.Logger
	// Ctx context.Context
	Cl         *rpc.Client
	Contract   *contract.Contract
	blockLevel int64
}

func (c *Client) SignTransaction() rpc.CallOptions {
	pK := tezos.MustParsePrivateKey("edskRz1HoD3cWkmWhCNS5LjBrJNWChGuKWB4HnVoN5UqVsUCpcNJR67ZxKs965u8RgRwptrtGc2ufYZoeECgB77RKm1gTbQ6eB")
	opts := rpc.DefaultOptions
	opts.Signer = signer.NewFromKey(pK)
	return opts
}

func (c *Client) SendTransaction(ctx context.Context, connection *contract.Contract, parameters micheline.Parameters, sender tezos.Address) (*rpc.Receipt, error) {
	args := contract.NewTxArgs()

	args.WithParameters(parameters)

	opts := c.SignTransaction()

	argument := args.WithSource(sender).WithDestination(connection.Address())

	result, err := connection.Call(ctx, argument, &opts)

	if err != nil {
		return nil, err
	}

	err = PrettyEncode(result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetLastBlock(ctx context.Context, connection *rpc.Client) (*rpc.Block, error) {
	block, err := connection.GetHeadBlock(ctx)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *Client) GetBlockByHeight(ctx context.Context, connection *rpc.Client, blockLevel int64) (*rpc.Block, error) {
	block, err := connection.GetBlock(ctx, rpc.BlockLevel(blockLevel))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *Client) GetBlockHeightByHash(ctx context.Context, connection *rpc.Client, hash tezos.BlockHash) (uint64, error) {
	block, err := connection.GetBlock(ctx, hash)
	if err != nil {
		return 0, err
	}
	return uint64(block.Header.Level), nil
}

func (c *Client) GetBlockHeaderByHeight(ctx context.Context, connection *rpc.Client, blockLevel int64) (*rpc.BlockHeader, error) {
	block, err := connection.GetBlockHeader(ctx, rpc.BlockLevel(blockLevel))
	if err != nil {
		return nil, err
	}
	return block, nil
}

// func (c *Client) MonitorBlock(ctx context.Context, blockLevel int64, verifier IVerifier, callback func(v []*chain.Receipt) error) error {
// 	fmt.Println("reached in monitor block")
// 	relayTicker := time.NewTicker(DefaultBlockWaitInterval)
// 	defer relayTicker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return fmt.Errorf("Context done")
// 		case <-relayTicker.C:
// 			fmt.Println("*************************************************************")
// 			fmt.Print("Trying to fetch block for blockLevel ")
// 			fmt.Println(blockLevel)

// 			block, err := c.GetBlockByHeight(ctx, c.Cl, blockLevel)

// 			if err != nil {
// 				fmt.Println(err)
// 				fmt.Println("reducing the block level")
// 				blockLevel--
// 				fmt.Print("Trying to Fetch for block level ")
// 				fmt.Println(blockLevel)
// 				continue
// 			}

// 			header, err := c.GetBlockHeaderByHeight(ctx, c.Cl, blockLevel)
// 			if err != nil {
// 				return err
// 			}
// 			fmt.Println(block.Metadata.ProposerConsensusKey)

// 			err = verifier.Verify(ctx, header, block.Metadata.ProposerConsensusKey, c.Cl, header)

// 			if err != nil {
// 				fmt.Println(err)
// 				return err
// 			}
// 			c.blockLevel = blockLevel

// 			// err = verifier.Update(header, )

// 			if err != nil {
// 				fmt.Println(err)
// 				return err
// 			}

// 			PrettyEncode(header)

// 			blockOperations := block.Operations

// 			for i := 0; i < len(blockOperations); i++ {
// 				for j := 0; j < len(blockOperations[i]); j++ {
// 					for _, operation := range blockOperations[i][j].Contents {
// 						switch operation.Kind() {
// 						case tezos.OpTypeTransaction:
// 							tx := operation.(*rpc.Transaction)
// 							receipt, err := returnTxMetadata(tx, c.Contract.Address())
// 							if err != nil {
// 								return err
// 							}
// 							if len(receipt) != 0 {
// 								fmt.Println("found for block level ", block.Header.Level)
// 								fmt.Println("callback start")
// 								err := callback(receipt)
// 								fmt.Println("call back end")
// 								if err != nil {
// 									return err
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 		blockLevel++
// 	}
// }

func returnTxMetadata(tx *rpc.Transaction, contractAddress tezos.Address) ([]*chain.Receipt, error) {
	// _, err := fmt.Println(tx.Destination)
	// if err != nil {
	// 	return nil, err
	// }
	address := tx.Destination

	var receipts []*chain.Receipt
	if address.ContractAddress() == contractAddress.ContractAddress() {
		fmt.Println("Address matched")
		fmt.Println("****************")
		fmt.Println("****************")
		fmt.Println("****")
		fmt.Println("****")
		fmt.Println("****")
		fmt.Println("****************")
		fmt.Println("****************")
		fmt.Println("****")
		fmt.Println("****")
		fmt.Println("****")
		fmt.Println("****")
		fmt.Println("****")
		fmt.Println("****")

		if tx.Metadata.InternalResults[0].Tag == "TransferStart" {
			var events []*chain.Event

			events = append(events, &chain.Event{
				Message: []byte(tx.Metadata.InternalResults[0].Payload.String),
			})
			receipts = append(receipts, &chain.Receipt{
				Events: events,
			})
		}
	}
	return receipts, nil
}

func returnTxMetadata3(tx *rpc.Transaction, contractAddress tezos.Address, height uint64) (*chain.Receipt, error) {
	fmt.Println("reache to return tx metadata3", height)
	receipt := &chain.Receipt{}

	for i := 0; i < len(tx.Metadata.InternalResults); i++ {
		fmt.Println("reached in for")
		internalResults := tx.Metadata.InternalResults[i]
		if internalResults.Kind.String() == "event" && internalResults.Source.ContractAddress() == "KT1XamvZ9WgAmxq4eBGKi6ZRbLGuGJpYcqaj" {
			fmt.Println("Address matched")
			if internalResults.Tag == "Message" {
				message := internalResults.Payload.Args[0].Bytes
				next := internalResults.Payload.Args[1].Args[0].String
				seq := internalResults.Payload.Args[1].Args[1].Int

				var events []*chain.Event
				events = append(events, &chain.Event{
					Message:  message,
					Next:     chain.BTPAddress(next),
					Sequence: seq.Uint64(),
				})

				receipt.Index = uint64(i)
				receipt.Height = height
				receipt.Events = events
				fmt.Println(message, next, seq)
			}

		}
	}
	return receipt, nil
}

func (c *Client) GetClient() *rpc.Client {
	return c.Cl
}

func (c *Client) GetBalance(ctx context.Context, connection *rpc.Client, account tezos.Address, blockLevel int64) (*big.Int, error) {
	balance, err := connection.GetContractBalance(ctx, account, rpc.BlockLevel(blockLevel))
	if err != nil {
		return nil, err
	}
	return balance.Big(), nil
}

func (c *Client) GetStatus(ctx context.Context, contr *contract.Contract, link string) (TypesLinkStats, error) {

	fmt.Println("reached in get status of tezos")
	prim := micheline.Prim{}

	in := "{ \"string\": \"" + link + "\" }"
	fmt.Println(in)
	fmt.Println(contr.Address().ContractAddress())

	if err := prim.UnmarshalJSON([]byte(in)); err != nil {
		fmt.Println("couldnot unmarshall empty string")
		fmt.Println(err)
		return *new(TypesLinkStats), err
	}

	result, err := contr.RunView(ctx, "get_status", prim)
	if err != nil {
		fmt.Println(err)
		return *new(TypesLinkStats), err
	}
	linkStats := &TypesLinkStats{}

	linkStats.CurrentHeight = result.Args[0].Args[0].Int
	linkStats.RxHeight = result.Args[0].Args[1].Int
	linkStats.RxSeq = result.Args[1].Int
	linkStats.TxSeq = result.Args[2].Int

	return *linkStats, nil
}

func (c *Client) GetOperationByHash(ctx context.Context, clinet *rpc.Client, blockHash tezos.BlockHash, list int, pos int) (*rpc.Operation, error) {
	operation, err := clinet.GetBlockOperation(ctx, blockHash, list, pos)
	if err != nil {
		return nil, err
	}
	return operation, nil
}

func (c *Client) HandleRelayMessage(ctx context.Context, callArgs contract.CallArguments, opts *rpc.CallOptions) (*rpc.Receipt, error) {
	fmt.Println("handling relay message")
	PrintU()
	result, err := c.Contract.Call(ctx, callArgs, opts)
	if err != nil {
		fmt.Println(err)
		fmt.Println("because error")
		return nil, err
	}
	fmt.Println(result)
	return result, nil
}

func NewClient(uri string, src tezos.Address, l log.Logger) (*Client, error) {

	fmt.Println("uri is : " + uri)

	c, err := rpc.NewClient(uri, nil)

	conn := contract.NewContract(src, c)

	if err != nil {
		return nil, err
	}

	return &Client{Log: l, Cl: c, Contract: conn}, nil
}

func PrettyEncode(data interface{}) error {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	fmt.Println(buffer.String())
	return nil
}

func returnTxMetadata2(block *rpc.Block, contractAddress tezos.Address, blockHeight int64, cl *Client) (bool, []*chain.Receipt, error) {
	blockOperations := block.Operations

	var tx *rpc.Transaction
	var receipt []*chain.Receipt
	for i := 0; i < len(blockOperations); i++ {
		for j := 0; j < len(blockOperations[i]); j++ {
			for _, operation := range blockOperations[i][j].Contents {
				switch operation.Kind() {
				case tezos.OpTypeTransaction:
					tx = operation.(*rpc.Transaction)
					r, err := returnTxMetadata3(tx, contractAddress, uint64(blockHeight))
					if err != nil {
						return false, nil, err
					}
					if len(r.Events) != 0 {
						fmt.Println("r is not nil for ", uint64(blockHeight))
						receipt = append(receipt, r)
					}
				}
			}
		}
	}
	// var transaction *rpc.Transaction

	if len(receipt) == 0 {
		return false, nil, nil
	}
	fmt.Println("found Message")
	return true, receipt, nil
}

func PrintU() {
	for i := 0; i < 100; i++ {
		fmt.Println("U")
	}
}
