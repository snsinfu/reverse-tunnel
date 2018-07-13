package main

import (
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/snsinfu/reverse-tunnel/config"
	"github.com/snsinfu/reverse-tunnel/taskch"
)

var dialer = websocket.DefaultDialer

func startAgent(conf config.Agent) error {
	tasks := taskch.New()

	for i := range conf.Forwards {
		forw := conf.Forwards[i]

		switch forw.Port.Protocol {
			case "tcp":
				tasks.Go(func() error {
					return remoteListenTCP(forw, conf)
				})

			case "udp":
				tasks.Go(func() error {
					return remoteListenUDP(forw, conf)
				})

			default:
				log.Printf("unrecognized protocol: %s", forw.Port.Protocol)
		}
	}

	return tasks.Wait()
}

func remoteListenTCP(forw config.Forward, conf config.Agent) error {
	url := conf.GatewayURL + "/tcp/" + strconv.Itoa(forw.Port.Port)

	header := http.Header{}
	header.Add("Authorization", "Bearer "+conf.AuthKey)

	ws, _, err := dialer.Dial(url, header)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		event := map[string]string{}
		if err := ws.ReadJSON(&event); err != nil {
			return err
		}

		log.Printf("forwarding remote connection from %s to %s", event["peer_address"], forw.Destination)
		sessionID := event["session_id"]

		go func() error {
			addr, err := net.ResolveTCPAddr("tcp", forw.Destination)
			if err != nil {
				return err
			}

			conn, err := net.DialTCP("tcp", nil, addr)
			if err != nil {
				return err
			}
			defer conn.Close()

			url := conf.GatewayURL + "/tcp/session/" + sessionID
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

					n, err := r.Read(buf)
					if err != nil {
						return err
					}

					if _, err := conn.Write(buf[:n]); err != nil {
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

	return nil
}

func remoteListenUDP(forw config.Forward, conf config.Agent) error {
	url := conf.GatewayURL + "/udp/" + strconv.Itoa(forw.Port.Port)

	header := http.Header{}
	header.Add("Authorization", "Bearer "+conf.AuthKey)

	ws, _, err := dialer.Dial(url, header)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		event := map[string]string{}
		if err := ws.ReadJSON(&event); err != nil {
			return err
		}

		log.Printf("remote connection from %s", event["peer_address"])
		sessionID := event["session_id"]

		go func() error {
			addr, err := net.ResolveUDPAddr("udp", forw.Destination)
			if err != nil {
				return err
			}

			conn, err := net.DialUDP("udp", nil, addr)
			if err != nil {
				return err
			}
			defer conn.Close()

			url := conf.GatewayURL + "/udp/session/" + sessionID
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

					n, err := r.Read(buf)
					if err != nil {
						return err
					}

					if _, err := conn.Write(buf[:n]); err != nil {
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

	return nil
}
