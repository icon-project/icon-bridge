package tezos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	// "io"
	"log"
	"time"

	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/micheline"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/signer"
	"blockwatch.cc/tzgo/tezos"
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
	// GetBlockHeaderByHeight(ctx context.Context, connection *rpc.Client, blockLevel int64)
	// GetBlockMetadataByHash(ctx context.Context, connection *rpc.Client, blockHash tezos.Hash)

	MonitorBlock(ctx context.Context, client *rpc.Client, connection *contract.Contract, blockLevel int64) (*rpc.Block, error)
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
	Ctx context.Context
	Cl *rpc.Client
	Contract *contract.Contract
	blockLevel int64
}

func (c *Client) SignTransaction() rpc.CallOptions{
	pK := tezos.MustParsePrivateKey("")
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

func (c *Client) GetBlockHeaderByHeight(ctx context.Context, connection *rpc.Client, blockLevel int64)(*rpc.BlockHeader, error){
	block, err := connection.GetBlockHeader(ctx, rpc.BlockLevel(blockLevel))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (c *Client) MonitorBlock(blockLevel int64, verifier IVerifier) (error) {
	fmt.Println("reached in monitor block")
	relayTicker := time.NewTicker(DefaultBlockWaitInterval)
	defer relayTicker.Stop()

	for {
		select {
		case <- c.Ctx.Done():
			return fmt.Errorf("Context done")
		case <- relayTicker.C:
			fmt.Println("*************************************************************")
			fmt.Print("Trying to fetch block for blockLevel ")
			fmt.Println(blockLevel)

			block, err := c.GetBlockByHeight(c.Ctx, c.Cl, blockLevel)
			
			if err != nil {
				fmt.Println(err)
				fmt.Println("reducing the block level")
				blockLevel--
				fmt.Print("Trying to Fetch for block level ")
				fmt.Println(blockLevel)
				continue
			}
			
			header, err := c.GetBlockHeaderByHeight(c.Ctx, c.Cl, blockLevel)
			if err != nil {
				return err
			}

			err = verifier.Verify(header, &block.Metadata.Baker)

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
							returnTxMetadata(tx, c.Contract.Address())
						}
					}
				}
			}
		}
		blockLevel++
	}
}

func returnTxMetadata(tx *rpc.Transaction, contractAddress tezos.Address) error {
	_, err := fmt.Println(tx.Destination)
	if err != nil {
		return err
	}
	address := tx.Destination
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

		if tx.Metadata.InternalResults[0].Tag == "TokenMinted" {
			fmt.Println("Payload is")
			fmt.Println(tx.Metadata.InternalResults[0].Payload.Int)
		}
	}
	return nil
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

func (c *Client) HandleRelayMessage(ctx context.Context, callArgs contract.CallArguments) (*rpc.Receipt, error) {
	return nil, nil 
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
