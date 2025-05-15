package store

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	rdb *redis.Client
	ctx context.Context
}

// NewStore returns a Store wrapping the given Redis client
func NewStore(rdb *redis.Client) *Store {
	return &Store{
		rdb: rdb,
		ctx: context.Background(),
	}
}

// AddUser appends a username to the chat's user list
func (s *Store) AddUser(chatID int64, username string) error {
	key := fmt.Sprintf("chat:%d:users", chatID)
	return s.rdb.RPush(s.ctx, key, username).Err()
}

// GetUsers returns all usernames in the chat's user list
func (s *Store) GetUsers(chatID int64) ([]string, error) {
	key := fmt.Sprintf("chat:%d:users", chatID)
	return s.rdb.LRange(s.ctx, key, 0, -1).Result()
}

// SetActiveIndex sets the current index for the chat
func (s *Store) SetActiveIndex(chatID int64, index int) error {
	key := fmt.Sprintf("chat:%d:active_index", chatID)
	return s.rdb.Set(s.ctx, key, index, 0).Err()
}

// GetActiveIndex returns the current index for the chat
func (s *Store) GetActiveIndex(chatID int64) (int, error) {
	key := fmt.Sprintf("chat:%d:active_index", chatID)
	return s.rdb.Get(s.ctx, key).Int()
}

// ClearUsers deletes the user list for the chat
func (s *Store) ClearUsers(chatID int64) error {
	key := fmt.Sprintf("chat:%d:users", chatID)
	return s.rdb.Del(s.ctx, key).Err()
}

// SetEstablish clears and sets a new ordered user list for the chat
func (s *Store) SetEstablish(chatID int64, users []string) error {
	if err := s.ClearUsers(chatID); err != nil {
		return err
	}
	if err := s.SetActiveIndex(chatID, 0); err != nil {
		return err
	}
	for _, u := range users {
		if err := s.AddUser(chatID, u); err != nil {
			return err
		}
	}
	return nil
}

// Next advances the active index and returns the new username
func (s *Store) Next(chatID int64) (string, error) {
	users, err := s.GetUsers(chatID)
	if err != nil || len(users) == 0 {
		return "", err
	}
	idx, err := s.GetActiveIndex(chatID)
	if err != nil {
		return "", err
	}
	idx = (idx + 1) % len(users)
	if err := s.SetActiveIndex(chatID, idx); err != nil {
		return "", err
	}
	return users[idx], nil
}

// Prev moves back the active index and returns the new username
func (s *Store) Prev(chatID int64) (string, error) {
	users, err := s.GetUsers(chatID)
	if err != nil || len(users) == 0 {
		return "", err
	}
	idx, err := s.GetActiveIndex(chatID)
	if err != nil {
		return "", err
	}
	if idx == 0 {
		idx = len(users) - 1
	} else {
		idx--
	}
	if err := s.SetActiveIndex(chatID, idx); err != nil {
		return "", err
	}
	return users[idx], nil
}

// Who returns the current active username
func (s *Store) Who(chatID int64) (string, error) {
	users, err := s.GetUsers(chatID)
	if err != nil || len(users) == 0 {
		return "", err
	}
	idx, err := s.GetActiveIndex(chatID)
	if err != nil {
		return "", err
	}
	return users[idx], nil
}
