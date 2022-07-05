package stat

import (
	"context"
	"time"

	"github.com/icon-project/icon-bridge/common/log"
)

type statCollector struct {
	disable  bool
	cfg      *StatConfig
	verbose  bool
	log      log.Logger
	stopChan chan struct{}
}

type StatCollector interface {
	Start(ctx context.Context) error
	Stop()
}

func NewService(cfg *StatConfig, l log.Logger) (StatCollector, error) {
	s := &statCollector{cfg: cfg, log: l, stopChan: make(chan struct{}), disable: false, verbose: false}
	s.cfg = s.ensureConfig(cfg)
	s.verbose = s.cfg.Verbose
	return s, nil
}

func (s *statCollector) ensureConfig(cfg *StatConfig) *StatConfig {
	if cfg == nil { // if config is not provided, service is disabled; provide at least an empty config
		s.disable = true
		s.log.Warn("System Metrics Collector Service is disabled as configuration has not been provided in relay config file")
		return &defaultConfig
	}
	if cfg.LoggingInterval == nil {
		s.log.Info("Using default config for logging Interval")
		cfg.LoggingInterval = defaultConfig.LoggingInterval
	} else {
		if cfg.LoggingInterval.HeartBeat == nil {
			s.log.Infof("Using default config for HeartBeat logging Interval:  %d seconds", DefaultHeartBeatLoggingInterval)
			cfg.LoggingInterval.HeartBeat = &DefaultHeartBeatLoggingInterval
		} else if cfg.LoggingInterval.HeartBeat != nil && *cfg.LoggingInterval.HeartBeat > 0 && *cfg.LoggingInterval.HeartBeat < MinimumHeartBeatLoggingInterval {
			s.log.Infof("HeartBeat logging Interval should be at least %d seconds; Using this minimum value", MinimumHeartBeatLoggingInterval)
			cfg.LoggingInterval.HeartBeat = &MinimumHeartBeatLoggingInterval
		} else if cfg.LoggingInterval.HeartBeat != nil && *cfg.LoggingInterval.HeartBeat <= 0 {
			s.log.Info("HeartBeat Logging has been disabled")
			cfg.LoggingInterval.HeartBeat = nil
		} else {
			s.log.Infof("HeartBeat interval set from config is %d seconds", *cfg.LoggingInterval.HeartBeat)
		}
		if cfg.LoggingInterval.SystemMetrics == nil {
			s.log.Infof("Using default config for SystemMetrics logging Interval: %d seconds", DefaultSystemMetricsLoggingInterval)
			cfg.LoggingInterval.SystemMetrics = &DefaultSystemMetricsLoggingInterval
		} else if cfg.LoggingInterval.SystemMetrics != nil && *cfg.LoggingInterval.SystemMetrics > 0 && *cfg.LoggingInterval.SystemMetrics < MinimumSystemMetricsLoggingInterval {
			s.log.Infof("SystemMetrics logging Interval should be at least %d seconds; Using this minimum value", MinimumSystemMetricsLoggingInterval)
			cfg.LoggingInterval.SystemMetrics = &MinimumSystemMetricsLoggingInterval
		} else if cfg.LoggingInterval.SystemMetrics != nil && *cfg.LoggingInterval.SystemMetrics <= 0 {
			s.log.Info("System Metrics Logging has been disabled")
			cfg.LoggingInterval.SystemMetrics = nil
		} else {
			s.log.Infof("SystemMetrics interval set from config is %d seconds", *cfg.LoggingInterval.SystemMetrics)
		}
	}
	if cfg.Trigger == nil {
		s.log.Info("Using default config for trigger criteria")
		cfg.Trigger = defaultConfig.Trigger
	}
	return cfg
}

func (s *statCollector) Start(ctx context.Context) error {
	if s.disable {
		return nil
	}

	sysTicker := time.NewTicker(time.Duration(DefaultSystemMetricsLoggingInterval) * time.Second)
	heartTicker := time.NewTicker(time.Duration(DefaultHeartBeatLoggingInterval) * time.Second)
	if s.cfg.LoggingInterval.SystemMetrics != nil && *s.cfg.LoggingInterval.SystemMetrics > 0 {
		sysTicker = time.NewTicker(time.Duration(*s.cfg.LoggingInterval.SystemMetrics) * time.Second)
	}
	if s.cfg.LoggingInterval.HeartBeat != nil && *s.cfg.LoggingInterval.HeartBeat > 0 {
		heartTicker = time.NewTicker(time.Duration(*s.cfg.LoggingInterval.HeartBeat) * time.Second)
	}

	go func() {
		defer sysTicker.Stop()
		defer heartTicker.Stop()

		var err error
		var metMap map[string]interface{}
		for {
			select {
			case <-ctx.Done():
				return
			case <-sysTicker.C:
				if s.cfg.LoggingInterval.SystemMetrics == nil { // continue if config sysMetrics logging was disabled (nil in config)
					continue
				}
				metMap, err = getFilteredMetrics(s.cfg.Trigger, s.cfg.Verbose) // can send both result and error
				if metMap != nil && len(metMap) > 0 {
					s.log.WithFields(metMap).Warn("System Alert")
				}
				if err != nil {
					s.log.Error("getFilteredMetricsFunc; SysMetrics; Error ", err)
				}
			case <-heartTicker.C:
				if s.cfg.LoggingInterval.HeartBeat == nil {
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
			case <-s.stopChan:
				s.log.Warn("Stopping Service StatCollector")
				break
			}
		}
	}()

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
