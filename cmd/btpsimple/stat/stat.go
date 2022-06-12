package stat

import (
	"time"

	"github.com/icon-project/btp/common/log"
)

type statCollector struct {
	disable  bool
	cfg      *StatConfig
	verbose  bool
	log      log.Logger
	stopChan chan struct{}
}

type StatCollector interface {
	Start() error
	Stop()
}

func NewService(cfg *StatConfig, l log.Logger) StatCollector {
	s := &statCollector{cfg: cfg, log: l, stopChan: make(chan struct{}), disable: false, verbose: false}
	s.cfg = s.ensureConfig(cfg)
	s.verbose = s.cfg.Verbose
	return s
}

func (s *statCollector) ensureConfig(cfg *StatConfig) *StatConfig {
	if cfg == nil { // if config is not provided, service is disabled; provide at least an empty config
		s.disable = true
		s.log.Warn("System Metrics Collector Service is disabled as configuration has not been provided in relay config file")
		return &defaultConfig
	}
	if cfg.LoggingInterval == nil {
		s.log.Debug("Using default config for logging Interval")
		cfg.LoggingInterval = defaultConfig.LoggingInterval
	} else {
		if cfg.LoggingInterval.HeartBeat > 0 && cfg.LoggingInterval.HeartBeat < MIN_HEARTBEAT_SEC {
			s.log.Debugf("Interval should be greater than %d", MIN_HEARTBEAT_SEC)
			cfg.LoggingInterval.HeartBeat = MIN_HEARTBEAT_SEC
		} else if cfg.LoggingInterval.HeartBeat == 0 {
			s.log.Debug("HeartBeat interval is zero. Disabling it")
		}
		if cfg.LoggingInterval.SystemMetrics > 0 && cfg.LoggingInterval.SystemMetrics < MIN_SYSTEMMETRICS_SEC {
			s.log.Debugf("Interval should be greater than %d", MIN_SYSTEMMETRICS_SEC)
			cfg.LoggingInterval.SystemMetrics = MIN_SYSTEMMETRICS_SEC
		} else if cfg.LoggingInterval.SystemMetrics == 0 {
			s.log.Debug("SystemMetrics Logging interval is zero. Disabling it")
		}
	}
	if cfg.Trigger == nil {
		s.log.Debug("Using default config for trigger criteria")
		cfg.Trigger = defaultConfig.Trigger
	}
	return cfg
}

func (s *statCollector) Start() error {
	if s.disable {
		return nil
	}
	if s.cfg.LoggingInterval.SystemMetrics != 0 {
		sysTicker := time.NewTicker(time.Duration(s.cfg.LoggingInterval.SystemMetrics) * time.Second)
		go func() {
			defer sysTicker.Stop()

			var err error
			var metMap map[string]interface{}
			// The loop hasn't been exited despite error; The message is logged instead
			for {
				select {
				case <-sysTicker.C:
					if s.cfg.LoggingInterval.SystemMetrics == 0 {
						continue
					}
					metMap, err = getFilteredMetrics(s.cfg.Trigger, s.cfg.Verbose) // can send both result and error
					if metMap != nil && len(metMap) > 0 {
						s.log.WithFields(metMap).Warn("System Alert")
					}
					if err != nil {
						s.log.Error("getFilteredMetricsFunc; SysMetrics; Error ", err)
					}

				case <-s.stopChan:
					s.log.Debug("Stopping Service StatCollector")
					break
				}
			}
		}()
	}

	// TODO: StopChan has only been used one go-routine
	if s.cfg.LoggingInterval.HeartBeat != 0 {
		heartTicker := time.NewTicker(time.Duration(s.cfg.LoggingInterval.HeartBeat) * time.Second)
		go func() {
			defer heartTicker.Stop()
			var err error
			var metMap map[string]interface{}
			// The loop hasn't been exited despite error; The message is logged instead
			for {
				select {
				case <-heartTicker.C:
					if s.cfg.LoggingInterval.HeartBeat == 0 {
						continue
					}
					s.log.Info("HeartBeat Message")
					nt := getAlwaysTrigger(s.cfg.Trigger)
					metMap, err = getFilteredMetrics(nt, s.cfg.Verbose) // can send both result and error
					if metMap != nil && len(metMap) > 0 {
						s.log.WithFields(metMap).Info("System Info")
					}
					if err != nil {
						s.log.Error("getFilteredMetricsFunc; HeartBeat; Error ", err)
					}
				}
			}
		}()
	}
	return nil
}

func (s *statCollector) Stop() {
	if s.stopChan != nil {
		s.stopChan <- struct{}{}
	}
}

func getAlwaysTrigger(t []*Trigger) []*Trigger {
	if len(t) == 0 {
		t = defaultConfig.Trigger
	}
	n := make([]*Trigger, len(t))
	for i, v := range t {
		// Set TriggerThreshold: Value to zero to always trigger
		n[i] = &Trigger{Field: v.Field, Sign: v.Sign, Measurement: v.Measurement, Value: 0}
	}
	return n
}
