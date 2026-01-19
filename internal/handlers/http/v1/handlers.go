package apiv1

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/6ermvH/trash-bot/internal/repository"
	"github.com/6ermvH/trash-bot/internal/services/trashmanager"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Chats(ctx context.Context) ([]repository.Chat, error)
	Chat(ctx context.Context, chatID int64) (*repository.Chat, error)
	Stats(ctx context.Context) (trashmanager.Stats, error)

	Who(ctx context.Context, chatID int64) (string, error)
	Next(ctx context.Context, chatID int64) (string, error)
	Prev(ctx context.Context, chatID int64) (string, error)
	SetEstablish(ctx context.Context, chatID int64, users []string) error
	Subscribe(ctx context.Context, chatID int64) error
	Unsubscribe(ctx context.Context, chatID int64) error
}

type HandlerM struct {
	service Service
}

func New(service Service) *HandlerM {
	return &HandlerM{service: service}
}

func (h *HandlerM) Chats(c *gin.Context) {
	chats, err := h.service.Chats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

func (h *HandlerM) ChatByID(c *gin.Context) {
	idStr := c.Param("id")
	chatID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
		return
	}

	chat, err := h.service.Chat(c.Request.Context(), chatID)
	if err != nil {
		if errors.Is(err, repository.ErrChatIsNotInitialize) {
			c.JSON(http.StatusNotFound, gin.H{"error": "chat not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load chat"})
		return
	}

	c.JSON(http.StatusOK, chat)
}

func (h *HandlerM) Stats(c *gin.Context) {
	stats, err := h.service.Stats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
