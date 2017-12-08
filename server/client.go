package server

import (
	pb "grpc-go-chat/chat"
)

type Client struct {
	User   pb.User
	Stream pb.ChatService_StreamServer
}

func NewClient(user pb.User, stream pb.ChatService_StreamServer) *Client {
	return &Client{
		User:   user,
		Stream: stream,
	}
}
