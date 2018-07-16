package service

import (
	"time"

	"github.com/gorilla/websocket"
)

// Watch periodically sends ping to websocket peer and invokes handler on the
// first timeout.
func Watch(ws *websocket.Conn, timeout time.Duration, handler func() error) error {
	for tick := range time.NewTicker(timeout).C {
		if err := ws.WriteControl(websocket.PingMessage, []byte(""), tick.Add(timeout)); err != nil {
			break
		}
	}
	return handler()
}
