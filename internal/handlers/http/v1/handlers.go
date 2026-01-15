package apiv1

import (
	"context"
	"net/http"

	"github.com/6ermvH/trash-bot/internal/repository/inmemory"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Chats(ctx context.Context) []inmemory.Chat

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
	ids := h.service.Chats(c.Request.Context())

	c.JSON(http.StatusOK, ids)
}
