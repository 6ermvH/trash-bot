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
	Subscribe(ctx context.Context, chatID int64, notifyTime string) error
	Unsubscribe(ctx context.Context, chatID int64) error
}

type HandlerM struct {
	service Service
}

func New(service Service) *HandlerM {
	return &HandlerM{service: service}
}

func (h *HandlerM) Chats(ctx *gin.Context) {
	chats, err := h.service.Chats(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load chats"})

		return
	}

	ctx.JSON(http.StatusOK, chats)
}

func (h *HandlerM) ChatByID(ctx *gin.Context) {
	idStr := ctx.Param("id")

	chatID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})

		return
	}

	chat, err := h.service.Chat(ctx.Request.Context(), chatID)
	if err != nil {
		if errors.Is(err, repository.ErrChatIsNotInitialize) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "chat not found"})

			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load chat"})

		return
	}

	ctx.JSON(http.StatusOK, chat)
}

func (h *HandlerM) Stats(ctx *gin.Context) {
	stats, err := h.service.Stats(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stats"})

		return
	}

	ctx.JSON(http.StatusOK, stats)
}
