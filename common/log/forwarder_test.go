package log

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestForwarder(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	args := []struct {
		vendor string
		addr   string
	}{
		{HookVendorFluentd, "tcp://ci.arch.iconloop.com:24224"},
		{HookVendorLogstash, "tcp://ci.arch.iconloop.com:5045"},
	}
	c := &ForwarderConfig{
		Level: "info",
	}
	for _, arg := range args {
		t.Run(arg.vendor, func(t *testing.T) {
			c.Vendor = arg.vendor
			c.Address = arg.addr
			err := AddForwarder(c)
			assert.NoError(t, err, "error on AddForwarder", c)
			GlobalLogger().Println("TestForwarder", c.Vendor, time.Now())
		})
	}
}

func TestSlackHook(t *testing.T) {
	const URL = "https://hooks.slack.com/services/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"
	sh, _ := NewSlackClient(URL, logrus.AllLevels)
	l := New()
	SetGlobalLogger(l)
	globalLogger.addHook(sh)
	for i := 0; i < 2; i++ {
		WithFields(Fields{"Status": "Check " + fmt.Sprint(i)}).Error("Hello Again!!")
	}
	fmt.Println("Done")
	<-time.After(time.Second * 10)
}

//bind to /fluentd/etc/fluent.conf
var fluentdConfig = `
<source>
  @type  forward
  port   24224
  source_address_key host
</source>

#<filter **>
#  @type stdout
#</filter>

<match *.**>
  @type elasticsearch
  host elasticsearch
  port 9200
  logstash_format true
  logstash_prefix fluentd
  include_tag_key true
</match>
`

//bind to /usr/share/logstash/pipeline/logstash.conf
var logstashConfig = `
input {
  tcp {
    port => 9600
    codec => json_lines
  }
}

output {
  elasticsearch { hosts => ["http://elasticsearch:9200"] }
#  stdout { codec => rubydebug }
}
`
