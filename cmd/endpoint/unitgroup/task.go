package unitgroup

import (
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/tenv"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

type tEnvTask struct {
	id           int64
	tu           tenv.TEnv
	isolateAddrs bool
	tfunc        TaskFunc
}

type TaskFunc struct {
	PreRun  func(tu tenv.TEnv) error
	Run     func(tu tenv.TEnv) error
	PostRun func(tu tenv.TEnv) error
}

type tEnvTaskCache struct {
	mem       map[int64]tEnvTask
	mu        sync.RWMutex
	lastAdded int64
}

const (
	PRIVKEYPOS = 0
	PUBKEYPOS  = 1
)

func (ch *tEnvTaskCache) Add(task tEnvTask) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	if _, ok := ch.mem[task.id]; !ok {
		ch.mem[task.id] = task
	} else {
		ch.mem[task.id+1] = task
	}
	ch.lastAdded = task.id
}

func (ch *tEnvTaskCache) GetNew(latestRead int64) (retList map[int64]tEnvTask, latestTs int64) {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	if ch.lastAdded <= latestRead {
		latestTs = ch.lastAdded
		return
	}
	retList = map[int64]tEnvTask{}
	for ts := range ch.mem {
		if ts > latestRead { // if new, add
			retList[ts] = ch.mem[ts]
		}
	}
	latestTs = ch.lastAdded
	return
}

func (ch *tEnvTaskCache) Del(key int64) {
	ch.mu.Lock()
	defer ch.mu.Unlock()
	if _, ok := ch.mem[key]; ok {
		delete(ch.mem, key)
	}
}

func (ug *unitgroup) RegisterTestUnit(numAddrsPerChain map[chain.ChainType]int, task TaskFunc, isolateAddrs bool) (err error) {
	for key, num := range numAddrsPerChain {
		if num < 1 {
			delete(numAddrsPerChain, key)
		}
	}
	accountsPerChain, err := ug.createAccounts(numAddrsPerChain)
	if err != nil {
		return
	}
	newCfg := map[chain.ChainType]*chain.ChainConfig{}
	newGodKeys := map[chain.ChainType][2]string{}
	for name := range numAddrsPerChain {
		cfg, ok := ug.cfgPerChain[name]
		if !ok {
			err = errors.New("ChainType not known")
			return
		}
		newCfg[name] = cfg
		pair, ok := ug.godKeysPerChain[name]
		if !ok {
			err = errors.New("God wallet for chain not found")
		}
		newGodKeys[name] = pair
	}
	now := time.Now().UnixNano()

	tu, err := tenv.New(ug.log.WithFields(log.Fields{"id": strconv.Itoa(int(now))}), newCfg, accountsPerChain, newGodKeys)
	if err != nil {
		return
	}

	utask := tEnvTask{
		id:           now,
		tu:           tu,
		tfunc:        task,
		isolateAddrs: isolateAddrs,
	}
	ug.cache.Add(utask)

	return
}

