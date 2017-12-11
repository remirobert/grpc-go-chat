package grpc_go_chat

import (
	pb "grpc-go-chat/chat"
	. "grpc-go-chat/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientsProvider impl tests", func() {
	Context("When adding a new Client", func() {
		It("should be able to find the added user from the id", func() {
			cm := NewClientsProvider()
			s := NewMockServerStream(nil)

			c1 := NewClient(pb.User{Id: "1"}, s)
			c2 := NewClient(pb.User{Id: "2"}, s)
			cm.Add(*c1)
			cm.Add(*c2)

			clientFound1 := cm.Find(c1.User.Id)
			clientFound2 := cm.Find(c2.User.Id)

			Expect(clientFound1).ToNot(BeNil())
			Expect(clientFound1).To(Equal(c1))
			Expect(clientFound2).ToNot(BeNil())
			Expect(clientFound2).To(Equal(c2))
		})
		It("If the client doesn't exist should return nil", func() {
			cm := NewClientsProvider()
			s := NewMockServerStream(nil)

			c1 := NewClient(pb.User{Id: "1"}, s)
			cm.Add(*c1)

			clientFound := cm.Find("fake")

			Expect(clientFound).To(BeNil())
		})
	})
	Context("When removing a Client", func() {
		It("it should be removed in the clients manager", func() {
			cm := NewClientsProvider()
			s := NewMockServerStream(nil)

			c1 := NewClient(pb.User{Id: "1"}, s)
			c2 := NewClient(pb.User{Id: "2"}, s)
			cm.Add(*c1)
			cm.Add(*c2)

			cm.Remove(c1.User.Id)
			clientFound1 := cm.Find(c1.User.Id)
			clientFound2 := cm.Find(c2.User.Id)

			Expect(clientFound1).To(BeNil())
			Expect(clientFound2).ToNot(BeNil())
			Expect(clientFound2).To(Equal(c2))
		})
	})
	Context("when broadcasting a message", func() {
		It("the clients's stream should Receive it", func() {
			cm := NewClientsProvider()
			ms := NewMockServerStream(nil)

			c := Client{User: pb.User{Id: "1"}, Stream: ms}
			cm.Add(c)

			u := &pb.User{Id: "1", Username: ""}
			m1 := pb.ChatMessage{User: u, Type: pb.ChatMessage_USER_JOIN}
			m2 := pb.ChatMessage{
				User:    u,
				Type:    pb.ChatMessage_USER_LEAVE,
				Message: &pb.Message{Content: "hello"},
			}
			m3 := pb.ChatMessage{User: u, Type: pb.ChatMessage_USER_LEAVE}

			cm.BroadcastMessage(&m1)
			cm.BroadcastMessage(&m2)
			cm.BroadcastMessage(&m3)

			expectedMessages := []*pb.ChatMessage{&m1, &m2, &m3}

			Expect(ms.sentMessages).ToNot(BeEmpty())
			Expect(ms.sentMessages).To(Equal(expectedMessages))
		})
	})
})
