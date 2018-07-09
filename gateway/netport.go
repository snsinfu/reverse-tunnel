package main

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
