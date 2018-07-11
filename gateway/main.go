package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/snsinfu/reverse-tunnel/config"
)

// CLI usage in the docopt format.
const usage = `
Reverse tunnel gateway server

Usage:
  rt-gateway [-f <config>]

Options:
  -h, --help   Show usage information and exit.
  -f <config>  Specify gateway configuration file.
`

// Path to default configuration file. Can be nonexistent.
const defaultConfigPath = "gateway.yml"

// Default gateway configuration.
var defaultConfig = config.Gateway{
	ControlAddress: "127.0.0.1:9000",
}

func main() {
	options, _ := docopt.ParseDoc(usage)

	if err := run(options); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(options docopt.Opts) error {
	conf := defaultConfig

	if path, err := options.String("-f"); err == nil {
		if err := config.Load(path, &conf); err != nil {
			return err
		}
	} else {
		if err := config.Load(defaultConfigPath, &conf); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return startServer(conf)
}
