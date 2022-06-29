package executor

import (
	"math/big"
	"time"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chain"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/hmny"
	"github.com/icon-project/icon-bridge/cmd/endpoint/chain/icon"
	"github.com/icon-project/icon-bridge/common/errors"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	PRIVKEYPOS = 0
	PUBKEYPOS  = 1
)

type args struct {
	log             log.Logger
	clientsPerChain map[chain.ChainType]chain.ChainAPI
	godKeysPerChain map[chain.ChainType][2]string
}

func newArgs(l log.Logger, clientsPerChain map[chain.ChainType]*chain.ChainConfig, godKeysPerChain map[chain.ChainType][2]string) (t *args, err error) {
	tu := &args{log: l,
		clientsPerChain: map[chain.ChainType]chain.ChainAPI{},
		godKeysPerChain: godKeysPerChain,
	}
	for name, cfg := range clientsPerChain {
		if name == chain.HMNY {
			tu.clientsPerChain[name], err = hmny.NewApi(l, cfg)
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
				return nil, err
			}
		} else if name == chain.ICON {
			tu.clientsPerChain[name], err = icon.NewApi(l, cfg)
			if err != nil {
				err = errors.Wrap(err, "HMNY Err: ")
			}
		} else {
			return nil, errors.New("Unknown Chain Type supplied from config: " + string(name))
		}
	}
	return tu, nil
}

type callBackFunc func(args *args) error

