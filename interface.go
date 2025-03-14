package telegram

import "context"

type UserStorage interface {
	FindById(id UserId) (*User, error)
	Update(user *User) error
}

type Bot interface {
	// Start запускает бота
	Start(ctx context.Context) error
}
