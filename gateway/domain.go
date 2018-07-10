package main

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	sessionTimeout = 30 * time.Second
	agentTimeout   = 30 * time.Second
	bufferSize     = 1418
)

var (
	sessionMap  = map[string]Session{}
	sessionLock = sync.Mutex{}
)

type Session interface {
	Start(*websocket.Conn) error
	Close() error
}

// RegisterSession inserts session to internal pool and returns unique ID. The
// session will be closed if no one acquires the session before a deadline.
func RegisterSession(sess Session) string {
	sessionLock.Lock()
	defer sessionLock.Unlock()

	buf := make([]byte, 16)
	rand.Read(buf)
	id := hex.EncodeToString(buf)
	sessionMap[id] = sess

	go func() {
		<-time.After(sessionTimeout)

		if sess, _ := AcquireSession(id); sess != nil {
			sess.Close()
		}
	}()

	return id
}

// AcquireSession removes and returns session from internal pool.
func AcquireSession(id string) (Session, error) {
	sessionLock.Lock()
	defer sessionLock.Unlock()

	sess, ok := sessionMap[id]
	if !ok {
		return nil, ErrInvalidSessionID
	}
	delete(sessionMap, id)

	return sess, nil
}

// ping periodically sends ping message to websocket peer and invokes handler if
// ping is timed out.
func ping(ws *websocket.Conn, timeout time.Duration, handler func() error) error {
	for range time.NewTicker(timeout).C {
		err := ws.WriteControl(
			websocket.PingMessage,
			[]byte(""),
			time.Now().Add(timeout),
		)
		if err != nil {
			if isClosedError(err) {
				return handler()
			}
			return err
		}
	}
	return nil
}

// isClosedError checks if an error is about use of closed network connection.
func isClosedError(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
