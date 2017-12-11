package grpc_go_chat

import (
	pb "grpc-go-chat/chat"
	. "grpc-go-chat/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type spyClientsManager struct {
	removedClient       *string
	addedClient         *string
	foundClient         *string
	broadcastedMessages []*pb.ChatMessage
	resultFindClient    *Client
}

func (cm *spyClientsManager) Remove(id string) {
	cm.removedClient = &id
}

func (cm *spyClientsManager) Add(user Client) {
	cm.addedClient = &user.User.Id
}

func (cm *spyClientsManager) Find(id string) *Client {
	cm.foundClient = &id
	return cm.resultFindClient
}

func (cm *spyClientsManager) BroadcastMessage(message *pb.ChatMessage) {
	cm.broadcastedMessages = append(cm.broadcastedMessages, message)
}

func newSpyClientsManager(foundClient *Client) *spyClientsManager {
	return &spyClientsManager{
		broadcastedMessages: []*pb.ChatMessage{},
		resultFindClient:    foundClient,
	}
}

var _ = Describe("Server", func() {
	Context("when a new client first connecting to the chat", func() {
		It("the server should add it to the client manager and broadcast a join message", func() {
			scm := newSpyClientsManager(nil)
			s := NewServer(scm)

			user := &pb.User{Id: "123", Username: ""}
			m1 := pb.ChatMessage{
				Type: pb.ChatMessage_USER_JOIN,
				User: user,
			}
			m2 := pb.ChatMessage{
				Type:    pb.ChatMessage_USER_CHAT,
				User:    user,
				Message: &pb.Message{Content: "hello"},
			}
			m3 := pb.ChatMessage{
				Type: pb.ChatMessage_USER_LEAVE,
				User: user,
			}
			messages := []*pb.ChatMessage{&m1, &m2, &m3}
			msv := NewMockServerStream(messages)
			err := s.Stream(msv)

			Expect(err).To(BeNil())
			Expect(scm.addedClient).ToNot(BeNil())
			Expect(*scm.addedClient).To(Equal("123"))
			Expect(scm.broadcastedMessages).ToNot(BeEmpty())
			Expect(scm.broadcastedMessages).To(Equal(messages))
		})
		It("if the user is not logged, the server should return an error", func() {
			scm := newSpyClientsManager(nil)
			s := NewServer(scm)

			m := pb.ChatMessage{
				Type: pb.ChatMessage_USER_CHAT,
				User: &pb.User{Id: "123", Username: ""},
			}
			messages := []*pb.ChatMessage{&m}
			msv := NewMockServerStream(messages)
			err := s.Stream(msv)

			Expect(scm.addedClient).To(BeNil())
			Expect(scm.broadcastedMessages).To(BeEmpty())
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(NewAuthError(AuthMessageUserNotJoined)))
		})
		It("if the server has already joined the server should return an error", func() {
			client := &Client{User: pb.User{Id: "123"}}
			scm := newSpyClientsManager(client)
			s := NewServer(scm)

			m := pb.ChatMessage{
				Type: pb.ChatMessage_USER_JOIN,
				User: &pb.User{Id: "123"},
			}
			messages := []*pb.ChatMessage{&m}
			msv := NewMockServerStream(messages)
			err := s.Stream(msv)

			Expect(scm.foundClient).ToNot(BeNil())
			Expect(*scm.foundClient).To(Equal("123"))
			Expect(scm.addedClient).To(BeNil())
			Expect(scm.broadcastedMessages).To(BeEmpty())
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(NewAuthError(AuthMessageUserAlreadyJoined)))
		})
	})
})
