package config

import (
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
