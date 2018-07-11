package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/snsinfu/reverse-tunnel/config"
	"github.com/snsinfu/reverse-tunnel/gateway/tcp"
	"github.com/snsinfu/reverse-tunnel/gateway/udp"
)

var (
	ErrInvalidPort      = errors.New("invalid port number")
	ErrMissingKey       = errors.New("missing auth key")
	ErrInvalidKey       = errors.New("invalid auth key")
	ErrInvalidSessionID = errors.New("invalid session Id")
)

// Action interprets request and drives domain logic.
type Action struct {
	tcp tcp.Service
	udp udp.Service
}

// NewAction returns an Action object with given configuration.
func NewAction(conf config.Gateway) *Action {
	return &Action{
		tcp: tcp.NewService(conf),
		udp: udp.NewService(conf),
	}
}

// GetHome handles GET / request.
func (action *Action) GetHome(c echo.Context) error {
	return Data(c, map[string]interface{}{"version": 0})
}

// GetTCPPort handles GET /tcp/:port request.
func (action *Action) GetTCPPort(c echo.Context) error {
	port, err := parsePort(c.Param("port"))
	if err != nil {
		return Error(c, err)
	}

	key, err := extractAuthKey(c)
	if err != nil {
		return Error(c, err)
	}

	binder, err := action.tcp.NewBinder(key, port)
	if err != nil {
		return Error(c, err)
	}

	return WebSocket(c, binder.Start)
}

// GetUDPPort handles GET /udp/:port request.
func (action *Action) GetUDPPort(c echo.Context) error {
	port, err := parsePort(c.Param("port"))
	if err != nil {
		return Error(c, err)
	}

	key, err := extractAuthKey(c)
	if err != nil {
		return Error(c, err)
	}

	binder, err := action.udp.NewBinder(key, port)
	if err != nil {
		return Error(c, err)
	}

	return WebSocket(c, binder.Start)
}

// GetTCPSession handles GET /tcp/session/:id request.
func (action *Action) GetTCPSession(c echo.Context) error {
	id := c.Param("id")

	sess := action.tcp.ResolveSession(id)
	if sess == nil {
		return Error(c, ErrInvalidSessionID)
	}

	return WebSocket(c, sess.Start)
}

// GetUDPSession handles GET /udp/session/:id request.
func (action *Action) GetUDPSession(c echo.Context) error {
	id := c.Param("id")

	sess := action.udp.ResolveSession(id)
	if sess == nil {
		return Error(c, ErrInvalidSessionID)
	}

	return WebSocket(c, sess.Start)
}

// parsePort parses s as a port number.
func parsePort(s string) (int, error) {
	const (
		portMin = 1
		portMax = 65535
	)

	port, err := strconv.Atoi(s)
	if err != nil {
		return port, ErrInvalidPort
	}

	if port < portMin || portMax < port {
		return port, ErrInvalidPort
	}

	return port, nil
}

// extractAuthKey extracts Bearer key from request header.
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

	key := strings.TrimPrefix(auth, authScheme)
	if key == "" {
		return "", ErrInvalidKey
	}

	return key, nil
}
