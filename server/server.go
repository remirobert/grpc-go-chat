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
		return errors.New("Join First")
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
			return errors.New("Already joined")
		case pb.ChatMessage_USER_LEAVE:
			s.processLeaveUser(message)
		case pb.ChatMessage_USER_CHAT:
			s.processChatUser(message)
		}

		/*
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
		*/
	}
}

func (s *server) processAddUser(message *pb.ChatMessage, stream pb.ChatService_StreamServer) error {
	if _, ok := s.users[message.User.Username]; ok {
		return errors.New("User already exists")
	}
	s.users[message.User.Username] = stream
	s.broadcastUserJoin(message.User)
	return nil
}

func (s *server) processLeaveUser(message *pb.ChatMessage) {

}

func (s *server) processChatUser(message *pb.ChatMessage) {

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

func (s* server) broadcastUserJoin(user *pb.User) {
	message := pb.ChatMessage{Type:pb.ChatMessage_USER_JOIN, User:user}
	s.broadcastMessage(&message)
}

func (s *server) sendErrorRegistrationUser(stream pb.ChatService_StreamServer, user *pb.User) {
	message := pb.ChatMessage{Type:pb.ChatMessage_USER_JOIN, User:user}
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
