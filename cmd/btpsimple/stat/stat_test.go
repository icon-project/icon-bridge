package stat

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/icon-project/btp/common/log"
)

func TestSystemMetricsGetter(t *testing.T) {
	if met, err := getSystemMetrics(); err != nil {
		panic(err)
	} else {
		if metBytes, err := json.Marshal(met); err != nil {
			panic(err)
		} else {
			fmt.Println(string(metBytes))
		}
	}
}

func TestFilteredMetrics(t *testing.T) {
	if res, err := getFilteredMetrics(defaultConfig.Trigger, false); err != nil {
		panic(err)
	} else {
		fmt.Println("Result ", res)
	}
}

func TestStatService(t *testing.T) {
	const URL = "https://hooks.slack.com/services/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"
	l := log.New()
	log.SetGlobalLogger(l)
	log.AddForwarder(&log.ForwarderConfig{Vendor: log.HookVendorSlack, Address: URL, Level: "info"})
	sv := NewService(&StatConfig{LoggingInterval: &LoggingInterval{HeartBeat: 10, SystemMetrics: 20}, Trigger: nil}, l)
	fmt.Println("Starting")
	sv.Start()
	<-time.After(time.Minute)
	fmt.Println("Closing")
	sv.Stop()
}