var DefaultCallBacks = map[string]callBackFunc{
	"Demo": func(args *args) error {
		// fund demo wallets
		args.log.Info("Starting demo...")
		ienv, ok := args.clientsPerChain[chain.ICON]
		if !ok {
			return errors.New("Icon client not found")
		}
		henv, ok := args.clientsPerChain[chain.HMNY]
		if !ok {
			return errors.New("Hmny client not found")
		}
		igod, ok := args.godKeysPerChain[chain.ICON]
		if !ok {
			return errors.New("God Keys not found for ICON")
		}
		hgod, ok := args.godKeysPerChain[chain.HMNY]
		if !ok {
			return errors.New("God keys not found for Hmy")
		}
		tmp, err := ienv.GetKeyPairs(1)
		if err != nil {
			return errors.New("Couldn't create demo account for icon")
		}
		iDemo := tmp[0]
		tmp, err = henv.GetKeyPairs(1)
		if err != nil {
			return errors.New("Couldn't create demo account for hmny")
		}
		hDemo := tmp[0]
		args.log.Info("Creating Demo Icon Account ", iDemo)
		args.log.Info("Creating Demo Hmy Account ", hDemo)
		showBalance := func(log log.Logger, env chain.ChainAPI, addr string, tokens []chain.TokenType) error {
			factor := new(big.Int)
			factor.SetString("10000000000000000", 10)
			for _, token := range tokens {
				if amt, err := env.GetCoinBalance(addr, token); err != nil {
					return err
				} else {
					log.Infof("%v: %v", token, amt.Div(amt, factor).String())
				}
			}
			return nil
		}
		args.log.Info("Funding Demo Wallets ")
		amt := new(big.Int)
		amt.SetString("250000000000000000000", 10)
		_, err = ienv.Transfer(&chain.RequestParam{FromChain: chain.ICON, ToChain: chain.ICON, SenderKey: igod[PRIVKEYPOS], FromAddress: igod[PUBKEYPOS], ToAddress: iDemo[PUBKEYPOS], Amount: *amt, Token: chain.ICXToken})
		if err != nil {
			return err
		}
		amt = new(big.Int)
		amt.SetString("10000000000000000000", 10)
		_, err = ienv.Transfer(&chain.RequestParam{FromChain: chain.ICON, ToChain: chain.ICON, SenderKey: igod[PRIVKEYPOS], FromAddress: igod[PUBKEYPOS], ToAddress: iDemo[PUBKEYPOS], Amount: *amt, Token: chain.IRC2Token})
		if err != nil {
			return err
		}
		amt = new(big.Int)
		amt.SetString("10000000000000000000", 10)
		_, err = henv.Transfer(&chain.RequestParam{FromChain: chain.HMNY, ToChain: chain.HMNY, SenderKey: hgod[PRIVKEYPOS], FromAddress: hgod[PUBKEYPOS], ToAddress: hDemo[PUBKEYPOS], Amount: *amt, Token: chain.ONEToken})
		if err != nil {
			return err
		}
		amt = new(big.Int)
		amt.SetString("10000000000000000000", 10)
		_, err = henv.Transfer(&chain.RequestParam{FromChain: chain.HMNY, ToChain: chain.HMNY, SenderKey: hgod[PRIVKEYPOS], FromAddress: hgod[PUBKEYPOS], ToAddress: hDemo[PUBKEYPOS], Amount: *amt, Token: chain.ERC20Token})
		if err != nil {
			return err
		}
		args.log.Info("Done funding")
		time.Sleep(time.Second * 10)
		// args.log.Info("ICON:  ")
		// if err := showBalance(args.log, ienv, iDemo[PUBKEYPOS], []chain.TokenType{chain.ICXToken, chain.IRC2Token, chain.ONEToken}); err != nil {
		// 	return err
		// }
		// args.log.Info("HMNY:   ")
		// if err := showBalance(args.log, henv, hDemo[PUBKEYPOS], []chain.TokenType{chain.ONEToken, chain.ERC20Token, chain.ICXToken}); err != nil {
		// 	return err
		// }

		args.log.Info("Transfer Native ICX to HMY")
		amt = new(big.Int)
		amt.SetString("2000000000000000000", 10)
		_, err = ienv.Transfer(&chain.RequestParam{FromChain: chain.ICON, ToChain: chain.HMNY, SenderKey: iDemo[PRIVKEYPOS], FromAddress: iDemo[PUBKEYPOS], ToAddress: *henv.GetBTPAddress(hDemo[PUBKEYPOS]), Amount: *amt, Token: chain.ICXToken})
		if err != nil {
			return err
		}
		args.log.Info("Transfer Native ONE to ICX")
		amt = new(big.Int)
		amt.SetString("2000000000000000000", 10)
		_, err = henv.Transfer(&chain.RequestParam{FromChain: chain.HMNY, ToChain: chain.ICON, SenderKey: hDemo[PRIVKEYPOS], FromAddress: hDemo[PUBKEYPOS], ToAddress: *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), Amount: *amt, Token: chain.ONEToken})
		if err != nil {
			return err
		}
		args.log.Info("Approve")
		time.Sleep(time.Second * 10)

		amt = new(big.Int)
		amt.SetString("100000000000000000000000", 10)
		_, err = ienv.Approve(iDemo[PRIVKEYPOS], *amt)
		if err != nil {
			return err
		}
		amt = new(big.Int)
		amt.SetString("100000000000000000000000", 10)
		_, err = henv.Approve(hDemo[PRIVKEYPOS], *amt)
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)

		args.log.Info("Transfer Wrapped")
		amt = new(big.Int)
		amt.SetString("1000000000000000000", 10)
		_, err = henv.Transfer(&chain.RequestParam{FromChain: chain.HMNY, ToChain: chain.ICON, SenderKey: hDemo[PRIVKEYPOS], FromAddress: hDemo[PUBKEYPOS], ToAddress: *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), Amount: *amt, Token: chain.ICXToken})
		if err != nil {
			return err
		}
		amt = new(big.Int)
		amt.SetString("1000000000000000000", 10)
		_, err = ienv.Transfer(&chain.RequestParam{FromChain: chain.ICON, ToChain: chain.HMNY, SenderKey: iDemo[PRIVKEYPOS], FromAddress: iDemo[PUBKEYPOS], ToAddress: *henv.GetBTPAddress(hDemo[PUBKEYPOS]), Amount: *amt, Token: chain.ONEToken})
		if err != nil {
			return err
		}
		time.Sleep(10 * time.Second)

		args.log.Info("Transfer Irc2 to HMY")
		amt = new(big.Int)
		amt.SetString("1000000000000000000", 10)
		_, err = ienv.Transfer(&chain.RequestParam{FromChain: chain.ICON, ToChain: chain.HMNY, SenderKey: iDemo[PRIVKEYPOS], FromAddress: iDemo[PUBKEYPOS], ToAddress: *henv.GetBTPAddress(hDemo[PUBKEYPOS]), Amount: *amt, Token: chain.IRC2Token})
		if err != nil {
			return err
		}
		args.log.Info("Transfer Erc20 to ICon")
		amt = new(big.Int)
		amt.SetString("1000000000000000000", 10)
		_, err = henv.Transfer(&chain.RequestParam{FromChain: chain.HMNY, ToChain: chain.ICON, SenderKey: hDemo[PRIVKEYPOS], FromAddress: hDemo[PUBKEYPOS], ToAddress: *ienv.GetBTPAddress(iDemo[PUBKEYPOS]), Amount: *amt, Token: chain.ERC20Token})
		if err != nil {
			return err
		}
		time.Sleep(15 * time.Second)
		args.log.Info("ICON:  ")
		if err := showBalance(args.log, ienv, iDemo[PUBKEYPOS], []chain.TokenType{chain.ICXToken, chain.IRC2Token, chain.ONEToken}); err != nil {
			return err
		}
		args.log.Info("HMNY:   ")
		if err := showBalance(args.log, henv, hDemo[PUBKEYPOS], []chain.TokenType{chain.ONEToken, chain.ERC20Token, chain.ICXToken}); err != nil {
			return err
		}
		args.log.Info("Done")
		return nil
	},
}

