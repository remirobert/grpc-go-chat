package grpc_go_chat

import (
	"google.golang.org/grpc"
	pb "grpc-go-chat/chat"
)

type MockStreamServer struct {
	sentMessages []*pb.ChatMessage
	respMessages []*pb.ChatMessage
	grpc.ServerStream
}

func (ms *MockStreamServer) Send(m *pb.ChatMessage) error {
	ms.sentMessages = append(ms.sentMessages, m)
	return nil
}

func (ms *MockStreamServer) Recv() (*pb.ChatMessage, error) {
	len := len(ms.respMessages)
	if len == 0 {
		return nil, nil
	}
	m := ms.respMessages[0]
	ms.respMessages = ms.respMessages[1:len]
	return m, nil
}

func NewMockServerStream(respMessages []*pb.ChatMessage) *MockStreamServer {
	return &MockStreamServer{
		sentMessages:[]*pb.ChatMessage{},
		respMessages:respMessages,
	}
}
