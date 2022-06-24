package backend

import "github.com/pkg/errors"

func (be *backend) Transfer(param *RequestParam) (txHash string, err error) {
	// what requestMatchingFns does this param match
	// If any add corresponding eventMatchingFn to processing Queue
	// then send the transaction
	// If error, pull watchers from queue, stop subscription provider go routine, then return error message
	// If result, a go routine will check queue and try to send to the receiver until timeout or end
	// If context is cancelled from the subscriber, stop trying to send

	addedIndex, err := be.findRequestMatchingCriteria(param)
	if err != nil {
		be.removeBulkFromEventListenerQueue(addedIndex)
		return "", err
	}
	reqApi := be.requestAPIPerChain[param.FromChain]
	if param.Token == ICXToken || param.Token == ONEToken { // Native Coin
		if param.FromChain == param.ToChain {
			txHash, err = reqApi.TransferCoin(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			txHash, err = reqApi.TransferCoinCrossChain(param.SenderKey, param.Amount, param.ToAddress)
		}
	} else if param.Token == IRC2Token || param.Token == ERC20Token { // EthToken
		if param.FromChain == param.ToChain {
			txHash, err = reqApi.TransferEthToken(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			_, txHash, err = reqApi.TransferEthTokenCrossChain(param.SenderKey, param.Amount, param.ToAddress)
		}
	} else if param.Token == ONEWrappedToken || param.Token == ICXWrappedToken { // WrappedToken
		if param.FromChain != param.ToChain {
			txHash, err = reqApi.TransferWrappedCoinCrossChain(param.SenderKey, param.Amount, param.ToAddress)
		} else {
			err = errors.New("Wrapped Coins can be transferred only cross chain")
		}
	} else {
		err = errors.New("Unrecognized token type ")
	}

	if err != nil {
		be.removeBulkFromEventListenerQueue(addedIndex)
	}

	return

}

// func (be *backend) Approve(param *TxnApproveParam) (txnHash string, err error) {
// 	reqApi := be.requestAPIPerChain[param.Chain]
// 	if (param.AccessToken == ONEToken && param.Chain == ICONChain) || (param.AccessToken == ICXToken && param.Chain == HMNYChain) {
// 		txnHash, _, err = reqApi.ApproveContractToAccessCrossCoin(param.OwnerKey, param.Amount)
// 	}
// 	err = errors.New("Invalid Request")
// 	return
// }

func (be *backend) RegisterCriterion(fn MatchingFn) error {
	be.availableCBMutex.Lock()
	defer be.availableCBMutex.Unlock()
	be.availableCBList = append(be.availableCBList, fn)
	return nil
}

func (be *backend) findRequestMatchingCriteria(req *RequestParam) ([]int, error) {
	be.availableCBMutex.RLock()
	defer be.availableCBMutex.RUnlock()
	addedIndex := []int{}
	for _, cb := range be.availableCBList {
		if cb.ReqFn(req) {
			wIndex, err := be.addToEventListenerQueue(cb.EventFn, req) // add eventFunction and inputRequest to another list while we're locking the current list
			if err != nil {
				return addedIndex, err
			}
			addedIndex = append(addedIndex, wIndex)
		}
	}
	return addedIndex, nil
}

func (be *backend) addToEventListenerQueue(fn EventMatchingFn, req *RequestParam) (writeIndex int, err error) {
	be.eventListenerMutex.Lock()
	defer be.eventListenerMutex.Unlock()
	if be.eventListenerCursor == EventListenerCapacity {
		err = errors.New("Event Listener Queue Full")
		return
	}
	be.eventListenerCursor += 1
	be.eventListenerList = append(be.eventListenerList, &eventListener{ip: req, fn: fn})
	writeIndex = int(be.eventListenerCursor) - 1
	return
}

func (be *backend) removeFromEventListenerQueue(arrayIndex int) {
	be.eventListenerMutex.Lock()
	defer be.eventListenerMutex.Unlock()
	be.eventListenerList[arrayIndex] = be.eventListenerList[be.eventListenerCursor-1] // fill empty array index with the last element
	be.eventListenerCursor -= 1                                                       // reduce cursor position
}

func (be *backend) removeBulkFromEventListenerQueue(toRemoveList []int) {
	be.eventListenerMutex.Lock()
	defer be.eventListenerMutex.Unlock()
	for _, arrayIndex := range toRemoveList {
		be.eventListenerList[arrayIndex] = be.eventListenerList[be.eventListenerCursor-1] // fill empty array index with the last element
		be.eventListenerCursor -= 1
	}
}

func (be *backend) checkIfLogMatchesListener(el *ReceiptEvent) (int, bool) {
	be.eventListenerMutex.RLock()
	be.eventListenerMutex.RUnlock()
	for evi, ev := range be.eventListenerList {
		if ev.fn(el, ev.ip) {
			return evi, true
		}
	}
	return 0, false
}

func (be *backend) Accounts(num int) error {
	return nil
}
