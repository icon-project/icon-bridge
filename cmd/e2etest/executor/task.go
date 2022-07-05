package executor

import (
	"context"
	"math/big"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	PRIVKEYPOS = 0
	PUBKEYPOS  = 1
)

type evt struct {
	msg  *chain.EventLogInfo
	name chain.ChainType
}

type args struct {
	watchRequestID uint64
	log            log.Logger
	src            chain.SrcAPI
	dst            chain.DstAPI
	srcKey         string
	srcAddr        string
	dstAddr        string
	godKey         string
	sinkChan       <-chan *evt
}

func newArgs(id uint64, l log.Logger,
	src chain.SrcAPI, dst chain.DstAPI,
	godKey string, srcKey string, srcAddr string, dstAddr string,
	sinkChan <-chan *evt,
) (t *args, err error) {
	tu := &args{watchRequestID: id, log: l,
		src: src, dst: dst,
		godKey: godKey,
		srcKey: srcKey, srcAddr: srcAddr, dstAddr: dstAddr,
		sinkChan: sinkChan,
	}
	return tu, nil
}

type callBackFunc func(ctx context.Context, args *args) error

var showBalance = func(args *args) error {
	args.log.Infof("Balance of Src Addr %v", args.srcAddr)
	for _, coin := range []string{"ICX", "ETH", "ONE"} {
		if amt, err := args.src.GetCoinBalance(coin, args.srcAddr); err != nil {
			return errors.Wrapf(err, "GetCoinBalance(%v, %v)", coin, args.srcAddr)
		} else {
			args.log.Infof(" %v %v", amt.String(), coin)
		}
	}
	args.log.Infof("Balance of Dst Addr %v", args.dstAddr)
	for _, coin := range []string{"ICX", "ETH", "ONE"} {
		if amt, err := args.dst.GetCoinBalance(coin, args.dstAddr); err != nil {
			return errors.Wrapf(err, "GetCoinBalance(%v, %v)", coin, args.dstAddr)
		} else {
			args.log.Infof(" %v %v", amt.String(), coin)
		}
	}
	return nil
}

