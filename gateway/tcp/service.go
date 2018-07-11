package tcp

import (
	"errors"
	"net"
	"strconv"

	"github.com/snsinfu/reverse-tunnel/config"
	"github.com/snsinfu/reverse-tunnel/ports"
)

var (
	ErrUnauthorizedKey   = errors.New("unauthorized key")
	ErrInsufficientScope = errors.New("insufficient scope")
)

// Service represents TCP tunneling service.
type Service struct {
	authorities map[string]ports.Set
	sessions    SessionStore
}

// NewService returns a UDP Service object with given configuration.
func NewService(conf config.Gateway) Service {
	auths := map[string]ports.Set{}

	for _, agent := range conf.Agents {
		set := ports.Set{}

		for _, np := range agent.Ports {
			if np.Protocol == "tcp" {
				set.Append(np.Port)
			}
		}

		auths[agent.AuthKey] = set
	}

	return Service{
		authorities: auths,
	}
}

// NewBinder returns a Binder for specified port.
func (service *Service) NewBinder(key string, port int) (*Binder, error) {
	set, ok := service.authorities[key]
	if !ok {
		return nil, ErrUnauthorizedKey
	}

	if !set.Has(port) {
		return nil, ErrInsufficientScope
	}

	addr, err := net.ResolveTCPAddr(network, ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	binder := &Binder{
		addr:     addr,
		sessions: &service.sessions,
	}
	return binder, nil
}

// ResolveSession resolves session ID to session object. The session ID cannot
// be reused later.
func (service *Service) ResolveSession(id string) *Session {
	return service.sessions.Resolve(id)
}