var DefaultTaskFunctions = map[string]TaskFunc{
	"DemoTransaction": {
		Run: func(tu tenv.TEnv) error {
			showIconBalance := func(ienv *chain.EnvVariables) error {
				tu.Logger().Info("ICON Balance ++++++++++++++")
				if amt, err := ienv.Client.GetCoinBalance(ienv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					return err
				} else {
					tu.Logger().Info("ICX ", amt.String())
				}
				if amt, err := ienv.Client.GetEthToken(ienv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					return err
				} else {
					tu.Logger().Info("Eth ", amt.String())
				}
				if amt, err := ienv.Client.GetWrappedCoin(ienv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					return err
				} else {
					tu.Logger().Info("WrappedONE ", amt.String())
				}
				return nil
			}
			showHmnyBalance := func(henv *chain.EnvVariables) error {
				tu.Logger().Info("HMNY Balance +++++++++++++")
				if amt, err := henv.Client.GetCoinBalance(henv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					return err
				} else {
					tu.Logger().Info("ONE ", amt.String())
				}

				if amt, err := henv.Client.GetEthToken(henv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					return err
				} else {
					tu.Logger().Info("Eth ", amt.String())
				}
				if amt, err := henv.Client.GetWrappedCoin(henv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					return err
				} else {
					tu.Logger().Info("WrappedICX ", amt.String())
				}
				return nil
			}
			ienv, err := tu.GetEnvVariables(chain.ICON)
			if err != nil {
				return err
			}
			henv, err := tu.GetEnvVariables(chain.HMNY)
			if err != nil {
				return err
			}
			time.Sleep(time.Second * 5)
			tu.Logger().Info("Step 2. Transfer Native ICX (ICON -> HMNY): ")
			i2h_nativecoin_transfer_amount := new(big.Int)
			i2h_nativecoin_transfer_amount.SetString("2000000000000000000", 10)
			if _, err := ienv.Client.TransferCoinCrossChain(
				ienv.AccountsKeys[0][PRIVKEYPOS],
				*i2h_nativecoin_transfer_amount,
				*henv.Client.GetBTPAddress(henv.AccountsKeys[0][PUBKEYPOS])); err != nil {
				return errors.Wrap(err, "Transfer ICX to HMNY ")
			}
			time.Sleep(time.Second * 5)
			showIconBalance(ienv)
			showHmnyBalance(henv)

			tu.Logger().Info("Step 3. Transfer Native ONE (HMNY -> ICON): ")
			h2i_nativecoin_transfer_amount := new(big.Int)
			h2i_nativecoin_transfer_amount.SetString("2000000000000000000", 10)
			rxAddr := *ienv.Client.GetBTPAddress(ienv.AccountsKeys[0][PUBKEYPOS])
			if _, err := henv.Client.TransferCoinCrossChain(
				henv.AccountsKeys[0][PRIVKEYPOS],
				*h2i_nativecoin_transfer_amount, rxAddr); err != nil {
				return errors.Wrap(err, "Transfer ONE to ICON ")
			}
			time.Sleep(time.Second * 5)
			showIconBalance(ienv)
			showHmnyBalance(henv)

			tu.Logger().Info("Step 4. Approve ICON NativeCoinBSH ")
			allowMount := new(big.Int)
			allowMount.SetString("100000000000000000000000", 10)
			if _, amt, err := ienv.Client.ApproveContractToAccessCrossCoin(ienv.AccountsKeys[0][PRIVKEYPOS], *allowMount); err != nil {
				return errors.Wrap(err, " Approve ICON ")
			} else {
				tu.Logger().Info("ICON Allowed Amount ", amt.String())
			}
			time.Sleep(time.Second * 5)
			tu.Logger().Info("Step 5. Approve HMNY BSHCore to access ")
			allowMount = new(big.Int)
			allowMount.SetString("100000000000000000000000", 10)
			if _, amt, err := henv.Client.ApproveContractToAccessCrossCoin(henv.AccountsKeys[0][PRIVKEYPOS], *allowMount); err != nil {
				return errors.Wrap(err, " Approve HMNY ")
			} else {
				tu.Logger().Info("HMNY Allowed Amount ", amt.String())
			}
			time.Sleep(time.Second * 5)
			tu.Logger().Info("Step 6. Transfer Wrapped ICX (HMNY -> ICON):")
			h2i_wrapped_ICX_transfer_amount := new(big.Int)
			h2i_wrapped_ICX_transfer_amount.SetString("1000000000000000000", 10)
			if _, err := henv.Client.TransferWrappedCoinCrossChain(
				henv.AccountsKeys[0][PRIVKEYPOS],
				*h2i_wrapped_ICX_transfer_amount,
				*ienv.Client.GetBTPAddress(ienv.AccountsKeys[0][PUBKEYPOS]),
			); err != nil {
				return errors.Wrap(err, " Transfer Wrapped ICX ")
			}
			time.Sleep(time.Second * 5)
			showIconBalance(ienv)
			showHmnyBalance(henv)

			tu.Logger().Info("Step 7. Transfer Wrapped ONE (ICON -> HMNY):")
			i2h_wrapped_ONE_transfer_amount := new(big.Int)
			i2h_wrapped_ONE_transfer_amount.SetString("1000000000000000000", 10)
			if _, err := ienv.Client.TransferWrappedCoinCrossChain(
				ienv.AccountsKeys[0][PRIVKEYPOS],
				*i2h_wrapped_ONE_transfer_amount,
				*henv.Client.GetBTPAddress(henv.AccountsKeys[0][PUBKEYPOS]),
			); err != nil {
				return errors.Wrap(err, " Transfer Wrapped ONE ")
			}
			time.Sleep(time.Second * 5)
			showIconBalance(ienv)
			showHmnyBalance(henv)

			tu.Logger().Info("Step 8. Transfer irc2.ETH (ICON -> HMNY):")
			i2h_irc2_ETH_transfer_amount := new(big.Int)
			i2h_irc2_ETH_transfer_amount.SetString("1000000000000000000", 10)
			if _, _, err := ienv.Client.TransferEthTokenCrossChain(
				ienv.AccountsKeys[0][PRIVKEYPOS],
				*i2h_irc2_ETH_transfer_amount,
				*henv.Client.GetBTPAddress(henv.AccountsKeys[0][PUBKEYPOS])); err != nil {
				return errors.Wrap(err, " Transfer irc2.ETH to HMNY ")
			}
			time.Sleep(time.Second * 5)
			showIconBalance(ienv)
			showHmnyBalance(henv)

			tu.Logger().Info("Step 9. Transfer erc20.ETH (HMNY -> ICON):")
			h2i_erc20_ETH_transfer_amount := new(big.Int)
			h2i_erc20_ETH_transfer_amount.SetString("1000000000000000000", 10)
			if _, _, err := henv.Client.TransferEthTokenCrossChain(
				henv.AccountsKeys[0][PRIVKEYPOS],
				*h2i_erc20_ETH_transfer_amount,
				*ienv.Client.GetBTPAddress(ienv.AccountsKeys[0][PUBKEYPOS])); err != nil {
				return errors.Wrap(err, " Transfer erc20.ETH to ICON ")
			}
			time.Sleep(time.Second * 15)
			showIconBalance(ienv)
			showHmnyBalance(henv)

			tu.Logger().Info("Step 10: DONE")

			return nil
		},
		PreRun: func(tu tenv.TEnv) (err error) {
			tu.Logger().Info("Starting test unit to show demo transactions")
			ienv, err := tu.GetEnvVariables(chain.ICON)
			if err != nil {
				return err
			}
			henv, err := tu.GetEnvVariables(chain.HMNY)
			if err != nil {
				return err
			}
			if len(ienv.AccountsKeys) != 1 || len(henv.AccountsKeys) != 1 {
				return errors.New("This demo constrains a single demo wallet. Found > 1")
			}
			showIconBalance := func(ienv *chain.EnvVariables) error {
				if amt, err := ienv.Client.GetCoinBalance(ienv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					tu.Logger().Error(err)
					return err
				} else {
					tu.Logger().Info("ICX ", amt.String())
				}
				if amt, err := ienv.Client.GetEthToken(ienv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					tu.Logger().Error(err)
					return err
				} else {
					tu.Logger().Info("Eth ", amt.String())
				}
				if amt, err := ienv.Client.GetWrappedCoin(ienv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					tu.Logger().Error(err)
					return err
				} else {
					tu.Logger().Info("WrappedONE ", amt.String())
				}
				return nil
			}
			showHmnyBalance := func(henv *chain.EnvVariables) error {
				if amt, err := henv.Client.GetCoinBalance(henv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					tu.Logger().Error(err)
					return err
				} else {
					tu.Logger().Info("ONE ", amt.String())
				}

				if amt, err := henv.Client.GetEthToken(henv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					tu.Logger().Error(err)
					return err
				} else {
					tu.Logger().Info("Eth ", amt.String())
				}
				if amt, err := henv.Client.GetWrappedCoin(henv.AccountsKeys[0][PUBKEYPOS]); err != nil {
					tu.Logger().Error(err)
					return err
				} else {
					tu.Logger().Info("WrappedICX ", amt.String())
				}
				return nil
			}

			tu.Logger().Info("Step 1. Funding demo wallets ...")
			tu.Logger().Info("AC ", ienv.AccountsKeys)
			tu.Logger().Info("HC ", henv.AccountsKeys)

			icx_target := new(big.Int)
			icx_target.SetString("250000000000000000000", 10)
			if _, err := ienv.Client.TransferCoin(ienv.GodKeys[PRIVKEYPOS], *icx_target, ienv.AccountsKeys[0][PUBKEYPOS]); err != nil {
				errors.Wrap(err, "Transfer ICX ")
				tu.Logger().Error(err)
				return err
			}
			irc2_target := new(big.Int)
			irc2_target.SetString("10000000000000000000", 10)
			if _, err := ienv.Client.TransferEthToken(ienv.GodKeys[PRIVKEYPOS], *irc2_target, ienv.AccountsKeys[0][PUBKEYPOS]); err != nil {
				errors.Wrap(err, "Transfer IRC2 ")
				tu.Logger().Error(err)
				return err
			}

			one_target := new(big.Int)
			one_target.SetString("10000000000000000000", 10)
			if _, err := henv.Client.TransferCoin(henv.GodKeys[PRIVKEYPOS], *one_target, henv.AccountsKeys[0][PUBKEYPOS]); err != nil {
				err = errors.Wrap(err, "Transfer One ")
				tu.Logger().Error(err)
				return err
			}
			erc20_target := new(big.Int)
			erc20_target.SetString("10000000000000000000", 10)
			if _, err := henv.Client.TransferEthToken(henv.GodKeys[PRIVKEYPOS], *erc20_target, henv.AccountsKeys[0][PUBKEYPOS]); err != nil {
				errors.Wrap(err, "Transfer Erc20 ")
				tu.Logger().Error(err)
				return err
			}
			tu.Logger().Info("Showing new balance")
			showIconBalance(ienv)
			showHmnyBalance(henv)
			return
		},
	},
}
