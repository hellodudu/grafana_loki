package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/hellodudu/grafana_loki/loki_conn"
	"github.com/judwhite/go-svc/svc"
	"github.com/mreiferson/go-options"
)

func lokiConnFlagSet(opts *loki_conn.Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("loki_conn", flag.ExitOnError)

	flagSet.String("url", "loki:3100/api/prom/push", "loki's url")

	flagSet.Duration("interval", opts.Interval, "interval seconds of connecting to loki")

	return flagSet
}

type program struct {
	once     sync.Once
	lokiConn *loki_conn.LokiConn
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal("%s", err)
	}
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}

func (p *program) Start() error {
	opts := loki_conn.NewOptions()

	flagSet := lokiConnFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	var cfg map[string]interface{}
	configFlag := flagSet.Lookup("config")
	if configFlag != nil {
		configFile := configFlag.Value.String()
		if configFile != "" {
			_, err := toml.DecodeFile(configFile, &cfg)
			if err != nil {
				fmt.Errorf("failed to load config file %s - %s", configFile, err)
			}
		}
	}

	options.Resolve(opts, flagSet, cfg)
	lokiConn, err := loki_conn.New(opts)
	if err != nil {
		fmt.Errorf("failed to instantiate nsqlookupd", err)
	}
	p.lokiConn = lokiConn

	go func() {
		err := p.lokiConn.Main()
		if err != nil {
			p.Stop()
			os.Exit(1)
		}
	}()

	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		p.lokiConn.Exit()
	})
	return nil
}
