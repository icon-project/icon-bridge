package algo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/bmizerany/assert"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const sandboxAddress = "http://localhost:4001"
const sandboxToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const (
	KMD_ADDRESS         = "http://localhost:4002"
	KMD_TOKEN           = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	KMD_WALLET_NAME     = "unencrypted-default-wallet"
	KMD_WALLET_PASSWORD = ""
)
const approvalPath = "bmc/approval.teal"
const clearPath = "bmc/clear.teal"

func getAccounts() ([]crypto.Account, error) {
	client, err := kmd.MakeClient(KMD_ADDRESS, KMD_TOKEN)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %+v", err)
	}

	resp, err := client.ListWallets()
	if err != nil {
		return nil, fmt.Errorf("Failed to list wallets: %+v", err)
	}

	var walletId string
	for _, wallet := range resp.Wallets {
		if wallet.Name == KMD_WALLET_NAME {
			walletId = wallet.ID
		}
	}

	if walletId == "" {
		return nil, fmt.Errorf("No wallet named %s", KMD_WALLET_NAME)
	}

	whResp, err := client.InitWalletHandle(walletId, KMD_WALLET_PASSWORD)
	if err != nil {
		return nil, fmt.Errorf("Failed to init wallet handle: %+v", err)
	}

	addrResp, err := client.ListKeys(whResp.WalletHandleToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to list keys: %+v", err)
	}

	var accts []crypto.Account
	for _, addr := range addrResp.Addresses {
		expResp, err := client.ExportKey(whResp.WalletHandleToken, KMD_WALLET_PASSWORD, addr)
		if err != nil {
			return nil, fmt.Errorf("Failed to export key: %+v", err)
		}

		acct, err := crypto.AccountFromPrivateKey(expResp.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("Failed to create account from private key: %+v", err)
		}

		accts = append(accts, acct)
	}

	return accts, nil
}

func createTestSender() (chain.Sender, error) {
	accts, err := getAccounts()
	if err != nil {
		return nil, fmt.Errorf("Error generating KMD account: %v", err)
	}
	account := accts[0]

	algodAccess := []string{sandboxAddress, sandboxToken}
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	appId, err := deployContract(ctx, algodAccess, [2]string{approvalPath, clearPath}, account)
	if err != nil {
		return nil, fmt.Errorf("Error deploying BMC: %v", err)
	}
	opts := map[string]interface{}{"app_id": appId}
	rawOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("Marshalling opts: %v", err)
	}
	w, err := wallet.NewAvmWalletFromPrivateKey(&account.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create wallet: %v", err)
	}

	s, err := NewSender(
		chain.BTPAddress(icon_bmc), chain.BTPAddress(algo_bmc),
		algodAccess, w,
		rawOpts, log.New())
	if err != nil {
		return nil, fmt.Errorf("Error creating new sender: %v", err)
	}
	return s, nil
}

func Test_Abi(t *testing.T) {
	s, err := createTestSender()
	if err != nil {
		t.Logf("Failed creting new sender:%v", err)
		t.FailNow()
	}
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	ret, err := s.(*sender).callAbi(ctx, "concat_strings",
		[]interface{}{[]string{"this", "string", "is", "joined"}})

	if err != nil {
		t.Logf("Failed calling abi:%v", err)
		t.FailNow()
	}
	fmt.Println(ret)
	concatString := ret.MethodResults[0].ReturnValue.(string)
	assert.Equal(t, concatString, "thisstringisjoined")
}

func Test_Segment(t *testing.T) {
	s, err := createTestSender()
	if err != nil {
		t.Logf("Failed creting new sender:%v", err)
		t.FailNow()
	}
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	msg := &chain.Message{
		From: chain.BTPAddress(icon_bmc),
		Receipts: []*chain.Receipt{{
			Index:  0,
			Height: 1,
			Events: []*chain.Event{
				{
					Message: []byte{97, 98, 99, 100, 101, 102},
				},
				{
					Message: []byte{44, 32, 33, 4, 101, 255},
				},
			},
		},
			{
				Index:  0,
				Height: 2,
				Events: []*chain.Event{
					{
						Message: []byte{55, 56, 222, 34, 6, 3},
					},
					{
						Message: []byte{64, 2, 4, 111, 55, 23},
					},
				},
			}},
	}
	tx, newmsg, err := s.Segment(ctx, msg)

	if err != nil {
		t.Logf("Couldn't segment message:%v", err)
		t.FailNow()
	}
	fmt.Println(tx)
	fmt.Println("......................")
	fmt.Println(newmsg)

	sss := tx.(*relayTx).msg

	recovered_pay_bytes := make([]byte, 1000000)
	base64.StdEncoding.Decode(recovered_pay_bytes, sss)
	rm := &chain.RelayMessage{}
	msgpack.Decode(recovered_pay_bytes, &rm)

	recSli := make([]chain.Receipt, 0)
	for _, r := range rm.Receipts {
		decodedReceipt := make([]byte, 1000)
		base64.StdEncoding.Decode(decodedReceipt, r)

		var finalRcp chain.Receipt
		msgpack.Decode(decodedReceipt, &finalRcp)

		recSli = append(recSli, finalRcp)
	}
	fmt.Println("......................")

}
