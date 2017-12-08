package main

import (
	. "grpc-go-chat/server"
)

func main() {
	server := NewServer()
	server.Start()
}

