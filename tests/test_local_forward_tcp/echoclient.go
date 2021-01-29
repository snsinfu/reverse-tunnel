package main

import (
	"io"
	"net"
	"os"
)

func main() {
	address := os.Args[1]

	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	n, err := io.Copy(conn, os.Stdin)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, n + 1)

	m, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	os.Stdout.Write(buf[:m])
}
