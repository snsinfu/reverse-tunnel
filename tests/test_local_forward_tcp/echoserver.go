package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	address := os.Args[1]

	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	fmt.Println("Listening on", address)

	conn, err := l.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Echoing")

	if _, err := io.Copy(conn, conn); err != nil {
		panic(err)
	}

	fmt.Println("Finished")
}
