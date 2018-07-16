package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/snsinfu/reverse-tunnel/config"
)

// Start starts tunneling server with given configuration.
func Start(conf config.Server) error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	action := NewAction(conf)
	e.GET("/tcp/:port", action.GetTCPPort)
	e.GET("/udp/:port", action.GetUDPPort)
	e.GET("/session/:id", action.GetSession)

	return e.Start(conf.ControlAddress)
}
