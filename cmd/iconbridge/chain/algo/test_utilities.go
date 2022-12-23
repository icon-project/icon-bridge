package algo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/algorand/go-algorand-sdk/client/kmd"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const (
	sandboxAddress      = "http://localhost:4001"
	sandboxToken        = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	KMD_ADDRESS         = "http://localhost:4002"
	KMD_TOKEN           = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	KMD_WALLET_NAME     = "unencrypted-default-wallet"
	KMD_WALLET_PASSWORD = ""
	algo_bmc            = "btp://0x14.algo/0x293b2D1B12393c70fCFcA0D9cb99889fFD4A23a8"
	icon_bmc            = "btp://0x1.icon/cx06f42ea934731b4867fca00d37c25aa30bc3e3d7"
)

var (
	testnetAddress = os.Getenv("ALGO_TEST_ADR")
	testnetToken   = os.Getenv("ALGO_TEST_TOK")
	testnetAccess  = []string{testnetAddress, testnetToken}
	sandboxAccess  = []string{sandboxAddress, sandboxToken}
)

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

func createTestReceiver(algodAccess []string) (chain.Receiver, error) {
	opts := map[string]interface{}{"syncConcurrency": 2}
	rawOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling options: %v", err)
	}

	rcv, err := NewReceiver(chain.BTPAddress(icon_bmc), chain.BTPAddress(algo_bmc),
		algodAccess, rawOpts, log.New())
	if err != nil {
		return nil, fmt.Errorf("Error creating new receiver: %v", err)
	}
	return rcv, nil
}

func createTestSender(algodAccess []string) (chain.Sender, error) {
	accts, err := getAccounts()
	if err != nil {
		return nil, fmt.Errorf("Error generating KMD account: %v", err)
	}
	account := accts[0]

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)

	appId, err := deployContract(ctx, algodAccess, [2]string{"approval.teal", "clear.teal"}, account)
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