var TestScripts []callBackFunc = []callBackFunc{
	func(ctx context.Context, args *args) error {
		args.log.Info("Test Script to check wrapped coin transfer")
		args.log.Info("Showing Initial Balance")
		if err := showBalance(args); err != nil {
			return errors.Wrap(err, "showBalance ")
		}
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		coin := "ONE"
		if hash, err := args.src.Approve(coin, args.srcKey, *amt); err != nil {
			return errors.Wrap(err, "Approve ")
		} else {
			if _, _, err := args.src.WaitForTxnResult(context.TODO(), hash); err != nil {
				return errors.Wrapf(err, "WaitForTxnResult %v %v", coin, hash)
			} else {
				args.log.Infof("Approve %v with %v %v. Hash: %v", args.dstAddr, amt.String(), coin, hash)
			}
		}
		if hash, err := args.src.Transfer(coin, args.srcKey, args.dstAddr, *amt); err != nil {
			return errors.Wrap(err, "Transfer ")
		} else {
			if _, _, err := args.src.WaitForTxnResult(context.TODO(), hash); err != nil {
				return errors.Wrapf(err, "WaitForTxnResult %v %v", coin, hash)
			} else {
				args.log.Infof("Transfer %v with %v %v. Hash: %v", args.dstAddr, amt.String(), coin, hash)
			}
		}
		time.Sleep(time.Second * 15)
		args.log.Info("Showing Final Balance")
		if err := showBalance(args); err != nil {
			return errors.Wrap(err, "showBalance ")
		}
		return nil
	},
	func(ctx context.Context, args *args) error {
		args.log.Info("Test Script to check token transfer")
		args.log.Info("Showing Initial Balance")
		if err := showBalance(args); err != nil {
			return errors.Wrap(err, "showBalance ")
		}
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		coin := "ETH"
		if hash, err := args.src.Approve(coin, args.srcKey, *amt); err != nil {
			return errors.Wrap(err, "Approve ")
		} else {
			if _, _, err := args.src.WaitForTxnResult(context.TODO(), hash); err != nil {
				return errors.Wrapf(err, "WaitForTxnResult %v %v", coin, hash)
			} else {
				args.log.Infof("Approve %v with %v %v. Hash: %v", args.dstAddr, amt.String(), coin, hash)
			}
		}
		if hash, err := args.src.Transfer(coin, args.srcKey, args.dstAddr, *amt); err != nil {
			return errors.Wrap(err, "Transfer ")
		} else {
			if _, _, err := args.src.WaitForTxnResult(context.TODO(), hash); err != nil {
				return errors.Wrapf(err, "WaitForTxnResult %v %v", coin, hash)
			} else {
				args.log.Infof("Transfer %v with %v %v. Hash: %v", args.dstAddr, amt.String(), coin, hash)
			}
		}
		time.Sleep(time.Second * 15)
		args.log.Info("Showing Final Balance")
		if err := showBalance(args); err != nil {
			return errors.Wrap(err, "showBalance ")
		}
		return nil
	},
	func(ctx context.Context, args *args) error {
		args.log.Info("Test Script to check native coin event logs")
		args.log.Info("Showing Initial Balance")
		if err := showBalance(args); err != nil {
			return errors.Wrap(err, "showBalance ")
		}
		amt := new(big.Int)
		amt.SetString("2000000000000000000", 10)
		coin := "ICX"
		if hash, err := args.src.Transfer(coin, args.srcKey, args.dstAddr, *amt); err != nil {
			return errors.Wrap(err, "Transfer ")
		} else {
			if _, els, err := args.src.WaitForTxnResult(context.TODO(), hash); err != nil {
				return errors.Wrapf(err, "WaitForTxnResult %v %v", coin, hash)
			} else {
				args.log.Infof("Transfer %v with %v %v. Hash: %v", args.dstAddr, amt.String(), coin, hash)
				for _, el := range els {
					if el.EventType != chain.TransferStart {
						continue
					}
					seq, _ := el.GetSeq()
					args.log.Infof("WatchForTransferEnd %v, %v, %v", args.watchRequestID, coin, seq)
					if err := args.src.WatchForTransferEnd(args.watchRequestID, coin, seq); err != nil {
						return errors.Wrap(err, "WatchForTransferEnd")
					}
					args.log.Infof("WatchForTransferReceived %v, %v, %v", args.watchRequestID, coin, seq)
					if err := args.dst.WatchForTransferReceived(args.watchRequestID, coin, seq); err != nil {
						return errors.Wrap(err, "WatchForTransferEnd")
					}
				}
				args.log.Infof("Waiting for events ")
				counter := 0
				for {
					select {
					case <-ctx.Done():
						args.log.Warn("Context Cancelled Exiting task")
						return nil
					case res := <-args.sinkChan:
						args.log.Infof("%v: %+v", res.name, res.msg)
						counter += 1
						if counter >= 2 {
							args.log.Infof("Received all events. Closing...")
							return nil
						}
					}
				}
			}
		}
		return nil
	},
}

// var DemoSubCallback callBackFunc = func(ctx context.Context, args *args) error {

// 	// fund demo wallets
// 	args.log.Info("Starting demo...")
// 	ienv, ok := args.clientsPerChain[chain.ICON]
// 	if !ok {
// 		return errors.New("Icon client not found")
// 	}
// 	henv, ok := args.clientsPerChain[chain.HMNY]
// 	if !ok {
// 		return errors.New("Hmny client not found")
// 	}
// 	igod, ok := args.godKeysPerChain[chain.ICON]
// 	if !ok {
// 		return errors.New("God Keys not found for ICON")
// 	}
// 	hgod, ok := args.godKeysPerChain[chain.HMNY]
// 	if !ok {
// 		return errors.New("God keys not found for Hmy")
// 	}
// 	tmp, err := ienv.GetKeyPairs(1)
// 	if err != nil {
// 		return errors.New("Couldn't create demo account for icon")
// 	}
// 	iDemo := tmp[0]
// 	tmp, err = henv.GetKeyPairs(1)
// 	if err != nil {
// 		return errors.New("Couldn't create demo account for hmny")
// 	}
// 	hDemo := tmp[0]
// 	args.log.Info("Creating Demo Icon Account ", iDemo)
// 	args.log.Info("Creating Demo Hmy Account ", hDemo)
// 	// findAddrForContract := func(inputName chain.ContractName) (retAddr string, ok bool) {
// 	// 	for addr, name := range args.addrToName {
// 	// 		if name == inputName {
// 	// 			return addr, true
// 	// 		}
// 	// 	}
// 	// 	return "", false
// 	// }

