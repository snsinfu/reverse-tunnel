package tcp

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/snsinfu/reverse-tunnel/server/service"
)

// Timeout value used for checking websocket connection loss.
const connTimeout = 3 * time.Second

// Binder implements service.Binder for TCP tunneling service.
type Binder struct {
	addr *net.TCPAddr
}

// Start binds to a TCP port and creates tcp.Session for each client connection.
func (binder Binder) Start(ws *websocket.Conn, store *service.SessionStore) error {
	ln, err := net.ListenTCP("tcp", binder.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	go service.Watch(ws, connTimeout, ln.Close)

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			return err
		}

		sess := NewSession(conn)
		id := store.Add(sess)

		err = ws.WriteJSON(service.BinderAcceptMessage{
			Event:       "accept",
			SessionID:   id,
			PeerAddress: conn.RemoteAddr().String(),
		})
		if err != nil {
			return err
		}
	}
}
