package config

// GatewayConfig represents gateway server configuration file.
type GatewayConfig struct {
	ControlAddress string      `yaml:"control_address"`
	Authorities    []Authority `yaml:"authorities"`
}

// Authority holds access information for a single authorized agent.
type Authority struct {
	AccessKey string    `yaml:"access_key"`
	Ports     []NetPort `yaml:"ports"`
}
