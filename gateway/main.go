package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/snsinfu/reverse-tunnel/config"
)

const usage = `
Reverse tunnel gateway server

Usage:
  rt-gateway [-f <config>]

Options:
  -h, --help   Show usage information and exit.
  -f <config>  Load server configuration from YAML file.
`

var defaultConf = config.GatewayConfig{
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
	conf := defaultConf

	if path, err := options.String("-f"); err == nil {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := config.Load(file, &conf); err != nil {
			return err
		}
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	Mount(e)

	return e.Start(conf.ControlAddress)
}
