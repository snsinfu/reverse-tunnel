package config

import (
	"github.com/snsinfu/reverse-tunnel/ports"
)

type Gateway struct {
	ControlAddress string      `yaml:"control_address"`
	Agents         []AgentAuth `yaml:"agents"`
}

type AgentAuth struct {
	AuthKey string          `yaml:"auth_key"`
	Ports   []ports.NetPort `yaml:"ports"`
}
