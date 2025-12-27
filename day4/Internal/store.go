package internal

import (
 "fmt"
)

func NewUserStore() *UserStore {
	return &UserStore{
		users: map[string]User{},
	}
}

func (us *UserStore) Create(user User) error {
	us.mu.Lock()
	defer us.mu.Unlock()

	if _, exists := us.users[user.Id]; exists {
		return fmt.Errorf("user with id %s already exists", user.Id)
	}
	
	us.users[user.Id] = user
	return nil
}

func (us *UserStore) GetByID(id string) (User, error) {
	us.mu.RLock()
	defer us.mu.RUnlock()

    if userVal, ok := us.users[id]; ok {
        return userVal, nil
    }

    return User{}, fmt.Errorf("user with id %s not found", id)
}

func (u *UserStore) GetAllUser() []User{
	u.mu.RLock()
	defer u.mu.RUnlock()
	
	users := make([]User, 0, len(u.users))

	for _,v := range u.users {
		users = append(users, v)
	}

	return users
}