package tcp

import (
	"net"

	"github.com/gorilla/websocket"
	"github.com/snsinfu/reverse-tunnel/config"
	"github.com/snsinfu/reverse-tunnel/taskch"
)

// Session
type Session struct {
	conn *net.TCPConn
}

// NewSession returns a Session object.
func NewSession(conn *net.TCPConn) *Session {
	return &Session{conn}
}

// Start starts bidirectional tunneling via websocket.
func (sess *Session) Start(ws *websocket.Conn) error {
	defer sess.conn.Close()

	tasks := taskch.New()

	// Uplink
	tasks.Go(func() error {
		buf := make([]byte, config.BufferSize)

		for {
			n, err := sess.conn.Read(buf)
			if err != nil {
				return err
			}

			if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				return err
			}
		}
	})

	// Downlink
	tasks.Go(func() error {
		buf := make([]byte, config.BufferSize)

		for {
			_, r, err := ws.NextReader()
			if err != nil {
				return err
			}

			n, err := r.Read(buf)
			if err != nil {
				return err
			}

			if _, err := sess.conn.Write(buf[:n]); err != nil {
				return err
			}
		}
	})

	return tasks.Wait()
}
