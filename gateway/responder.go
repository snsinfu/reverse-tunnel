package main

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

const (
	wsCloseTimeout = 5 * time.Second
)

// WebSocket upgrades a request to a websocket session.
func WebSocket(c echo.Context, handler func(*websocket.Conn) error) error {
	upgrader := websocket.Upgrader{}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return Error(c, err)
	}

	go func() {
		if err := handler(ws); err != nil {
			c.Logger().Errorf("websocket: %s", err)
		}

		ws.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(wsCloseTimeout),
		)
		ws.Close()
	}()

	return nil
}

// Error responds to request with an error.
func Error(c echo.Context, err error) error {
	return err
}
