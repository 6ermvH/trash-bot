package telegram

import (
	"context"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service interface {
	Who(ctx context.Context, chatID int64) (string, error)
	Next(ctx context.Context, chatID int64) (string, error)
	Prev(ctx context.Context, chatID int64) (string, error)
	SetEstablish(ctx context.Context, chatID int64, users []string) error
	Subscribe(ctx context.Context, chatID int64) error
	Unsubscribe(ctx context.Context, chatID int64) error
}

type TgBotHandler struct {
	service Service
}

func New(service Service) *TgBotHandler {
	return &TgBotHandler{
		service: service,
	}
}

func (t *TgBotHandler) Start(ctx context.Context, botApi *bot.Bot, update *models.Update) {
	_, err := botApi.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Привет!",
		ReplyMarkup: t.getKeyboardOnStart(botApi),
	})
	if err != nil {
		log.Printf("Start. send message: %v", err.Error())
	}
}

func (t *TgBotHandler) SetEstablish(ctx context.Context, botApi *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	parts := strings.Fields(update.Message.Text)
	if len(parts) < 2 {
		t.sendMessage(
			ctx,
			botApi,
			chatID,
			"Provide at least one username after /set",
			"SetEstablish send message error",
		)
		return
	}
	users := parts[1:]

	if err := t.service.SetEstablish(ctx, chatID, users); err != nil {
		t.sendServiceError(ctx, botApi, chatID, err, "SetEstablish")
		return
	}

	_, err := botApi.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Establish success",
	})
	if err != nil {
		log.Printf("SetEstablish. send message: %v", err.Error())
	}
}

func (t *TgBotHandler) Next(ctx context.Context, botApi *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	username, err := t.service.Next(ctx, chatID)
	if err != nil {
		t.sendServiceError(ctx, botApi, chatID, err, "Next")
		return
	}

	if _, err := botApi.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   username,
	}); err != nil {
		log.Printf("Next. send message: %v", err.Error())
	}
}

func (t *TgBotHandler) Prev(ctx context.Context, botApi *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	username, err := t.service.Prev(ctx, chatID)
	if err != nil {
		t.sendServiceError(ctx, botApi, chatID, err, "Prev")
		return
	}

	if _, err := botApi.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   username,
	}); err != nil {
		log.Printf("Prev. send message: %v", err.Error())
	}
}

func (t *TgBotHandler) Who(ctx context.Context, botApi *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	username, err := t.service.Who(ctx, chatID)
	if err != nil {
		t.sendServiceError(ctx, botApi, chatID, err, "Who")
		return
	}

	if _, err := botApi.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   username,
	}); err != nil {
		log.Printf("Who. send message: %v", err.Error())
	}
}
