package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Stat(ctx context.Context) []int64
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

func (h *HandlerM) Stat(c *gin.Context) {
	ids := h.service.Stat(c.Request.Context())

	c.JSON(200, ids)
}
