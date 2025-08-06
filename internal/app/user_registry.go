package app

import (
	"sync"
)

type UserRegistry struct {
	users map[string]interface{}
	mu    sync.RWMutex
}

func NewUserRegistry() *UserRegistry {
	return &UserRegistry{
		users: make(map[string]interface{}),
	}
}

func (r *UserRegistry) RegisterUser(username string, userData interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[username] = userData
}

func (r *UserRegistry) GetUser(username string) (interface{}, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[username]
	return user, ok
}
