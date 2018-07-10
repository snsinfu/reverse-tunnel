package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	httpListenAddr = ":3000"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	Mount(e)

	e.Logger.Fatal(e.Start(httpListenAddr))
}
