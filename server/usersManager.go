package main

import "sync"

type UsersManager struct {
	users map[string]User
	mutex *sync.Mutex
}

func (u *UsersManager) Remove(id string) {
	u.mutex.Lock()
	delete(u.users, id)
	u.mutex.Unlock()
}

func (u *UsersManager) Add(user User) {
	u.mutex.Lock()
	u.users[user.id] = user
	u.mutex.Unlock()
}

func (u* UsersManager) Find(id string) *User {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	for _, user := range u.users {
		if user.id == id {
			return &user
		}
	}
	return nil
}

func NewUsers() *UsersManager {
	return &UsersManager{
		users: make(map[string]User),
		mutex: &sync.Mutex{},
	}
}