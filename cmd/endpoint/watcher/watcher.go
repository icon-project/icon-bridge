package watcher

import (
	"context"

	"github.com/icon-project/icon-bridge/common/log"

	"github.com/icon-project/icon-bridge/cmd/endpoint/chainAPI/chain"
)

type Watcher interface {
	Start(ctx context.Context) error
}

type watcher struct {
	log     log.Logger
	subChan <-chan *chain.SubscribedEvent
	errChan <-chan error
}

func New(log log.Logger, subChan <-chan *chain.SubscribedEvent, errChan <-chan error) (Watcher, error) {
	w := &watcher{log: log, subChan: subChan, errChan: errChan}
	return w, nil
}

func (w *watcher) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.log.Warn("Watcher; Context Cancelled")
				return
			case msg := <-w.subChan:
				w.log.Info(msg)
			case err := <-w.errChan:
				w.log.Error(err)
				return
			}
		}
	}()
	return nil
}

func (w *watcher) process() {

}