/*
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
			if _, _, err := ienv.Client.TransferCoinCrossChain(
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
			if _, _, err := henv.Client.TransferCoinCrossChain(
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
			if _, _, amt, err := ienv.Client.ApproveContractToAccessCrossCoin(ienv.AccountsKeys[0][PRIVKEYPOS], *allowMount); err != nil {
				return errors.Wrap(err, " Approve ICON ")
			} else {
				tu.Logger().Info("ICON Allowed Amount ", amt.String())
			}
			time.Sleep(time.Second * 5)
			tu.Logger().Info("Step 5. Approve HMNY BSHCore to access ")
			allowMount = new(big.Int)
			allowMount.SetString("100000000000000000000000", 10)
			if _, _, amt, err := henv.Client.ApproveContractToAccessCrossCoin(henv.AccountsKeys[0][PRIVKEYPOS], *allowMount); err != nil {
				return errors.Wrap(err, " Approve HMNY ")
			} else {
				tu.Logger().Info("HMNY Allowed Amount ", amt.String())
			}
			time.Sleep(time.Second * 5)
			tu.Logger().Info("Step 6. Transfer Wrapped ICX (HMNY -> ICON):")
			h2i_wrapped_ICX_transfer_amount := new(big.Int)
			h2i_wrapped_ICX_transfer_amount.SetString("1000000000000000000", 10)
			if _, _, err := henv.Client.TransferWrappedCoinCrossChain(
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
			if _, _, err := ienv.Client.TransferWrappedCoinCrossChain(
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
			if _, _, _, _, err := ienv.Client.TransferEthTokenCrossChain(
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
			if _, _, _, _, err := henv.Client.TransferEthTokenCrossChain(
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
			if _, _, err := ienv.Client.TransferCoin(ienv.GodKeys[PRIVKEYPOS], *icx_target, ienv.AccountsKeys[0][PUBKEYPOS]); err != nil {
				errors.Wrap(err, "Transfer ICX ")
				tu.Logger().Error(err)
				return err
			}
			irc2_target := new(big.Int)
			irc2_target.SetString("10000000000000000000", 10)
			if _, _, err := ienv.Client.TransferEthToken(ienv.GodKeys[PRIVKEYPOS], *irc2_target, ienv.AccountsKeys[0][PUBKEYPOS]); err != nil {
				errors.Wrap(err, "Transfer IRC2 ")
				tu.Logger().Error(err)
				return err
			}

			one_target := new(big.Int)
			one_target.SetString("10000000000000000000", 10)
			if _, _, err := henv.Client.TransferCoin(henv.GodKeys[PRIVKEYPOS], *one_target, henv.AccountsKeys[0][PUBKEYPOS]); err != nil {
				err = errors.Wrap(err, "Transfer One ")
				tu.Logger().Error(err)
				return err
			}
			erc20_target := new(big.Int)
			erc20_target.SetString("10000000000000000000", 10)
			if _, _, err := henv.Client.TransferEthToken(henv.GodKeys[PRIVKEYPOS], *erc20_target, henv.AccountsKeys[0][PUBKEYPOS]); err != nil {
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
*/
