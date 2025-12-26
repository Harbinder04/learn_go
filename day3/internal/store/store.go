package store

import (
	"fmt"
	"log"
	"sync"
)

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type UserStore struct {
	mu sync.RWMutex
	users map[string]User
}

func NewUserStore() *UserStore {
    return &UserStore{
        users: make(map[string]User),
    }
} 


func (us *UserStore) Create(user User) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exists := us.users[user.ID]; exists {
		return fmt.Errorf("user with id %s already exists", user.ID)
	}
	log.Printf("User created with %s", user.ID)
	us.users[user.ID] = user
	return nil
}

func (us *UserStore) GetByID(id string) (User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

    if userVal, ok := us.users[id]; ok {
        return userVal, nil
    }
    
    // log.Printf("User with id [%v] doesn't exist", id)
    return User{}, fmt.Errorf("user with id %s not found", id)
}

func (us *UserStore) List() []User {
	us.mu.RLock()
	defer us.mu.RUnlock()

	values := make([]User, 0, len(us.users))

	for _, user := range us.users {
		values = append(values, user)
	}

	return values
}

