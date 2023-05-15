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
	DefaultSendTransactionRetryInterval = 30 * time.Second
	DefaultGetTransactionResultPollingInterval = 15 * time.Second
	DefaultBlockWaitInterval = 15 * time.Second
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
	Cl *rpc.Client
	Contract *contract.Contract
	blockLevel int64
}

func (c *Client) SignTransaction() rpc.CallOptions{
	pK := tezos.MustParsePrivateKey("edskRz1HoD3cWkmWhCNS5LjBrJNWChGuKWB4HnVoN5UqVsUCpcNJR67ZxKs965u8RgRwptrtGc2ufYZoeECgB77RKm1gTbQ6eB")
	opts := rpc.DefaultOptions
	opts.Signer = signer.NewFromKey(pK)
	return opts
}

func (c *Client) SendTransaction(ctx context.Context, connection *contract.Contract, parameters micheline.Parameters, sender tezos.Address) (*rpc.Receipt, error){
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

func (c *Client) GetLastBlock(ctx context.Context, connection *rpc.Client) (*rpc.Block, error){
	block, err := connection.GetHeadBlock(ctx)
	if err != nil{
		return nil, err
	}
	return block, nil
}

func (c *Client) GetBlockByHeight(ctx context.Context, connection *rpc.Client, blockLevel int64)(*rpc.Block, error){
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

func (c *Client) GetBlockHeaderByHeight(ctx context.Context, connection *rpc.Client, blockLevel int64)(*rpc.BlockHeader, error){
	block, err := connection.GetBlockHeader(ctx, rpc.BlockLevel(blockLevel))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *Client) MonitorBlock(ctx context.Context, blockLevel int64, verifier IVerifier, callback func(v []*chain.Receipt) error) (error) {
	fmt.Println("reached in monitor block")
	relayTicker := time.NewTicker(DefaultBlockWaitInterval)
	defer relayTicker.Stop()

	for {
		select {
		case <- ctx.Done():
			return fmt.Errorf("Context done")
		case <- relayTicker.C:
			fmt.Println("*************************************************************")
			fmt.Print("Trying to fetch block for blockLevel ")
			fmt.Println(blockLevel)

			block, err := c.GetBlockByHeight(ctx, c.Cl, blockLevel)
			
			if err != nil {
				fmt.Println(err)
				fmt.Println("reducing the block level")
				blockLevel--
				fmt.Print("Trying to Fetch for block level ")
				fmt.Println(blockLevel)
				continue
			}
			
			header, err := c.GetBlockHeaderByHeight(ctx, c.Cl, blockLevel)
			if err != nil {
				return err
			}
			fmt.Println(block.Metadata.ProposerConsensusKey)

			err = verifier.Verify(ctx, header, block.Metadata.ProposerConsensusKey, c.Cl, header.Hash)

			if err != nil {
				fmt.Println(err)
				return err
			}
			c.blockLevel = blockLevel
			
			err = verifier.Update(header)

			if err != nil {
				fmt.Println(err)
				return err 
			}


			PrettyEncode(header)
			

			blockOperations := block.Operations
			
			for i := 0; i < len(blockOperations); i++ {
				for j := 0; j < len(blockOperations[i]); j ++{ 
					for _, operation := range blockOperations[i][j].Contents {
						switch operation.Kind() {
						case tezos.OpTypeTransaction:
							tx := operation.(*rpc.Transaction)
							receipt, err := returnTxMetadata(tx, c.Contract.Address())
							if err != nil {
								return err
							}
							if len(receipt) != 0 {
								fmt.Println("found for block level ", block.Header.Level)
								fmt.Println("callback start")
								err := callback(receipt)
								fmt.Println("call back end")
								if err != nil {
									return err
								}
							}
						}
					}
				}
			}
		}
		blockLevel++
	}
}

func returnTxMetadata(tx *rpc.Transaction, contractAddress tezos.Address) ([]*chain.Receipt, error) {
	_, err := fmt.Println(tx.Destination)
	if err != nil {
		return nil, err
	}
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

func (c *Client) GetClient()(*rpc.Client) {
	return c.Cl
}

func (c *Client) GetBalance(ctx context.Context, connection *rpc.Client, account tezos.Address, blockLevel int64)(*big.Int, error){
	balance, err := connection.GetContractBalance(ctx, account, rpc.BlockLevel(blockLevel))
	if err != nil {
		return nil, err
	}
	return balance.Big(), nil 
}

func (c *Client) GetStatus(ctx context.Context, contr *contract.Contract) (TypesLinkStats, error){
	prim := micheline.Prim{}
	status, err := contr.RunCallback(ctx, "getStatus", prim)
	if err != nil {
		return *new(TypesLinkStats), err
	}
	var stats TypesLinkStats 
	err = status.Decode(stats)
	if err != nil {
		return *new(TypesLinkStats), err 
	}
	return stats, nil
}

func (c *Client) GetOperationByHash(ctx context.Context, clinet *rpc.Client, blockHash tezos.BlockHash, list int, pos int) (*rpc.Operation, error){
	operation, err := clinet.GetBlockOperation(ctx, blockHash, list, pos)
	if err != nil {
		return nil, err
	}
	return operation, nil 
}

func (c *Client) HandleRelayMessage(ctx context.Context, callArgs contract.CallArguments, opts *rpc.CallOptions) (*rpc.Receipt, error) {
	fmt.Println("handling relay message")
	result, err := c.Contract.Call(ctx, callArgs, opts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

func NewClient(uri string, src tezos.Address, l log.Logger) (*Client, error){

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
		for j := 0; j < len(blockOperations[i]); j ++{ 
			for _, operation := range blockOperations[i][j].Contents {
				switch operation.Kind() {
				case tezos.OpTypeTransaction:
					tx = operation.(*rpc.Transaction)
					r, err := returnTxMetadata(tx, contractAddress)
					if err != nil {
						return false, nil, err
					}
					receipt = r
				}
			}
		}
	}
	// var transaction *rpc.Transaction
	
	if len(receipt) == 0 {
		return false, receipt, nil 
	} 

	return true, receipt, nil
}


