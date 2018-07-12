package config

import (
	"github.com/snsinfu/reverse-tunnel/ports"
)

type Agent struct {
	GatewayURL string    `yaml:"gateway_url"`
	AuthKey    string    `yaml:"auth_key"`
	Forwards   []Forward `yaml:"forwards"`
}

type Forward struct {
	Port        ports.NetPort `yaml:"port"`
	Destination string        `yaml:"destination"`
}
