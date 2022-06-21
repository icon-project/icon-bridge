package hmny

import (
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	bshcore "github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny/abi/bsh/bshcore"
	erc20 "github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny/abi/bsh/erc20tradable"
	bep20tkn "github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny/abi/tokenbsh/bep20tkn"
	bshproxy "github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny/abi/tokenbsh/bshproxy"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	coinName        = "ICX"
	DefaultGasLimit = 80000000
)

// grouped rpc api clients
type client struct {
	ethCl           *ethclient.Client
	log             log.Logger
	bshc            *bshcore.Bshcore
	erc             *erc20.Erc20tradable
	bep             *bep20tkn.BEP
	tokbsh          *bshproxy.TokenBSH
	contractAddress *contractAddress
	networkID       string
}

type contractAddress struct {
	btp_hmny_erc20               string
	btp_hmny_nativecoin_bsh_core string
	btp_hmny_token_bsh_proxy     string
}

func (cAddr *contractAddress) FromMap(contractAddrsMap map[string]string) {
	cAddr.btp_hmny_erc20 = contractAddrsMap["btp_hmny_erc20"]
	cAddr.btp_hmny_nativecoin_bsh_core = contractAddrsMap["btp_hmny_nativecoin_bsh_core"]
	cAddr.btp_hmny_token_bsh_proxy = contractAddrsMap["btp_hmny_token_bsh_proxy"]
}

func New(url string, l log.Logger, contractAddrsMap map[string]string, networkID string) (chain.Client, error) {
	cAddr := &contractAddress{}
	cAddr.FromMap(contractAddrsMap)

	return newClient(url, l, cAddr, networkID)
}

func newClient(url string, l log.Logger, cAddress *contractAddress, networkID string) (*client, error) {

	clrpc, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	cleth := ethclient.NewClient(clrpc)

	bshc, err := bshcore.NewBshcore(common.HexToAddress(cAddress.btp_hmny_nativecoin_bsh_core), cleth)
	if err != nil {
		return nil, err
	}
	coinAddress, err := bshc.CoinId(&bind.CallOpts{Pending: false, Context: nil}, coinName)
	if err != nil {
		return nil, err
	}
	bep, err := bep20tkn.NewBEP(common.HexToAddress(cAddress.btp_hmny_erc20), cleth)
	if err != nil {
		return nil, err
	}
	erc, err := erc20.NewErc20tradable(coinAddress, cleth)
	if err != nil {
		return nil, err
	}
	tokbsh, err := bshproxy.NewTokenBSH(common.HexToAddress(cAddress.btp_hmny_token_bsh_proxy), cleth)
	if err != nil {
		return nil, err
	}
	c := &client{
		ethCl:           cleth,
		log:             l,
		bshc:            bshc,
		erc:             erc,
		bep:             bep,
		tokbsh:          tokbsh,
		contractAddress: cAddress,
		networkID:       networkID,
	}
	return c, nil
}

func (c *client) GetCoinBalance(addr string) (*big.Int, error) {
	return c.GetHmnyBalance(addr)
}

func (c *client) GetEthToken(addr string) (val *big.Int, err error) {
	return c.GetHmnyErc20Balance(addr)
}

func (c *client) GetWrappedCoin(addr string) (val *big.Int, err error) {
	return c.GetHmnyWrappedICX(addr)
}

func (c *client) TransferCoin(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	if txnHash, err = c.TransferHmnyOne(senderKey, amount, recepientAddress); err != nil && err.Error() == core.ErrReplaceUnderpriced.Error() {
		duration := time.Millisecond * 500 * time.Duration(rand.Intn(11)+1) // Delay of [500ms, 6 seconds]
		c.log.Warn("Retrying Hmny One Transaction after (ms) ", duration.Milliseconds())
		time.Sleep(duration)
		return c.TransferHmnyOne(senderKey, amount, recepientAddress)
	}
	return
}

func (c *client) TransferEthToken(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return c.TransferErc20(senderKey, amount, recepientAddress)
}

func (c *client) TransferCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return c.TransferOneToIcon(senderKey, recepientAddress, amount)
}

func (c *client) TransferWrappedCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return c.TransferWrappedICXFromHmnyToIcon(senderKey, amount, recepientAddress)
}

func (c *client) TransferEthTokenCrossChain(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, transferTxnHash string, err error) {
	return c.TransferERC20ToIcon(senderKey, amount, recepientAddress)
}
func (c *client) ApproveContractToAccessCrossCoin(ownerKey string, amount big.Int) (approveTxnHash string, allowanceAmount *big.Int, err error) {
	return c.ApproveHmnyNativeBSHCoreToAccessICX(ownerKey, amount)
}

func (c *client) GetAddressFromPrivKey(key string) (*string, error) {
	return getAddressFromPrivKey(key)
}

func (c *client) GetFullAddress(addr string) *string {
	fullAddr := "btp://" + c.networkID + ".hmny/" + addr
	return &fullAddr
}
