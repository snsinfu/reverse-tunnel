package service

import (
	"net"

	"github.com/gorilla/websocket"
)

// Session is an abstraction of a tunneling session.
type Session interface {
	// PeerAddr returns the network address of a connected peer. This must be
	// invariant throughout the entire lifetime of a session.
	PeerAddr() net.Addr

	// Start starts a tunneling session through given websocket channel.
	Start(ws *websocket.Conn) error

	// Close closes client connection. This cancels tunneling session.
	Close() error
}
