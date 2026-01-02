package repository

import "errors"

var (
	ErrChatIsEmpty         = errors.New("chat don`t have someone user in list")
	ErrChatIsNotInitialize = errors.New("chat don`t initialize manager")
)
