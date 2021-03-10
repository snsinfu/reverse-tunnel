package service

import (
	"time"

	"github.com/gorilla/websocket"
)

// Watch watches for broken WebSocket connection. This function periodically
// sends ping message to the websocket peer and invokes `handler` on first
// timeout. The caller must continuously read something from `ws` to allow
// pong messages to be received.
func Watch(ws *websocket.Conn, timeout time.Duration, handler func() error) error {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	pong := make(chan bool)
	ws.SetPongHandler(func(_ string) error {
		pong <- true
		return nil
	})

	for tick := range ticker.C {
		if err := ws.WriteControl(websocket.PingMessage, []byte(""), tick.Add(timeout)); err != nil {
			break
		}

		select {
		case <-pong:
			continue

		case <-ticker.C:
		}
		break
	}

	return handler()
}
