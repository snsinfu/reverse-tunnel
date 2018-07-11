package udp

import (
	"net"
	"sync"
)

// SessionStore holds active tunneling sessions.
type SessionStore struct {
	sync.Mutex

	sessions map[string]*Session
}

// Add adds an unresolved (not yet active) session. The session must be resolved
// later with sess.token.
func (store *SessionStore) Add(sess *Session) {
	store.Lock()
	defer store.Unlock()

	if store.sessions == nil {
		store.sessions = map[string]*Session{}
	}

	sess.store = store
	store.sessions[sess.peer.String()] = sess
}

// Get returns session for given UDP peer. Returns nil if session is not
// established yet for the peer.
func (store *SessionStore) Get(peer *net.UDPAddr) *Session {
	store.Lock()
	defer store.Unlock()

	sess, ok := store.sessions[peer.String()]
	if !ok {
		return nil
	}

	return sess
}

// Remove removes sess.
func (store *SessionStore) Remove(sess *Session) {
	store.Lock()
	defer store.Unlock()

	delete(store.sessions, sess.peer.String())
}

// Resolve searches for unresolved session with given token.
func (store *SessionStore) Resolve(token string) *Session {
	store.Lock()
	defer store.Unlock()

	for _, sess := range store.sessions {
		if sess.token == token {
			sess.token = ""
			return sess
		}
	}

	return nil
}
