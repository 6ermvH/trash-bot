package main

import (
	"context"
	"log"

	"github.com/6ermvH/trash-bot/cmd/bot"
	"github.com/6ermvH/trash-bot/cmd/panel"
	"github.com/6ermvH/trash-bot/internal/config"
	"github.com/6ermvH/trash-bot/internal/repository/inmemory"
	"github.com/6ermvH/trash-bot/internal/repository/sqlite"
	"github.com/6ermvH/trash-bot/internal/services/trashmanager"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := config.LoadEnvFile(".env"); err != nil {
		log.Fatalf("load .env file: %v", err)
	}

	cfg, err := config.NewFromFile("config/base.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	group, ctx := errgroup.WithContext(context.Background())

	var trashm *trashmanager.Service

	switch cfg.Database.Type {
	case "sqlite":
		repo, err := sqlite.New(cfg.Database.Path)
		if err != nil {
			log.Fatalf("failed to open sqlite db: %v", err)
		}

		defer func() {
			if err := repo.Close(); err != nil {
				log.Printf("close sqlite db: %v", err)
			}
		}()

		trashm = trashmanager.New(repo)

		log.Printf("Using SQLite database: %s\n", cfg.Database.Path)
	default:
		repo := inmemory.New()

		trashm = trashmanager.New(repo)

		log.Println("Using in-memory database")
	}

	if cfg.Server.Enabled {
		group.Go(func() error {
			return panel.Start(ctx, cfg, trashm)
		})
		log.Printf("Server started on port: %s\n", cfg.Server.Port)
	}

	group.Go(func() error {
		return bot.Start(ctx, cfg, trashm)
	})
	log.Printf("Bot started\n")

	if err := group.Wait(); err != nil {
		log.Printf("application ended with error: %v\n", err)

		return
	}
}
