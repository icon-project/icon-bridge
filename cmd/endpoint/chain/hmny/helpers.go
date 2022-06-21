package hmny

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/icon-bridge/common/wallet"
)

func GetWalletFromPrivKey(privKey string) (wal *wallet.EvmWallet, pKey *ecdsa.PrivateKey, err error) {
	privBytes, err := hex.DecodeString(privKey)
	if err != nil {
		return
	}
	ethPrivKey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return
	}
	wal = &wallet.EvmWallet{
		Skey: ethPrivKey,
		Pkey: &ethPrivKey.PublicKey,
	}
	return wal, ethPrivKey, nil
}

func getAddressFromPrivKey(pKey string) (*string, error) {
	privBytes, err := hex.DecodeString(pKey)
	if err != nil {
		return nil, err
	}
	privKey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(privKey.PublicKey).String()
	return &addr, nil
}

func (c *client) GetHmnyBalance(addr string) (*big.Int, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return c.ethCl.BalanceAt(ctx, common.HexToAddress(addr), nil)
}

func (c *client) GetHmnyErc20Balance(addr string) (val *big.Int, err error) {
	//common.HexToAddress("0x8fc668275b4fa032342ea3039653d841f069a83b")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	opts := bind.CallOpts{Pending: false, Context: ctx}
	return c.bep.BalanceOf(&opts, common.HexToAddress(addr))
}

func (c *client) GetHmnyWrappedICX(addr string) (val *big.Int, err error) {
	coinName := "ICX"
	v, err := c.bshc.GetBalanceOf(&bind.CallOpts{Pending: false, Context: context.Background()}, common.HexToAddress(addr), coinName)
	if err != nil {
		return
	}
	return v.UsableBalance, nil
}

func (c *client) getTransactionRequest(senderKey string) (*bind.TransactOpts, error) {
	_, senderPrivKey, err := GetWalletFromPrivKey(senderKey)
	if err != nil {
		return nil, err
	}
	chainID, err := c.ethCl.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	txo, err := bind.NewKeyedTransactorWithChainID(senderPrivKey, chainID)
	if err != nil {
		return nil, err
	}
	txo.GasPrice, err = c.ethCl.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	txo.GasLimit = uint64(DefaultGasLimit)
	return txo, nil
}

func (c *client) waitForResults(ctx context.Context, txHash common.Hash) (txr *types.Receipt, err error) {
	const DefaultGetTransactionResultPollingInterval = 1500 * time.Millisecond //1.5sec
	ticker := time.NewTicker(time.Duration(DefaultGetTransactionResultPollingInterval) * time.Nanosecond)
	retryLimit := 10
	retryCounter := 0
	for {
		defer ticker.Stop()
		select {
		case <-ctx.Done():
			err = errors.New("Context Cancelled")
			return
		case <-ticker.C:
			if retryCounter >= retryLimit {
				err = errors.New("Retry Limit Exceeded while waiting for results of transaction")
				return
			}
			retryCounter++
			//c.log.Debugf("GetTransactionResult Attempt: %d", retryCounter)
			txr, err = c.ethCl.TransactionReceipt(context.TODO(), txHash)
			if err != nil && err == ethereum.NotFound {
				continue
			}
			//c.log.Debugf("GetTransactionResult hash:%v, txr:%+v, err:%+v", txHash, txr, err)
			return
		}
	}
}
func (c *client) TransferHmnyOne(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	senderWallet, senderPrivKey, err := GetWalletFromPrivKey(senderKey)
	if err != nil {
		err = errors.Wrap(err, "SenderKey decode ")
		return
	}
	nonce, err := c.ethCl.PendingNonceAt(context.Background(), common.HexToAddress(senderWallet.Address()))
	if err != nil {
		return
	}
	gasPrice, err := c.ethCl.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}
	chainID, err := c.ethCl.ChainID(context.Background())
	if err != nil {
		return
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(recepientAddress), &amount, uint64(DefaultGasLimit), gasPrice, []byte{})
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), senderPrivKey)
	if err != nil {
		return
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err = c.ethCl.SendTransaction(ctx, signedTx); err != nil {
		return
	}
	txnHash = signedTx.Hash().String()
	_, err = c.waitForResults(context.TODO(), signedTx.Hash())
	return
}

