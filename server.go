package main

import (
	"fmt"
	"net"
	"os"

	"github.com/yigithanbalci/go-http-server/pkg/workers"
)

type Server struct {
	HttpVer string
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		workers.QueueConn(conn)
	}
}
