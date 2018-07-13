package tcp

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

const (
	binderTimeout = 3 * time.Second
	network       = "tcp4"
)

type Binder struct {
	addr     *net.TCPAddr
	sessions *SessionStore
}

func (binder Binder) Start(ws *websocket.Conn) error {
	ln, err := net.ListenTCP(network, binder.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	go watchConnLoss(ws, binderTimeout, ln.Close)

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			return err
		}

		token := binder.sessions.Add(NewSession(conn))
		err = ws.WriteJSON(map[string]interface{}{
			"event":        "accept",
			"session_id":   token,
			"peer_address": conn.RemoteAddr().String(),
		})
		if err != nil {
			return err
		}
	}
}
