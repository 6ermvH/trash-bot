package main

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gitlab.com/6ermvH/trash_bot/internal/config"
	"gitlab.com/6ermvH/trash_bot/internal/server"
	"gitlab.com/6ermvH/trash_bot/internal/store"
	"gitlab.com/6ermvH/trash_bot/internal/telegram"
)

func main() {
	// Load configuration from environment
	cfg, err := config.Load(os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("redis connect error: %v", err)
	}

	// Start healthcheck HTTP server
	go server.Start(cfg.Server)

	// Initialize storage manager
	storeMgr := store.NewStore(rdb)

	// Initialize and run Telegram bot
	botService := telegram.NewService(cfg.Telegram.Token, storeMgr, cfg.OpenRouterAPIKey)
	if err := botService.Run(context.Background()); err != nil {
		log.Fatalf("telegram bot error: %v", err)
	}
}
