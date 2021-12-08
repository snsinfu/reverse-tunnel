package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/snsinfu/reverse-tunnel/server"
)

const usage = `
Reverse tunnel gateway server

Usage:
  rtun-server [-f <config>]

Options:
  -h, --help   Show usage information and exit.
  -f <config>  Specify gateway configuration file.
`

const defaultConfigPath = "rtun-server.yml"

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
	if path, err := options.String("-f"); err == nil {
        return server.Start(path, true)
	}
    return server.Start(defaultConfigPath, false)
}
