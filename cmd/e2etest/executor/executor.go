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
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/bsc"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/hmny"
	"github.com/icon-project/icon-bridge/cmd/e2etest/chain/icon"
	cicon "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"
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
		} else if name == chain.BSC {
			ex.clientsPerChain[name], err = bsc.NewApi(l, cfg)
			if err != nil {
				err = errors.Wrap(err, "bsc.NewApi ")
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
		//ex.log.Debugf("Cleaning up after function returns. Removing channel of id %v", id)
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
							dst <- &evt{chainType: chains[chosen], msg: res}
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

func (ex *executor) Execute(ctx context.Context, srcChainName, dstChainName chain.ChainType, coinName string, amount *big.Int, scr Script) (err error) {
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

	srcKeys, err := src.GetKeyPairs(1)
	if err != nil {
		return errors.Wrapf(err, "GetKeyPairs for src %v", srcChainName)
	}
	srcAddress := src.GetBTPAddress(srcKeys[0][PUBKEYPOS])
	dstKeys, err := dst.GetKeyPairs(1)
	if err != nil {
		return errors.Wrapf(err, "GetKeyPairs for dst %v", dstChainName)
	}
	dstAddress := dst.GetBTPAddress(dstKeys[0][PUBKEYPOS])

	args := &args{
		watchRequestID: id,
		log:            log,
		src:            src,
		dst:            dst,
		srcKey:         srcKeys[0][PRIVKEYPOS],
		srcAddr:        srcAddress,
		dstAddr:        dstAddress,
		sinkChan:       sinkChan,
		coinName:       coinName,
	}

	amountPerCoin := map[string]*big.Int{coinName: amount}
	if src.GetNativeCoin() != coinName { // Add native coin to cut gas fees
		amountPerCoin[src.GetNativeCoin()] = amount
	}
	// Funding accounts
	nativeCoinsPerChain := map[string]chain.ChainType{"ICX": chain.ICON, "TICX": chain.ICON, "BNB": chain.BSC, "TBNB": chain.BSC, "ONE": chain.HMNY, "TONE": chain.HMNY}
	for coinName, amt := range amountPerCoin {
		if chainType, ok := nativeCoinsPerChain[coinName]; !ok {
			return fmt.Errorf("Unexpected coinName %v", coinName)
		} else {
			chainClient, ok := ex.clientsPerChain[chainType]
			if !ok {
				return fmt.Errorf("Client for chain %v doesnot exist in config ", chainType)
			}
			chainGodKey, ok := ex.godKeysPerChain[chainType]
			if !ok {
				return fmt.Errorf("GodKeys for chain %v doesnot exist in config ", chainType)
			}
			if err := ex.fund(chainClient, chainGodKey[PRIVKEYPOS], srcAddress, coinName, amt); err != nil {
				return errors.Wrapf(err, "Fund Dst: %v %v", dstChainName, srcAddress)
			} else {
				ex.log.Infof("Chain %v Funded amount %v %v to Src Chain %v Addr %v", chainType, amt, coinName, srcChainName, srcAddress)
			}
		}
	}

	time.Sleep(time.Second * 15)
	ex.log.Infof("Run ID %v %v, Transfer %v From %v To %v", id, scr.Name, coinName, src.GetChainType(), dst.GetChainType())
	if scr.Callback == nil {
		return errors.New("Callback function was nil")
	}
	if err := scr.Callback(ctx, args); err != nil {
		return errors.Wrap(err, "CallBackFunc ")
	}
	ex.log.Infof("Completed Succesfully. ID %v %v, Transfer %v From %v To %v", id, scr.Name, coinName, src.GetChainType(), dst.GetChainType())
	// CleanupFunc removeChan() is called after cb() on function return
	// so make sure cb() returns only after all the test logic is finished
	return
}

func (ex *executor) fund(api chain.ChainAPI, senderKey string, recepientAddr string, coin string, amount *big.Int) error {
	if api.GetNativeCoin() != coin {
		if _, err := api.Approve(coin, senderKey, *amount); err != nil {
			return errors.Wrapf(err, "Approve(%v,%v,%v)", coin, senderKey, amount.String())
		}
	}
	if hash, err := api.Transfer(coin, senderKey, recepientAddr, *amount); err != nil {
		return errors.Wrapf(err, "Transfer(%v,%v,%v,%v)", coin, senderKey, recepientAddr, amount.String())
	} else {
		if res, _, err := api.WaitForTxnResult(context.TODO(), hash); err != nil {
			return errors.Wrapf(err, "WaitForTxnResult %v %v", coin, hash)
		} else {
			if api.GetChainType() == chain.BSC {
				txResult, ok := res.(*ethTypes.Receipt)
				if !ok {
					return fmt.Errorf("TransactionResult Expected Type *ethtypes.Receipt Got Type %T Hash %v", res, hash)
				} else if ok && txResult.Status != 1 {
					return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Hash %v", txResult.Status, hash)
				}
			} else if api.GetChainType() == chain.ICON {
				iconTxResult, ok := res.(*cicon.TransactionResult)
				if !ok {
					return fmt.Errorf("TransactionResult Expected Type *icon.TransactionResult Got Type %T Hash %v", res, hash)
				}
				if status, err := iconTxResult.Status.Int(); err != nil {
					return errors.Wrapf(err, "Transaction Result; Error: %v Hash %v", err, hash)
				} else if status != 1 {
					return errors.Wrapf(err, "Transaction Result Expected Status 1. Got %v Failure %+v Hash %v", status, iconTxResult.Failure, hash)
				}
			}
		}
	}

	return nil
}

func (ex *executor) GetFundedAddresses(addressMap map[chain.ChainType]uint) (map[chain.ChainType][][2]string, error) {
	retMap := map[chain.ChainType][][2]string{}
	fundAmount := new(big.Int)
	fundAmount.SetString("100000000000000000000", 10)
	nativeCoinsPerChain := map[string]chain.ChainType{"ICX": chain.ICON, "TICX": chain.ICON, "BNB": chain.BSC, "TBNB": chain.BSC}
	tokensToFund := []string{"ICX", "BNB", "TICX", "TBNB"}
	for chainName, numAddr := range addressMap {
		cl, ok := ex.clientsPerChain[chainName]
		if !ok {
			return nil, fmt.Errorf("Client for chain %v doesnot exist in config ", chainName)
		}
		keys, err := cl.GetKeyPairs(int(numAddr))
		if err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairs for src %v", chainName)
		}
		for i := 0; i < len(keys); i++ {
			keys[i][PUBKEYPOS] = cl.GetBTPAddress(keys[i][PUBKEYPOS])
		}
		retMap[chainName] = keys
		// Fund
		for _, coinName := range tokensToFund {
			if chainType, ok := nativeCoinsPerChain[coinName]; !ok {
				return nil, fmt.Errorf("Unexpected coinName %v", coinName)
			} else {
				godCl, ok := ex.clientsPerChain[chainType]
				if !ok {
					return nil, fmt.Errorf("Client for chain %v doesnot exist in config ", chainType)
				}
				godKey, ok := ex.godKeysPerChain[chainType]
				if !ok {
					return nil, fmt.Errorf("GodKeys for chain %v doesnot exist in config ", chainType)
				}
				for _, key := range keys {
					if err := ex.fund(godCl, godKey[PRIVKEYPOS], key[PUBKEYPOS], coinName, fundAmount); err != nil {
						return nil, errors.Wrapf(err, "Fund %v Dst: %v", coinName, key[PUBKEYPOS])
					} else {
						time.Sleep(time.Second)
						ex.log.Infof("Funded %v %v on chain %v by chain %v", fundAmount.String(), coinName, chainName, godCl.GetChainType())
					}
				}

			}
		}
	}
	time.Sleep(time.Second * 10)
	return retMap, nil
}

func (ex *executor) ExecuteOnAddr(ctx context.Context, srcChainName, dstChainName chain.ChainType, coinName string, srcKeys, dstKeys [2]string, scr Script) (err error) {
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
	args := &args{
		watchRequestID: id,
		log:            log,
		src:            src,
		dst:            dst,
		srcKey:         srcKeys[PRIVKEYPOS],
		srcAddr:        srcKeys[PUBKEYPOS],
		dstAddr:        dstKeys[PUBKEYPOS],
		sinkChan:       sinkChan,
		coinName:       coinName,
	}

	ex.log.Infof("Run ID %v %v, Transfer %v From %v To %v", id, scr.Name, coinName, src.GetChainType(), dst.GetChainType())
	if scr.Callback == nil {
		return errors.New("Callback function was nil")
	}
	if err := scr.Callback(ctx, args); err != nil {
		return errors.Wrapf(err, "CallBackFunc ID %v Err: %v", id, err)
	}
	ex.log.Infof("Completed Succesfully. ID %v %v, Transfer %v From %v To %v", id, scr.Name, coinName, src.GetChainType(), dst.GetChainType())
	// CleanupFunc removeChan() is called after cb() on function return
	// so make sure cb() returns only after all the test logic is finished
	return nil
}
