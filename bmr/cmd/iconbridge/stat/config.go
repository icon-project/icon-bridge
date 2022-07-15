package stat

/*
ExampleConfig:
"stat": {
    "verbose": false,
    "logging_interval":{
      "heartbeat":60,
      "system_metrics":20
    },
    "trigger":[{
      "measurement":"LoadAverage",
      "field":"LoadAvg1",
      "value":0.2,
      "sign":">="
    },
    {
      "measurement":"MemoryUsage",
      "field":"UsedPercent",
      "value":20,
      "sign":">="
    }
  ]
  }

1. Stat service does not run if "stat" key-value pair is not mentioned in config
2. Stat service runs with default config if "stat" key is present but value is empty i.e. "stat":{}
3. If fields Verbose, LoggingInterval, Trigger is not specified in config, then default value for those configs are used
4. For fields: heartbeat and system_metrics of logging_interval, following holds true for each
  	If unspecified, use default config value
	If specified and set to 0, disable logging
	If specified and >0 && < minimum, enable logging and set interval to minimum
	Else , enable logging and set interval to that provided in config

CHECK ensureConfig() for detail
*/
var (
	DefaultHeartBeatLoggingInterval     uint = 5 * 60 // 5 minutes
	DefaultSystemMetricsLoggingInterval uint = 5 * 60 // 5 minutes
	MinimumHeartBeatLoggingInterval     uint = 10     // 10 seconds
	MinimumSystemMetricsLoggingInterval uint = 10     // 10 seconds
)

var defaultConfig = StatConfig{
	Verbose:         false,
	LoggingInterval: &LoggingInterval{HeartBeat: &DefaultHeartBeatLoggingInterval, SystemMetrics: &DefaultSystemMetricsLoggingInterval},
	Trigger: []*Trigger{
		{Measurement: "LoadAverage", Field: "LoadAvg5", Value: 1.5, Sign: ">"},
		{Measurement: "MemoryUsage", Field: "UsedPercent", Value: 90, Sign: ">"},
		{Measurement: "DiskUsage", Field: "UsedPercent", Value: 90, Sign: ">"},
	},
}

type StatConfig struct {
	Verbose         bool             `json:"verbose,omitempy"`           // whether to display all fields or just the one used in trigger criteria
	LoggingInterval *LoggingInterval `json:"logging_interval,omitempty"` // check every X seconds
	Trigger         []*Trigger       `json:"trigger,omitempty"`          // defines threshold for alert to trigger
}

type LoggingInterval struct {
	HeartBeat     *uint `json:"heartbeat,omitempty"`
	SystemMetrics *uint `json:"system_metrics,omitempty"`
}

type Trigger struct { // Trigger Criterion: if Memory.UsedPercent > 90, then trigger alert
	Measurement MeasurementType `json:"measurement"` // Measurement: eg Meemory
	Field       string          `json:"field"`       // A field of measurement eg UsedPercent
	Value       float64         `json:"value"`       // Threshold Value: eg 90
	Sign        string          `json:"sign"`        // Relational Operator: eg >
}

type MeasurementType string

const (
	LOADAVERAGE MeasurementType = "LoadAverage"
	MEMORYUSAGE MeasurementType = "MemoryUsage"
	DISKUSAGE   MeasurementType = "DiskUsage"
)
