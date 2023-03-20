package tezos

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"strconv"
	// "blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/contract"
	"blockwatch.cc/tzgo/tezos"
)

type receiver struct {
	log log.Logger
	src tezos.Address
	dst tezos.Address
	client Client
}

type Receiver struct {
	log log.Logger
	src tezos.Address
	dst tezos.Address
	Client *Client
}

func (r *Receiver) Subscribe(ctx context.Context, msgCh chan<- *chain.Message, opts chain.SubscribeOptions) (errCh <-chan error, err error) {
	r.Client.Contract = contract.NewContract(r.src, r.Client.Cl)
	r.Client.Ctx = ctx

	opts.Seq++

	_errCh := make(chan error)

	verifier, err := r.NewVerifier(int64(opts.Height) - 1) 

	if err != nil {
		_errCh <- err
		return _errCh, err 
	}

	go func() {
		defer close(_errCh)
		err := r.Client.MonitorBlock(int64(opts.Height), verifier)
		
		if err != nil {
			_errCh <- err 
		}

		fmt.Println("Printing from inside the receiver")
	}()
	
	return _errCh, nil
}

func NewReceiver(src, dst tezos.Address, urls []string, rawOpts json.RawMessage, l log.Logger) (chain.Receiver, error){
	var client *Client
	var err error

	if len(urls) == 0 {
		return nil, fmt.Errorf("Empty urls")
	}

	client, err = NewClient(urls[0], src, l)
	if err != nil {
		return nil, err
	}

	receiver := &Receiver{
		log: l,
		src: src,
		dst: dst,
		Client: client,
	}

	return receiver, nil
}

func (r *Receiver) NewVerifier(previousHeight int64) (vri IVerifier, err error) {
	header, err := r.Client.GetBlockHeaderByHeight(r.Client.Ctx, r.Client.Cl, previousHeight)
	if err != nil {
		return nil, err
	}

	fittness, err := strconv.ParseInt(string(header.Fitness[1].String()), 16, 64)
	if err != nil {
		return nil, err
	}

	chainIdHash, err := r.Client.Cl.GetChainId(r.Client.Ctx)
	if err != nil {
		return nil, err
	}

	id := chainIdHash.Uint32()

	if err != nil {
		return nil, err 
	} 

	vr := &Verifier{
		mu: sync.RWMutex{},
		next: header.Level,
		parentHash: header.Hash,
		parentFittness: fittness,
		chainID: id,
	}

	return vr, nil
}