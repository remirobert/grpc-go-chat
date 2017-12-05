package main

import (
	pb "grpc-go-chat/chat"
	"net"
	"fmt"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"errors"
)

type server struct {
	users map[string]pb.ChatService_StreamServer
}

func (s *server) Stream(stream pb.ChatService_StreamServer) error {
	request, err := stream.Recv()
	if err != nil {
		return err
	}

	if request.Type != pb.ChatMessage_USER_JOIN {
		return errors.New("join first")
	}

	if err := s.processAddUser(request, stream); err != nil {
		return err
	}

	for {
		message, err := stream.Recv()
		if err != nil {
			s.processLeaveUser(message)
			return err
		}

		switch message.Type {
		case pb.ChatMessage_USER_JOIN:
			return errors.New("already joined")
		case pb.ChatMessage_USER_LEAVE:
			s.processLeaveUser(message)
		case pb.ChatMessage_USER_CHAT:
			s.processChatUser(message)
		}
	}
}

func (s *server) processAddUser(message *pb.ChatMessage, stream pb.ChatService_StreamServer) error {
	if message.User == nil {
		return errors.New("no user found in the request")
	}
	if _, ok := s.users[message.User.Username]; ok {
		return errors.New("user already exists")
	}
	s.users[message.User.Username] = stream
	log.Print("new user joined the channel : ", message.User.Username)
	s.broadcastUserJoin(message.User)
	return nil
}

func (s *server) processLeaveUser(message *pb.ChatMessage) {
	log.Print("remove the user : ", message)
	if message == nil || message.User == nil {
		return
	}
	delete(s.users, message.User.Username)
	log.Print("new user left the channel : ", message.User.Username)
	s.broadcastUserLeave(message.User)
}

func (s *server) processChatUser(message *pb.ChatMessage) {
	s.broadcastMessage(message)
}

func (s *server) broadcastMessage(message *pb.ChatMessage) {
	for _, stream := range s.users {
		stream.Send(message)
	}
}

func (s* server) broadcastUserJoin(user *pb.User) {
	message := pb.ChatMessage{Type:pb.ChatMessage_USER_JOIN, User:user}
	s.broadcastMessage(&message)
}

func (s* server) broadcastUserLeave(user *pb.User) {
	message := pb.ChatMessage{Type:pb.ChatMessage_USER_LEAVE, User:user}
	s.broadcastMessage(&message)
}

func main() {
	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		fmt.Print(err)
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	server := &server{users: make(map[string]pb.ChatService_StreamServer)}

	pb.RegisterChatServiceServer(s, server)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		fmt.Print(err)
		log.Fatalf("failed to serve: %v", err)
	}
}
