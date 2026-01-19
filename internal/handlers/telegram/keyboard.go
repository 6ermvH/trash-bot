package telegram

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func (t *TgBotHandler) getKeyboardOnStart(botApi *bot.Bot) *inline.Keyboard {
	keyboard := inline.New(botApi).Row().
		Button("Кто выносит", []byte(""), t.handleWho).
		Row().
		Button("Следующий", []byte(""), t.handleNext).
		Row().
		Button("Предыдущий", []byte(""), t.handlePrev)

	return keyboard
}

func (t *TgBotHandler) handleWho(
	ctx context.Context,
	botApi *bot.Bot,
	mes models.MaybeInaccessibleMessage,
	data []byte,
) {
	username, err := t.service.Who(ctx, mes.Message.Chat.ID)
	if err != nil {
		t.sendServiceError(ctx, botApi, mes.Message.Chat.ID, err, "Inline Who")
		return
	}

	t.sendMessage(
		ctx,
		botApi,
		mes.Message.Chat.ID,
		"Мусор выносит: "+username,
		"Handler Who send message error",
	)
}

func (t *TgBotHandler) handleNext(
	ctx context.Context,
	botApi *bot.Bot,
	mes models.MaybeInaccessibleMessage,
	data []byte,
) {
	username, err := t.service.Next(ctx, mes.Message.Chat.ID)
	if err != nil {
		t.sendServiceError(ctx, botApi, mes.Message.Chat.ID, err, "Inline Next")
		return
	}

	t.sendMessage(
		ctx,
		botApi,
		mes.Message.Chat.ID,
		"Мусор выносит: "+username,
		"Handler Next send message error",
	)
}

func (t *TgBotHandler) handlePrev(
	ctx context.Context,
	botApi *bot.Bot,
	mes models.MaybeInaccessibleMessage,
	data []byte,
) {
	username, err := t.service.Prev(ctx, mes.Message.Chat.ID)
	if err != nil {
		t.sendServiceError(ctx, botApi, mes.Message.Chat.ID, err, "Inline Prev")
		return
	}

	t.sendMessage(
		ctx,
		botApi,
		mes.Message.Chat.ID,
		"Мусор выносит: "+username,
		"Handler Prev send message error",
	)
}

func (t *TgBotHandler) sendMessage(
	ctx context.Context,
	botApi *bot.Bot,
	chatID int64,
	text string,
	errPrefix string,
) {
	if _, err := botApi.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		log.Printf("%s: %v", errPrefix, err.Error())
	}
}
