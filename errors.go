package telegram

// NewUnresolvedCommandError создает новую ошибку для случая не допустимого действия пользователя.
func NewUnresolvedCommandError(err error) *UnresolvedCommandError {
	return &UnresolvedCommandError{err: err}
}

// UnresolvedCommandError ошибка для случая недопустимого действия пользователя на его позиции в боте(place). Такие ошибки обычно требуется
// по-особенному обрабатывать, и благодаря такой ошибке появляется такая возможность.
//
// Важно, данную ошибку стоит применить лишь тогда, когда пользователь пытается вызвать действие, которое не предусмотренно выполнять из его места.
// Например:
// - пользователь ввёл номер телефона, чтобы авторизоваться, и бот ожидает, что пользователь введёт код из смс, но пользователь нажимает кнопку "Главное меню"
type UnresolvedCommandError struct {
	err error
}

// Error возвращает текст ошибки
func (e *UnresolvedCommandError) Error() string {
	return e.err.Error()
}
