package main

import (
	"github.com/labstack/echo"
	"github.com/snsinfu/reverse-tunnel/gateway/udp"
)

// GetUDPPort handles GET /udp/:port request.
func GetUDPPort(c echo.Context) error {
	port, err := ParsePort(c.Param("port"))
	if err != nil {
		return Error(c, err)
	}

	key, err := ExtractAuthKey(c)
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
