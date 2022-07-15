package executor

import (
	"context"
	"math/big"
	"time"

	"github.com/icon-project/icon-bridge/bmr/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/bmr/common/errors"
	"github.com/icon-project/icon-bridge/bmr/common/log"
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
	for _, coin := range []string{"ICX", "ONE", "TICX", "TONE"} {
		if amt, err := args.src.GetCoinBalance(coin, args.srcAddr); err != nil {
			return errors.Wrapf(err, "GetCoinBalance(%v, %v)", coin, args.srcAddr)
		} else {
			args.log.Infof(" %v %v", amt.String(), coin)
		}
	}
	args.log.Infof("Balance of Dst Addr %v", args.dstAddr)
	for _, coin := range []string{"ICX", "ONE", "TICX", "TONE"} {
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
		coin := "TICX"
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
