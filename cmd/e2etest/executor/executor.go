package executor

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/icon"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

type executor struct {
	log             log.Logger
	godKeysPerChain map[chain.ChainType][2]string
	cfgPerChain     map[chain.ChainType]*chain.ChainConfig
	clientsPerChain map[chain.ChainType]chain.ChainAPI
	sinkChanPerID   map[uint64]chan *evt
	syncChanMtx     sync.RWMutex
	stoppedChan     chan struct{}
}

func New(l log.Logger, cfgPerChain map[chain.ChainType]*chain.ChainConfig) (ex *executor, err error) {
	getKeyPairFromFile := func(walFile string, password string) (pair [2]string, err error) {
		keyReader, err := os.Open(walFile)
		if err != nil {
			err = errors.Wrapf(err, "os.Open file %v", walFile)
			return
		}
		defer keyReader.Close()

		keyStore, err := ioutil.ReadAll(keyReader)
		if err != nil {
			err = errors.Wrapf(err, "ioutil.ReadAll %v", walFile)
			return
		}
		key, err := keystore.DecryptKey(keyStore, password)
		if err != nil {
			err = errors.Wrapf(err, "keystore.Decrypt %v", walFile)
			return
		}
		privBytes := ethcrypto.FromECDSA(key.PrivateKey)
		privString := hex.EncodeToString(privBytes)
		addr := ethcrypto.PubkeyToAddress(key.PrivateKey.PublicKey)
		pair = [2]string{privString, addr.String()}
		return
	}
	ex = &executor{
		log:             l,
		cfgPerChain:     cfgPerChain,
		godKeysPerChain: make(map[chain.ChainType][2]string),
		clientsPerChain: make(map[chain.ChainType]chain.ChainAPI),
		sinkChanPerID:   make(map[uint64]chan *evt),
		syncChanMtx:     sync.RWMutex{},
		stoppedChan:     make(chan struct{}),
	}
	for name, cfg := range cfgPerChain {
		// GodKeys
		if pair, err := getKeyPairFromFile(cfg.GodWallet.Path, cfg.GodWallet.Password); err != nil {
			return nil, errors.Wrapf(err, "getKeyPairFromFile(%v)", cfg.GodWallet.Path)
		} else {
			ex.godKeysPerChain[name] = pair
		}
		//Clients
		if name == chain.HMNY {
			ex.clientsPerChain[name], err = hmny.NewApi(l, cfg)
			if err != nil {
				err = errors.Wrap(err, "hmny.NewApi ")
				return nil, err
			}
		} else if name == chain.ICON {
			ex.clientsPerChain[name], err = icon.NewApi(l, cfg)
			if err != nil {
				err = errors.Wrap(err, "icon.NewApi ")
				return nil, err
			}

		} else {
			return nil, errors.New("Unknown Chain Type supplied from config: " + string(name))
		}
	}

	return
}

func (ex *executor) Done() <-chan struct{} {
	return ex.stoppedChan
}

func (ex *executor) getID() (uint64, error) {
	ex.syncChanMtx.RLock()
	defer ex.syncChanMtx.RUnlock()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 3; i++ { // Try 5 times to get a random number not already used
		id := uint64(rand.Intn(10000)) // human readable range
		if _, ok := ex.sinkChanPerID[id]; !ok {
			return id, nil
		}
	}
	for i := 0; i < 2; i++ {
		id := rand.Uint64() // larger search space
		if _, ok := ex.sinkChanPerID[id]; !ok {
			return id, nil
		}
	}
	return 0, errors.New("Couldn't find a unique ID")
}

func (ex *executor) addChan(id uint64, ch chan *evt) {
	ex.syncChanMtx.Lock()
	defer ex.syncChanMtx.Unlock()
	ex.sinkChanPerID[id] = ch
}

func (ex *executor) removeChan(id uint64) {
	ex.syncChanMtx.Lock()
	defer ex.syncChanMtx.Unlock()
	if ch, ok := ex.sinkChanPerID[id]; ok {
		ex.log.Debugf("Cleaning up after function returns. Removing channel of id %v", id)
		if ch != nil {
			close(ex.sinkChanPerID[id])
		}
		delete(ex.sinkChanPerID, id)
	}
	if len(ex.sinkChanPerID) == 0 {
		ex.log.Info("All test scripts have been completed")
		ex.stoppedChan <- struct{}{}
	}
}

func (ex *executor) getChan(id uint64) chan *evt {
	ex.syncChanMtx.RLock()
	defer ex.syncChanMtx.RUnlock()
	if _, ok := ex.sinkChanPerID[id]; ok {
		return ex.sinkChanPerID[id]
	} else {
		ex.log.Warnf("Message Target id %v does not exist", id)
	}
	return nil
}

