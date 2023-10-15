package config

import (
	"fmt"

	"github.com/snsinfu/reverse-tunnel/ports"
)

// ServerDefault is a default server configuration.
var ServerDefault = Server{
	ControlAddress: "localhost:9000",
	TLSConf:        TLSConf{},
}

// Server represents a configuration of a reverse tunnel server program.
type Server struct {
	ControlAddress string      `yaml:"control_address"`
	TLSConf        TLSConf     `yaml:"tls"`
	Agents         []AgentAuth `yaml:"agents"`
}

// TLSConf represents TLS configuration for the server.
type TLSConf struct {
	CertPath string `yaml:"cert_path"`
	KeyPath  string `yaml:"key_path"`
}

// AgentAuth represents an agent and its access rights authorized in a reverse
// tunnel server program.
type AgentAuth struct {
	AuthKey string          `yaml:"auth_key"`
	Ports   []ports.NetPort `yaml:"ports"`
}

// Check checks agent config for obvious mistakes. Returns a non-nil error if a
// bad configuration is found.
func (conf *Server) Check() error {
	if conf.ControlAddress == "" {
		return fmt.Errorf("control_address is empty")
	}

	for _, agent := range conf.Agents {
		if agent.AuthKey == "" {
			return fmt.Errorf("auth_key is unconfigured")
		}
	}

	return nil
}
