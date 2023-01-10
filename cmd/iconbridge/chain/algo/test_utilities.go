package algo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"github.com/icon-project/icon-bridge/common/log"
	"github.com/icon-project/icon-bridge/common/wallet"
)

const (
	sandboxAddress = "http://localhost:4001"
	sandboxToken   = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testAccountSk  = "Vfscmda+xG0c9OQGVAgTd6mny016riYjml/RW5AkhkKP12Aujgpm7kCRGgYaColPK8PRRCibwWuJldYelaHo+Q=="
	algoBmc        = "btp://0x14.algo/0x293b2D1B12393c70fCFcA0D9cb99889fFD4A23a8"
	iconBmc        = "btp://0x1.icon/cx06f42ea934731b4867fca00d37c25aa30bc3e3d7"
)

var (
	testnetAddress = os.Getenv("ALGO_TEST_ADR")
	testnetToken   = os.Getenv("ALGO_TEST_TOK")
	testnetAccess  = []string{testnetAddress, testnetToken}
	sandboxAccess  = []string{sandboxAddress, sandboxToken}
)

func createTestReceiver(algodAccess []string, round uint64, hash [32]byte) (chain.Receiver, error) {
	opts := map[string]interface{}{"syncConcurrency": 2,
		"Verifier": Verifier{
			round,
			hash,
		},
	}
	rawOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling options: %v", err)
	}

	rcv, err := NewReceiver(chain.BTPAddress(iconBmc), chain.BTPAddress(algoBmc),
		algodAccess, rawOpts, log.New())
	if err != nil {
		return nil, fmt.Errorf("Error creating new receiver: %v", err)
	}
	return rcv, nil
}

func createTestSender(algodAccess []string) (chain.Sender, error) {
	privateKey, err := base64.StdEncoding.DecodeString(testAccountSk)
	if err != nil {
		return nil, fmt.Errorf("Error decoding private key: %s", err)
	}
	account, err := crypto.AccountFromPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("Can't get account from private key: %s", err)
	}

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
		chain.BTPAddress(iconBmc), chain.BTPAddress(algoBmc),
		algodAccess, w,
		rawOpts, log.New())
	if err != nil {
		return nil, fmt.Errorf("Error creating new sender: %v", err)
	}
	return s, nil
}

func genAlgoAccount() {
	acc := crypto.GenerateAccount()
	log.Printf("Private key: %s\n", base64.StdEncoding.EncodeToString(acc.PrivateKey))
	log.Printf("Address:     %s\n", acc.Address)
}
