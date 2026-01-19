package panel

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/6ermvH/trash-bot/internal/config"
	handlers "github.com/6ermvH/trash-bot/internal/handlers/http/v1"
	"github.com/6ermvH/trash-bot/internal/services/trashmanager"
	"github.com/gin-gonic/gin"
)

//go:embed all:web
var webFS embed.FS

func Start(ctx context.Context, cfg *config.Config, trashm *trashmanager.Service) error {
	router := gin.Default()

	// API routes
	api := router.Group("/api")

	authHandler := handlers.NewAuthHandler(
		cfg.Server.AdminLogin,
		cfg.Server.AdminPassword,
		cfg.Server.JWTSecret,
	)
	api.POST("/login", authHandler.Login)

	handle := handlers.New(trashm)

	protected := api.Group("/")
	protected.Use(handlers.AuthMiddleware(cfg.Server.JWTSecret))
	{
		protected.GET("/stats", handle.Stats)
		protected.GET("/chats", handle.Chats)
		protected.GET("/chats/:id", handle.ChatByID)
	}

	// Static files
	webContent, err := fs.Sub(webFS, "web")
	if err != nil {
		return fmt.Errorf("failed to get web fs: %w", err)
	}
	router.GET("/", func(c *gin.Context) {
		c.FileFromFS("index.html", http.FS(webContent))
	})
	router.GET("/style.css", func(c *gin.Context) {
		c.FileFromFS("style.css", http.FS(webContent))
	})
	router.GET("/app.js", func(c *gin.Context) {
		c.FileFromFS("app.js", http.FS(webContent))
	})

	port := ":" + cfg.Server.Port

	if err := router.Run(port); err != nil {
		return fmt.Errorf("Start server on port %s : %w", port, err)
	}

	return nil
}
