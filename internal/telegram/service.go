package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab.com/6ermvH/trash_bot/internal/store"
)

// --- Constants for user-facing messages and API details ---

const (
	// --- Command-related text ---
	msgStart           = "Выберите действие:"
	msgHey             = "Exactly!!!"
	msgEstablishSet    = "Список пользователей обновлён"
	msgEstablishErr    = "Ошибка при установке списка пользователей"
	msgUserAdded       = "Пользователь %s добавлен в список"
	msgUserAddErr      = "Ошибка при добавлении: %v"
	msgUserAddUsage    = "Укажите имя пользователя. Пример: /add_user @nickname"
	msgUserRemoved     = "Пользователь %s удалён из списка"
	msgUserRemoveErr   = "Ошибка при удалении: %v"
	msgUserRemoveUsage = "Укажите имя пользователя. Пример: /remove_user @nickname"
	msgNext            = "Теперь мусор выносит: %s"
	msgNextErr         = "Невозможно перейти на следующего"
	msgPrev            = "Теперь мусор выносит: %s"
	msgPrevErr         = "Невозможно перейти на предыдущего"
	msgWho             = "Мусор выносит: %s"
	msgWhoErr          = "Не удалось получить текущее имя"
	msgUnknownCmd      = "Неизвестная команда. Используйте /help"

	// --- Callback-related text ---
	msgCallbackWhoErr  = "Ошибка при получении пользователя"
	msgCallbackNextErr = "Ошибка при переходе к следующему"

	// --- Chat-related text ---
	msgChatDisabled    = "Ключ OpenRouter API не настроен."
	msgChatRequestErr  = "Ошибка при формировании запроса к чату"
	msgChatClientErr   = "Ошибка при обращении к чату"
	msgChatDecodeErr   = "Ошибка при разборе ответа чата"
	msgChatEmptyResp   = "Пустой ответ от чата"
	msgChatPayloadErr  = "Ошибка при подготовке запроса к чату"

	// --- OpenRouter API ---
	openRouterAPIURL     = "https://openrouter.ai/api/v1/chat/completions"
	openRouterModel      = "meta-llama/llama-3.3-8b-instruct:free"
	openRouterPromptTmpl = "Ты чёрный браток, отвечай как реальный Homie, вот запрос:'%s'"
)

// --- Structs for OpenRouter API ---

type openRouterRequest struct {
	Model    string          `json:"model"`
	Messages []openRouterMsg `json:"messages"`
}

type openRouterMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Service holds bot API, store manager, HTTP client, and OpenRouter key
type Service struct {
	logger        *slog.Logger
	bot           *tgbotapi.BotAPI
	store         *store.Store
	httpClient    *http.Client
	openRouterKey string
}

