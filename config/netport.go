package config

import (
	"fmt"
)

// NetPort is a port number paired with protocol name (tcp or udp).
type NetPort struct {
	Port     int
	Protocol string
}

// ParseNetPort parses s as a NetPort in a form like "80/tcp".
func ParseNetPort(s string) (NetPort, error) {
	np := NetPort{}
	_, err := fmt.Sscanf(s, "%d/%s", &np.Port, &np.Protocol)

	return np, err
}

// String returns a string representation of a NetPort as a form like "80/tcp".
func (np NetPort) String() string {
	return fmt.Sprintf("%d/%s", np.Port, np.Protocol)
}

// MarshalYAML implements yaml.Marshaler interface.
func (np NetPort) MarshalYAML() (interface{}, error) {
	return np.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (np *NetPort) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var repr string

	if err := unmarshal(&repr); err != nil {
		return err
	}

	p, err := ParseNetPort(repr)
	if err != nil {
		return err
	}
	*np = p

	return nil
}
