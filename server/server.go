package server

import (
	pb "grpc-go-chat/chat"
	"net"
	"fmt"
	"log"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	cm ClientManagerProvider
}

func (s *Server) Stream(stream pb.ChatService_StreamServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok != true {
		log.Print("impossible to get the md")
	} else {
		log.Print("md from the context : ", md)
	}
	request, err := stream.Recv()
	if err != nil {
		return err
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

		switch message.Type {
		case pb.ChatMessage_USER_JOIN:
			return NewAuthError(AuthMessageUserAlreadyJoined)
		case pb.ChatMessage_USER_LEAVE:
			s.processLeaveUser(message)
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
	log.Print("new user joined the channel : ", message.User.Username)
	s.broadcastUserJoin(message.User)
	return nil
}

func (s *Server) processLeaveUser(message *pb.ChatMessage) {
	log.Print("remove the user : ", message)
	if message == nil || message.User == nil {
		return
	}
	s.cm.Remove(message.User.Id)
	log.Print("new user left the channel : ", message.User.Username)
	s.broadcastUserLeave(message.User)
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

func (s *Server) Start() {
	lis, err := net.Listen("tcp", ":8083")
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

func NewServer() *Server {
	return &Server{cm: NewClientsProvider()}
}