// NewService creates a new Telegram service with dependencies
func NewService(logger *slog.Logger, token string, storeMgr *store.Store, openRouterKey string) *Service {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Error("failed to create bot", "error", err)
		os.Exit(1)
	}
	bot.Debug = false
	return &Service{
		logger:        logger,
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

// --- Command Handlers ---

// handleCommand processes incoming slash commands by dispatching to specific handlers
func (s *Service) handleCommand(msg *tgbotapi.Message) {
	s.logger.Info("handling command", "command", msg.Command(), "from", msg.From.UserName, "chat_id", msg.Chat.ID)
	var text string
	var markup interface{} // Can be keyboard or nil

	switch msg.Command() {
	case "start":
		text = msgStart
		markup = commandKeyboard()
	case "help":
		text = helpText()
	case "hey":
		text = msgHey
	case "set_establish":
		text = s.handleSetEstablish(msg)
	case "add_user":
		text = s.handleAddUser(msg)
	case "remove_user":
		text = s.handleRemoveUser(msg)
	case "next":
		text = s.handleNext(msg)
	case "prev":
		text = s.handlePrev(msg)
	case "who":
		text = s.handleWho(msg)
	case "chat":
		text = s.handleChat(msg.CommandArguments())
	default:
		text = msgUnknownCmd
	}

	s.sendMessage(msg.Chat.ID, text, markup)
}

func (s *Service) handleSetEstablish(msg *tgbotapi.Message) string {
	users := strings.Fields(msg.CommandArguments())
	if err := s.store.SetEstablish(msg.Chat.ID, users); err != nil {
		return msgEstablishErr
	}
	return msgEstablishSet
}

func (s *Service) handleAddUser(msg *tgbotapi.Message) string {
	user := msg.CommandArguments()
	if user == "" {
		return msgUserAddUsage
	}
	if err := s.store.AddUser(msg.Chat.ID, user); err != nil {
		return fmt.Sprintf(msgUserAddErr, err)
	}
	return fmt.Sprintf(msgUserAdded, user)
}

func (s *Service) handleRemoveUser(msg *tgbotapi.Message) string {
	user := msg.CommandArguments()
	if user == "" {
		return msgUserRemoveUsage
	}
	if err := s.store.RemoveUser(msg.Chat.ID, user); err != nil {
		return fmt.Sprintf(msgUserRemoveErr, err)
	}
	return fmt.Sprintf(msgUserRemoved, user)
}

func (s *Service) handleNext(msg *tgbotapi.Message) string {
	name, err := s.store.Next(msg.Chat.ID)
	if err != nil {
		return msgNextErr
	}
	return fmt.Sprintf(msgNext, name)
}

func (s *Service) handlePrev(msg *tgbotapi.Message) string {
	name, err := s.store.Prev(msg.Chat.ID)
	if err != nil {
		return msgPrevErr
	}
	return fmt.Sprintf(msgPrev, name)
}

func (s *Service) handleWho(msg *tgbotapi.Message) string {
	name, err := s.store.Who(msg.Chat.ID)
	if err != nil {
		return msgWhoErr
	}
	return fmt.Sprintf(msgWho, name)
}

// --- Chat Handler ---

// handleChat sends user prompt to OpenRouter API and returns the response text
func (s *Service) handleChat(prompt string) string {
	if s.openRouterKey == "" {
		return msgChatDisabled
	}

	requestBody := openRouterRequest{
		Model: openRouterModel,
		Messages: []openRouterMsg{
			{
				Role:    "user",
				Content: fmt.Sprintf(openRouterPromptTmpl, prompt),
			},
		},
	}

	payloadBytes, err := json.Marshal(requestBody)
	if err != nil {
		s.logger.Error("chat payload marshal error", "error", err)
		return msgChatPayloadErr
	}

	req, err := http.NewRequest("POST", openRouterAPIURL, bytes.NewReader(payloadBytes))
	if err != nil {
		s.logger.Error("chat http request create error", "error", err)
		return msgChatRequestErr
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.openRouterKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("chat http client error", "error", err)
		return msgChatClientErr
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Error("response body close error", "error", err)
		}
	}()

	var orResp openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&orResp); err != nil {
		s.logger.Error("chat response decode error", "error", err)
		return msgChatDecodeErr
	}
	if len(orResp.Choices) > 0 && orResp.Choices[0].Message.Content != "" {
		return orResp.Choices[0].Message.Content
	}
	return msgChatEmptyResp
}

// --- Callback Handler ---

// handleCallback processes inline button callbacks
func (s *Service) handleCallback(cb *tgbotapi.CallbackQuery) {
	s.logger.Info("handling callback", "data", cb.Data, "from", cb.From.UserName, "chat_id", cb.Message.Chat.ID)
	var text string

	switch cb.Data {
	case "who":
		name, err := s.store.Who(cb.Message.Chat.ID)
		if err != nil {
			text = msgCallbackWhoErr
		} else {
			text = fmt.Sprintf(msgWho, name)
		}
	case "next":
		name, err := s.store.Next(cb.Message.Chat.ID)
		if err != nil {
			text = msgCallbackNextErr
		} else {
			text = fmt.Sprintf(msgNext, name)
		}
	default:
		// Do nothing for unknown callbacks, just acknowledge
		if _, err := s.bot.Request(tgbotapi.NewCallback(cb.ID, "")); err != nil {
			s.logger.Error("callback ack error", "error", err)
		}
		return
	}

	// Acknowledge the callback
	if _, err := s.bot.Request(tgbotapi.NewCallback(cb.ID, cb.Data)); err != nil {
		s.logger.Error("callback ack error", "error", err)
	}

	// Edit the message with the result
	if text != "" {
		edit := tgbotapi.NewEditMessageText(cb.Message.Chat.ID, cb.Message.MessageID, text)
		if _, err := s.bot.Send(edit); err != nil {
			s.logger.Error("message edit error", "error", err)
		}
	}
}

// --- Helpers ---

// sendMessage wraps sending a text message with optional reply markup
func (s *Service) sendMessage(chatID int64, text string, markup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	if markup != nil {
		msg.ReplyMarkup = markup
	}
	s.logger.Info("sending message", "chat_id", chatID, "text", text)
	if _, err := s.bot.Send(msg); err != nil {
		s.logger.Error("send message error", "error", err)
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
	return `/help - показать помощь
/start - начать работу бота
/set_establish [users] - задать/перезаписать список пользователей
/add_user [user] - добавить пользователя в конец списка
/remove_user [user] - удалить пользователя из списка
/who - кто сейчас
/next - следующий пользователь
/prev - предыдущий пользователь
/chat [text] - общение через AI
/hey - ответ от бота
`
}
