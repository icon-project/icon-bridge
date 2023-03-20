package tezos

import (
	"github.com/icon-project/icon-bridge/cmd/iconbridge/chain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/icon-project/icon-bridge/common/log"
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
	client *Client
}

func (r *receiver) Subscribe(
	ctx context.Context, msgCh chan<- *chain.Message,
	opts chain.SubscribeOptions) (errCh <-chan error, err error) {

	r.client.Contract = contract.NewContract(r.src, r.client.Cl)
	r.client.Ctx = ctx

	opts.Seq++

	_errCh := make(chan error)

	verifier, err := r.NewVerifier(int64(opts.Height) - 1) 

	if err != nil {
		_errCh <- err
		return _errCh, err 
	}

	go func() {
		defer close(_errCh)
		err := r.client.MonitorBlock(int64(opts.Height), verifier)
		
		if err != nil {
			_errCh <- err 
		}

		fmt.Println("Printing from inside the receiver")
	}()
	
	return _errCh, nil
}

func NewReceiver(
	src, dst chain.BTPAddress, urls []string,
	rawOpts json.RawMessage, l log.Logger) (chain.Receiver, error){

	var client *Client
	var err error

	if len(urls) == 0 {
		return nil, fmt.Errorf("Empty urls")
	}

	srcAddr := tezos.MustParseAddress(src.String())

	dstAddr := tezos.MustParseAddress(dst.String())

	client, err = NewClient(urls[0], srcAddr, l)
	if err != nil {
		return nil, err
	}

	r := &receiver{
		log: l,
		src: srcAddr,
		dst: dstAddr,
		client: client,
	}

	return r, nil
}

func (r *receiver) NewVerifier(previousHeight int64) (vri IVerifier, err error) {
	header, err := r.client.GetBlockHeaderByHeight(r.client.Ctx, r.client.Cl, previousHeight)
	if err != nil {
		return nil, err
	}

	fittness, err := strconv.ParseInt(string(header.Fitness[1].String()), 16, 64)
	if err != nil {
		return nil, err
	}

	chainIdHash, err := r.client.Cl.GetChainId(r.client.Ctx)
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