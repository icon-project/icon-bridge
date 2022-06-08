package stat

/*
ExampleConfig:
"stat": {
    "Verbose": false,
    "LoggingInterval":{
      "HeartBeat":60,
      "SystemMetrics":20
    },
    "Trigger":[{
      "Measurement":"LoadAverage",
      "Field":"LoadAvg1",
      "Value":0.2,
      "Sign":">="
    },
    {
      "Measurement":"MemoryUsage",
      "Field":"UsedPercent",
      "Value":20,
      "Sign":">="
    }
  ]
  }

1. Stat service does not run if "stat" key-value pair is not mentioned in config
2. Stat service runs with default config if "stat" key is present but value is empty i.e. "stat":{}
3. If fields Verbose, LoggingInterval, Trigger is not specified in config, then default value for those configs are used
4. If these fields are specified, then that specified field is used
  For example, using "LoggingInterval":{"HeartBeat":60} means HeartBeat is logged every 60 seconds while "SystemMetrics" is not logged at all
  For example, using "LoggingInterval":{} means neither HeartBeat nor "SystemMetrics" is logged
  For example, using "Trigger":[] means SystemMetrics does not have any trigger to create alert

CHECK ensureConfig() for detail
*/

var defaultConfig = StatConfig{
	LoggingInterval: &LoggingInterval{HeartBeat: 60 * 10, SystemMetrics: 60 * 10}, //every 10 minutes
	Trigger: []*Trigger{
		{Measurement: "LoadAverage", Field: "LoadAvg5", Value: 1.5, Sign: ">"},
		{Measurement: "MemoryUsage", Field: "UsedPercent", Value: 90, Sign: ">"},
		{Measurement: "DiskUsage", Field: "UsedPercent", Value: 90, Sign: ">"},
	},
}

const (
	MIN_HEARTBEAT_SEC     = 10 // 10 seconds
	MIN_SYSTEMMETRICS_SEC = 10 // 10 seconds
)

type StatConfig struct {
	Verbose         bool             `json:"Verbose,omitempy"`          // whether to display all fields or just the one used in trigger criteria
	LoggingInterval *LoggingInterval `json:"LoggingInterval,omitempty"` // check every X seconds
	Trigger         []*Trigger       `json:"Trigger,omitempty"`         // defines threshold for alert to trigger
}

type LoggingInterval struct {
	HeartBeat     int `json:"HeartBeat,omitempty"`
	SystemMetrics int `json:"SystemMetrics,omitempty"`
}

type Trigger struct { // Trigger Criterion: if Memory.UsedPercent > 90, then trigger alert
	Measurement MeasurementType `json:"Measurement"` // Measurement: eg Meemory
	Field       string          `json:"Field"`       // A field of measurement eg UsedPercent
	Value       float64         `json:"Value"`       // Threshold Value: eg 90
	Sign        string          `json:"Sign"`        // Relational Operator: eg >
}

type MeasurementType string

const (
	LOADAVERAGE MeasurementType = "LoadAverage"
	MEMORYUSAGE MeasurementType = "MemoryUsage"
	DISKUSAGE   MeasurementType = "DiskUsage"
)
