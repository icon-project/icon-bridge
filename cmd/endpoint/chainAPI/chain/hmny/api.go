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
	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	bshcore "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/hmny/abi/bsh/bshcore"
	erc20 "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/hmny/abi/bsh/erc20tradable"
	bep20tkn "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/hmny/abi/tokenbsh/bep20tkn"
	bshproxy "github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain/hmny/abi/tokenbsh/bshproxy"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	coinName        = "ICX"
	DefaultGasLimit = 80000000
)

type requestAPI struct {
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

func NewRequestAPI(url string, l log.Logger, contractAddrsMap map[string]string, networkID string) (chain.RequestAPI, error) {
	return newRequestAPI(url, l, contractAddrsMap, networkID)
}

type contractAddress struct {
	btp_hmny_erc20               string
	btp_hmny_nativecoin_bsh_core string
	btp_hmny_token_bsh_proxy     string
}

func newRequestAPI(url string, l log.Logger, contractAddrsMap map[string]string, networkID string) (*requestAPI, error) {
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
	a := &requestAPI{
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

func (r *requestAPI) GetCoinBalance(addr string) (*big.Int, error) {
	return r.GetHmnyBalance(addr)
}

func (r *requestAPI) GetEthToken(addr string) (val *big.Int, err error) {
	return r.GetHmnyErc20Balance(addr)
}

func (r *requestAPI) GetWrappedCoin(addr string) (val *big.Int, err error) {
	return r.GetHmnyWrappedICX(addr)
}

func (r *requestAPI) TransferCoin(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	if txnHash, err = r.TransferHmnyOne(senderKey, amount, recepientAddress); err != nil && err.Error() == core.ErrReplaceUnderpriced.Error() {
		duration := time.Millisecond * 500 * time.Duration(rand.Intn(11)+1) // Delay of [500ms, 6 seconds]
		r.log.Warn("Retrying Hmny One Transaction after (ms) ", duration.Milliseconds())
		time.Sleep(duration)
		return r.TransferHmnyOne(senderKey, amount, recepientAddress)
	}
	return
}

func (r *requestAPI) TransferEthToken(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return r.TransferErc20(senderKey, amount, recepientAddress)
}

func (r *requestAPI) TransferCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return r.TransferOneToIcon(senderKey, recepientAddress, amount)
}

func (r *requestAPI) TransferWrappedCoinCrossChain(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	return r.TransferWrappedICXFromHmnyToIcon(senderKey, amount, recepientAddress)
}

func (r *requestAPI) TransferEthTokenCrossChain(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, transferTxnHash string, err error) {
	return r.TransferERC20ToIcon(senderKey, amount, recepientAddress)
}
func (r *requestAPI) ApproveContractToAccessCrossCoin(ownerKey string, amount big.Int) (approveTxnHash string, allowanceAmount *big.Int, err error) {
	return r.ApproveHmnyNativeBSHCoreToAccessICX(ownerKey, amount)
}

func (r *requestAPI) GetAddressFromPrivKey(key string) (*string, error) {
	return getAddressFromPrivKey(key)
}

func (r *requestAPI) GetBTPAddress(addr string) *string {
	fullAddr := "btp://" + r.networkID + ".hmny/" + addr
	return &fullAddr
}
