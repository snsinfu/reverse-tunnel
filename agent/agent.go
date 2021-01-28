package agent

import (
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/snsinfu/go-taskch"
	"github.com/snsinfu/reverse-tunnel/config"
	"github.com/snsinfu/reverse-tunnel/ports"
	"github.com/snsinfu/reverse-tunnel/server/service"
)

const retryInterval = 10 * time.Second

// dialer is the websocket dialer used to connect to a gateway server.
var dialer = websocket.DefaultDialer

// Agent remote-listens on a port on a gateway server.
type Agent struct {
	gatewayURL  string
	key         string
	service     ports.NetPort
	destination string
}

// Start starts agents with given configuration.
func Start(conf config.Agent) error {
	tasks := taskch.New()

	for _, forw := range conf.Forwards {
		agent := Agent{
			gatewayURL:  conf.GatewayURL,
			key:         conf.AuthKey,
			service:     forw.Port,
			destination: forw.Destination,
		}

		tasks.Go(func() error {
			delay := time.Tick(retryInterval)
			for {
				if err := agent.Start(); err != nil {
					log.Printf("Agent error %q - recovering...", err)
					<-delay
				}
			}
		})
	}

	return tasks.Wait()
}

// Start starts remote-listening.
func (agent Agent) Start() error {
	url := agent.gatewayURL + "/" + agent.service.Protocol + "/" + strconv.Itoa(agent.service.Port)

	header := http.Header{}
	header.Add("Authorization", "Bearer "+agent.key)

	ws, _, err := dialer.Dial(url, header)
	if err != nil {
		return err
	}
	defer ws.Close()

	log.Printf("Listening on remote port: %s", agent.service)

	for {
		accept := service.BinderAcceptMessage{}
		if err := ws.ReadJSON(&accept); err != nil {
			return err
		}

		log.Printf("Tunneling remote connection from %s to %s", accept.PeerAddress, agent.destination)

		go func() error {
			defer log.Printf("Closing connection from %s", accept.PeerAddress)

			conn, err := net.Dial(agent.service.Protocol, agent.destination)
			if err != nil {
				return err
			}
			defer conn.Close()

			url := agent.gatewayURL + "/session/" + accept.SessionID
			ws, _, err := dialer.Dial(url, nil)
			if err != nil {
				return err
			}
			defer ws.Close()

			tasks := taskch.New()

			// Uplink
			tasks.Go(func() error {
				buf := make([]byte, config.BufferSize)

				for {
					_, r, err := ws.NextReader()
					if err != nil {
						return err
					}

					if _, err := io.CopyBuffer(conn, r, buf); err != nil {
						return err
					}
				}
			})

			// Downlink
			tasks.Go(func() error {
				buf := make([]byte, config.BufferSize)

				for {
					n, err := conn.Read(buf)
					if err != nil {
						return err
					}

					if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
						return err
					}
				}
			})

			return tasks.Wait()
		}()
	}
}
