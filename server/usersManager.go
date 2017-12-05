package main

import (
	"sync"
	"log"
	pb "grpc-go-chat/chat"
)

type UsersManager struct {
	users map[string]User
	mutex *sync.Mutex
}

func (u *UsersManager) Remove(id string) {
	log.Print("[USER manager] remove new user: ", id)
	u.mutex.Lock()
	delete(u.users, id)
	u.mutex.Unlock()
}

func (u *UsersManager) Add(user User) {
	log.Print("[USER manager] add new user: ", user)
	u.mutex.Lock()
	u.users[user.id] = user
	u.mutex.Unlock()
}

func (u* UsersManager) Find(id string) *User {
	log.Print("[USER manager] find user: ", id)
	u.mutex.Lock()
	defer u.mutex.Unlock()
	for _, user := range u.users {
		if user.id == id {
			return &user
		}
	}
	return nil
}

func (u *UsersManager) BroadcastMessage(message *pb.ChatMessage) {
	log.Print("[User manager] broadcast : ", *message)
	u.mutex.Lock()
	for _, user := range u.users {
		user.stream.Send(message)
	}
	u.mutex.Lock()
}

func NewUsersManager() *UsersManager {
	return &UsersManager{
		users: make(map[string]User),
		mutex: &sync.Mutex{},
	}
}