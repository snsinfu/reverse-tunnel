package config

import (
	"fmt"

	"github.com/snsinfu/reverse-tunnel/ports"
)

// AgentDefault is a default agent configuration.
var AgentDefault = Agent{
	GatewayURL: "ws://localhost:9000",
}

// Agent represnets a configuration of a reverse tunnel agent program.
type Agent struct {
	GatewayURL string    `yaml:"gateway_url"`
	AuthKey    string    `yaml:"auth_key"`
	Forwards   []Forward `yaml:"forwards"`
}

// Forward represents a configuration of a single port forwarding.
type Forward struct {
	Port        ports.NetPort `yaml:"port"`
	Destination string        `yaml:"destination"`
}

// Check checks agent config for obvious mistakes. Returns a non-nil error if a
// bad configuration is found.
func (conf *Agent) Check() error {
	if conf.GatewayURL == "" {
		return fmt.Errorf("gateway_url is empty")
	}

	if conf.AuthKey == "" {
		return fmt.Errorf("auth_key is empty")
	}

	for _, forw := range conf.Forwards {
		if forw.Port.Protocol == "" || forw.Port.Port == 0 {
			return fmt.Errorf("port is unconfigured")
		}

		if forw.Destination == "" {
			return fmt.Errorf("port %s destination is empty", forw.Port)
		}
	}

	return nil
}
