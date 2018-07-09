package main

import (
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// mount adds application routes to echo instance.
func mount(e *echo.Echo) {
	e.GET("/tcp/:port", GetTCPPort)
	e.GET("/udp/:port", GetUDPPort)
	e.GET("/session/:id", GetSession)
}

// GetTCPPort handles /tcp/:port request.
func GetTCPPort(c echo.Context) error {
	port, err := parsePort(c.Param("port"))
	if err != nil {
		return err
	}

	key, err := extractAuthKey(c)
	if err != nil {
		return err
	}

	if err := checkAccess(key, "tcp4", port); err != nil {
		return err
	}

	addr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}

	return WebSocket(func(ws *websocket.Conn) error {
		ln, err := net.ListenTCP("tcp4", addr)
		if err != nil {
			return err
		}
		defer ln.Close()
	})
}

// GetUDPPort handles /udp/:port request.
func GetUDPPort(c echo.Context) error {
}

// extractAuthKey extracts Bearer authorization key from request header.
func extractAuthKey(c echo.Context) (string, error) {
	const (
		authHeader = "Authorization"
		authScheme = "Bearer "
	)

	auth, ok := c.Request().Header[authHeader]
	if !ok {
		return "", ErrMissingKey
	}

	if !strings.HasPrefix(auth, authScheme) {
		return "", ErrInvalidKey
	}

	return strings.TrimPrefix(auth, authScheme), nil
}

// parsePort parses string as a port number, which must be between 1 and 65535.
func parsePort(s string) (int, error) {
	const (
		portMin = 1
		portMax = 65535
	)

	port, err := strconv.Atoi(s)
	if err != nil {
		return 0, ErrInvalidPort
	}

	if port < portMin || portMax < port {
		return port, ErrInvalidPort
	}

	return port, nil
}