// 	watchTransferStart := func(elInfo []*chain.EventLogInfo) error {
// 		// for _, el := range elInfo {
// 		// 	if el.EventType != chain.TransferStart {
// 		// 		continue
// 		// 	}
// 		// 	seq, err := el.GetSeq()
// 		// 	if err != nil {
// 		// 		return err
// 		// 	}
// 		// 	ctrName, ok := args.addrToName[el.ContractAddress]
// 		// 	if !ok {
// 		// 		return fmt.Errorf("Event %v generated by %v is not in config", el.EventType, el.ContractAddress)
// 		// 	}
// 		// 	args.log.Infof("Generated event %v contractName %v SeqNo %v", el.EventType, ctrName, seq)
// 		// 	if ctrName == chain.NativeBSHIcon {
// 		// 		if ctr, ok := findAddrForContract(chain.NativeBSHPeripheryHmy); ok {
// 		// 			henv.WatchForTransferReceived(args.id, seq, ctr)
// 		// 			ienv.WatchForTransferEnd(args.id, seq, el.ContractAddress)
// 		// 		} else {
// 		// 			return errors.New("NativeBSHPeripheryHmy does not exist in config")
// 		// 		}
// 		// 	} else if ctrName == chain.NativeBSHPeripheryHmy {
// 		// 		if ctr, ok := findAddrForContract(chain.NativeBSHIcon); ok {
// 		// 			henv.WatchForTransferEnd(args.id, seq, el.ContractAddress)
// 		// 			ienv.WatchForTransferReceived(args.id, seq, ctr)
// 		// 		} else {
// 		// 			return errors.New("NativeBSHIcon does not exist in config")
// 		// 		}
// 		// 	} else if ctrName == chain.TokenBSHIcon {
// 		// 		if ctr, ok := findAddrForContract(chain.TokenBSHImplHmy); ok {
// 		// 			henv.WatchForTransferReceived(args.id, seq, ctr)
// 		// 			ienv.WatchForTransferEnd(args.id, seq, el.ContractAddress)
// 		// 		} else {
// 		// 			return errors.New("TokenBSHImplHmy does not exist in config")
// 		// 		}
// 		// 	} else if ctrName == chain.TokenBSHImplHmy {
// 		// 		if ctr, ok := findAddrForContract(chain.TokenBSHIcon); ok {
// 		// 			henv.WatchForTransferEnd(args.id, seq, el.ContractAddress)
// 		// 			ienv.WatchForTransferReceived(args.id, seq, ctr)
// 		// 		} else {
// 		// 			return errors.New("NativeBSHIcon does not exist in config")
// 		// 		}
// 		// 	} else {
// 		// 		args.log.Warnf("Unexpected contract name %v ", ctrName)
// 		// 	}
// 		// }
// 		return nil
// 	}

// 	args.log.Info("Funding Demo Wallets ")
// 	amt := new(big.Int)
// 	amt.SetString("250000000000000000000", 10)
// 	_, err = ienv.Transfer("ICX", igod[PRIVKEYPOS], *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	amt = new(big.Int)
// 	amt.SetString("10000000000000000000", 10)
// 	_, err = ienv.Transfer("ETH", igod[PRIVKEYPOS], *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	amt = new(big.Int)
// 	amt.SetString("10000000000000000000", 10)
// 	_, err = henv.Transfer("ONE", hgod[PRIVKEYPOS], *henv.GetBTPAddress(hDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	amt = new(big.Int)
// 	amt.SetString("10000000000000000000", 10)
// 	_, err = henv.Transfer("ETH", hgod[PRIVKEYPOS], *henv.GetBTPAddress(hDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	args.log.Info("Done funding")
// 	time.Sleep(time.Second * 10)
// 	go func(ctx context.Context) {
// 		args.log.Info("Starting Watch")
// 		defer args.closeFunc()
// 		counter := 0
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				args.log.Warn("Context Cancelled Exiting task")
// 				return
// 			case res := <-args.sinkChan:
// 				args.log.Infof("%v: %+v", res.name, res.msg)
// 				counter += 1
// 				if counter >= 4 { // 2 Watch calls * 2 TxEvents{End,Rx}
// 					args.log.Infof("Received all events. Closing...")
// 					return
// 				}
// 			}
// 		}

