package stat

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/icon-project/icon-bridge/common/log"
)


func TestStatService(t *testing.T) {
	const URL = "https://hooks.slack.com/services/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"
	l := log.New()
	log.SetGlobalLogger(l)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log.AddForwarder(&log.ForwarderConfig{Vendor: log.HookVendorSlack, Address: URL, Level: "info"})
	var h uint = 10
	//s := 20
	sv, _ := NewService(&StatConfig{LoggingInterval: &LoggingInterval{HeartBeat: &h, SystemMetrics: nil}, Trigger: nil}, l)
	fmt.Println("Starting")
	sv.Start(ctx)
	<-time.After(time.Second * 25)
	fmt.Println("Closing")
	sv.Stop()
}
