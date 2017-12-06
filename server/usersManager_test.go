package main

import (
	pb "grpc-go-chat/chat"
	"testing"
	"google.golang.org/grpc"
)

type mockStream struct {
	rcvMessage *pb.ChatMessage
	grpc.ServerStream
}

func (s *mockStream) Send(m *pb.ChatMessage) error {
	s.rcvMessage = m
	return nil
}

func (s *mockStream) Recv() (*pb.ChatMessage, error) {
	return nil, nil
}

func TestUsers_Add(t *testing.T) {
	um := NewUsers()
	s := new(mockStream)
	newUser := User{stream: s}
	um.Add(newUser)
	userFound := um.Find(newUser.id)
	if userFound == nil {
		t.Errorf("The user [%s - %s] wasn't added.", newUser.username, newUser.id)
	}
}

func TestUsers_BroadcastMessage(t *testing.T) {
	um := NewUsers()
	s := new(mockStream)
	newUser := User{stream: s}
	um.Add(newUser)
	m := pb.ChatMessage{User: &pb.User{Id:"123", Username:"remi"}, Type: pb.ChatMessage_USER_JOIN}
	um.BroadcastMessage(&m)
	if s.rcvMessage == nil {
		t.Errorf("The stream didn't receive the message.", newUser.username, newUser.id)
	}
}
