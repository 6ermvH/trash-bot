package bot

import (
	"context"
	"fmt"

	"github.com/6ermvH/trash-bot/internal/config"
	"github.com/6ermvH/trash-bot/internal/handlers/telegram"
	"github.com/6ermvH/trash-bot/internal/repository/inmemory"
	"github.com/6ermvH/trash-bot/internal/services/trashmanager"
	"github.com/go-telegram/bot"
)

func Start(ctx context.Context, cfg *config.Config) error {

	opts := []bot.Option{}

	b, err := bot.New(cfg.Telegram.BotKey, opts...)
	if err != nil {
		return fmt.Errorf("init bot: %v", err)
	}

	repo := inmemory.New()
	trashm := trashmanager.New(repo)
	handlers := telegram.New(trashm)

	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, handlers.Start)
	b.RegisterHandler(bot.HandlerTypeMessageText, "set", bot.MatchTypeCommand, handlers.SetEstablish)
	b.RegisterHandler(bot.HandlerTypeMessageText, "next", bot.MatchTypeCommand, handlers.Next)
	b.RegisterHandler(bot.HandlerTypeMessageText, "prev", bot.MatchTypeCommand, handlers.Prev)
	b.RegisterHandler(bot.HandlerTypeMessageText, "who", bot.MatchTypeCommand, handlers.Who)

	b.Start(ctx)

	return nil
}
