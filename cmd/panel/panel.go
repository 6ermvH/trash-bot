package panel

import (
	"context"
	"fmt"

	"github.com/6ermvH/trash-bot/internal/config"
	handlers "github.com/6ermvH/trash-bot/internal/handlers/http/v1"
	"github.com/6ermvH/trash-bot/internal/services/trashmanager"
	"github.com/gin-gonic/gin"
)

func Start(ctx context.Context, cfg *config.Config, trashm *trashmanager.Service) error {
	router := gin.Default()

	authHandler := handlers.NewAuthHandler(
		cfg.Server.AdminLogin,
		cfg.Server.AdminPassword,
		cfg.Server.JWTSecret,
	)
	router.POST("/login", authHandler.Login)

	handle := handlers.New(trashm)

	protected := router.Group("/")
	protected.Use(handlers.AuthMiddleware(cfg.Server.JWTSecret))
	{
		protected.GET("/stats", handle.Stats)
		protected.GET("/chats", handle.Chats)
		protected.GET("/chats/:id", handle.ChatByID)
	}

	port := ":" + cfg.Server.Port

	if err := router.Run(port); err != nil {
		return fmt.Errorf("Start server on port %s : %w", port, err)
	}

	return nil
}
