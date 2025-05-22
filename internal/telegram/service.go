package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/6ermvH/trash_bot/internal/store"
)

// Service holds bot API, store manager, HTTP client, and OpenRouter key
type Service struct {
	bot           *tgbotapi.BotAPI
	store         *store.Store
	httpClient    *http.Client
	openRouterKey string
}

// NewService creates a new Telegram service with dependencies
func NewService(token string, storeMgr *store.Store, openRouterKey string) *Service {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panicf("failed to create bot: %v", err)
	}
	bot.Debug = false
	return &Service{
		bot:           bot,
		store:         storeMgr,
		httpClient:    &http.Client{},
		openRouterKey: openRouterKey,
	}
}

// Run starts receiving updates and dispatching handlers
func (s *Service) Run(ctx context.Context) error {
	updates := s.bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			s.handleCommand(update.Message)
		} else if update.CallbackQuery != nil {
			s.handleCallback(update.CallbackQuery)
		}
	}
	return nil
}

// handleCommand processes incoming slash commands
func (s *Service) handleCommand(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	var text string

	switch msg.Command() {
	case "start":
		text = "Выберите действие:"
		s.sendMessage(chatID, text, commandKeyboard())
		return
	case "help":
		text = helpText()
	case "hey":
		text = "Exactly!!!"
	case "set_establish":
		users := strings.Fields(msg.CommandArguments())
		if err := s.store.SetEstablish(chatID, users); err != nil {
			text = "Ошибка при установке списка пользователей"
		} else {
			text = "Список пользователей обновлён"
		}
	case "next":
		name, err := s.store.Next(chatID)
		if err != nil {
			text = "Невозможно перейти на следующего"
		} else {
			text = "Теперь мусор выносит: " + name
		}
	case "prev":
		name, err := s.store.Prev(chatID)
		if err != nil {
			text = "Невозможно перейти на предыдущего"
		} else {
			text = "Теперь мусор выносит: " + name
		}
	case "who":
		name, err := s.store.Who(chatID)
		if err != nil {
			text = "Не удалось получить текущее имя"
		} else {
			text = "Мусор выносит: " + name
		}
	case "chat":
		text = s.handleChat(msg.CommandArguments())
	default:
		text = "Неизвестная команда. Используйте /help"
	}

	s.sendMessage(chatID, text, tgbotapi.InlineKeyboardMarkup{})
}

// handleChat sends user prompt to OpenRouter API and returns the response text
func (s *Service) handleChat(prompt string) string {
    payload := fmt.Sprintf(
			`{"model":"meta-llama/llama-3.3-8b-instruct:free","messages":[{"role":"user","content":"Ты чёрный браток, отвечай как реальный Homie, вот запрос:'%s'"}]}`,
        prompt,
    )
    req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", strings.NewReader(payload))
    if err != nil {
        return "Ошибка при формировании запроса к чату"
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+s.openRouterKey)

    resp, err := s.httpClient.Do(req)
    if err != nil {
        return "Ошибка при обращении к чату"
    }
    defer resp.Body.Close()

    var orResp struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
        return "Ошибка при разборе ответа чата"
    }
    if len(orResp.Choices) > 0 {
        return orResp.Choices[0].Message.Content
    }
    return "Пустой ответ от чата"
}
// handleCallback processes inline button callbacks
func (s *Service) handleCallback(cb *tgbotapi.CallbackQuery) {
	chatID := cb.Message.Chat.ID
	var text string

	switch cb.Data {
	case "who":
		name, err := s.store.Who(chatID)
		if err != nil {
			text = "Ошибка при получении пользователя"
		} else {
			text = "Мусор выносит: " + name
		}
	case "next":
		name, err := s.store.Next(chatID)
		if err != nil {
			text = "Ошибка при переходе к следующему"
		} else {
			text = "Теперь мусор выносит: " + name
		}
	default:
		text = ""
	}

	if _, err := s.bot.Request(tgbotapi.NewCallback(cb.ID, cb.Data)); err != nil {
		log.Printf("callback ack error: %v", err)
	}
	edit := tgbotapi.NewEditMessageText(chatID, cb.Message.MessageID, text)
	if _, err := s.bot.Send(edit); err != nil {
		log.Printf("message edit error: %v", err)
	}
}

// sendMessage wraps sending a text message with optional reply markup with optional reply markup
func (s *Service) sendMessage(chatID int64, text string, markup tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	if len(markup.InlineKeyboard) > 0 {
		msg.ReplyMarkup = &markup
	}
	if _, err := s.bot.Send(msg); err != nil {
		log.Printf("send message error: %v", err)
	}
}

// commandKeyboard builds the inline keyboard for /start
func commandKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Кто выносит мусор", "who"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вынес мусор", "next"),
		),
	)
}

// helpText returns help message string
func helpText() string {
	return "/help - показать помощь\n" +
		"/start - начать работу бота\n" +
		"/set_establish [users] - задать список пользователей\n" +
		"/who - кто сейчас\n" +
		"/next - следующий пользователь\n" +
		"/prev - предыдущий пользователь\n" +
		"/chat [text] - общение через AI\n" +
		"/hey - ответ от бота\n"
}
