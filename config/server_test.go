package config

import (
	"testing"

	"github.com/snsinfu/reverse-tunnel/ports"
)

func TestServer_Check_BadConfig(t *testing.T) {
	badConfigs := []Server{
		// Empty control address
		{
			ControlAddress: "",
		},

		// No Auth key
		{
			ControlAddress: ":9000",
			Agents: []AgentAuth{
				{
					Ports: []ports.NetPort{{"tcp", 8080}},
				},
			},
		},
	}

	for _, badConfig := range badConfigs {
		err := badConfig.Check()
		if err == nil {
			t.Fatalf("error not caught: %v", badConfig)
		}
	}
}
