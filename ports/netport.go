package ports

import (
	"fmt"
)

// NetPort is a pair of IP protocol name and port number.
type NetPort struct {
	Protocol string
	Port     int
}

// ParseNetPort parses s as NetPort in a form like "80/tcp".
func ParseNetPort(s string) (NetPort, error) {
	np := NetPort{}
	_, err := fmt.Sscanf(s, "%d/%s", &np.Port, &np.Protocol)
	return np, err
}

// String formats np as string in a form like "80/tcp".
func (np NetPort) String() string {
	return fmt.Sprintf("%d/%s", np.Port, np.Protocol)
}

// MarshalYAML implements yaml.Marshaler interface.
func (np NetPort) MarshalYAML() (interface{}, error) {
	return np.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler interface.
func (np *NetPort) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	val, err := ParseNetPort(s)
	if err != nil {
		return err
	}

	*np = val

	return nil
}
