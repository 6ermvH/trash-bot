package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func addUser(chatID int64, username string) error {
	key := fmt.Sprintf("chat:%d:users", chatID)
	return rdb.RPush(ctx, key, username).Err()
}

func getUsers(chatID int64) ([]string, error) {
	key := fmt.Sprintf("chat:%d:users", chatID)
	return rdb.LRange(ctx, key, 0, -1).Result()
}

func setActiveIndex(chatID int64, index int) error {
	key := fmt.Sprintf("chat:%d:active_index", chatID)
	return rdb.Set(ctx, key, index, 0).Err()
}

func getActiveIndex(chatID int64) (int, error) {
	key := fmt.Sprintf("chat:%d:active_index", chatID)
	return rdb.Get(ctx, key).Int()
}

func clearUserList(chatID int64) error {
	key := fmt.Sprintf("chat:%d:users", chatID)
	return rdb.Del(ctx, key).Err()
}
