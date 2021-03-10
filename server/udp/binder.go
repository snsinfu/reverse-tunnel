package udp

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/snsinfu/reverse-tunnel/config"
	"github.com/snsinfu/reverse-tunnel/server/service"
)

// connTimeout is the timeout used for checking websocket connection loss.
const connTimeout = 3 * time.Second

// Binder implements service.Binder for UDP tunneling service.
type Binder struct {
	addr *net.UDPAddr
}

// Start binds to a UDP port and routes incoming packets to udp.Session objects.
func (binder Binder) Start(ws *websocket.Conn, store *service.SessionStore) error {
	conn, err := net.ListenUDP("udp", binder.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Forcifully close connection (thus end session) if the agent does not
	// respond to ping.
	go service.Watch(ws, connTimeout, conn.Close)

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

	buf := make([]byte, config.BufferSize)

	for {
		n, peer, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		if sess, ok := store.Get(peer).(*Session); ok {
			sess.SendToAgent(buf[:n])
		} else {
			sess := NewSession(conn, peer)
			id := store.Add(sess)

			err = ws.WriteJSON(service.BinderAcceptMessage{
				Event:       "accept",
				SessionID:   id,
				PeerAddress: peer.String(),
			})
			if err != nil {
				return err
			}

			// NOTE: The message is dropped here, which is acceptable since it
			// is UDP. But it makes a noticeable delay, for example, on a mosh
			// handshake. Maybe the message should be queued.
		}
	}
}
