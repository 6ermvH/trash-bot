package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/6ermvH/trash-bot/internal/repository"
	"github.com/go-telegram/bot"
)

type Service interface {
	GetSubscribedChats(ctx context.Context) ([]repository.Chat, error)
	Who(ctx context.Context, chatID int64) (string, error)
}

type Scheduler struct {
	service Service
	botAPI  *bot.Bot
}

func New(service Service, botAPI *bot.Bot) *Scheduler {
	return &Scheduler{
		service: service,
		botAPI:  botAPI,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	// Выравниваем запуск на начало минуты
	now := time.Now()
	nextMinute := now.Truncate(time.Minute).Add(time.Minute)
	time.Sleep(time.Until(nextMinute))

	// Проверяем сразу после выравнивания
	s.checkAndNotify(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checkAndNotify(ctx)
		}
	}
}

func (s *Scheduler) checkAndNotify(ctx context.Context) {
	currentTime := time.Now().Format("15:04")

	chats, err := s.service.GetSubscribedChats(ctx)
	if err != nil {
		log.Printf("scheduler: get subscribed chats: %v", err)

		return
	}

	for _, chat := range chats {
		if chat.NotifyTime == nil {
			continue
		}

		if *chat.NotifyTime != currentTime {
			continue
		}

		s.sendNotification(ctx, chat.ID)
	}
}

func (s *Scheduler) sendNotification(ctx context.Context, chatID int64) {
	username, err := s.service.Who(ctx, chatID)
	if err != nil {
		log.Printf("scheduler: get who for chat %d: %v", chatID, err)

		return
	}

	_, err = s.botAPI.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "🗑 Напоминание: сегодня мусор выносит " + username,
	})
	if err != nil {
		log.Printf("scheduler: send notification to chat %d: %v", chatID, err)
	}
}
