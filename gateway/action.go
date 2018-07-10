package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/snsinfu/reverse-tunnel/gateway2/udp"
)

var (
	ErrInvalidPort = errors.New("invalid port")
	ErrMissingKey  = errors.New("missing key")
	ErrInvalidKey  = errors.New("invalid key")
)

// GetUDPPort handles GET /udp/:port request.
func GetUDPPort(c echo.Context) error {
	port, err := parsePort(c.Param("port"))
	if err != nil {
		return Error(c, err)
	}

	key, err := extractAuthKey(c)
	if err != nil {
		return Error(c, err)
	}

	listener, err := udp.GetListener(key, port)
	if err != nil {
		return Error(c, err)
	}

	return WebSocket(c, listener.Start)
}

// GetSession handles GET /session/:id request.
func GetSession(c echo.Context) error {
	id := c.Param("id")

	sess, err := udp.ResolveSession(id)
	if err != nil {
		return Error(c, err)
	}

	return WebSocket(c, sess.Start)
}

// parsePort parses s as a valid port number.
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

// extractAuthKey extracts Bearer token from request header.
func extractAuthKey(c echo.Context) (string, error) {
	const (
		authHeader = "Authorization"
		authScheme = "Bearer "
	)

	auth := c.Request().Header.Get(authHeader)
	if auth == "" {
		return "", ErrMissingKey
	}

	if !strings.HasPrefix(auth, authScheme) {
		return "", ErrInvalidKey
	}

	return strings.TrimPrefix(auth, authScheme), nil
}
