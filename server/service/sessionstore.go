package service

import (
	"fmt"
	"net"
	"sync"

	"github.com/snsinfu/reverse-tunnel/hexid"
	"github.com/snsinfu/reverse-tunnel/ports"
)

// sessionIDSize is the number of random bytes encoded in a session ID.
const sessionIDSize = 8

// SessionStore is a concurrent storage of sessions.
type SessionStore struct {
	sync.Mutex

	sessions map[string]Session
	tokens   map[string]string
}

// Add adds a new session to store and returns a token ID for the session. The
// session can later be retrieved using the token for only a single time.
func (store *SessionStore) Add(sess Session) string {
	store.Lock()
	defer store.Unlock()

	if store.sessions == nil {
		store.sessions = map[string]Session{}
		store.tokens = map[string]string{}
	}

	peer := encodeAddr(sess.PeerAddr())
	id := hexid.New(sessionIDSize)

	store.sessions[peer] = sess
	store.tokens[id] = peer

	return id
}

// Resolve returns a session associated to given token ID. Nil is returned if
// the token is invalid. Token is invalidated after calling this function.
func (store *SessionStore) Resolve(id string) Session {
	store.Lock()
	defer store.Unlock()

	peer, ok := store.tokens[id]
	if !ok {
		return nil
	}
	delete(store.tokens, id)

	return store.sessions[peer]
}

// Get returns the session connected with given peer. Nil is returned if no
// such session in store.
func (store *SessionStore) Get(peer net.Addr) Session {
	store.Lock()
	defer store.Unlock()

	sess, ok := store.sessions[encodeAddr(peer)]
	if !ok {
		return nil
	}

	return sess
}

// Remove removes sess from store.
func (store *SessionStore) Remove(sess Session) {
	store.Lock()
	defer store.Unlock()

    addr := encodeAddr(sess.PeerAddr())
    if compare, ok := store.sessions[addr]; ok && compare == sess {
        delete(store.sessions, addr)
    }
}

func (store *SessionStore) Close(port ports.NetPort, key string) {
    for _, sess := range store.sessions {
        fmt.Println(sess)
        if sess.GetPort() == port && sess.GetKey() == key {
            sess.Close()
        }
    }
}

// encodeAddr encodes network address as a string.
func encodeAddr(addr net.Addr) string {
	return addr.String() + "/" + addr.Network()
}
