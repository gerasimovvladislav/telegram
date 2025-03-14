package telegram

import (
	"errors"
	"sync"
)

type Users struct {
	mu sync.RWMutex

	users map[UserId]*User
}

func NewUsers(users ...[]*User) *Users {
	usersMap := make(map[UserId]*User)
	for _, user := range users {
		for _, u := range user {
			usersMap[u.ID] = u
		}
	}

	return &Users{
		users: usersMap,
	}
}

func (u *Users) FindById(id UserId) (*User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	user, ok := u.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (u *Users) Update(user *User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.users[user.ID] = user

	return nil
}
