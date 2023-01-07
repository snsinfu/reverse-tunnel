package tcp

import (
	"errors"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/snsinfu/reverse-tunnel/server/service"
)

const (
	// Timeout value used for checking websocket connection loss.
	connTimeout = 3 * time.Second

	// Time to wait before retrying a failed Accept().
	acceptRetryWait = 100 * time.Millisecond
)

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

	// Forcefully close connection (thus end session) if the agent does not
	// respond to ping.
	go service.Watch(ws, connTimeout, ln.Close)

	go func() {
		for {
			// Agent does not send message to this channel in the current
			// protocol, but it is required to drain the channel to check
			// for ping responses.
			if _, _, err := ws.NextReader(); err != nil {
				break
			}
		}
	}()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			var nerr net.Error
			if errors.As(err, &nerr) && nerr.Temporary() {
				time.Sleep(acceptRetryWait)
				continue
			}
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
