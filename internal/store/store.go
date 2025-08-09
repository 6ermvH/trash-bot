package store

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	rdb    *redis.Client
	ctx    context.Context
	logger *slog.Logger
}

// NewStore returns a Store wrapping the given Redis client
func NewStore(rdb *redis.Client, logger *slog.Logger) *Store {
	return &Store{
		rdb:    rdb,
		ctx:    context.Background(),
		logger: logger,
	}
}

// AddUser appends a username to the chat's user list if it's not already present
func (s *Store) AddUser(chatID int64, username string) error {
	key := fmt.Sprintf("chat:%d:users", chatID)
	s.logger.Info("adding user", "chat_id", chatID, "username", username)

	// Check if user already exists to prevent duplicates
	users, err := s.GetUsers(chatID)
	if err != nil {
		s.logger.Error("failed to get users", "error", err)
		return err
	}
	for _, u := range users {
		if u == username {
			err := fmt.Errorf("user %s already exists in the list", username)
			s.logger.Warn("user already exists", "error", err)
			return err
		}
	}

	return s.rdb.RPush(s.ctx, key, username).Err()
}

// RemoveUser removes a username from the chat's user list
func (s *Store) RemoveUser(chatID int64, username string) error {
	key := fmt.Sprintf("chat:%d:users", chatID)
	s.logger.Info("removing user", "chat_id", chatID, "username", username)
	// LREM 0 means remove all occurrences of the value
	return s.rdb.LRem(s.ctx, key, 0, username).Err()
}

// GetUsers returns all usernames in the chat's user list
func (s *Store) GetUsers(chatID int64) ([]string, error) {
	key := fmt.Sprintf("chat:%d:users", chatID)
	return s.rdb.LRange(s.ctx, key, 0, -1).Result()
}

// SetActiveIndex sets the current index for the chat
func (s *Store) SetActiveIndex(chatID int64, index int) error {
	key := fmt.Sprintf("chat:%d:active_index", chatID)
	s.logger.Info("setting active index", "chat_id", chatID, "index", index)
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
	s.logger.Info("clearing users", "chat_id", chatID)
	return s.rdb.Del(s.ctx, key).Err()
}

// SetEstablish clears and sets a new ordered user list for the chat
func (s *Store) SetEstablish(chatID int64, users []string) error {
	s.logger.Info("setting establish", "chat_id", chatID, "users", users)
	key := fmt.Sprintf("chat:%d:users", chatID)

	// Use a transaction to ensure atomicity
	pipe := s.rdb.TxPipeline()
	pipe.Del(s.ctx, key)
	if len(users) > 0 {
		// Convert []string to []interface{} for RPush
		userInterfaces := make([]interface{}, len(users))
		for i, u := range users {
			userInterfaces[i] = u
		}
		pipe.RPush(s.ctx, key, userInterfaces...)
	}
	pipe.Set(s.ctx, fmt.Sprintf("chat:%d:active_index", chatID), 0, 0)

	_, err := pipe.Exec(s.ctx)
	if err != nil {
		s.logger.Error("failed to set establish in transaction", "error", err)
		return err
	}

	return nil
}

// Next advances the active index and returns the new username
func (s *Store) Next(chatID int64) (string, error) {
	users, err := s.GetUsers(chatID)
	if err != nil || len(users) == 0 {
		s.logger.Error("failed to get users for next", "error", err)
		return "", err
	}
	idx, err := s.GetActiveIndex(chatID)
	if err != nil {
		s.logger.Error("failed to get active index for next", "error", err)
		return "", err
	}
	idx = (idx + 1) % len(users)
	if err := s.SetActiveIndex(chatID, idx); err != nil {
		s.logger.Error("failed to set active index for next", "error", err)
		return "", err
	}
	return users[idx], nil
}

// Prev moves back the active index and returns the new username
func (s *Store) Prev(chatID int64) (string, error) {
	users, err := s.GetUsers(chatID)
	if err != nil || len(users) == 0 {
		s.logger.Error("failed to get users for prev", "error", err)
		return "", err
	}
	idx, err := s.GetActiveIndex(chatID)
	if err != nil {
		s.logger.Error("failed to get active index for prev", "error", err)
		return "", err
	}
	if idx == 0 {
		idx = len(users) - 1
	} else {
		idx--
	}
	if err := s.SetActiveIndex(chatID, idx); err != nil {
		s.logger.Error("failed to set active index for prev", "error", err)
		return "", err
	}
	return users[idx], nil
}

// Who returns the current active username
func (s *Store) Who(chatID int64) (string, error) {
	users, err := s.GetUsers(chatID)
	if err != nil || len(users) == 0 {
		s.logger.Error("failed to get users for who", "error", err)
		return "", err
	}
	idx, err := s.GetActiveIndex(chatID)
	if err != nil {
		s.logger.Error("failed to get active index for who", "error", err)
		return "", err
	}
	return users[idx], nil
}
