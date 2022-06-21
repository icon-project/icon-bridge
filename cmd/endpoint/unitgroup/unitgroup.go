package unitgroup

import (
	"context"
	"sync"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	DefaultTestUnitTimeoutSeconds = 60
	DefaultFundingAmount          = 100000000000000
)

type unitgroup struct {
	godKeysPerChain map[chain.ChainType][2]string
	cfgPerChain     map[chain.ChainType]*chain.ChainConfig
	log             log.Logger
	cache           *tEnvTaskCache
}

func New(l log.Logger, cfgPerChain map[chain.ChainType]*chain.ChainConfig) (ug *unitgroup, err error) {

	ug = &unitgroup{
		log:             l,
		cfgPerChain:     cfgPerChain,
		cache:           &tEnvTaskCache{mem: map[int64]tEnvTask{}, mu: sync.RWMutex{}, lastAdded: 0},
		godKeysPerChain: make(map[chain.ChainType][2]string),
	}
	for name, cfg := range cfgPerChain {
		if pair, err := GetKeyPairFromFile(cfg.GodWallet.Path, cfg.GodWallet.Password); err != nil {
			return nil, err
		} else {
			ug.godKeysPerChain[name] = pair
		}
	}
	return
}

func (ug *unitgroup) Start(ctx context.Context) error {
	cachePoller := time.NewTicker(time.Duration(2) * time.Second)
	lastChecked := int64(0)
	res := map[int64]tEnvTask{}
	errChan := make(chan error)
	go func() { // Poll for newly added tasks and feed it to executor
		defer cachePoller.Stop()
		for {
			select {
			case <-cachePoller.C:
				res, lastChecked = ug.cache.GetNew(lastChecked)
				//ug.log.Warn("Poll ", len(res), lastChecked)
				for ts, r := range res {
					ug.log.Warn("Spawn processing go routine")
					go ug.process(ctx, ts, r, errChan)
				}
			case <-ctx.Done():
				break
			case err := <-errChan:
				ug.log.Error("UnitGroup; Start Error ", err)
			}
		}
	}()
	return nil
}

func (ug *unitgroup) process(ctx context.Context, ts int64, r tEnvTask, errChan chan error) {
	var err error
	if r.tfunc.PreRun != nil {
		if err = r.tfunc.PreRun(r.tu); err != nil {
			errChan <- err
		}
	}
	if r.tfunc.Run != nil {
		if err = r.tfunc.Run(r.tu); err != nil {
			errChan <- err
		}
	} else {
		err = errors.New("Run Function should be given")
		errChan <- err
	}
	if r.tfunc.PostRun != nil {
		if err = r.tfunc.PostRun(r.tu); err != nil {
			errChan <- err
		}
	}
	ug.cache.Del(ts)
}
