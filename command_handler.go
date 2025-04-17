package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Кто выносит мусор", "who"),
		tgbotapi.NewInlineKeyboardButtonData("Вынес мусор", "next"),
	),
)

func InitTelegramAPI() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil && update.Message.IsCommand() {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			switch update.Message.Command() {
			case "start":
				msg.ReplyMarkup = numericKeyboard
			case "prev":
				msg.Text = "Теперь мусор выносит: " + prev(update.Message.Chat.ID)
			case "next":
				msg.Text = "Теперь мусор выносит: " + next(update.Message.Chat.ID)
			case "set_establish":
				usersString := update.Message.CommandArguments()
				users := strings.Fields(usersString)
				setEstablish(update.Message.Chat.ID, users)
				msg.Text = "Список пользователей обновлен"
			case "help":
				var builder strings.Builder
				builder.WriteString("/help просмотреть все возможности\n")
				builder.WriteString("/start запустить меню бота\n")
				builder.WriteString("-> Команды для отладки <-\n")
				builder.WriteString("/set_establish [пользователи] задать пользователей\n")
				builder.WriteString("/next переход на следующего\n")
				builder.WriteString("/prev переход на предыдущего")
				msg.Text = builder.String()

			default:
				msg.Text = "I don't know that command"
			}

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}

		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			data := ""
			if update.CallbackQuery.Data == "who" {
				data = "Мусор выносит: " + who(update.CallbackQuery.Message.Chat.ID)
			} else if update.CallbackQuery.Data == "next" {
				data = "Теперь мусор выносит: " + next(update.CallbackQuery.Message.Chat.ID)
			}

			edit := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				data,
			)

			bot.Send(edit)
		}
	}
}