func (ex *executor) Subscribe(ctx context.Context) {
	go func() {
		lenCls := len(ex.clientsPerChain)
		chains := make([]chain.ChainType, lenCls)
		cases := make([]reflect.SelectCase, 1+(lenCls*2))
		i := 0
		for name, cl := range ex.clientsPerChain {
			ex.log.Debugf("Start Subscription %v", name)
			sinkChan, errChan, err := cl.Subscribe(ctx)
			if err != nil {
				ex.log.Error(errors.Wrapf(err, "%v: Subscribe()", name))
			}
			chains[i] = name
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(sinkChan)}
			cases[i+lenCls] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(errChan)}
			i++
		}
		cases[len(cases)-1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())}
		for {
			chosen, value, ok := reflect.Select(cases)
			if !ok {
				if chosen == len(cases)-1 {
					ex.log.Error("Context cancelled. Exiting Executor")
					return
				}
				if chosen >= lenCls {
					ex.log.Debugf("Sender Closed ErrMessage ChannelID %v Client %v", chosen, chains[chosen-lenCls])
				} else {
					ex.log.Debugf("Sender Closed EvtMessage ChannelID %v Client %v", chosen, chains[chosen])
				}
				cases[chosen].Chan = reflect.ValueOf(nil)
				continue
			}

			if chosen < lenCls { // [0, lenCapi-1] is message
				res, dok := value.Interface().(*chain.EventLogInfo)
				if !dok {
					ex.log.Errorf("Got interface of type %T; Expected errorType", value)
					break
				}
				if len(res.IDs) > 0 {
					for _, id := range res.IDs {
						if dst := ex.getChan(id); dst != nil {
							dst <- &evt{name: chains[chosen], msg: res}
						}
					}
				} else {
					ex.log.Warnf("Message without target received %+v", res)
				}
			} else if chosen >= lenCls && chosen < 2*len(cases) {
				res, eok := value.Interface().(error)
				if !eok {
					ex.log.Errorf("Got interface of type %T; Expected errorType", value)
					break
				}
				ex.log.Errorf("ErrMessage %v %+v", chains[chosen-lenCls], res)
			} else {
				ex.log.Error("Context cancelled. Exiting Executor")
				return
			}
		}
	}()
}

func (ex *executor) Execute(ctx context.Context, srcChainName, dstChainName chain.ChainType, amount *big.Int, cb callBackFunc) (err error) {
	id, err := ex.getID()
	if err != nil {
		return errors.Wrap(err, "getID ")
	}
	log := ex.log.WithFields(log.Fields{"pid": id})
	sinkChan := make(chan *evt)
	ex.addChan(id, sinkChan)
	defer ex.removeChan(id) // should defer be called by cb() instead to make sure cb() was done

	src, ok := ex.clientsPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("Client for chain %v doesnot exist in config ", srcChainName)
	}
	dst, ok := ex.clientsPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("Client for chain %v doesnot exist in config ", dstChainName)
	}
	srcGKey, ok := ex.godKeysPerChain[srcChainName]
	if !ok {
		return fmt.Errorf("GodKeys for chain %v doesnot exist in config", srcChainName)
	}
	dstGKey, ok := ex.godKeysPerChain[dstChainName]
	if !ok {
		return fmt.Errorf("GodKeys for chain %v doesnot exist in config", dstChainName)
	}
	srcKeys, err := src.GetKeyPairs(1)
	if err != nil {
		return errors.Wrapf(err, "GetKeyPairs for src %v", srcChainName)
	}
	srcAddress := src.GetBTPAddress(srcKeys[0][1])
	dstKeys, err := dst.GetKeyPairs(1)
	if err != nil {
		return errors.Wrapf(err, "GetKeyPairs for dst %v", dstChainName)
	}
	dstAddress := dst.GetBTPAddress(dstKeys[0][1])

	coinsToFund := map[chain.ChainType][]string{chain.ICON: {"ICX", "TICX"}, chain.HMNY: {"ONE", "TONE"}}
	//wrappedCointsToFund := map[chain.ChainType][]string{chain.ICON: {"ONE", "TONE"}, chain.HMNY: {"ICX", "TICX"}}

	log.Info("Funding src and dst accounts with and native coin and token")
	if err := ex.fund(log, src, srcGKey[0], srcAddress, coinsToFund[srcChainName], amount); err != nil {
		return errors.Wrapf(err, "Fund Src: %v %v", srcChainName, srcAddress)
	}
	// if err := ex.fund(log, dst, dstGKey[0], srcAddress, wrappedCointsToFund[srcChainName], amount); err != nil {
	// 	return errors.Wrapf(err, "Fund Src: %v %v", srcChainName, srcAddress)
	// }
	if err := ex.fund(log, dst, dstGKey[0], dstAddress, coinsToFund[dstChainName], amount); err != nil {
		return errors.Wrapf(err, "Fund Dst: %v %v", dstChainName, srcAddress)
	}
	// if err := ex.fund(log, src, srcGKey[0], dstAddress, wrappedCointsToFund[dstChainName], amount); err != nil {
	// 	return errors.Wrapf(err, "Fund Dst: %v %v", dstChainName, srcAddress)
	// }

	time.Sleep(time.Second * 10)
	args, err := newArgs(
		id, log,
		src, dst, srcGKey[0],
		srcKeys[0][0], srcAddress, dstAddress,
		sinkChan,
	)

	if err != nil {
		return errors.Wrap(err, "newArgs ")
	}

	if cb != nil {
		if err := cb(ctx, args); err != nil {
			return errors.Wrap(err, "CallBackFunc ")
		}
	} else {
		return errors.New("Callback function was nil")
	}

	// CleanupFunc removeChan() is called after cb() on function return
	// so make sure cb() returns only after all the test logic is finished
	return
}

func (ex *executor) fund(log log.Logger, api chain.ChainAPI, senderKey string, recepientAddr string, coins []string, amount *big.Int) error {
	for _, coin := range coins {
		if hash, err := api.Transfer(coin, senderKey, recepientAddr, *amount); err != nil {
			return errors.Wrapf(err, "Transfer(%v,%v,%v,%v)", coin, senderKey, recepientAddr, amount.String())
		} else {
			if _, _, err := api.WaitForTxnResult(context.TODO(), hash); err != nil {
				return errors.Wrapf(err, "WaitForTxnResult %v %v", coin, hash)
			} else {
				log.Infof("Funded %v with %v %v. Hash: %v", recepientAddr, amount.String(), coin, hash)
			}
		}

	}
	return nil
}
