package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/snsinfu/reverse-tunnel/config"
	"golang.org/x/crypto/acme/autocert"
)

// Start starts tunneling server with given configuration.
func Start(conf config.Server) error {
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
	e.GET("/tcp/:port", action.GetTCPPort)
	e.GET("/udp/:port", action.GetUDPPort)
	e.GET("/session/:id", action.GetSession)

	if useTLS {
		return e.StartAutoTLS(conf.ControlAddress)
	}
	return e.Start(conf.ControlAddress)
}
