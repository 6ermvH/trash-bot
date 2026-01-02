package telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func (t *TgBotHandler) getKeyboardOnStart(b *bot.Bot) *inline.Keyboard {
	kb := inline.New(b).Row().
		Button("Кто выносит", []byte(""), func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			username, err := t.service.Who(ctx, mes.Message.Chat.ID)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: mes.Message.Chat.ID,
					Text:   "Во время выполнения запроса произошла ошибка",
				})
				return
			}

			text := fmt.Sprintf("Мусор выносит: %s", username)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text:   text,
			})
		}).Row().
		Button("Следующий", []byte(""), func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			username, err := t.service.Next(ctx, mes.Message.Chat.ID)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: mes.Message.Chat.ID,
					Text:   "Во время выполнения запроса произошла ошибка",
				})
				return
			}

			text := fmt.Sprintf("Мусор выносит: %s", username)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text:   text,
			})
		}).Row().
		Button("Предыдущий", []byte(""), func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
			username, err := t.service.Prev(ctx, mes.Message.Chat.ID)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: mes.Message.Chat.ID,
					Text:   "Во время выполнения запроса произошла ошибка",
				})
				return
			}

			text := fmt.Sprintf("Мусор выносит: %s", username)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: mes.Message.Chat.ID,
				Text:   text,
			})
		})

	return kb
}
