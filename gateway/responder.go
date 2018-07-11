package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/snsinfu/reverse-tunnel/gateway/tcp"
	"github.com/snsinfu/reverse-tunnel/gateway/udp"
)

const wsCloseTimeout = 5 * time.Second

// errorMap maps error object to HTTP status code.
var errorMap = map[error]int{
	ErrInvalidPort:           http.StatusBadRequest,
	ErrMissingKey:            http.StatusUnauthorized,
	ErrInvalidKey:            http.StatusBadRequest,
	ErrInvalidSessionID:      http.StatusNotFound,
	udp.ErrUnauthorizedKey:   http.StatusUnauthorized,
	udp.ErrInsufficientScope: http.StatusForbidden,
	tcp.ErrUnauthorizedKey:   http.StatusUnauthorized,
	tcp.ErrInsufficientScope: http.StatusForbidden,
}

// Data sends JSON response to the client.
func Data(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, data)
}

// Error interprets err and sends error response to the client.
func Error(c echo.Context, err error) error {
	status, ok := errorMap[err]
	if !ok {
		status = http.StatusInternalServerError
	}

	return c.JSON(status, map[string]string{"error": err.Error()})
}

// WebSocket starts websocket session. Handler is invoked in a goroutine.
func WebSocket(c echo.Context, handler func(ws *websocket.Conn) error) error {
	upgrader := websocket.Upgrader{}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return Error(c, err)
	}

	go func() {
		if err := handler(ws); err != nil {
			c.Logger().Error(err)
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
