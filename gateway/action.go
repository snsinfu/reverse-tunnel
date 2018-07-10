package main

import (
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

// Mount adds application routes to echo instance.
func Mount(e *echo.Echo) {
	e.GET("/tcp/:port", GetTCPPort)
	e.GET("/udp/:port", GetUDPPort)
	e.GET("/session/:id", GetSession)
}

// GetTCPPort handles /tcp/:port request. It starts listening on specified TCP
// port and notifies client connections to the requester via WebSocket.
func GetTCPPort(c echo.Context) error {
	port, err := parsePort(c.Param("port"))
	if err != nil {
		return err
	}

	key, err := extractAuthKey(c)
	if err != nil {
		return err
	}

	agent, err := NewTCPAgent(port, key)
	if err != nil {
		return err
	}

	return WebSocket(c, agent.Start)
}

// GetUDPPort handles /udp/:port request. It starts listening on specified UDP
// port and notifies client connections to the requester via WebSocket.
func GetUDPPort(c echo.Context) error {
	port, err := parsePort(c.Param("port"))
	if err != nil {
		return err
	}

	key, err := extractAuthKey(c)
	if err != nil {
		return err
	}

	agent, err := NewUDPAgent(port, key)
	if err != nil {
		return err
	}

	return WebSocket(c, agent.Start)
}

// GetSession handles /session/:id request. It starts tunneling TCP/UDP packets
// to the requester via WebSocket.
func GetSession(c echo.Context) error {
	id := c.Param("id")

	sess, err := AcquireSession(id)
	if err != nil {
		return err
	}

	return WebSocket(c, sess.Start)
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

// extractAuthKey extracts Bearer authorization key from request header.
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
