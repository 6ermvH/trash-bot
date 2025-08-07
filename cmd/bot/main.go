package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/redis/go-redis/v9"
	"gitlab.com/6ermvH/trash_bot/internal/config"
	"gitlab.com/6ermvH/trash_bot/internal/server"
	"gitlab.com/6ermvH/trash_bot/internal/store"
	"gitlab.com/6ermvH/trash_bot/internal/telegram"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load configuration from environment
	cfg, err := config.Load(os.Getenv("CONFIG_PATH"))
	if err != nil {
		logger.Error("config load error", "error", err)
		os.Exit(1)
	}

	logger.Info("config loaded")

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logger.Error("redis connect error", "error", err)
		os.Exit(1)
	}
	logger.Info("redis connected")

	// Start healthcheck HTTP server
	go server.Start(cfg.Server)
	logger.Info("healthcheck server started")

	// Initialize storage manager
	storeMgr := store.NewStore(rdb, logger)
	logger.Info("storage manager initialized")

	// Initialize and run Telegram bot
	botService := telegram.NewService(logger, cfg.Telegram.Token, storeMgr, cfg.OpenRouterAPIKey)
	logger.Info("telegram bot service initialized, starting...")
	if err := botService.Run(context.Background()); err != nil {
		logger.Error("telegram bot error", "error", err)
		os.Exit(1)
	}
}