// 	}(ctx)
// 	time.Sleep(time.Second * 15)

// 	args.log.Info("Transfer Native ICX to HMY")
// 	amt = new(big.Int)
// 	amt.SetString("2000000000000000000", 10)
// 	_, err = ienv.Transfer("ICX", iDemo[PRIVKEYPOS], *henv.GetBTPAddress(hDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	args.log.Info("Transfer Native ONE to ICX")
// 	amt = new(big.Int)
// 	amt.SetString("2000000000000000000", 10)
// 	_, err = henv.Transfer("ONE", hDemo[PRIVKEYPOS], *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	args.log.Info("Approve")
// 	time.Sleep(time.Second * 10)

// 	amt = new(big.Int)
// 	amt.SetString("100000000000000000000000", 10)
// 	_, err = ienv.Approve("ONE", iDemo[PRIVKEYPOS], *amt)
// 	if err != nil {
// 		return err
// 	}
// 	amt = new(big.Int)
// 	amt.SetString("100000000000000000000000", 10)
// 	_, err = henv.Approve("ICX", hDemo[PRIVKEYPOS], *amt)
// 	if err != nil {
// 		return err
// 	}
// 	time.Sleep(5 * time.Second)

// 	args.log.Info("Transfer Wrapped")
// 	amt = new(big.Int)
// 	amt.SetString("1000000000000000000", 10)
// 	hash, err := ienv.Transfer("ONE", iDemo[PRIVKEYPOS], *henv.GetBTPAddress(hDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	_, elInfo, err := ienv.WaitForTxnResult(ctx, hash)
// 	if err != nil {
// 		return err
// 	}
// 	watchTransferStart(elInfo)

// 	amt = new(big.Int)
// 	amt.SetString("1000000000000000000", 10)
// 	hash, err = henv.Transfer("ICX", hDemo[PRIVKEYPOS], *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	_, elInfo, err = henv.WaitForTxnResult(ctx, hash)
// 	if err != nil {
// 		return err
// 	}
// 	watchTransferStart(elInfo)

// 	return nil
// }

// var DemoRequestCallback callBackFunc = func(ctx context.Context, args *args) error {
// 	defer args.closeFunc()
// 	// fund demo wallets
// 	args.log.Info("Starting demo...")
// 	ienv, ok := args.clientsPerChain[chain.ICON]
// 	if !ok {
// 		return errors.New("Icon client not found")
// 	}
// 	henv, ok := args.clientsPerChain[chain.HMNY]
// 	if !ok {
// 		return errors.New("Hmny client not found")
// 	}
// 	igod, ok := args.godKeysPerChain[chain.ICON]
// 	if !ok {
// 		return errors.New("God Keys not found for ICON")
// 	}
// 	hgod, ok := args.godKeysPerChain[chain.HMNY]
// 	if !ok {
// 		return errors.New("God keys not found for Hmy")
// 	}
// 	tmp, err := ienv.GetKeyPairs(1)
// 	if err != nil {
// 		return errors.New("Couldn't create demo account for icon")
// 	}
// 	iDemo := tmp[0]
// 	tmp, err = henv.GetKeyPairs(1)
// 	if err != nil {
// 		return errors.New("Couldn't create demo account for hmny")
// 	}
// 	hDemo := tmp[0]
// 	args.log.Info("Creating Demo Icon Account ", iDemo)
// 	args.log.Info("Creating Demo Hmy Account ", hDemo)
// 	showBalance := func(log log.Logger, env chain.ChainAPI, addr string, tokens []string) error {
// 		factor := new(big.Int)
// 		factor.SetString("10000000000000000", 10)
// 		for _, token := range tokens {
// 			if amt, err := env.GetCoinBalance(token, addr); err != nil {
// 				return err
// 			} else {
// 				log.Infof("%v: %v", token, amt.Div(amt, factor).String())
// 			}
// 		}
// 		return nil
// 	}
// 	args.log.Info("Funding Demo Wallets ")
// 	amt := new(big.Int)
// 	amt.SetString("250000000000000000000", 10)
// 	_, err = ienv.Transfer("ICX", igod[PRIVKEYPOS], *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	amt = new(big.Int)
// 	amt.SetString("10000000000000000000", 10)
// 	_, err = ienv.Transfer("ETH", igod[PRIVKEYPOS], *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	amt = new(big.Int)
// 	amt.SetString("10000000000000000000", 10)
// 	_, err = henv.Transfer("ONE", hgod[PRIVKEYPOS], *henv.GetBTPAddress(hDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	amt = new(big.Int)
// 	amt.SetString("10000000000000000000", 10)
// 	_, err = henv.Transfer("ETH", hgod[PRIVKEYPOS], *henv.GetBTPAddress(hDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	args.log.Info("Done funding")
// 	time.Sleep(time.Second * 10)
// 	// args.log.Info("ICON:  ")
// 	// if err := showBalance(args.log, ienv, iDemo[PUBKEYPOS], []chain.TokenType{chain.ICXToken, chain.IRC2Token, chain.ONEToken}); err != nil {
// 	// 	return err
// 	// }
// 	// args.log.Info("HMNY:   ")
// 	// if err := showBalance(args.log, henv, hDemo[PUBKEYPOS], []chain.TokenType{chain.ONEToken, chain.ERC20Token, chain.ICXToken}); err != nil {
// 	// 	return err
// 	// }

