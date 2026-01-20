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

func (t *TgBotHandler) getTimeSelectionKeyboard(botApi *bot.Bot) *inline.Keyboard {
	keyboard := inline.New(botApi).
		Row().
		Button("06:00", []byte("06:00"), t.handleTimeSelection).
		Button("07:00", []byte("07:00"), t.handleTimeSelection).
		Button("08:00", []byte("08:00"), t.handleTimeSelection).
		Button("09:00", []byte("09:00"), t.handleTimeSelection).
		Row().
		Button("10:00", []byte("10:00"), t.handleTimeSelection).
		Button("11:00", []byte("11:00"), t.handleTimeSelection).
		Button("12:00", []byte("12:00"), t.handleTimeSelection).
		Button("13:00", []byte("13:00"), t.handleTimeSelection).
		Row().
		Button("14:00", []byte("14:00"), t.handleTimeSelection).
		Button("15:00", []byte("15:00"), t.handleTimeSelection).
		Button("16:00", []byte("16:00"), t.handleTimeSelection).
		Button("17:00", []byte("17:00"), t.handleTimeSelection).
		Row().
		Button("18:00", []byte("18:00"), t.handleTimeSelection).
		Button("19:00", []byte("19:00"), t.handleTimeSelection).
		Button("20:00", []byte("20:00"), t.handleTimeSelection).
		Button("21:00", []byte("21:00"), t.handleTimeSelection).
		Row().
		Button("22:00", []byte("22:00"), t.handleTimeSelection).
		Button("23:00", []byte("23:00"), t.handleTimeSelection)

	return keyboard
}

func (t *TgBotHandler) handleTimeSelection(
	ctx context.Context,
	botApi *bot.Bot,
	mes models.MaybeInaccessibleMessage,
	data []byte,
) {
	chatID := mes.Message.Chat.ID
	selectedTime := string(data)

	if err := t.service.Subscribe(ctx, chatID, selectedTime); err != nil {
		t.sendServiceError(ctx, botApi, chatID, err, "TimeSelection Subscribe")

		return
	}

	// Удаляем сообщение с клавиатурой
	_, _ = botApi.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatID,
		MessageID: mes.Message.ID,
	})

	t.sendMessage(
		ctx,
		botApi,
		chatID,
		"✅ Вы подписались на ежедневные напоминания в "+selectedTime,
		"TimeSelection send message error",
	)
}
