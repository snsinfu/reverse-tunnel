package tcp

import (
	"sync"

	"github.com/snsinfu/reverse-tunnel/uid"
)

// SessionStore holds unresolved tunneling sessions.
type SessionStore struct {
	sync.Mutex

	sessions map[string]*Session
}

// Add adds an unresolved (not yet active) session. Returns token.
func (store *SessionStore) Add(sess *Session) string {
	store.Lock()
	defer store.Unlock()

	if store.sessions == nil {
		store.sessions = map[string]*Session{}
	}

	token := uid.New()
	store.sessions[token] = sess

	return token
}

// Resolve removes session and returns it.
func (store *SessionStore) Resolve(token string) *Session {
	store.Lock()
	defer store.Unlock()

	if sess, ok := store.sessions[token]; ok {
		return sess
	}

	return nil
}
