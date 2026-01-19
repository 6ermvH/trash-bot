package telegram

import (
	"context"
	"errors"
	"log"

	"github.com/6ermvH/trash-bot/internal/services/trashmanager"
	"github.com/go-telegram/bot"
)

func userErrorMessage(err error) string {
	switch {
	case errors.Is(err, trashmanager.ErrTryToAddUsers),
		errors.Is(err, trashmanager.ErrTryToInitialize):
		return err.Error()
	default:
		return "Request failed. Try again later."
	}
}

func (t *TgBotHandler) sendServiceError(
	ctx context.Context,
	botApi *bot.Bot,
	chatID int64,
	err error,
	logPrefix string,
) {
	if err == nil {
		return
	}

	log.Printf("%s: %v", logPrefix, err)
	t.sendMessage(
		ctx,
		botApi,
		chatID,
		userErrorMessage(err),
		logPrefix+" send message error",
	)
}
