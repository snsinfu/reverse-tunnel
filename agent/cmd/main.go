package main

import (
	"fmt"
	"os"

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

const defaultConfigPath = "rtun.yml"

func main() {
	options, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}

	if err := run(options); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
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

	return agent.Start(conf)
}
