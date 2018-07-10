package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	controlAddress = "127.0.0.1:9000"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/udp/:port", GetUDPPort)
	e.GET("/session/:id", GetSession)

	e.Logger.Fatal(e.Start(controlAddress))
}
