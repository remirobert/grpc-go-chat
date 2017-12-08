package main

import (
	pb "grpc-go-chat/chat"
)

type Client struct {
	user   pb.User
	stream pb.ChatService_StreamServer
}

func NewClient(user pb.User, stream pb.ChatService_StreamServer) *Client {
	return &Client{
		user:   user,
		stream: stream,
	}
}
