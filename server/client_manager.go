package server

import (
	"sync"
	pb "grpc-go-chat/chat"
)

type ClientManager interface {
	Remove(id string)
	Add(user Client)
	Find(id string) *Client
	BroadcastMessage(message *pb.ChatMessage)
}

type Clients struct {
	clients map[string]Client
	mutex   *sync.Mutex
}

func (cm *Clients) Remove(id string) {
	cm.mutex.Lock()
	delete(cm.clients, id)
	cm.mutex.Unlock()
}

func (cm *Clients) Add(client Client) {
	cm.mutex.Lock()
	cm.clients[client.User.Id] = client
	cm.mutex.Unlock()
}

func (cm *Clients) Find(id string) *Client {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for _, c := range cm.clients {
		if c.User.Id == id {
			return &c
		}
	}
	return nil
}

func (cm *Clients) BroadcastMessage(message *pb.ChatMessage) {
	cm.mutex.Lock()
	for _, c := range cm.clients {
		c.Stream.Send(message)
	}
	cm.mutex.Unlock()
}

func NewClientsManager() *Clients {
	return &Clients{
		clients: make(map[string]Client),
		mutex:   &sync.Mutex{},
	}
}
