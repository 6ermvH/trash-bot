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
	users := strings.Split(update.Message.Text, " ")[1:]

	if err := t.service.SetEstablish(ctx, chatID, users); err != nil {
		log.Printf("SetEstablish. call service: %v", err.Error())
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
		log.Printf("Next. call service: %v", err.Error())
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
		log.Printf("Prev. call service: %v", err.Error())
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
		log.Printf("Who. call service: %v", err.Error())
	}

	if _, err := botApi.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   username,
	}); err != nil {
		log.Printf("Who. send message: %v", err.Error())
	}
}
