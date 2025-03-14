package telegram

type User struct {
	ID            UserId
	Place         *Place
	originalState *User
}

func NewAnonUser(userId UserId) *User {
	b := User{
		ID:    userId,
		Place: NewPlace(PlaceIdEmpty),
	}

	originalBotUser := b
	b.originalState = &originalBotUser

	return &b
}

// Reset сбрасывает все данные пользователя до состояния анонимного пользователя
func (u *User) Reset() {
	u.Place = NewPlace(PlaceIdEmpty)
}
func (u *User) Original() *User {
	return u.originalState
}
