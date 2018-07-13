package udp

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

const (
	binderTimeout = 3 * time.Second
	network       = "udp4"
)

// Binder accepts UDP packets and starts session.
type Binder struct {
	addr     *net.UDPAddr
	sessions *SessionStore
}

func (binder Binder) Start(ws *websocket.Conn) error {
	conn, err := net.ListenUDP(network, binder.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	go watchConnLoss(ws, binderTimeout, conn.Close)

	buf := make([]byte, bufferSize)

	for {
		n, peer, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		if sess := binder.sessions.Get(peer); sess != nil {
			sess.SendToServer(buf[:n])
		} else {
			sess := NewSession(conn, peer)
			binder.sessions.Add(sess)

			err = ws.WriteJSON(map[string]interface{}{
				"event":        "accept",
				"session_id":   sess.token,
				"peer_address": peer.String(),
			})
			if err != nil {
				return err
			}
		}
	}
}
