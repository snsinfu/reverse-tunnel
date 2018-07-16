package service

import (
	"github.com/gorilla/websocket"
)

// Service is an abstraction of a tunneling service.
type Service interface {
	GetBinder(key string, port int) (Binder, error)
}

// Binder creates tunnleing session for each client connection.
type Binder interface {
	Start(ws *websocket.Conn, store *SessionStore) error
}
