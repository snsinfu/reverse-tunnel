package main

import (
	"net"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	udpSessionTimeout = 30 * time.Second
)

type UDPAgent struct {
	port int
}

func NewUDPAgent(port int, key string) (*UDPAgent, error) {
	return &UDPAgent{port}, nil
}

func (agent *UDPAgent) Start(ws *websocket.Conn) error {
	addr, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(agent.port))
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	go ping(ws, agentTimeout, conn.Close)

	routes := map[string]*UDPSession{}
	buf := make([]byte, bufferSize)

	for {
		n, client, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}
		clientID := client.String()

		if sess, ok := routes[clientID]; ok && sess.ws != nil {
			sess.ws.WriteMessage(websocket.BinaryMessage, buf[:n])
		} else {
			sess := &UDPSession{conn, client, nil}
			id := RegisterSession(sess)
			routes[clientID] = sess

			err = ws.WriteJSON(map[string]interface{}{
				"event":          "accept",
				"session_id":     id,
				"client_address": client,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type UDPSession struct {
	conn   *net.UDPConn
	client *net.UDPAddr
	ws     *websocket.Conn
}

func (sess *UDPSession) Start(ws *websocket.Conn) error {
	go ping(ws, udpSessionTimeout, ws.Close)

	buf := make([]byte, bufferSize)
	sess.ws = ws

	for {
		_, r, err := ws.NextReader()
		if err != nil {
			return err
		}

		n, err := r.Read(buf)
		if err != nil {
			return err
		}

		if _, err := sess.conn.WriteToUDP(buf[:n], sess.client); err != nil {
			return err
		}
	}

	return nil
}

func (sess *UDPSession) Close() error {
	return sess.conn.Close()
}
