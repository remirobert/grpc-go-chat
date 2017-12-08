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
	cm := NewClientsProvider()
	s := new(mockStream)
	c := Client{stream: s}
	cm.Add(c)
	userFound := cm.Find(c.user.Id)
	if userFound == nil {
		t.Errorf("The user [%s - %s] wasn't added.", c.user.Username, c.user.Id)
	}2
}

func TestUsers_Remove(t *testing.T) {
	cm := NewClientsProvider()
	s := new(mockStream)
	c := Client{stream: s}
	cm.Add(c)
	cm.Remove(c.user.Id)
	userFound := cm.Find(c.user.Id)
	if userFound != nil {
		t.Errorf("The user [%s - %s] should be removed.", c.user.Username, c.user.Id)
	}
}

func TestUsers_Find(t *testing.T) {
	cm := NewClientsProvider()
	userFound := cm.Find("fake")
	if userFound != nil {
		t.Errorf("No user should be find.")
	}
}


func TestUsers_BroadcastMessage(t *testing.T) {
	cm := NewClientsProvider()
	s := new(mockStream)
	c := Client{stream: s}
	cm.Add(c)
	m := pb.ChatMessage{User: &pb.User{Id:"123", Username:"remi"}, Type: pb.ChatMessage_USER_JOIN}
	cm.BroadcastMessage(&m)
	if s.rcvMessage == nil {
		t.Errorf("The stream didn't receive the message.", c.user.Username, c.user.Id)
	}
}
