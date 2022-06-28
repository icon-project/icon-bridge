package hmny

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	bshcore "github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny/abi/bsh/bshcore"
	erc20 "github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny/abi/bsh/erc20tradable"
	bep20tkn "github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny/abi/tokenbsh/bep20tkn"
	bshproxy "github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny/abi/tokenbsh/bshproxy"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/pkg/errors"
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

func NewRequestAPI(url string, l log.Logger, contractAddrsMap map[string]string, networkID string) (*requestAPI, error) {
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

func (r *requestAPI) GetCoinBalance(addr string, coinType chain.TokenType) (*big.Int, error) {
	if coinType == chain.ONEToken {
		return r.getHmnyBalance(addr)
	} else if coinType == chain.ERC20Token {
		return r.getHmnyErc20Balance(addr)
	} else if coinType == chain.ICXToken {
		return r.getHmnyWrappedICX(addr)
	}
	return nil, errors.New("Unsupported Token Type ")
}

func (r *requestAPI) Transfer(param *chain.RequestParam) (txnHash string, err error) {
	if param.FromChain != chain.HMNY {
		err = errors.New("Source Chan should be Hmny")
		return
	}
	if param.ToChain == chain.HMNY {
		if param.Token == chain.ONEToken {
			txnHash, _, err = r.transferHmnyOne(param.SenderKey, param.Amount, param.ToAddress)
		} else if param.Token == chain.ERC20Token {
			txnHash, _, err = r.transferErc20(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			err = errors.New("For intra chain transfer; unsupported token type ")
		}
	} else if param.ToChain == chain.ICON {
		if param.Token == chain.ONEToken {
			txnHash, _, err = r.transferOneToIcon(param.SenderKey, param.ToAddress, param.Amount)
		} else if param.Token == chain.ERC20Token {
			_, _, txnHash, _, err = r.transferERC20ToIcon(param.SenderKey, param.Amount, param.ToAddress)
		} else if param.Token == chain.ICXToken {
			txnHash, _, err = r.transferWrappedICXFromHmnyToIcon(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			err = errors.New("For intra chain transfer; unsupported token type ")
		}
	} else {
		err = errors.New("Unsupport Transaction Parameters ")
	}
	return
}

func (r *requestAPI) Approve(ownerKey string, amount big.Int) (approveTxnHash string, logs interface{}, allowanceAmount *big.Int, err error) {
	return r.approveHmnyNativeBSHCoreToAccessICX(ownerKey, amount)
}

func (r *requestAPI) WaitForTxnResult(hash string) (txr interface{}, err error) {
	txr, err = r.waitForResults(context.TODO(), common.HexToHash(hash))
	return
}

func (r *requestAPI) GetBTPAddress(addr string) *string {
	fullAddr := "btp://" + r.networkID + ".hmny/" + addr
	return &fullAddr
}

func (r *requestAPI) GetKeyPairs(num int) ([][2]string, error) {
	var err error
	res := make([][2]string, num)
	for i := 0; i < num; i++ {
		res[i], err = generateKeyPair()
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
