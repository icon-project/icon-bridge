package main

import (
	"context"
	"encoding/json"
	"flag"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/icon-project/icon-bridge/cmd/iconbridge/relay"
	"github.com/icon-project/icon-bridge/cmd/iconbridge/stat"
	"github.com/icon-project/icon-bridge/common/config"
	"github.com/icon-project/icon-bridge/common/log"

	_ "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/bsc"
	_ "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/hmny"
	_ "github.com/icon-project/icon-bridge/cmd/iconbridge/chain/icon"
)

var (
	cfgFile string
)

func init() {
	flag.StringVar(&cfgFile, "config", "", "multi-relay config.json file")
}

type Config struct {
	config.FileConfig `json:",squash"`
	relay.Config      `json:",squash"`
	LogLevel          string               `json:"log_level"`
	ConsoleLevel      string               `json:"console_level"`
	LogWriter         *log.WriterConfig    `json:"log_writer,omitempty"`
	LogForwarder      *log.ForwarderConfig `json:"log_forwarder,omitempty"`
	StatConfig        *stat.StatConfig     `json:"stat_collector,omitempty"`
}

func main() {
	flag.Parse()

	cfg, err := loadConfig(cfgFile)
	if err != nil {
		log.Fatalf("failed to load config: file=%q, err=%q", cfgFile, err)
	}

	l := setLogger(cfg)
	relay, err := relay.NewMultiRelay(&cfg.Config, l)
	if err != nil {
		log.Fatalf("failed to create MultiRelay: %v", err)
	}
	scollector, err := stat.NewService(
		cfg.StatConfig,
		l.WithFields(log.Fields{
			log.FieldKeyService: "BMR-BSC",
		}))
	if err != nil {
		log.Error("failed to create StatCollector for MultiRelay: %v", err)
	}
	// for net/http/pprof
	go func() { http.ListenAndServe("0.0.0.0:6060", nil) }()
	runRelay(relay, scollector)
}

func runRelay(relay relay.Relay, sc stat.StatCollector) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigCh)
		cancel()
	}()

	go func() {
		select {
		case <-sigCh: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-sigCh // second signal, hard exit
		os.Exit(2)
	}()

	if err := sc.Start(ctx); err != nil {
		log.Error(err)
	}

	if err := relay.Start(ctx); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func loadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func setLogger(cfg *Config) log.Logger {
	l := log.New()
	log.SetGlobalLogger(l)
	stdlog.SetOutput(l.WriterLevel(log.WarnLevel))
	if cfg.LogWriter != nil {
		if cfg.LogWriter.Filename == "" {
			log.Fatalln("Empty LogWriterConfig filename!")
		}
		var lwCfg log.WriterConfig
		lwCfg = *cfg.LogWriter
		lwCfg.Filename = cfg.ResolveAbsolute(lwCfg.Filename)
		w, err := log.NewWriter(&lwCfg)
		if err != nil {
			log.Panicf("Fail to make writer err=%+v", err)
		}
		err = l.SetFileWriter(w)
		if err != nil {
			log.Panicf("Fail to set file l err=%+v", err)
		}
	}

	if lv, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.Panicf("Invalid log_level=%s", cfg.LogLevel)
	} else {
		l.SetLevel(lv)
	}
	if lv, err := log.ParseLevel(cfg.ConsoleLevel); err != nil {
		log.Panicf("Invalid console_level=%s", cfg.ConsoleLevel)
	} else {
		l.SetConsoleLevel(lv)
	}

	if cfg.LogForwarder != nil {
		if cfg.LogForwarder.Vendor == "" && cfg.LogForwarder.Address == "" {
			log.Fatalln("Empty LogForwarderConfig vendor and address!")
		}
		if err := log.AddForwarder(cfg.LogForwarder); err != nil {
			log.Fatalf("Invalid log_forwarder err:%+v", err)
		}
	}

	return l
}
