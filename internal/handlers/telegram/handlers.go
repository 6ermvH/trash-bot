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

func (t *TgBotHandler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Привет!",
		ReplyMarkup: t.getKeyboardOnStart(b),
	})
	if err != nil {
		log.Printf("Handler Start end with: %v", err.Error())
	}
}

func (t *TgBotHandler) SetEstablish(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	users := strings.Split(update.Message.Text, " ")[1:]

	if err := t.service.SetEstablish(ctx, chatID, users); err != nil {
		log.Printf("Handler Start end with: %v", err.Error())
	}
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Establish success",
	})
	if err != nil {
		log.Printf("Handler Start end with: %v", err.Error())
	}
}

func (t *TgBotHandler) Next(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	username, err := t.service.Next(ctx, chatID)
	if err != nil {
		log.Printf("Handler Next end with: %v", err.Error())
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   username,
	}); err != nil {
		log.Printf("Handler Start end with: %v", err.Error())
	}
}

func (t *TgBotHandler) Prev(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	username, err := t.service.Prev(ctx, chatID)
	if err != nil {
		log.Printf("Handler Prev end with: %v", err.Error())
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   username,
	}); err != nil {
		log.Printf("Handler Prev end with: %v", err.Error())
	}
}

func (t *TgBotHandler) Who(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	username, err := t.service.Who(ctx, chatID)
	if err != nil {
		log.Printf("Handler Who end with: %v", err.Error())
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   username,
	}); err != nil {
		log.Printf("Handler Start end with: %v", err.Error())
	}
}
