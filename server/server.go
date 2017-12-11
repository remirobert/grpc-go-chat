package server

import (
	pb "grpc-go-chat/chat"
	"net"
	"fmt"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
)

type Server struct {
	cm ClientManagerProvider
}

func (s *Server) Stream(stream pb.ChatService_StreamServer) error {
	request, err := stream.Recv()
	if err != nil {
		return err
	}
	if request == nil {
		return NewRequestError(RequestInvalid)
	}

	if request.Type != pb.ChatMessage_USER_JOIN {
		return NewAuthError(AuthMessageUserNotJoined)
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
		if message == nil {
			return NewRequestError(RequestInvalid)
		}

		switch message.Type {
		case pb.ChatMessage_USER_JOIN:
			return NewAuthError(AuthMessageUserAlreadyJoined)
		case pb.ChatMessage_USER_LEAVE:
			return s.processLeaveUser(message)
		case pb.ChatMessage_USER_CHAT:
			s.broadcastChatUser(message)
		}
	}
}

func (s *Server) processAddUser(message *pb.ChatMessage, stream pb.ChatService_StreamServer) error {
	if message.User == nil {
		return NewRequestError(RequestMessageNoUser)
	}
	if u := s.cm.Find(message.User.Id); u != nil {
		return NewAuthError(AuthMessageUserAlreadyJoined)
	}
	newUser := NewClient(*message.User, stream)
	s.cm.Add(*newUser)
	s.broadcastUserJoin(message.User)
	return nil
}

func (s *Server) processLeaveUser(message *pb.ChatMessage) error {
	if message.User == nil {
		return NewRequestError(RequestUserMissing)
	}
	s.cm.Remove(message.User.Id)
	s.broadcastUserLeave(message.User)
	return nil
}

func (s *Server) broadcastChatUser(message *pb.ChatMessage) {
	s.cm.BroadcastMessage(message)
}

func (s *Server) broadcastUserJoin(user *pb.User) {
	message := pb.ChatMessage{Type: pb.ChatMessage_USER_JOIN, User: user}
	s.cm.BroadcastMessage(&message)
}

func (s *Server) broadcastUserLeave(user *pb.User) {
	message := pb.ChatMessage{Type: pb.ChatMessage_USER_LEAVE, User: user}
	s.cm.BroadcastMessage(&message)
}

const (
	protocol = "tcp"
	port = ":8083"
)

func (s *Server) Start() {
	lis, err := net.Listen(protocol, port)
	if err != nil {
		fmt.Print(err)
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	pb.RegisterChatServiceServer(server, s)
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		fmt.Print(err)
		log.Fatalf("failed to serve: %v", err)
	}
}

func NewDefaultServer() *Server {
	return &Server{cm: NewClientsProvider()}
}

func NewServer(cm ClientManagerProvider) *Server {
	return &Server{cm: cm}
}