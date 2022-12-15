package algo

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const _algodAddress = "http://localhost:4001"
const _algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const (
	KMD_ADDRESS         = "http://localhost:4002"
	KMD_TOKEN           = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	KMD_WALLET_NAME     = "unencrypted-default-wallet"
	KMD_WALLET_PASSWORD = ""
)
const approvalPath = "bmc/approval.teal"
const clearPath = "bmc/clear.teal"

func GetAccounts() ([]crypto.Account, error) {
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

func Test_NewSender(t *testing.T) {
	accts, err := GetAccounts()
	if err != nil {
		t.Logf("Error generating KMD account: %v", err)
		t.FailNow()
	}
	account := accts[0]

	algodAccess := []string{_algodAddress, _algodToken}

	appId, err := deployContract(algodAccess, [2]string{approvalPath, clearPath}, account)
	if err != nil {
		t.Logf("Error deploying BMC: %v", err)
		t.FailNow()
	}
	opts := map[string]interface{}{"app_id": appId}
	rawOpts, err := json.Marshal(opts)
	if err != nil {
		t.Logf("Marshalling opts: %v", err)
		t.FailNow()
	}
	w, err := wallet.NewAvmWalletFromPrivateKey(&account.PrivateKey)
	if err != nil {
		t.Logf("Couldn't create wallet: %v", err)
		t.FailNow()
	}

	s, err := NewSender(
		chain.BTPAddress(icon_bmc), chain.BTPAddress(algo_bmc),
		algodAccess, w,
		rawOpts, log.New())
	if err != nil {
		t.Logf("Error creating new sender: %v", err)
		t.FailNow()
	}
	kk, err := s.(*sender).callAbi("concat_strings",
		[]interface{}{[]string{"this", "string", "is", "joined"}})
	if err != nil {
		t.Logf("Error using abi: %v", err)
		t.FailNow()
	}
	fmt.Print(kk)

}
