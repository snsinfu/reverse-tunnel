package server

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/snsinfu/reverse-tunnel/config"
	"github.com/snsinfu/reverse-tunnel/server/service"
	"github.com/snsinfu/reverse-tunnel/server/tcp"
	"github.com/snsinfu/reverse-tunnel/server/udp"
)

// ErrInvalidPort is returned when
var ErrInvalidPort = errors.New("invalid port number")

// ErrMissingKey is returned when a request to an privileged resource does not
// contain an authorization key.
var ErrMissingKey = errors.New("missing auth key")

// ErrInvalidKey is returned when an authorization is is invalid or does not
// have sufficient privilege.
var ErrInvalidKey = errors.New("invalid auth key")

// ErrInvalidSessionID is returned when a requested session ID is not found.
var ErrInvalidSessionID = errors.New("invalid session ID")

// Action takes HTTP requests and serves tunneling service.
type Action struct {
	tcp   tcp.Service
	udp   udp.Service
	store service.SessionStore
}

// NewAction creates a new Action with given server configuration.
func NewAction(conf config.Server) Action {
	return Action{
		tcp: tcp.NewService(conf),
		udp: udp.NewService(conf),
	}
}

// GetTCPPort handles GET /tcp/:port request, backed by tcp.Binder.
func (action *Action) GetTCPPort(c echo.Context) error {
	return action.GetServicePort(c, &action.tcp)
}

// GetUDPPort handles GET /udp/:port request, backed by udp.Binder.
func (action *Action) GetUDPPort(c echo.Context) error {
	return action.GetServicePort(c, &action.udp)
}

// GetServicePort handles port binding request, backed by given service.
func (action *Action) GetServicePort(c echo.Context, serv service.Service) error {
	port, err := parsePort(c.Param("port"))
	if err != nil {
		return Error(c, err)
	}

	key, err := extractAuthKey(c)
	if err != nil {
		return Error(c, err)
	}

	binder, err := serv.GetBinder(key, port)
	if err != nil {
		return Error(c, err)
	}

	return WebSocket(c, func(ws *websocket.Conn) error {
		return binder.Start(ws, &action.store)
	})
}

// GetSession handles GET /session/:id request.
func (action *Action) GetSession(c echo.Context) error {
	id := c.Param("id")

	sess := action.store.Resolve(id)
	if sess == nil {
		return Error(c, ErrInvalidSessionID)
	}

	return WebSocket(c, func(ws *websocket.Conn) error {
		defer action.store.Remove(sess)
		return sess.Start(ws)
	})
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
