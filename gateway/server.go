package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/snsinfu/reverse-tunnel/config"
)

func startServer(conf config.Gateway) error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	action := NewAction(conf)
	e.GET("/", action.GetHome)
	e.GET("/tcp/:port", action.GetTCPPort)
	e.GET("/udp/:port", action.GetUDPPort)
	e.GET("/tcp/session/:id", action.GetTCPSession)
	e.GET("/udp/session/:id", action.GetUDPSession)

	return e.Start(conf.ControlAddress)
}
