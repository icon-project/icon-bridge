package tezos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	// "io"
	"time"

	"github.com/icon-project/icon-bridge/common/log"

	"blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/micheline"
	"blockwatch.cc/tzgo/rpc"
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
	GetBalance(ctx context.Context, connection *rpc.Client, account tezos.Address, blockLevel int64)
	GetBlockByHeight(ctx context.Context, connection *rpc.Client, blockLevel int64) (*rpc.Block, error)
	GetBlockHeightByHash(ctx context.Context, connection *rpc.Client, hash tezos.BlockHash) (int64, error)
	MonitorBlock(ctx context.Context, client *rpc.Client, connection *contract.Contract, blockLevel int64, callback func(v *types.BlockNotification) error) (*rpc.Block, error)
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
	Cl            *rpc.Client
	Contract      *contract.Contract
	blockLevel    int64
	BmcManagement tezos.Address
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

func filterMessageEvents(tx *rpc.Transaction, contractAddress tezos.Address, height uint64, dst string) (*chain.Receipt, error) {
	receipt := &chain.Receipt{}
	var events []*chain.Event

	for i := 0; i < len(tx.Metadata.InternalResults); i++ {
		internalResults := tx.Metadata.InternalResults[i]
		if internalResults.Kind.String() == "event" && internalResults.Source.ContractAddress() == contractAddress.ContractAddress() {
			if internalResults.Tag == "Message" {
				message := internalResults.Payload.Args[0].Bytes
				next := internalResults.Payload.Args[1].Args[0].String
				seq := internalResults.Payload.Args[1].Args[1].Int

				if next == dst {
					fmt.Println("found it")
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

func (c *Client) GetBMCManangement(ctx context.Context, contr *contract.Contract, account tezos.Address) (string, error) {
	fmt.Println("reached in getting bmc Management")
	result, err := contr.RunView(ctx, "get_bmc_periphery", micheline.Prim{})
	if err != nil {
		return "", err
	}
	return result.String, nil
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

func (c *Client) GetConsensusKey(ctx context.Context, bakerConsensusKey tezos.Address) (tezos.Key, error) {
	fmt.Println("baker consensus key", bakerConsensusKey.String())
	var exposedPublicKey tezos.Key
	for i := 0; i < 5; i++ {
		url := c.Cl.BaseURL.String() + "/chains/main/blocks/head/context/raw/json/contracts/index/" + bakerConsensusKey.String() + "/consensus_key/active"

		resp, err := http.Get(url)
		if err != nil {
			return tezos.Key{}, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return tezos.Key{}, err
		}
		//Convert the body to type string
		sb := string(body)

		exposedPublicKey, err = tezos.ParseKey(sb[1 : len(sb)-2])
		if err != nil {
			fmt.Println("continued to refetch again")
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	return exposedPublicKey, nil
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

func (c *Client) CustomCall(ctx context.Context, args []contract.CallArguments, opts *rpc.CallOptions) (*rpc.Receipt, error) {
	if opts == nil {
		opts = &rpc.DefaultOptions
	}

	// assemble batch transaction
	op := codec.NewOp().WithTTL(opts.TTL)
	for _, arg := range args {
		if arg == nil {
			continue
		}
		op.WithContents(arg.Encode())
	}

	var limits []tezos.Limits
	limit := tezos.Limits{
		GasLimit:     tezos.MumbainetParams.HardGasLimitPerOperation,
		StorageLimit: tezos.MumbainetParams.HardStorageLimitPerOperation,
	}

	limits = append(limits, limit)

	op.WithLimits(limits, 0).WithMinFee()

	// prepare, sign and broadcast
	return c.Cl.Send(ctx, op, opts)
}

func NewClient(uri string, src tezos.Address, bmcManagement tezos.Address, l log.Logger) (*Client, error) {

	fmt.Println("uri is : " + uri)

	c, err := rpc.NewClient(uri, nil)

	conn := contract.NewContract(src, c)

	if err != nil {
		return nil, err
	}

	return &Client{Log: l, Cl: c, Contract: conn, BmcManagement: bmcManagement}, nil
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

func filterTransactionOperations(block *rpc.Block, contractAddress tezos.Address, blockHeight int64, cl *Client, dst string) (bool, []*chain.Receipt, error) {
	blockOperations := block.Operations
	var tx *rpc.Transaction
	var receipt []*chain.Receipt
	for i := 0; i < len(blockOperations); i++ {
		for j := 0; j < len(blockOperations[i]); j++ {
			for _, operation := range blockOperations[i][j].Contents {
				switch operation.Kind() {
				case tezos.OpTypeTransaction:
					tx = operation.(*rpc.Transaction)
					r, err := filterMessageEvents(tx, cl.BmcManagement, uint64(blockHeight), dst)
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
