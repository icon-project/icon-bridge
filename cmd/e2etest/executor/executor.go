package executor

import (
	"context"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/icon-project/icon-bridge/cmd/e2etest/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	FEE_NUMERATOR   = 100
	FEE_DENOMINATOR = 10000
	FIXED_PRICE     = 5000
)

type NewApiCaller func(l log.Logger, cfg *chain.Config) (chain.ChainAPI, error)

var (
	APICallerFunc = map[chain.ChainType]NewApiCaller{}
)

type executor struct {
	env              string
	log              log.Logger
	godKeysPerChain  map[chain.ChainType]keypair
	demoKeysPerChain map[chain.ChainType][]keypair
	cfgPerChain      map[chain.ChainType]*chain.Config
	clientsPerChain  map[chain.ChainType]chain.ChainAPI
	sinkChanPerID    map[uint64]chan *evt
	syncChanMtx      sync.RWMutex
	stoppedChan      chan struct{}
}

func (ex *executor) Clients() map[chain.ChainType]chain.ChainAPI {
	return ex.clientsPerChain
}

func New(l log.Logger, cfg *Config) (ex *executor, err error) {
	ex = &executor{
		env:              cfg.Env,
		log:              l,
		cfgPerChain:      make(map[chain.ChainType]*chain.Config),
		godKeysPerChain:  make(map[chain.ChainType]keypair),
		demoKeysPerChain: make(map[chain.ChainType][]keypair),
		clientsPerChain:  make(map[chain.ChainType]chain.ChainAPI),
		sinkChanPerID:    make(map[uint64]chan *evt),
		syncChanMtx:      sync.RWMutex{},
		stoppedChan:      make(chan struct{}),
	}
	for _, chainCfg := range cfg.Chains {
		apiFunc, ok := APICallerFunc[chainCfg.Name]
		if !ok {
			err = errors.Wrapf(err, "%v NewApi Func does not exist", chainCfg.Name)
			return nil, err
		} else if apiFunc == nil {
			err = errors.Wrapf(err, "%v NewApi Func is nil", chainCfg.Name)
			return nil, err
		}
		ex.cfgPerChain[chainCfg.Name] = chainCfg
		ex.clientsPerChain[chainCfg.Name], err = apiFunc(l, chainCfg)
		if err != nil {
			err = errors.Wrap(err, "hmny.NewApi ")
			return nil, err
		}
		if priv, pub, err := ex.clientsPerChain[chainCfg.Name].GetKeyPairFromKeystore(chainCfg.GodWalletKeystorePath, chainCfg.GodWalletSecretPath); err != nil {
			return nil, errors.Wrapf(err, "GetKeyPairFromKeystore %v", err)
		} else {
			ex.godKeysPerChain[chainCfg.Name] = keypair{PrivKey: priv, PubKey: pub}
		}
		//Demo
		files, err := ioutil.ReadDir(chainCfg.DemoWalletKeystorePath)
		if err != nil {
			return nil, err
		}
		tmpFile, err := ioutil.TempFile("./", "tmp")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		for _, f := range files {
			fpath := filepath.Join(chainCfg.DemoWalletKeystorePath, f.Name())
			if priv, pub, err := ex.clientsPerChain[chainCfg.Name].GetKeyPairFromKeystore(fpath, filepath.Join("./", tmpFile.Name())); err != nil {
				return nil, errors.Wrapf(err, "DemoWallets GetKeyPairFromKeystore %v", err)
			} else {
				ex.demoKeysPerChain[chainCfg.Name] = append(ex.demoKeysPerChain[chainCfg.Name], keypair{PrivKey: priv, PubKey: pub})
			}
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
		id := uint64(rand.Intn(100000)) // human readable range
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

func (ex *executor) addChan(id uint64, ch chan *evt) error {
	ex.syncChanMtx.Lock()
	defer ex.syncChanMtx.Unlock()
	if _, ok := ex.sinkChanPerID[id]; !ok {
		ex.sinkChanPerID[id] = ch
	} else {
		return errors.New("Duplicate ID")
	}
	return nil
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
	// if len(ex.sinkChanPerID) == 0 {
	// 	ex.log.Info("All test scripts have been completed")
	// 	ex.stoppedChan <- struct{}{}
	// }
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

// func (ex *executor) Execute(ctx context.Context, srcChainName, dstChainName chain.ChainType, coinNames []string, scr Script) (err error) {
// 	id, err := ex.getID()
// 	if err != nil {
// 		return errors.Wrap(err, "getID ")
// 	}
// 	log := ex.log.WithFields(log.Fields{"pid": id})
// 	sinkChan := make(chan *evt)
// 	ex.addChan(id, sinkChan)
// 	defer ex.removeChan(id) // should defer be called by cb() instead to make sure cb() was done

// 	srcCl, ok := ex.clientsPerChain[srcChainName]
// 	if !ok {
// 		return fmt.Errorf("Client for chain %v not found", srcChainName)
// 	}
// 	dstCl, ok := ex.clientsPerChain[dstChainName]
// 	if !ok {
// 		return fmt.Errorf("Client for chain %v not found", dstChainName)
// 	}
// 	srcGod, ok := ex.godKeysPerChain[srcChainName]
// 	if !ok {
// 		return fmt.Errorf("GodKeys for chain %v not found", srcChainName)
// 	}
// 	dstGod, ok := ex.godKeysPerChain[dstChainName]
// 	if !ok {
// 		return fmt.Errorf("GodKeys for chain %v not found", dstChainName)
// 	}
// 	srcDemo, ok := ex.demoKeysPerChain[srcChainName]
// 	if !ok {
// 		return fmt.Errorf("DemoKeys for chain %v not found", srcChainName)
// 	}
// 	dstDemo, ok := ex.demoKeysPerChain[dstChainName]
// 	if !ok {
// 		return fmt.Errorf("DemoKeys for chain %v not found", dstChainName)
// 	}
// 	srcCfg, ok := ex.cfgPerChain[srcChainName]
// 	if !ok {
// 		return fmt.Errorf("Cfg for chain %v not found", srcChainName)
// 	}
// 	dstCfg, ok := ex.cfgPerChain[dstChainName]
// 	if !ok {
// 		return fmt.Errorf("Cfg for chain %v not found", srcChainName)
// 	}
// 	btsAddressPerChain := map[chain.ChainType]string{
// 		srcChainName: srcCfg.ContractAddresses[chain.BTS],
// 		dstChainName: dstCfg.ContractAddresses[chain.BTS],
// 	}
// 	gasLimitPerChain := map[chain.ChainType]int64{
// 		srcChainName: srcCfg.GasLimit,
// 		dstChainName: dstCfg.GasLimit,
// 	}

// 	ts := &testSuite{
// 		id:                 id,
// 		logger:             log,
// 		env:                ex.env,
// 		subChan:            sinkChan,
// 		btsAddressPerChain: btsAddressPerChain,
// 		gasLimitPerChain:   gasLimitPerChain,
// 		clsPerChain:        map[chain.ChainType]chain.ChainAPI{srcChainName: srcCl, dstChainName: dstCl},
// 		godKeysPerChain:    map[chain.ChainType]keypair{srcChainName: srcGod, dstChainName: dstGod},
// 		demoKeysPerChain:   map[chain.ChainType][]keypair{srcChainName: srcDemo, dstChainName: dstDemo},
// 		fee:                fee{numerator: big.NewInt(FEE_NUMERATOR), denominator: big.NewInt(FEE_DENOMINATOR), fixed: big.NewInt(FIXED_PRICE)},
// 	}

// 	ex.log.Infof("Run ID %v %v, Transfer %v From %v To %v", id, scr.Name, coinNames, srcChainName, dstChainName)
// 	if scr.Callback == nil {
// 		return errors.New("Callback function was nil")
// 	}
// 	if _, err := scr.Callback(ctx, srcChainName, dstChainName, coinNames, ts); err != nil {
// 		return errors.Wrap(err, "CallBackFunc ")
// 	}
// 	ex.log.Infof("Completed Succesfully. ID %v %v, Transfer %v From %v To %v", id, scr.Name, coinNames, srcChainName, dstChainName)
// 	// CleanupFunc removeChan() is called after cb() on function return
// 	// so make sure cb() returns only after all the test logic is finished
// 	return
// }
