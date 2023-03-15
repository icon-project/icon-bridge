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
	algodAddress     = "http://localhost:4001"
	kmdAddress       = "http://localhost:4002"
	algoToken        = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testAddress      = "5SOSEYHUIAGIFPOHCV6ANF7KYSAU5AXHF5TNM4BQMIVU4DNWZOQJ427XRE"
	testAccountSk    = "5pw5+iRVb91t/RfIHF4c1RDMiIQ434PPVDN+Vv2qv9LsnSJg9EAMgr3HFXwGl+rEgU6C5y9m1nAwYitODbbLoA=="
	algoBmc          = "btp://0x14.algo/0x293b2D1B12393c70fCFcA0D9cb99889fFD4A23a8"
	iconBmc          = "btp://0x2.icon/cx04d4cc5ee639aa2fc5f2ededa7b50df6044dd325"
	bmcCompilePyPath = "../../../../pyteal/bmc/builder.py"
)

var (
	testnetAddress = os.Getenv("ALGO_TEST_ADR")
	testnetToken   = os.Getenv("ALGO_TEST_TOK")
	testnetAccess  = []string{testnetAddress, testnetToken}
	sandboxAccess  = []string{algodAddress, algoToken}
)

func createTestReceiver(algodAccess []string, round uint64, hash string) (chain.Receiver, error) {
	opts := map[string]interface{}{"syncConcurrency": 2,
		"Verifier": VerifierOptions{
			round,
			hash,
		},
	}
	rawOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling options: %v", err)
	}

	rcv, err := NewReceiver(chain.BTPAddress(algoBmc), chain.BTPAddress(iconBmc),
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
	fmt.Println(account.Address)

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
	w, err := wallet.NewAvmWalletFromPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create wallet: %v", err)
	}

	s, err := NewSender(
		chain.BTPAddress(algoBmc), chain.BTPAddress(iconBmc),
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
