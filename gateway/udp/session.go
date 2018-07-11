package udp

import (
	"errors"
	"net"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/gorilla/websocket"
	"github.com/snsinfu/reverse-tunnel/uid"
)

const sessionTimeout = 100 * time.Second

var ErrTunnelNotReady = errors.New("tunnel is not yet established")

// Session
type Session struct {
	token  string
	tunnel *websocket.Conn
	conn   *net.UDPConn
	peer   *net.UDPAddr
	idle   int32
	store  *SessionStore
}

// NewSession returns a Session object with random token and nil tunnel. The
// tunnel must be set to a valid websocket connection with ResolveSession.
func NewSession(conn *net.UDPConn, peer *net.UDPAddr) *Session {
	return &Session{
		token:  uid.New(),
		tunnel: nil,
		conn:   conn,
		peer:   peer,
		idle:   0,
	}
}

// SendToServer sends msg to the other end of the tunnel. Calling this function
// resets internal idle counter.
func (sess *Session) SendToServer(msg []byte) error {
	atomic.StoreInt32(&sess.idle, 0)

	ws := (*websocket.Conn)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&sess.tunnel))))
	if ws == nil {
		return ErrTunnelNotReady
	}

	return ws.WriteMessage(websocket.BinaryMessage, msg)
}

// SendToClient sends msg to UDP connected client. Calling this function resets
// internal idle counter.
func (sess *Session) SendToClient(msg []byte) error {
	atomic.StoreInt32(&sess.idle, 0)

	_, err := sess.conn.WriteToUDP(msg, sess.peer)
	return err
}

// Start starts uplink tunneling. Downlink needs to be handled in the listener
// loop due to the connection-less nature of UDP.
func (sess *Session) Start(ws *websocket.Conn) error {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&sess.tunnel)), unsafe.Pointer(ws))

	// Kill idle session to prevent potential resource leak.
	go func() {
		for range time.NewTicker(sessionTimeout / 2).C {
			if atomic.AddInt32(&sess.idle, 1) == 2 {
				break
			}
		}
		sess.store.Remove(sess)

		ws.Close()
	}()

	buf := make([]byte, bufferSize)

	for {
		_, r, err := ws.NextReader()
		if err != nil {
			return err
		}

		n, err := r.Read(buf)
		if err != nil {
			return err
		}

		if err := sess.SendToClient(buf[:n]); err != nil {
			return err
		}
	}
}
