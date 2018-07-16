package config

import (
	"github.com/snsinfu/reverse-tunnel/ports"
)

// Server represents a configuration of a reverse tunnel server program.
type Server struct {
	ControlAddress string      `yaml:"control_address"`
	Agents         []AgentAuth `yaml:"agents"`
}

// AgentAuth represents an agent and its access rights authorized in a reverse
// tunnel server program.
type AgentAuth struct {
	AuthKey string          `yaml:"auth_key"`
	Ports   []ports.NetPort `yaml:"ports"`
}
