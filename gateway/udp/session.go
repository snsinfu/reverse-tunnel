package udp

import (
	"crypto/rand"
	"encoding/hex"
	"net"

	"github.com/gorilla/websocket"
)

type Session struct {
	id     string
	tunnel *websocket.Conn
	conn   *net.UDPConn
	peer   *net.UDPAddr
}

func (sess *Session) Start(ws *websocket.Conn) error {
	sess.tunnel = ws

	buf := make([]byte, 1500)

	for {
		_, r, err := ws.NextReader()
		if err != nil {
			return err
		}

		n, err := r.Read(buf)
		if err != nil {
			return err
		}

		if _, err := sess.conn.WriteToUDP(buf[:n], sess.peer); err != nil {
			return err
		}
	}
}

func generateSessionID() string {
	buf := [4]byte{}
	rand.Read(buf[:])
	return hex.EncodeToString(buf[:])
}
