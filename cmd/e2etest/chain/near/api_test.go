package near

import (
	"testing"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	RPC_URI      = "https://rpc.testnet.near.org"
	TokenGodKey  = ""
	TokenGodAddr = ""
	GodKey       = ""
	GodAddr      = ""
	DemoSrcKey   = ""
	DemoSrcAddr  = ""
	DemoDstAddr  = ""
	GodDstAddr   = ""
	NID          = "0x1.near"
	BtsOwner     = ""
)

func TestGetCoinNames(t *testing.T) {
	api, err := getNewApi()
	if err != nil {
		t.Fatal(err)
		return
	}
	i := api.NativeCoin()
	t.Log(i)
	// assert.Equal(t, 1, 0)
}

// func TestGetOwners(t *testing.T) {
// 	api, err := getNewApi()
// 	if err != nil {
// 		t.Fatalf("%+v", err)
// 		return
// 	}
// 	owner, err := api.CallBTS("get_owners", nil)
// 	if err != nil {
// 		t.Fatalf("%+v", err)
// 		return
// 	}
// 	if data, ok := (owner).(types.CallFunctionResponse); ok {
// 		var r []string
// 		err = json.Unmarshal(data.Result, &r)
// 		fmt.Println(data.BlockHash)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println(r)

// 		// assert.Equal(t, 1, 0)
// 	}

// }

// func TestIsUserBlackListed(t *testing.T) {
// 	rpi, err := getNewApi()
// 	if err != nil {
// 		t.Fatalf("%+v", err)
// 	}
// 	res, err := rpi.CallBTS(chain.IsUserBlackListed, []interface{}{
// 		"0x61.bsc",
// 		GodDstAddr,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println("Res ", res)
// }

func getNewApi() (chain.ChainAPI, error) {
	srcEndpoint := RPC_URI
	addrToName := map[chain.ContractName]string{
		chain.BTS: "bts.iconbridge.testnet",
	}
	l := log.New()
	log.SetGlobalLogger(l)
	return NewApi(l, &chain.Config{
		Name:              chain.NEAR,
		URL:               srcEndpoint,
		ContractAddresses: addrToName,
		NativeTokens:      []string{},
		WrappedCoins:      []string{"btp-0x2.icon-ICX"},
		NativeCoin:        "btp-0x1.near-NEAR",
		NetworkID:         NID,
		// GasLimit:          300000000000000,
	})
}
