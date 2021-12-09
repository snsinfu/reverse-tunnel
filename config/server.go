package config

import (
	"fmt"

	"github.com/snsinfu/reverse-tunnel/ports"
)

// ServerDefault is a default server configuration.
var ServerDefault = Server{
	ControlAddress: "localhost:9000",
	LetsEncrypt:    LetsEncrypt{CacheDir: ".autocert_cache"},
}

// Server represents a configuration of a reverse tunnel server program.
type Server struct {
	ControlAddress string      `yaml:"control_address"`
	LetsEncrypt    LetsEncrypt `yaml:"lets_encrypt"`
	Agents         []AgentAuth `yaml:"agents"`
}

// LetsEncrypt represents autocert configuration for the server.
type LetsEncrypt struct {
	Domain   string `yaml:"domain"`
	CacheDir string `yaml:"cache_dir"`
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

func (conf *Server) Has(port ports.NetPort, authKey string) bool {
    for _, agt := range conf.Agents {
        if agt.has(port, authKey) {
            return true
        }
    }
    return false
}

func (agentAuth AgentAuth) has(port ports.NetPort, authKey string) bool {
    if authKey != agentAuth.AuthKey {
        return false
    }
    for _, agtPort := range agentAuth.Ports {
        if agtPort == port {
            return true
        }
    }
    return false
}
