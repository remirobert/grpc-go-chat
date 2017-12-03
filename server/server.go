package main

import (
	pb "test-chat/chat"
	"net"
	"fmt"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
)

type server struct {
	users map[string]pb.ChatService_StreamServer
}

func (s *server) Stream(stream pb.ChatService_StreamServer) error {
	for {
		message, err := stream.Recv()
		if err != nil {
			return err
		}

		if message.Register {
			if message.User != nil {
				log.Print("register new user : " + message.User.Username)
				s.addUser(message.User, stream)
			}
		}
		if s.checkRegistrationUser(message.User) {
			s.broadcastMessage(message)
		} else {
			s.sendErrorRegistrationUser(stream)
		}
	}
}

func (s *server) broadcastMessage(message *pb.ChatMessage) {
	if message.Message != nil {
		log.Print("broadcast new message : " + message.Message.Content)
	} else {
		log.Print("user registration : " + message.User.Username)
	}

	for _, stream := range s.users {
		stream.Send(message)
	}
}

func (s *server) sendErrorRegistrationUser(stream pb.ChatService_StreamServer) {
	messageContent := pb.Message{Type: pb.Message_ERROR, Content: "Error authentification"}
	message := pb.ChatMessage{Message: &messageContent}
	stream.Send(&message)
}

func (s *server) addUser(user *pb.User, stream pb.ChatService_StreamServer) {
	s.users[user.Username] = stream
}

func (s *server) removeUser(user string) {
	delete(s.users, user)
}

func (s *server) checkRegistrationUser(user *pb.User) bool {
	for username := range s.users {
		if username == user.Username {
			return true
		}
	}
	return false
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
