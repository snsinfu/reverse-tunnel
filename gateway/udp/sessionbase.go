package udp

import (
	"net"
	"sync"
)

var sessions = SessionBase{}

type SessionBase struct {
	sessions map[string]*Session
	mutex    sync.Mutex
}

func (sb *SessionBase) Add(sess *Session) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()

	sb.sessions[sess.peer.String()] = sess
}

func (sb *SessionBase) Get(peer *net.UDPAddr) (*Session, error) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()

	if sess, ok := sb.sessions[peer.String()]; ok {
		return sess, nil
	}
	return nil, ErrNoSession
}

func (sb *SessionBase) Resolve(id string) (*Session, error) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()

	for key := range sb.sessions {
		sess := sb.sessions[key]

		if sess.id == id {
			sess.id = ""
			return sess, nil
		}
	}
	return nil, ErrNoSession
}

func ResolveSession(id string) (*Session, error) {
	return sessions.Resolve(id)
}