func (c *client) TransferErc20(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	txo, err := c.getTransactionRequest(senderKey)
	if err != nil {
		return
	}
	txo.Context = context.Background()
	txn, err := c.bep.Transfer(txo, common.HexToAddress(recepientAddress), &amount)
	if err != nil {
		return
	}
	txnHash = txn.Hash().String()
	_, err = c.waitForResults(context.TODO(), txn.Hash())
	return
}

func (c *client) TransferOneToIcon(senderKey string, recepientAddress string, amount big.Int) (txnHash string, err error) {
	txo, err := c.getTransactionRequest(senderKey)
	if err != nil {
		return
	}
	txo.Value = &amount
	txo.Context = context.Background()
	txn, err := c.bshc.TransferNativeCoin(txo, recepientAddress)
	if err != nil {
		return
	}
	txnHash = txn.Hash().String()
	_, err = c.waitForResults(context.TODO(), txn.Hash())
	return
}

func (c *client) TransferWrappedICXFromHmnyToIcon(senderKey string, amount big.Int, recepientAddress string) (txnHash string, err error) {
	txo, err := c.getTransactionRequest(senderKey)
	if err != nil {
		return
	}
	txo.Context = context.Background()
	txn, err := c.bshc.Transfer(txo, coinName, &amount, recepientAddress)
	if err != nil {
		return
	}
	txnHash = txn.Hash().String()
	_, err = c.waitForResults(context.TODO(), txn.Hash())
	return
}

func (c *client) TransferERC20ToIcon(senderKey string, amount big.Int, recepientAddress string) (approveTxnHash, transferTxnHash string, err error) {
	txo, err := c.getTransactionRequest(senderKey)
	if err != nil {
		return
	}
	txo.Context = context.Background()
	approveTxn, err := c.bep.Approve(txo, common.HexToAddress(c.contractAddress.btp_hmny_token_bsh_proxy), &amount)
	if err != nil {
		return
	}
	_, err = c.waitForResults(context.TODO(), approveTxn.Hash())
	if err != nil {
		return
	}
	transferTxn, err := c.tokbsh.Transfer(txo, "ETH", &amount, recepientAddress)
	if err != nil {
		return
	}
	approveTxnHash = approveTxn.Hash().String()
	transferTxnHash = transferTxn.Hash().String()
	_, err = c.waitForResults(context.TODO(), transferTxn.Hash())
	if err != nil {
		return
	}
	return
}

func (c *client) ApproveHmnyNativeBSHCoreToAccessICX(ownerKey string, amount big.Int) (approveTxnHash string, allowanceAmount *big.Int, err error) {
	txo, err := c.getTransactionRequest(ownerKey)
	if err != nil {
		return
	}
	txo.Context = context.Background()
	approveTxn, err := c.erc.Approve(txo, common.HexToAddress(c.contractAddress.btp_hmny_nativecoin_bsh_core), &amount)
	if err != nil {
		return
	}
	approveTxnHash = approveTxn.Hash().String()
	_, err = c.waitForResults(context.TODO(), approveTxn.Hash())
	if err != nil {
		return
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ownerWallet, _, err := GetWalletFromPrivKey(ownerKey)
	if err != nil {
		return
	}
	allowanceAmount, err = c.erc.Allowance(&bind.CallOpts{Pending: false, Context: ctx}, common.HexToAddress(ownerWallet.Address()), common.HexToAddress(c.contractAddress.btp_hmny_nativecoin_bsh_core))
	return
}

// func GetWalletAndPrivKey(walFile string, password string) (wal *wallet.EvmWallet, pKey *ecdsa.PrivateKey, err error) {
// 	keyReader, err := os.Open(walFile)
// 	if err != nil {
// 		return
// 	}
// 	defer keyReader.Close()

// 	keyStore, err := ioutil.ReadAll(keyReader)
// 	if err != nil {
// 		return
// 	}

// 	pKey, err = wallet.DecryptEvmKeyStore(keyStore, []byte(password))
// 	wal = &wallet.EvmWallet{
// 		Skey: pKey,
// 		Pkey: &pKey.PublicKey,
// 	}
// 	return
// }
