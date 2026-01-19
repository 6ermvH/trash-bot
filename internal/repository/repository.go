package repository

import "errors"

var (
	ErrChatIsEmpty         = errors.New("chat don`t have someone user in list")
	ErrChatIsNotInitialize = errors.New("chat don`t initialize manager")
)

type Chat struct {
	ID      int64    `json:"id"`
	Current int      `json:"currentUser"`
	Users   []string `json:"activeUsers"`
}
