package main

import (
	"sync"
	"log"
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
	log.Print("[USER manager] remove new user: ", id)
	cm.mutex.Lock()
	delete(cm.clients, id)
	cm.mutex.Unlock()
}

func (cm *ClientsManager) Add(client Client) {
	log.Print("[USER manager] add new user: ", client)
	cm.mutex.Lock()
	cm.clients[client.user.Id] = client
	cm.mutex.Unlock()
}

func (cm *ClientsManager) Find(id string) *Client {
	log.Print("[USER manager] find user: ", id)
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for _, c := range cm.clients {
		if c.user.Id == id {
			return &c
		}
	}
	return nil
}

func (cm *ClientsManager) BroadcastMessage(message *pb.ChatMessage) {
	log.Print("[User manager] broadcast : ", *message)
	cm.mutex.Lock()
	for _, c := range cm.clients {
		c.stream.Send(message)
	}
	cm.mutex.Unlock()
}

func NewClientsProvider() *ClientsManager {
	return &ClientsManager{
		clients: make(map[string]Client),
		mutex:   &sync.Mutex{},
	}
}
