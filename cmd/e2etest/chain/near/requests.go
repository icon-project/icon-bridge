package near

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain/near/types"
	"github.com/icon-project/icon-bridge/common/errors"
)

type requestAPI struct {
	contractNameToAddress map[chain.ContractName]string
	networkID             string
	cl                    *near.Client
	stepLimit             int64
	nativeCoin            string
	wrappedCoinsAddr      map[string]string
	nativeTokensAddr      map[string]string
}

type coinNames struct {
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Network string `json:"network"`
}

func newRequestAPI(cl *near.Client, cfg *chain.Config) (req *requestAPI, err error) {
	if !strings.Contains(cfg.NetworkID, ".near") {
		return nil, fmt.Errorf("Expected cfg.NetwrkID=0xnid.near Got %v", cfg.NetworkID)
	}
	req = &requestAPI{
		networkID:             strings.Split(cfg.NetworkID, ".")[0],
		contractNameToAddress: cfg.ContractAddresses,
		cl:                    cl,
		stepLimit:             cfg.GasLimit,
		nativeCoin:            cfg.NativeCoin,
	}
	req.nativeTokensAddr, req.wrappedCoinsAddr, err = req.getCoinAddresses(cfg.NativeTokens, cfg.WrappedCoins)
	return req, err
}

func (r *requestAPI) getCoinAddresses(nativeTokens, wrappedCoins []string) (tokenAddrMap, wrappedAddrMap map[string]string, err error) {
	btsaddr, ok := r.contractNameToAddress[chain.BTS]
	var coin_names []coinNames
	if !ok {
		err = fmt.Errorf("contractNameToAddress doesn't include name %v", chain.BTS)
		return
	}
	res, err := r.callContract(btsaddr, map[string]interface{}{}, "coins")
	if err != nil {
		err = errors.Wrap(err, "callContract coinNames ")
		return
	} else if res == nil {
		err = fmt.Errorf("Call to Method %v returned nil", "coinNames")
		return
	}
	resArr := res.(near.CallFunctionResult).Result
	err = json.Unmarshal(resArr, &coin_names)
	println(coin_names)
	if err != nil {
		err = fmt.Errorf("For method coinNames, Expected Type []interface{} Got %T", err)
		return
	}
	coinNames := []string{}
	for _, re := range coin_names {
		c := re.Name
		if c == r.nativeCoin {
			continue
		}
		coinNames = append(coinNames, c)
	}
	exists := func(arr []string, val string) bool {
		for _, a := range arr {
			if a == val {
				return true
			}
		}
		return false
	}

	// all registered coins have to be given in input config
	allInputCoins := append(nativeTokens, wrappedCoins...)
	for _, coinName := range coinNames {
		if !exists(allInputCoins, coinName) {
			err = fmt.Errorf("Registered coin %v not provided in input config ", coinName)
			return
		}
	}
	// all coins given in input config have to have been registered
	for _, inputCoin := range allInputCoins {
		if !exists(coinNames, inputCoin) {
			err = fmt.Errorf("Input coin %v does not exist among registered coins ", inputCoin)
			return
		}
	}
	getAddr := func(coin string) (coinId string, err error) {
		var res interface{}

		res, err = r.callContract(btsaddr, map[string]interface{}{"coin_name": coin}, "coin_id")
		if err != nil {
			err = errors.Wrap(err, "callContract coinId ")
			return
		} else if res == nil {
			err = fmt.Errorf("Call to Method %v returned nil for _coinName=%v", "coinId", coin)
			return
		}
		resArr := res.(near.CallFunctionResult).Result
		var coin_id []byte
		err = json.Unmarshal(resArr, &coin_id)
		if err != nil {
			err = errors.Wrap(err, "callContract coinId ")
			return
		}
		coinId = base64.StdEncoding.EncodeToString(coin_id)
		if !ok {
			err = fmt.Errorf("For method coinId, Expected Type string Got %T", res)
			return
		}
		return coinId, nil
	}

	tokenAddrMap = map[string]string{}
	for _, coin := range nativeTokens {
		tokenAddrMap[coin], err = getAddr(coin)
		if err != nil {
			return
		}
	}
	wrappedAddrMap = map[string]string{}
	for _, coin := range wrappedCoins {
		wrappedAddrMap[coin], err = getAddr(coin)
		if err != nil {
			return
		}
	}
	return
}

func (r *requestAPI) callContract(contractAddress string, args map[string]interface{}, method string) (interface{}, error) {
	methodParam, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	param := &types.CallFunction{
		RequestType:  "call_function",
		Finality:     "final",
		AccountId:    types.AccountId(contractAddress),
		MethodName:   method,
		ArgumentsB64: base64.URLEncoding.EncodeToString(methodParam),
	}

	var res near.CallFunctionResult
	_, err = r.cl.Call("query", param, &res)
	if err != nil {
		return nil, errors.Wrap(err, "Call ")
	}
	return res, nil
}
