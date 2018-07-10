package main

import (
	"net"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/snsinfu/reverse-tunnel/taskch"
)

type TCPAgent struct {
	port int
}

func NewTCPAgent(port int, key string) (*TCPAgent, error) {
	return &TCPAgent{port}, nil
}

func (agent *TCPAgent) Start(ws *websocket.Conn) error {
	addr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(agent.port))
	if err != nil {
		return err
	}

	ln, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	go ping(ws, agentTimeout, ln.Close)

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			return err
		}

		id := RegisterSession(&TCPSession{conn})

		err = ws.WriteJSON(map[string]interface{}{
			"event":          "accept",
			"session_id":     id,
			"client_address": conn.RemoteAddr(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type TCPSession struct {
	conn *net.TCPConn
}

func (sess *TCPSession) Start(ws *websocket.Conn) error {
	tasks := taskch.New()

	// Uplink
	tasks.Go(func() error {
		buf := make([]byte, bufferSize)

		for {
			n, err := sess.conn.Read(buf)
			if err != nil {
				return err
			}

			if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				return err
			}
		}

		return nil
	})

	// Downlink
	tasks.Go(func() error {
		buf := make([]byte, bufferSize)

		for {
			_, r, err := ws.NextReader()
			if err != nil {
				return err
			}

			n, err := r.Read(buf)
			if err != nil {
				return err
			}

			if _, err = sess.conn.Write(buf[:n]); err != nil {
				return err
			}
		}

		return nil
	})

	return tasks.Wait()
}

func (sess *TCPSession) Close() error {
	return sess.conn.Close()
}
