package server

import (
	"sync"
	pb "grpc-go-chat/chat"
)

type ClientManagerProvider interface {
	Remove(id string)
	Add(user Client)
	Find(id string) *Client
	BroadcastMessage(message *pb.ChatMessage)
}

type ClientsManager struct {
	clients map[string]Client
	mutex   *sync.Mutex
}

func (cm *ClientsManager) Remove(id string) {
	cm.mutex.Lock()
	delete(cm.clients, id)
	cm.mutex.Unlock()
}

func (cm *ClientsManager) Add(client Client) {
	cm.mutex.Lock()
	cm.clients[client.User.Id] = client
	cm.mutex.Unlock()
}

func (cm *ClientsManager) Find(id string) *Client {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for _, c := range cm.clients {
		if c.User.Id == id {
			return &c
		}
	}
	return nil
}

func (cm *ClientsManager) BroadcastMessage(message *pb.ChatMessage) {
	cm.mutex.Lock()
	for _, c := range cm.clients {
		c.Stream.Send(message)
	}
	cm.mutex.Unlock()
}

func NewClientsProvider() *ClientsManager {
	return &ClientsManager{
		clients: make(map[string]Client),
		mutex:   &sync.Mutex{},
	}
}
