package config

import (
	"testing"

	"github.com/snsinfu/reverse-tunnel/ports"
)

func TestAgent_Check_BadConfig(t *testing.T) {
	badConfigs := []Agent{
		// Empty gateway URL
		{
			AuthKey: "key",
		},

		// Empty auth key
		{
			GatewayURL: "ws://localhost:9000",
		},

		// No port
		{
			GatewayURL: "ws://localhost:9000",
			AuthKey:    "key",
			Forwards: []Forward{
				{
					Destination: "127.0.0.1:10000",
				},
			},
		},

		// No destination
		{
			GatewayURL: "ws://localhost:9000",
			AuthKey:    "key",
			Forwards: []Forward{
				{
					Port: ports.NetPort{"tcp", 10000},
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
