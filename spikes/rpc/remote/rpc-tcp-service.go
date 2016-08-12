package main

import (
	"net"
	"net/rpc"
	"os"
)

func main() {
	compose := new(Compose)

	rpc.Register(compose)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		os.Exit(1)
	}

	rpc.Accept(listener)
}
