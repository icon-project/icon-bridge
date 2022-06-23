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

type api struct {
	contractAddress *contractAddress
	networkID       string
	ethCl           *ethclient.Client
	log             log.Logger
	bshc            *bshcore.Bshcore
	erc             *erc20.Erc20tradable
	bep             *bep20tkn.BEP
	tokbsh          *bshproxy.TokenBSH
}

func (cAddr *contractAddress) FromMap(contractAddrsMap map[string]string) {
	cAddr.btp_hmny_erc20 = contractAddrsMap["btp_hmny_erc20"]
	cAddr.btp_hmny_nativecoin_bsh_core = contractAddrsMap["btp_hmny_nativecoin_bsh_core"]
	cAddr.btp_hmny_token_bsh_proxy = contractAddrsMap["btp_hmny_token_bsh_proxy"]
}

func New(url string, l log.Logger, contractAddrsMap map[string]string, networkID string) (chain.API, error) {
	return newAPI(url, l, contractAddrsMap, networkID)
}

type contractAddress struct {
	btp_hmny_erc20               string
	btp_hmny_nativecoin_bsh_core string
	btp_hmny_token_bsh_proxy     string
}

func newAPI(url string, l log.Logger, contractAddrsMap map[string]string, networkID string) (*api, error) {
	cAddress := &contractAddress{}
	cAddress.FromMap(contractAddrsMap)

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
	a := &api{
		log:             l,
		contractAddress: cAddress,
		networkID:       networkID,
		ethCl:           cleth,
		bshc:            bshc,
		erc:             erc,
		bep:             bep,
		tokbsh:          tokbsh,
	}
	return a, nil
}

func (a *api) GetCoinBalance(addr string) (*big.Int, error) {
	return a.GetHmnyBalance(addr)
}

func (a *api) GetEthToken(addr string) (val *big.Int, err error) {
	return a.GetHmnyErc20Balance(addr)
}

func (a *api) GetWrappedCoin(addr string) (val *big.Int, err error) {
	return a.GetHmnyWrappedICX(addr)
}

func (a *api) TransferCoin(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	if txnHash, err = a.TransferHmnyOne(senderKey, amount, recepientAddress); err != nil && err.Error() == core.ErrReplaceUnderpriced.Error() {
		duration := time.Millisecond * 500 * time.Duration(rand.Intn(11)+1) // Delay of [500ms, 6 seconds]
		a.log.Warn("Retrying Hmny One Transaction after (ms) ", duration.Milliseconds())
		time.Sleep(duration)
		return a.TransferHmnyOne(senderKey, amount, recepientAddress)
	}
	return
}

func (a *api) TransferEthToken(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return a.TransferErc20(senderKey, amount, recepientAddress)
}

func (a *api) TransferCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return a.TransferOneToIcon(senderKey, recepientAddress, amount)
}

func (a *api) TransferWrappedCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return a.TransferWrappedICXFromHmnyToIcon(senderKey, amount, recepientAddress)
}

func (a *api) TransferEthTokenCrossChain(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, transferTxnHash string, err error) {
	return a.TransferERC20ToIcon(senderKey, amount, recepientAddress)
}
func (a *api) ApproveContractToAccessCrossCoin(ownerKey string, amount big.Int) (approveTxnHash string, allowanceAmount *big.Int, err error) {
	return a.ApproveHmnyNativeBSHCoreToAccessICX(ownerKey, amount)
}

func (a *api) GetAddressFromPrivKey(key string) (*string, error) {
	return getAddressFromPrivKey(key)
}

func (a *api) GetBTPAddress(addr string) *string {
	fullAddr := "btp://" + a.networkID + ".hmny/" + addr
	return &fullAddr
}
