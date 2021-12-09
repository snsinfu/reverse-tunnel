package server

import (
	"fmt"
    "os"
    "os/signal"
    "syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/snsinfu/reverse-tunnel/config"
	"golang.org/x/crypto/acme/autocert"
)

// Start starts tunneling server with given configuration.
func Start(path string, required bool) error {

    conf, err := load(path, required)
    if err != nil {
        return err
    }

	e := echo.New()
	e.HideBanner = true

	// Enable TLS when Let's Encrypt domain is configured. Do not require the
	// control address port to be 443 because the port could be redirected.
	useTLS := (conf.LetsEncrypt.Domain != "")

	if useTLS {
		e.AutoTLSManager.Prompt = autocert.AcceptTOS
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(conf.LetsEncrypt.Domain)
		e.AutoTLSManager.Cache = autocert.DirCache(conf.LetsEncrypt.CacheDir)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	action := NewAction(conf)

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGHUP)
    go func() {
        for true {
            <-sigs
            fmt.Println("Reloading config " + path)
            conf, err := load(path, required)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
            action.Update(conf)
        }
    }()

	e.GET("/tcp/:port", action.GetTCPPort)
	e.GET("/udp/:port", action.GetUDPPort)
	e.GET("/session/:id", action.GetSession)

	if useTLS {
		return e.StartAutoTLS(conf.ControlAddress)
	}
	return e.Start(conf.ControlAddress)
}

func load(path string, required bool) (config.Server, error) {
	conf := config.ServerDefault
    err := config.Load(path, &conf)
    if required && err != nil {
        return conf, err
    }
	if err := conf.Check(); err != nil {
		return conf, fmt.Errorf("config error: %w", err)
	}
    return conf, nil
}
