package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/snsinfu/reverse-tunnel/agent"
	"github.com/snsinfu/reverse-tunnel/config"
)

const usage = `
Reverse tunnel agent

Usage:
  rtun [-f <config>]

Options:
  -h, --help   Show usage information and exit.
  -f <config>  Specify agent configuration file.
`

const (
	defaultConfigPath = "rtun.yml"
	cancelWait        = 3 * time.Second
)

func main() {
	options, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}

	if err := run(options); err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}

func run(options docopt.Opts) error {
	conf := config.AgentDefault

	if path, err := options.String("-f"); err == nil {
		if err := config.Load(path, &conf); err != nil {
			return err
		}
	} else {
		if err := config.Load(defaultConfigPath, &conf); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	ctx := withSignalCancel(context.Background())
	err := agent.Start(conf, ctx)

	if errors.Is(err, context.Canceled) {
		log.Print("waiting for agents to stop...")
		time.Sleep(cancelWait)
		return nil
	}

	return err
}

func withSignalCancel(ctx context.Context) context.Context {
	newCtx, cancel := context.WithCancel(ctx)

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		cancel()
	}()
	return newCtx
}
