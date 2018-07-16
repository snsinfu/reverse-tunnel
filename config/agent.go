package config

import (
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
