package server

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// wsCloseTimeout is the timeout of a WebSocket close message.
const wsCloseTimeout = 3 * time.Second

// errorMap maps error object to HTTP status code.
var errorMap = map[error]int{}

// Upgrader is the websocket upgrader.
var upgrader = websocket.Upgrader{}

// Error responds to client with an error. The error is logged and translated
// to proper HTTP status response.
func Error(c echo.Context, err error) error {
	c.Logger().Error(err)

	status, ok := errorMap[err]
	if !ok {
		status = http.StatusInternalServerError
	}

	return c.JSON(status, map[string]string{"error": err.Error()})
}

// WebSocket starts websocket session. The handler is invoked in a goroutine.
func WebSocket(c echo.Context, handler func(ws *websocket.Conn) error) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return Error(c, err)
	}

	go func() {
		var closeMessage []byte

		if err := handler(ws); err != nil {
			c.Logger().Error(err)

			closeMessage = websocket.FormatCloseMessage(
				websocket.ClosePolicyViolation, "error: "+err.Error(),
			)
		} else {
			closeMessage = websocket.FormatCloseMessage(
				websocket.CloseNormalClosure, "",
			)
		}

		ws.WriteControl(
			websocket.CloseMessage, closeMessage, time.Now().Add(wsCloseTimeout),
		)
		ws.Close()
	}()

	return nil
}
