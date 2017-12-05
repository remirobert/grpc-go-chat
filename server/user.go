package main

import (
	pb "grpc-go-chat/chat"
)

type User struct {
	id       string
	username string
	stream   pb.ChatService_StreamServer
}

func NewUser(userData pb.User, stream pb.ChatService_StreamServer) *User {
	return &User{
		id: userData.Id,
		username:userData.Username,
		stream:stream,
	}
}
