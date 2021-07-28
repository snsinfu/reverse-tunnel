package agent

import (
	"context"
	"errors"
	"fmt"
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

const (
	retryInterval  = 10 * time.Second
	wsCloseTimeout = 3 * time.Second
)

// dialer is the websocket dialer used to connect to a gateway server.
var dialer = websocket.DefaultDialer

// Agent tunnels remote port on a gateway server to local destination.
type Agent struct {
	gatewayURL  string
	key         string
	service     ports.NetPort
	destination string
}

// Start starts tunneling agents with given configurations. The agents and
// tunneled connections are cancelable via context.
func Start(conf config.Agent, ctx context.Context) error {
	if err := conf.Check(); err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	tasks := taskch.New()

	for _, forw := range conf.Forwards {
		agent := Agent{
			gatewayURL:  conf.GatewayURL,
			key:         conf.AuthKey,
			service:     forw.Port,
			destination: forw.Destination,
		}

		tasks.Go(func() error {
			retry := time.NewTicker(retryInterval)
			defer retry.Stop()
			for {
				err := agent.Start(ctx)
				if err == nil || !isRecoverable(err) {
					return err
				}
				log.Printf("Agent error %q - recovering...", err)
				<-retry.C
			}
		})
	}

	return tasks.Wait()
}

// Returns true if agent error is recoverable (by restarting agent).
func isRecoverable(err error) bool {
	if errors.Is(err, context.Canceled) {
		return false
	}
	return true
}

// hookCancel launches a goroutine for handling task cancellation. handler is
// called upon cancellation. The caller must invoke returned function after the
// cancelable task is finished.
func hookCancel(ctx context.Context, handler func()) func() {
	end := make(chan struct{})
	unhook := func() {
		close(end)
	}
	go func() {
		select {
		case <-ctx.Done():
			handler()
		case <-end:
		}
	}()
	return unhook
}

// Start connects to tunneling server and starts listening on the remote ports.
// The connection to the server and tunnels can be canceled via context.
func (agent *Agent) Start(ctx context.Context) error {

	// Connection on the remote port is first notified to us as a WebSocket
	// message. We then call back server to create another WebSocket channel
	// for tunneling the connection (in Agent.tunnel function).

	url := agent.gatewayURL + "/" + agent.service.Protocol + "/" + strconv.Itoa(agent.service.Port)
	header := http.Header{}
	header.Add("Authorization", "Bearer "+agent.key)
	ws, _, err := dialer.DialContext(ctx, url, header)
	if err != nil {
		return err
	}
	defer closeWebsocket(ws)

	unhookCancel := hookCancel(ctx, func() {
		closeWebsocket(ws)
	})
	defer unhookCancel()

	log.Printf("Listening on remote port: %s", agent.service)

	for {
		accept := service.BinderAcceptMessage{}
		if err := ws.ReadJSON(&accept); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Agent %s closed normally.", agent.service)
				return nil
			}
			if errors.Is(err, net.ErrClosed) {
				log.Printf("Agent %s is canceled.", agent.service)
				return context.Canceled
			}
			return err
		}

		go func() {
			if err := agent.tunnel(accept, ctx); err != nil {
				log.Printf("Tunneling error: %s", err)
			}
		}()
	}
}

// tunnel proxies accepted connection on the remote tunneling server to a local
// connection to the forwarding destination. The tunnel can be canceled via
// context.
func (agent *Agent) tunnel(accept service.BinderAcceptMessage, ctx context.Context) error {
	log.Printf(
		"Tunneling remote connection from %s to %s",
		accept.PeerAddress,
		agent.destination,
	)

	// Local connection to the forwarding destination.
	conn, err := net.Dial(agent.service.Protocol, agent.destination)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Remote connection proxied through WebSocket.
	url := agent.gatewayURL + "/session/" + accept.SessionID
	ws, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return err
	}
	defer closeWebsocket(ws)

	unhookCancel := hookCancel(ctx, func() {
		conn.Close()
		closeWebsocket(ws)
	})
	defer unhookCancel()

	tasks := taskch.New()

	// The tunnel looks like this:
	//
	// Client <------> Server <------> Agent <------> Destination
	//                           ws            conn
	//

	// Uplink: (Client --->) Server ---> Agent ---> Destination
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

	// Downlink: (Client <---) Server <--- Agent <--- Destination
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

	err = tasks.Wait()

	if errors.Is(err, io.EOF) {
		log.Printf(
			"Destination closed. Finishing session %s -> %s",
			accept.PeerAddress,
			agent.destination,
		)
		return nil
	}

	if errors.Is(err, net.ErrClosed) {
		log.Printf(
			"Canceled. Finishing session %s -> %s",
			accept.PeerAddress,
			agent.destination,
		)
		return nil
	}

	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		log.Printf(
			"Tunnel closed. Finishing session %s -> %s",
			accept.PeerAddress,
			agent.destination,
		)
		return nil
	}

	log.Printf(
		"Error %q. Killing session %s -> %s",
		err,
		accept.PeerAddress,
		agent.destination,
	)

	return err
}

// closeWebsocket attempts to close a websocket session normally. It is ok to
// call this function on a connection that has already been closed by the peer.
func closeWebsocket(ws *websocket.Conn) {
	ws.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(wsCloseTimeout),
	)
	ws.Close()
}