// 	args.log.Info("Transfer Native ICX to HMY")
// 	amt = new(big.Int)
// 	amt.SetString("2000000000000000000", 10)
// 	_, err = ienv.Transfer("ICX", iDemo[PRIVKEYPOS], *henv.GetBTPAddress(hDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	args.log.Info("Transfer Native ONE to ICX")
// 	amt = new(big.Int)
// 	amt.SetString("2000000000000000000", 10)
// 	_, err = henv.Transfer("ONE", hDemo[PRIVKEYPOS], *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	args.log.Info("Approve")
// 	time.Sleep(time.Second * 10)

// 	amt = new(big.Int)
// 	amt.SetString("100000000000000000000000", 10)
// 	_, err = ienv.Approve("ONE", iDemo[PRIVKEYPOS], *amt)
// 	if err != nil {
// 		return err
// 	}
// 	amt = new(big.Int)
// 	amt.SetString("100000000000000000000000", 10)
// 	_, err = henv.Approve("ICX", hDemo[PRIVKEYPOS], *amt)
// 	if err != nil {
// 		return err
// 	}
// 	time.Sleep(5 * time.Second)

// 	args.log.Info("Transfer Wrapped")
// 	amt = new(big.Int)
// 	amt.SetString("1000000000000000000", 10)
// 	_, err = ienv.Transfer("ONE", iDemo[PRIVKEYPOS], *henv.GetBTPAddress(hDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}

// 	amt = new(big.Int)
// 	amt.SetString("1000000000000000000", 10)
// 	_, err = henv.Transfer("ICX", hDemo[PRIVKEYPOS], *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	time.Sleep(10 * time.Second)

// 	args.log.Info("Transfer Irc2 to HMY")
// 	amt = new(big.Int)
// 	amt.SetString("1000000000000000000", 10)
// 	_, err = ienv.Approve("ETH", iDemo[PRIVKEYPOS], *amt)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = ienv.Transfer("ETH", iDemo[PRIVKEYPOS], *henv.GetBTPAddress(hDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	args.log.Info("Transfer Erc20 to ICon")
// 	amt = new(big.Int)
// 	amt.SetString("1000000000000000000", 10)
// 	_, err = henv.Approve("ETH", hDemo[PRIVKEYPOS], *amt)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = henv.Transfer("ETH", hDemo[PRIVKEYPOS], *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), *amt)
// 	if err != nil {
// 		return err
// 	}
// 	time.Sleep(30 * time.Second)
// 	args.log.Info("ICON:  ")
// 	if err := showBalance(args.log, ienv, iDemo[PUBKEYPOS], []string{"ICX", "ETH", "ONE"}); err != nil {
// 		return err
// 	}
// 	args.log.Info("HMNY:   ")
// 	if err := showBalance(args.log, henv, hDemo[PUBKEYPOS], []string{"ONE", "ETH", "ICX"}); err != nil {
// 		return err
// 	}
// 	args.log.Info("Done")
// 	return nil
// }
