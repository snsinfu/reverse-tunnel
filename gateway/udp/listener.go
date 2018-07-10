package udp

import (
	"net"
	"strconv"

	"github.com/gorilla/websocket"
)

type Listener struct {
	Port int
}

func GetListener(key string, port int) (*Listener, error) {
	auth, ok := authorities[key]
	if !ok {
		return nil, ErrUnauthorizedKey
	}

	if !auth.Check(port) {
		return nil, ErrInsufficientScope
	}

	return &Listener{port}, nil
}

func (listener *Listener) Start(ws *websocket.Conn) error {
	addr, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(listener.Port))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := make([]byte, 1500)

	for {
		n, peer, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		if sess, err := sessions.Get(peer); err == nil {
			if sess.tunnel != nil {
				sess.tunnel.WriteMessage(websocket.BinaryMessage, buf[:n])
			}
		} else {
			sess := &Session{
				id:     generateSessionID(),
				tunnel: nil,
				conn:   conn,
				peer:   peer,
			}
			sessions.Add(sess)

			err = ws.WriteJSON(map[string]interface{}{
				"event":        "accept",
				"peer_address": peer,
				"session_id":   sess.id,
			})
			if err != nil {
				return err
			}
		}
	}
}
