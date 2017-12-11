package main

import (
	. "grpc-go-chat/server"
)

func main() {
	server := NewDefaultServer()
	server.Start()
}

