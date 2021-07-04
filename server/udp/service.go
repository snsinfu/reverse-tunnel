package udp

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/snsinfu/reverse-tunnel/config"
	"github.com/snsinfu/reverse-tunnel/ports"
	"github.com/snsinfu/reverse-tunnel/server/service"
)

// ErrUnauthorizedKey is returned when a key is not authorized.
var ErrUnauthorizedKey = errors.New("unauthorized key")

// ErrUnauthorizedPort is returned when a key is not allowed to bind to a
// requested port.
var ErrUnauthorizedPort = errors.New("unauthorized port")

// Service implements service.Service for UDP tunneling service.
type Service struct {
	authorities map[string]ports.Set
}

// NewService creates a udp.Service with given server configuration.
func NewService(conf config.Server) Service {
	auths := map[string]ports.Set{}

	for _, agent := range conf.Agents {
		set := ports.Set{}

		for _, np := range agent.Ports {
			if np.Protocol == "udp" {
				set.Add(np.Port)
			}
		}

		auths[agent.AuthKey] = set
	}

	return Service{authorities: auths}
}

// GetBinder returns a udp.Binder for an agent with given authorization key and
// given UDP port.
func (serv Service) GetBinder(key string, port int) (service.Binder, error) {
	set, ok := serv.authorities[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnauthorizedKey, key)
	}

	if !set.Has(port) {
		return nil, fmt.Errorf("%w: %d/udp (key: %s)", ErrUnauthorizedPort, port, key)
	}

	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	return &Binder{addr: addr}, nil
}
