package icon

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icon-project/icon-bridge/common/log"
)

const (
	TransferStartSignature    = "TransferStart(Address,str,int,bytes)"
	TransferEndSignature      = "TransferEnd(Address,int,int,bytes)"
	TransferReceivedSignature = "TransferReceived(str,Address,int,bytes)"
	btp_icon_nativecoin_bsh   = "cxe51b3bd6326df70806419707d0f77eda86711bf5"
	btp_icon_token_bsh        = "cx5924a147ae30091ed9c6fe0c153ef77de4132902"
)

// btp_icon_nativecoin_bsh

type notification struct {
	event     *EventNotification
	signature string
}

func NewReceiver() {
	url := "http://localhost:9080/api/v3/default"
	l := log.New()
	log.SetGlobalLogger(l)
	contractAddress := []string{btp_icon_nativecoin_bsh, btp_icon_token_bsh}
	streamChan := make(chan *notification)

	go func() {
		for {
			select {
			case n := <-streamChan:
				fmt.Println(n.event.Hash, n.signature, *n.event)
			}
		}
	}()

	for _, cAddr := range contractAddress {
		go getEventStreamForIcon(url, l, TransferStartSignature, cAddr, 0x438d, streamChan)
		go getEventStreamForIcon(url, l, TransferEndSignature, cAddr, 0x438d, streamChan)
		go getEventStreamForIcon(url, l, TransferReceivedSignature, cAddr, 0x438d, streamChan)
	}
	fmt.Println("Wait")
	time.Sleep(time.Hour)
}

func getEventStreamForIcon(url string, l log.Logger, eventSignature string, contractAddress string, height int64, sink chan<- *notification) {
	l = l.WithFields(log.Fields{"Signature": eventSignature, "Contract": contractAddress})
	client := newClient(url, l)
	ef := EventFilter{
		Addr:      Address(contractAddress),
		Signature: eventSignature,
	}
	evtReq := &EventRequest{
		EventFilter: ef,
		Height:      NewHexInt(height),
		Logs:        HexInt("0x1"),
	}

	client.MonitorEvent(context.TODO(), evtReq,
		func(conn *websocket.Conn, v *EventNotification) error {
			sink <- &notification{event: v, signature: eventSignature}
			return nil
		}, func(c *websocket.Conn, err error) { fmt.Println("Exiting", err); return })
	return
}
